package progress

import (
	"fmt"
)

type Sink interface {
	Start(resolution int)
	End()
	Tick()
}

type stdout struct {
	resolution, progress int
	step, threshold float32
}

func NewStdout() *stdout {
	return &stdout{}
}

func (s *stdout) Start(resolution int) {
	s.resolution = resolution
	s.progress = 0
	s.step = 10.0
	s.threshold = s.step

//	fmt.Println("0%")
}

func (s *stdout) End() {
	fmt.Printf("\n")
}

func (s *stdout) Tick() {
	if (s.progress >= s.resolution) {
		return
	}

	s.progress++

	p := (float32(s.progress) / float32(s.resolution)) * 100.0

	if p >= s.threshold {
		s.threshold += s.step 

		fmt.Printf("%d%%... ", int(p))		
	}
}
