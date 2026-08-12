package providers

import (
	"context"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"
)

type Provider interface {
	Name() string
	Plan(context.Context, contract.ApplicationIntent) (contract.Plan, error)
}
