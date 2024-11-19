package api

type PostTokenReq struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PostTokenResp struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	AccessToken  string `json:"access_token" binding:"required"`
	Name         string `json:"name" binding:"required"`
	ExpiresAt    int64  `json:"expires_at" binding:"required"`
}

type ValidateTokenReq struct {
	Token string `json:"token" binding:"required"`
}
type ValidateTokenResp struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" `
}

type DeleteTokenReq struct {
	Token string `json:"token" binding:"required"`
}
