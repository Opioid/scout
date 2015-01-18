package substitute

type Pool struct {
    pool chan *Sample
}

func NewPool(max int) *Pool {
    p := Pool{make(chan *Sample, max)}
    return &p
}    

func (p *Pool) Borrow() *Sample {
    var s *Sample

    select {
    case s = <- p.pool:
    default:
        s = new(Sample)
    }

    return s
}

func (p *Pool) Return(s *Sample) {
    select {
    case p.pool <- s:
    default:
        // let it go
    }
}