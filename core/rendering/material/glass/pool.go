package glass

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