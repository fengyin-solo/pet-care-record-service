package scopeaudit

import (
	"sort"
	"testing"
	"time"

	"petsmanagement/internal/adapter"
	"petsmanagement/internal/model"
	"petsmanagement/internal/scopepool"
	. "petsmanagement/internal/service"
	"petsmanagement/internal/store"
)

func TestAsyncAuditKeepsEachPetLookupIdentity(t *testing.T) {
	svc := New(store.NewMemoryStore(), nil, nil)
	pool := scopepool.New()
	audit := adapter.NewRequestAudit()
	svc.AuditPetLookup("owner-A", "pet-A", []string{"calm"}, pool, audit)
	svc.AuditPetLookup("owner-B", "pet-B", []string{"alert"}, pool, audit)
	audit.Release()
	deadline := time.Now().Add(time.Second)
	for len(audit.Entries()) < 2 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	entries := audit.Entries()
	if len(entries) != 2 {
		t.Fatalf("audit did not finish both requests: %+v", entries)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].OwnerID < entries[j].OwnerID })
	if entries[0].OwnerID != "owner-A" || entries[0].PetID != "pet-A" || entries[0].Label != "calm" {
		t.Fatalf("first lookup inherited the next request identity: %+v", entries)
	}
	if entries[1].OwnerID != "owner-B" || entries[1].PetID != "pet-B" || entries[1].Label != "alert" {
		t.Fatalf("second lookup audit was corrupted: %+v", entries)
	}

	labels := []string{"original"}
	scope := &model.RequestScope{}
	scope.Prepare("owner-C", "pet-C", labels)
	labels[0] = "reused"
	if scope.Labels[0] != "original" {
		t.Fatalf("request scope retained the caller label buffer: %+v", scope)
	}

	reused := pool.Get()
	if reused.OwnerID != "" || reused.PetID != "" || len(reused.Labels) != 0 {
		t.Fatalf("pool returned identity left by the prior request: %+v", reused)
	}
	pool.Put(reused)

	unownedAudit := adapter.NewRequestAudit()
	unownedDone := make(chan struct{})
	go func() {
		unownedAudit.Record(&model.RequestScope{OwnerID: "leaked", PetID: "leaked"})
		close(unownedDone)
	}()
	<-unownedAudit.Started()
	unownedAudit.Release()
	<-unownedDone
	if leaked := unownedAudit.Entries(); len(leaked) != 0 {
		t.Fatalf("audit accepted a pooled scope without ownership transfer: %+v", leaked)
	}
}
