package iprs

import "time"

type VerifyRequest struct {
	NationalID string `json:"national_id"`
}

type VerifyResponse struct {
	NationalID  string    `json:"national_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Verified    bool      `json:"verified"`
}
