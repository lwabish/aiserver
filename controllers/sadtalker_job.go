package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/handlers/sadtalker"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	batchv1 "k8s.io/api/batch/v1"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type SadTalkerJobReconciler struct {
	client.Client
	*handlers.BaseHandler
	logr.Logger
}

func (r *SadTalkerJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	job := &batchv1.Job{}
	if err := r.Get(ctx, req.NamespacedName, job); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !jobSucceeded(job) {
		return ctrl.Result{}, nil
	}
	r.Info("Job succeeded", "uid", job.Name)
	r.UpdateTaskStatus(job.Name, models.TaskStatusSuccess)

	logString, err := getPodLogString(r.C, job.Name, job.Namespace)
	if err != nil {
		r.UpdateTaskStatus(job.Name, models.TaskStatusResultMissing)
		return ctrl.Result{}, err
	}

	if result := extractResult(logString); result != "" {
		r.SaveTaskResult(job.Name, result)
	} else {
		r.UpdateTaskStatus(job.Name, models.TaskStatusResultMissing)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SadTalkerJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	p1 := predicate.NewPredicateFuncs(func(object client.Object) bool {
		if job, ok := object.(*batchv1.Job); ok && filterSadTalkerJob(job) {
			return true
		}
		return false
	})
	p2 := predicate.Funcs{
		CreateFunc: func(event event.CreateEvent) bool {
			return false
		},
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			return true
		},
		GenericFunc: func(genericEvent event.GenericEvent) bool {
			return false
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(predicate.And(p1, p2)).
		For(&batchv1.Job{}).
		Complete(r)
}

func filterSadTalkerJob(job *batchv1.Job) bool {
	return job.Annotations[utils.TaskTypeKey] == sadtalker.JobType
}

func extractResult(log string) string {
	re := regexp.MustCompile(`\./results/\d{4}_\d{2}_\d{2}_\d{2}\.\d{2}\.\d{2}\.mp4\n`)
	match := re.FindString(log)
	if match == "" {
		return match
	}
	// 去除匹配字符串两端的换行符和"./results/"
	result := match[:len(match)-1] // 移除尾部的换行符
	result = result[10:]           // 移除前面的"./results/"
	return result
}
