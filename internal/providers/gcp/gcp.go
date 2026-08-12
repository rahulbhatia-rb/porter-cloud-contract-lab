package gcp

import (
	"context"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"
)

type Adapter struct{}

func (Adapter) Name() string { return "gcp" }
func (Adapter) Plan(_ context.Context, in contract.ApplicationIntent) (contract.Plan, error) {
	class := "general-purpose"
	if in.GPU { class = "gpu" }
	return contract.Plan{Provider:"gcp", Runtime:"gke", InstanceClass:class,
		Capabilities: contract.CapabilitySet{Autoscaling:true, ManagedIngress:true, PersistentStorage:true, GPU:true, ZeroDowntime:true}}, nil
}
