package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/models"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type taskParam struct {
	photo string
	audio string
}

func (p taskParam) String() string {
	return p.photo + "|" + p.audio
}

func (s *handler) getParam(uid string) *taskParam {
	s.Lock()
	defer s.Unlock()
	return s.workerParam[uid]
}

func (s *handler) setParam(uid string, param *taskParam) {
	s.Lock()
	defer s.Unlock()
	s.workerParam[uid] = param
}

func (s *handler) Process(task *models.Task) {
	s.UpdateTaskStatus(task.Uid, models.TaskStatusRunning)
	var err error
	defer func() {
		if err != nil {
			s.L.Errorf("Process sad talker task failed: %s %s", task.Uid, err.Error())
			s.UpdateTaskStatus(task.Uid, models.TaskStatusFailed)
		}
	}()
	p := s.getParam(task.Uid)
	s.L.Infof("Processing sad talker task: %s %s", task.Uid, p)
	err = s.workerFunc(task, p)
}

func (s *handler) invokeSadTalker(task *models.Task, p *taskParam) error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	args := []string{
		"inference.py",
		"--driven_audio",
		path.Join(curDir, p.audio),
		"--source_image",
		path.Join(curDir, p.photo),
		"--result_dir",
		path.Join(curDir, "results"),
	}
	s.L.Debugf("sadtalker command: %s %v %v", s.pythonPath, args, s.extraArgs)
	cmd := exec.Command(s.pythonPath, append(args, s.extraArgs...)...)
	cmd.Dir = s.projectPath

	var output []byte
	if output, err = cmd.Output(); err != nil {
		return err
	}

	s.L.Debugf("sadtalker stdout: %s", output)
	if result := ParseResult(string(output)); result != "" {
		s.SaveTaskResult(task.Uid, result)
		s.UpdateTaskStatus(task.Uid, models.TaskStatusSuccess)
	} else {
		s.L.Warnf("sadtalker result not found: %s", task.Uid)
		s.UpdateTaskStatus(task.Uid, models.TaskStatusResultMissing)
	}
	return nil
}

func ParseResult(log string) string {
	re := regexp.MustCompile(`/results/\d{4}_\d{2}_\d{2}_\d{2}\.\d{2}\.\d{2}\.mp4\n`)
	match := re.FindString(log)
	if match == "" {
		return match
	}
	result := match[:len(match)-1] // 移除尾部的换行符
	splits := strings.Split(result, "/")
	return splits[len(splits)-1]
}
