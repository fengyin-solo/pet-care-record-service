package adapter

import (
	"sync"

	"petsmanagement/internal/model"
)

type AuditEntry struct{ OwnerID, PetID, Label string }

type RequestAudit struct {
	started chan struct{}
	release chan struct{}
	mu      sync.Mutex
	entries []AuditEntry
}

func NewRequestAudit() *RequestAudit {
	return &RequestAudit{started: make(chan struct{}, 8), release: make(chan struct{})}
}

func (a *RequestAudit) Record(scope *model.RequestScope) {
	a.started <- struct{}{}
	<-a.release
	entry := AuditEntry{OwnerID: scope.OwnerID, PetID: scope.PetID}
	if len(scope.Labels) > 0 {
		entry.Label = scope.Labels[0]
	}
	a.mu.Lock()
	a.entries = append(a.entries, entry)
	a.mu.Unlock()
}

func (a *RequestAudit) Started() <-chan struct{} { return a.started }
func (a *RequestAudit) Release()                 { close(a.release) }
func (a *RequestAudit) Entries() []AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]AuditEntry(nil), a.entries...)
}
