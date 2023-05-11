package serializers

type CreateEventReq struct {
	Name        string   `json:"name" binding:"required,_alpha-num,min=2,max=50"`
	Description string   `json:"description"`
	Variables   []string `json:"variables" binding:"variable_format"`
	AppId       uint     `json:"app_id" binding:"required"`
}

type UpdateEventReq struct {
	Name        string   `json:"name" binding:"_alpha-num,max=50"`
	Description string   `json:"description"`
	Variables   []string `json:"variables" binding:"variable_format"`
}
