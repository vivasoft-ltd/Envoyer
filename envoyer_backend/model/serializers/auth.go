package serializers

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTTokenResponse struct {
	ExpiredAt int64  `json:"expired_at"`
	Token     string `json:"access_token"`
	Refresh   string `json:"refresh_token"`
	Id        uint   `json:"id"`
	Role      string `json:"role,omitempty"`
	AppId     uint   `json:"app_id,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type PublishReq struct {
	AppKey    string `json:"app_key" binding:"required"`
	ClientKey string `json:"client_key" binding:"required"`
}
