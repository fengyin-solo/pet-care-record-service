package model

type FollowUpDelivery struct {
	FollowUpID   string
	Delivered    bool
	SlotReleased bool
}

func (d *FollowUpDelivery) MarkSlotReleased() {}
