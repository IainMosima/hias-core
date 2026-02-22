package schema

import "time"

type CreatePolicyRequest struct {
	PlanID            string    `json:"plan_id" binding:"required,uuid"`
	PolicyholderName  string    `json:"policyholder_name" binding:"required"`
	PolicyholderEmail string    `json:"policyholder_email" binding:"required,email"`
	PolicyholderPhone string    `json:"policyholder_phone" binding:"required"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
}

type EnrollMemberRequest struct {
	NationalID   string `json:"national_id"`
	Name         string `json:"name" binding:"required"`
	DateOfBirth  string `json:"date_of_birth" binding:"required"`
	Gender       string `json:"gender" binding:"required,oneof=male female other"`
	Relationship string `json:"relationship" binding:"required,oneof=principal spouse child parent"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

type ActivatePolicyRequest struct {
	PaymentReference string `json:"payment_reference" binding:"required"`
}
