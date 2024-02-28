package openvoice

import (
	"fmt"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"os"
	"os/exec"
	"path"
	"time"
)

type taskParam struct {
	text string
	// image with face
	audioPath string
}

func (p taskParam) String() string {
	return p.text + "|" + p.audioPath
}

func (h *handler) invoke(task *models.Task, tp handlers.TaskParam) error {
	var p *taskParam
	var ok bool
	if p, ok = tp.(*taskParam); !ok {
		return fmt.Errorf("task param type error")
	}
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}

	resultFileName := fmt.Sprintf("%s%s", time.Now().Format("2006_01_02_15.04.05"), ".wav")
	args := []string{
		"run.py",
		"--text",
		p.text,
		"--audio",
		path.Join(curDir, p.audioPath),
		"--output",
		path.Join(curDir, fmt.Sprintf("%s/%s", utils.ResultDir, resultFileName)),
	}
	h.L.Debugf("openvoice command: %s %v %v", h.pythonPath, args, h.extraArgs)
	cmd := exec.Command(h.pythonPath, append(args, h.extraArgs...)...)
	cmd.Dir = h.projectPath

	var output []byte
	if output, err = cmd.Output(); err != nil {
		return err
	}

	h.L.Debugf("openvoice stdout: %s", output)
	h.SaveTaskResult(task.Uid, resultFileName)
	h.UpdateTaskStatus(task.Uid, models.TaskStatusSuccess)
	return nil
}
