package adapter

import "errors"

var ErrPublisherBusy = errors.New("health publisher busy")

// Publisher 是健康状态发布器的抽象接口，便于在 service 层替换测试替身。
type Publisher interface {
	Publish(key string) error
}

// HealthPublisher 模拟健康状态发布器。
// effects 以 key 为维度做幂等：同一稳定标识重复到达下游只产生一次副作用。
type HealthPublisher struct {
	failFirst bool
	calls     int
	effects   map[string]int
}

func NewHealthPublisher(failFirst bool) *HealthPublisher {
	return &HealthPublisher{failFirst: failFirst, effects: make(map[string]int)}
}

func (p *HealthPublisher) Publish(key string) error {
	p.calls++
	if p.failFirst && p.calls == 1 {
		return ErrPublisherBusy
	}
	// 幂等：同一 key 已发布过则直接返回成功，不重复产生副作用。
	if p.effects[key] > 0 {
		return nil
	}
	p.effects[key] = 1
	return nil
}

func (p *HealthPublisher) Calls() int { return p.calls }
func (p *HealthPublisher) Effects() int {
	total := 0
	for _, count := range p.effects {
		total += count
	}
	return total
}
