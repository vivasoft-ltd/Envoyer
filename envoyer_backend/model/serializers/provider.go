package serializers

import (
	"gorm.io/datatypes"
)

type CreateProviderReq struct {
	AppId        uint           `json:"app_id" binding:"required"`
	Type         string         `json:"type" binding:"required"`
	ProviderType string         `json:"provider_type" binding:"required"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Config       datatypes.JSON `json:"config" binding:"required"`
	Priority     uint           `json:"priority"`
	Active       bool           `json:"active"`
	Policy       datatypes.JSON `json:"policy"`
}

type UpdateProviderReq struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Config      datatypes.JSON `json:"config"`
	Priority    uint           `json:"priority"`
	Active      bool           `json:"active"`
	Policy      datatypes.JSON `json:"policy"`
}

type UpdatePriorityReq struct {
	ProviderPriority []Priority `json:"priority"`
}

type Priority struct {
	Id       uint `json:"id"`
	Priority uint `json:"priority"`
}
