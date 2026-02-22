package auth

type TokenMaker interface {
	CreateToken(userID string, email string, role string, permissions []string, duration interface{}) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
