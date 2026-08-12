package cost

import "github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"

func Estimate(in contract.ApplicationIntent, provider string) float64 {
	base := float64(in.CPU)/1000*22 + float64(in.MemoryMB)/1024*6 + float64(in.StorageGB)*0.12
	if in.GPU { base += 450 }
	multiplier := 1.0
	if provider == "gcp" { multiplier = 0.97 }
	return base * float64(in.Replicas) * multiplier
}
