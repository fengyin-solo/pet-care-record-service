package model

import (
	"strings"
	"time"
)

// Owner 宠物主人档案。
type Owner struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验主人字段并规范化。
func (o *Owner) Validate() error {
	o.Name = strings.TrimSpace(o.Name)
	o.Phone = strings.TrimSpace(o.Phone)
	o.Email = strings.TrimSpace(o.Email)
	o.Address = strings.TrimSpace(o.Address)
	if o.Name == "" {
		return NewValidationError("name", "主人姓名不能为空")
	}
	if o.Phone == "" {
		return NewValidationError("phone", "联系电话不能为空")
	}
	return nil
}

// OwnerFilter 主人查询筛选条件。
type OwnerFilter struct {
	Keyword string
}

// Match 判断主人是否命中筛选条件。
func (f OwnerFilter) Match(o *Owner) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" &&
			!strings.Contains(strings.ToLower(o.Name), k) &&
			!strings.Contains(strings.ToLower(o.Phone), k) &&
			!strings.Contains(strings.ToLower(o.Email), k) {
			return false
		}
	}
	return true
}
