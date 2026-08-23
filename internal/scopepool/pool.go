package scopepool

import "petsmanagement/internal/model"

// Pool 复用 RequestScope 对象以减少分配。
//
// Put 回收时调用 Reset 清空 scope 并撤销所有权，使回收对象不再
// 残留上一任调用方的身份；下一次 Get 拿到的 scope 必须经过
// Prepare 交接所有权才能被审计写入。
type Pool struct{ items chan *model.RequestScope }

// New 创建容量为 1 的对象池，预置一个空 scope。
func New() *Pool {
	p := &Pool{items: make(chan *model.RequestScope, 1)}
	p.items <- &model.RequestScope{}
	return p
}

// Get 取出一个 scope。取出后必须经 Prepare 交接所有权才能用于审计。
func (p *Pool) Get() *model.RequestScope { return <-p.items }

// Put 回收 scope。回收前先 Reset，清除上一任身份并撤销所有权。
func (p *Pool) Put(scope *model.RequestScope) {
	scope.Reset()
	p.items <- scope
}
