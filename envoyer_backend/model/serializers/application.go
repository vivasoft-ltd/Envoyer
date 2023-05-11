package serializers

type CreateAppReq struct {
	Name        string `json:"name" binding:"required,_alpha-num,min=2,max=30"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type UpdateAppReq struct {
	Description string `json:"description"`
	Active      bool   `json:"active"`
}
