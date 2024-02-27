package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/models"
	"os/exec"
	"regexp"
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
			s.L.Errorf("Process sad talker task failed:%s %s", task.Uid, err.Error())
			s.UpdateTaskStatus(task.Uid, models.TaskStatusFailed)
		}
	}()
	p := s.getParam(task.Uid)
	s.L.Infof("Processing sad talker task:%s %s", task.Uid, p)
	err = s.workerFunc(task, p)
}

func (s *handler) invokeSadTalker(task *models.Task, p *taskParam) error {
	cmd := exec.Command(s.pythonPath, "inference.py", "--driven_audio", p.audio, "--source_image", p.photo)
	if err := cmd.Run(); err != nil {
		return err
	}
	stdout, err := cmd.Output()
	if err != nil {
		s.UpdateTaskStatus(task.Uid, models.TaskStatusResultMissing)
	}

	if result := ParseResult(string(stdout)); result != "" {
		s.SaveTaskResult(task.Uid, result)
	} else {
		s.UpdateTaskStatus(task.Uid, models.TaskStatusResultMissing)
	}
	return nil
}

func ParseResult(log string) string {
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
