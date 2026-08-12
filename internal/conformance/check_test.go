package conformance

import (
	"testing"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"
)

func TestGPURequirement(t *testing.T) {
	in := contract.ApplicationIntent{Replicas:1, GPU:true}
	plan := contract.Plan{Provider:"example", Capabilities:contract.CapabilitySet{GPU:false}}
	got := Check(in, plan)
	if got.Passed { t.Fatal("expected conformance failure") }
}
