package serializers

type CreateClientReq struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description"`
	AppId       uint   `json:"app_id" binding:"required"`
}

type UpdateClientReq struct {
	Name        string `json:"name" binding:"max=50"`
	Description string `json:"description"`
}
