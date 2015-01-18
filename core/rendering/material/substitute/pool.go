package substitute

import (
    "sync"
)

type stack struct {
    samples []*Sample
    top uint32
}

type Pool struct {
    pool sync.Pool

//    stacks []stack
}

func NewPool(numWorkers uint32) *Pool {
    p := Pool{}
    p.pool.New = NewSample

/*
    p.stacks = make([]stack, numWorkers)

    for i := range p.stacks {
        p.stacks[i].samples = make([]*Sample, 4)

        for j := range p.stacks[i].samples {
            p.stacks[i].samples[j] = new(Sample)
        }
    }
*/

    return &p
}    

func (p *Pool) Get(workerId uint32) *Sample {
    s := p.pool.Get() 
    return s.(*Sample)

//    s := p.stacks[workerId].samples[p.stacks[workerId].top]
//    p.stacks[workerId].top++
//    return s
}

func (p *Pool) Put(s interface{}, workerId uint32) {
    p.pool.Put(s)

//    p.stacks[workerId].top--
}

func NewSample() interface{} {
    return new(Sample)
}