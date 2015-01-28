package glass

import (
	
)

type stack struct {
	samples []*Sample
	top uint32
}

type BinnedStack struct {
	stacks []stack
}

func NewBinnedStack(numBins uint32) *BinnedStack {
	s := BinnedStack{}

	s.stacks = make([]stack, numBins)

	for i := range s.stacks {
		stack := &s.stacks[i]
		stack.samples = make([]*Sample, 4)

		for j := range stack.samples {
			stack.samples[j] = new(Sample)
		}
	}

	return &s
}    

func (s *BinnedStack) Pop(binId uint32) *Sample {
	stack := &s.stacks[binId]
    sample := stack.samples[stack.top]
    stack.top++
    return sample
}

func (s *BinnedStack) Push(binId uint32) {
    s.stacks[binId].top--
}

