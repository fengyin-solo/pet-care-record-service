package adapter

import (
	"sync"

	"petsmanagement/internal/model"
)

// AuditEntry 是已落库的审计条目，OwnerID/PetID/Label 各自独立，
// 不会被后续 scope 的 Prepare/Reset 影响。
type AuditEntry struct{ OwnerID, PetID, Label string }

// RequestAudit 协调审计写入的启停与记录。
//
// Record 在阻塞之前先拷贝 scope 的身份快照，因此调用方在
// <-Started 之后 Prepare 下一任主人/宠物，不会回头改写本次条目；
// 未通过 Prepare 交接所有权（Owned=false）的 scope 会被拒绝写入，
// 但仍正常走 Started/Release 流程，保证通道语义不变。
type RequestAudit struct {
	started chan struct{}
	release chan struct{}
	mu      sync.Mutex
	entries []AuditEntry
}

// NewRequestAudit 创建审计器。
//
// started 缓冲为 8，release 无缓冲：每次 Record 先发 started 通知
// 调用方可以继续（如 Put 回收 scope），再等待 Release 放行写入。
func NewRequestAudit() *RequestAudit {
	return &RequestAudit{
		started: make(chan struct{}, 8),
		release: make(chan struct{}),
	}
}

// Record 记录一次审计查询。
//
// 先拷贝 scope 的身份快照，再通过 started/release 与调用方同步。
// 这样即使调用方在 <-started 后立即 Prepare/reset 同一 scope 对象，
// 本次记录的内容也不会被篡改。Owned=false 的 scope 会被跳过写入，
// 但仍参与 started/release 握手，维持通道语义。
func (a *RequestAudit) Record(scope *model.RequestScope) {
	// 在阻塞前拷贝身份快照，避免后续 Prepare/Reset 改写本次条目。
	ownerID := scope.OwnerID
	petID := scope.PetID
	label := ""
	if len(scope.Labels) > 0 {
		label = scope.Labels[0]
	}
	owned := scope.Owned

	a.started <- struct{}{}
	<-a.release

	if !owned {
		// 未交接所有权的 scope 拒绝写入，但已走完 Started/Release。
		return
	}
	entry := AuditEntry{OwnerID: ownerID, PetID: petID, Label: label}
	a.mu.Lock()
	a.entries = append(a.entries, entry)
	a.mu.Unlock()
}

// Started 返回 started 通道，调用方借此得知 Record 已就绪。
func (a *RequestAudit) Started() <-chan struct{} { return a.started }

// Release 放行所有等待中的 Record 完成写入。
func (a *RequestAudit) Release() { close(a.release) }

// Entries 返回审计条目的拷贝。
func (a *RequestAudit) Entries() []AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]AuditEntry(nil), a.entries...)
}
