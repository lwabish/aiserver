package controllers

import (
	"context"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/handlers/sadtalker"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	batchv1 "k8s.io/api/batch/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type SadTalkerJobReconciler struct {
	client.Client
	*handlers.BaseHandler
}

func (r *SadTalkerJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	job := &batchv1.Job{}
	if err := r.Get(ctx, req.NamespacedName, job); err != nil {
		return ctrl.Result{}, err
	}

	if jobSucceed(job) {
		r.UpdateTaskStatus(job.Name, models.TaskStatusSuccess)
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

func jobSucceed(job *batchv1.Job) bool {
	complete := false
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete {
			complete = true
		}
	}
	return complete && *job.Spec.Completions == job.Status.Succeeded
}
