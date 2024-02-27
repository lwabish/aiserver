package roop

import (
	"fmt"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
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

	ext := filepath.Ext(p.target)
	resultFileName := fmt.Sprintf("%s%s", time.Now().Format("2006_01_02_15.04.05"), ext)
	args := []string{
		"run.py",
		"-s",
		path.Join(curDir, p.source),
		"-t",
		path.Join(curDir, p.target),
		"-o",
		path.Join(curDir, fmt.Sprintf("%s/%s", utils.ResultDir, resultFileName)),
	}
	h.L.Debugf("roop command: %s %v %v", h.pythonPath, args, h.extraArgs)
	cmd := exec.Command(h.pythonPath, append(args, h.extraArgs...)...)
	cmd.Dir = h.projectPath

	var output []byte
	if output, err = cmd.Output(); err != nil {
		return err
	}

	h.L.Debugf("roop stdout: %s", output)
	h.SaveTaskResult(task.Uid, resultFileName)
	h.UpdateTaskStatus(task.Uid, models.TaskStatusSuccess)
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
