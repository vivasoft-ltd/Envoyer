package serializers

type CreateUserReq struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
	AppId    uint   `json:"app_id" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
