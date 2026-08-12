package kubernetes

import (
	"context"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"
)

type Adapter struct{}

func (Adapter) Name() string { return "kubernetes" }
func (Adapter) Plan(_ context.Context, in contract.ApplicationIntent) (contract.Plan, error) {
	warnings := []string{}
	if in.GPU { warnings = append(warnings, "GPU availability depends on node pool configuration") }
	return contract.Plan{Provider:"kubernetes", Runtime:"kubernetes", InstanceClass:"portable",
		Capabilities: contract.CapabilitySet{Autoscaling:true, ManagedIngress:false, PersistentStorage:true, GPU:true, ZeroDowntime:true}, Warnings:warnings}, nil
}
