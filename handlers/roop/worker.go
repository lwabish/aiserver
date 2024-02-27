package roop

import (
	"github.com/lwabish/cloudnative-ai-server/models"
	"os"
	"os/exec"
	"path"
)

type taskParam struct {
	// image with face
	source string

	// image/video to be replaced
	target string
}

func (p taskParam) String() string {
	return p.source + "|" + p.target
}

func (h *handler) invoke(task *models.Task, p *taskParam) error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	args := []string{
		"run.py",
		"-s",
		path.Join(curDir, p.source),
		"-t",
		path.Join(curDir, p.target),
		"-o",
		path.Join(curDir, "results"),
	}
	h.L.Debugf("roop command: %s %v %v", h.pythonPath, args, h.extraArgs)
	cmd := exec.Command(h.pythonPath, append(args, h.extraArgs...)...)
	cmd.Dir = h.projectPath

	var output []byte
	if output, err = cmd.Output(); err != nil {
		return err
	}

	h.L.Debugf("roop stdout: %s", output)
	//if result := ParseResult(string(output)); result != "" {
	//	s.SaveTaskResult(task.Uid, result)
	//	s.UpdateTaskStatus(task.Uid, models.TaskStatusSuccess)
	//} else {
	//	s.L.Warnf("sadtalker result not found: %s", task.Uid)
	//	s.UpdateTaskStatus(task.Uid, models.TaskStatusResultMissing)
	//}
	return nil
}

func (h *handler) getParam(uid string) *taskParam {
	h.Lock()
	defer h.Unlock()
	return h.workerParam[uid]
}

func (h *handler) setParam(uid string, param *taskParam) {
	h.Lock()
	defer h.Unlock()
	h.workerParam[uid] = param
}
