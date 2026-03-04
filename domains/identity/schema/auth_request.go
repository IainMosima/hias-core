package schema

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	Name       string `json:"name" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	NationalID string `json:"national_id"`
	RoleName   string `json:"role_name"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
