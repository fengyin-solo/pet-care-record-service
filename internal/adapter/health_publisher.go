package adapter

import "errors"

var ErrPublisherBusy = errors.New("health publisher busy")

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
	p.effects[key]++
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
