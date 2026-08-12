package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/conformance"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/cost"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/providers"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/providers/aws"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/providers/gcp"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/providers/kubernetes"
)

func main() {
	if len(os.Args) != 2 { fmt.Fprintln(os.Stderr, "usage: contractlab <intent.json>"); os.Exit(2) }
	b, err := os.ReadFile(os.Args[1]); if err != nil { panic(err) }
	var in contract.ApplicationIntent
	if err := json.Unmarshal(b, &in); err != nil { panic(err) }
	adapters := []providers.Provider{aws.Adapter{}, gcp.Adapter{}, kubernetes.Adapter{}}
	for _, a := range adapters {
		plan, err := a.Plan(context.Background(), in); if err != nil { panic(err) }
		plan.EstimatedCost = cost.Estimate(in, a.Name())
		result := conformance.Check(in, plan)
		fmt.Printf("%-10s runtime=%-10s conformant=%-5v est_cost=$%.2f/mo warnings=%v failures=%v\n", a.Name(), plan.Runtime, result.Passed, plan.EstimatedCost, plan.Warnings, result.Failures)
	}
}
