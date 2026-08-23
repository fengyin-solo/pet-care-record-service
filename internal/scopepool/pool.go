package scopepool

import "petsmanagement/internal/model"

type Pool struct{ items chan *model.RequestScope }

func New() *Pool {
	p := &Pool{items: make(chan *model.RequestScope, 1)}
	p.items <- &model.RequestScope{}
	return p
}

func (p *Pool) Get() *model.RequestScope      { return <-p.items }
func (p *Pool) Put(scope *model.RequestScope) { p.items <- scope }
