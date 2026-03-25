package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/google/uuid"
)

type ContractService interface {
	CreateContract(ctx context.Context, providerID uuid.UUID, req providerSchema.CreateContractRequest) *schema.ServiceResponse[providerSchema.ContractResponse]
	ListContracts(ctx context.Context, providerID uuid.UUID) *schema.ServiceResponse[[]providerSchema.ContractResponse]
}
