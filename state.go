package docker

import (
	"fmt"
	"github.com/dotcloud/docker/api"
	"github.com/dotcloud/docker/utils"
	"sync"
	"time"
)

type State struct {
	sync.Mutex
	api.State
}

// String returns a human-readable description of the state
func (s *State) String() string {
	if s.Running {
		if s.Ghost {
			return fmt.Sprintf("Ghost")
		}
		return fmt.Sprintf("Up %s", utils.HumanDuration(time.Now().Sub(s.StartedAt)))
	}
	return fmt.Sprintf("Exit %d", s.ExitCode)
}

func (s *State) setRunning(pid int) {
	s.Running = true
	s.Ghost = false
	s.ExitCode = 0
	s.Pid = pid
	s.StartedAt = time.Now()
}

func (s *State) setStopped(exitCode int) {
	s.Running = false
	s.Pid = 0
	s.ExitCode = exitCode
}
