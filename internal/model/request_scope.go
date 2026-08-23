package model

// RequestScope 描述一次审计查询的主人、宠物与标签上下文。
//
// Owned 标记调用方是否已通过 Prepare 交接所有权：只有 Owned 为 true 时
// 审计才允许写入条目。Put 回收时通过 Reset 将 Owned 置 false，防止
// 回收对象残留上一任主人的身份而被误写。
type RequestScope struct {
	OwnerID string
	PetID   string
	Labels  []string
	Owned   bool
}

// Prepare 填充 scope 的主人、宠物与标签数据，拷贝标签切片以避免外部
// 修改穿透，并标记 Owned = true 表示所有权已交接给调用方。
func (s *RequestScope) Prepare(ownerID, petID string, labels []string) {
	s.OwnerID = ownerID
	s.PetID = petID
	s.Labels = make([]string, len(labels))
	copy(s.Labels, labels)
	s.Owned = true
}

// Reset 清空全部字段并撤销所有权，在 scope 回收到池之前调用，
// 使回收对象不再残留上一任调用方的身份。
func (s *RequestScope) Reset() {
	s.OwnerID = ""
	s.PetID = ""
	s.Labels = nil
	s.Owned = false
}
