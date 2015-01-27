package substitute
/*
import (
    "sync"
)

type Pool struct {
	pool sync.Pool
}

func NewPool(numWorkers uint32) *Pool {
    p := Pool{}
	p.pool.New = NewSample

    return &p
}    

func (p *Pool) Get(workerId uint32) *Sample {
	s := p.pool.Get() 
	return s.(*Sample)
}

func (p *Pool) Put(s interface{}, workerId uint32) {
	p.pool.Put(s)
}

func NewSample() interface{} {
    return new(Sample)
}
*/

import (
	"github.com/Opioid/scout/core/rendering/material"
)

type stack struct {
    samples []*Sample
    top uint32
}

type Pool struct {
    stacks []stack
}

func NewPool(numWorkers uint32) *Pool {
    p := Pool{}

    p.stacks = make([]stack, numWorkers)

    for i := range p.stacks {
        p.stacks[i].samples = make([]*Sample, 4)

        for j := range p.stacks[i].samples {
            p.stacks[i].samples[j] = new(Sample)
        }
    }

    return &p
}    

func (p *Pool) Get(workerId uint32) *Sample {
	stack := &p.stacks[workerId]

    s := stack.samples[stack.top]
    stack.top++
    return s
}

func (p *Pool) Put(s material.Sample, workerId uint32) {
    p.stacks[workerId].top--
}