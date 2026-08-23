package healthretry

import (
	"testing"

	"petsmanagement/internal/adapter"
	"petsmanagement/internal/model"
	. "petsmanagement/internal/service"
	"petsmanagement/internal/store"
)

func TestRetryStoresSingleHealthTransition(t *testing.T) {
	st := store.NewMemoryStore()
	svc := New(st, nil, nil)
	publisher := adapter.NewHealthPublisher(true)
	if err := svc.UpdateHealthWithRetry("pet-8", "recovering", publisher); err != nil {
		t.Fatalf("retry succeeded downstream but caller still received failure: %v", err)
	}
	history := svc.HealthUpdateHistory()
	if len(history) != 1 || history[0].Status != "recovering" {
		t.Fatalf("one logical update produced duplicate or wrong history: %+v", history)
	}
	if history[0].Key != "health-update:pet-8:recovering" {
		t.Fatalf("retry changed the logical update key: %q", history[0].Key)
	}
	if accepted := st.ApplyHealthUpdate(model.HealthUpdate{PetID: "pet-8", Status: "recovering", Key: "legacy-attempt-2"}); accepted {
		t.Fatal("store accepted the same logical update under an old attempt key")
	}
	if err := publisher.Publish(history[0].Key); err != nil {
		t.Fatalf("idempotent publisher replay failed: %v", err)
	}
	if publisher.Calls() != 3 || publisher.Effects() != 1 {
		t.Fatalf("publisher retry was not idempotent: calls=%d effects=%d", publisher.Calls(), publisher.Effects())
	}
}
