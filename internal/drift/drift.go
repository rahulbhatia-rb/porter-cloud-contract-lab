package drift

import (
	"strconv"
	"github.com/rahulbhatia-rb/porter-cloud-contract-lab/internal/contract"
)

type Finding struct { Field, Desired, Actual string }

func Compare(desired, actual contract.ApplicationIntent) []Finding {
	out := []Finding{}
	if desired.CPU != actual.CPU { out = append(out, Finding{"cpu_millicores", strconv.Itoa(desired.CPU), strconv.Itoa(actual.CPU)}) }
	if desired.MemoryMB != actual.MemoryMB { out = append(out, Finding{"memory_mb", strconv.Itoa(desired.MemoryMB), strconv.Itoa(actual.MemoryMB)}) }
	if desired.Replicas != actual.Replicas { out = append(out, Finding{"replicas", strconv.Itoa(desired.Replicas), strconv.Itoa(actual.Replicas)}) }
	return out
}
