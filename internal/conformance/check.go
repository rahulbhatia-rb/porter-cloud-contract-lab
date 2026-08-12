package conformance

import (
	"fmt"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"
)

type Result struct { Provider string; Passed bool; Failures []string }

func Check(in contract.ApplicationIntent, p contract.Plan) Result {
	failures := []string{}
	if in.MaxReplicas > in.MinReplicas && !p.Capabilities.Autoscaling { failures = append(failures, "autoscaling unsupported") }
	if in.Public && !p.Capabilities.ManagedIngress && p.Provider != "kubernetes" { failures = append(failures, "managed ingress unsupported") }
	if in.StorageGB > 0 && !p.Capabilities.PersistentStorage { failures = append(failures, "persistent storage unsupported") }
	if in.GPU && !p.Capabilities.GPU { failures = append(failures, "GPU unsupported") }
	if in.Replicas < 1 { failures = append(failures, fmt.Sprintf("replicas must be >= 1, got %d", in.Replicas)) }
	return Result{Provider:p.Provider, Passed:len(failures)==0, Failures:failures}
}
