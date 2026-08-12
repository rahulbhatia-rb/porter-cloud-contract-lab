package contract

type ApplicationIntent struct {
	Name string `json:"name"`
	Region string `json:"region"`
	CPU int `json:"cpu_millicores"`
	MemoryMB int `json:"memory_mb"`
	Replicas int `json:"replicas"`
	MinReplicas int `json:"min_replicas"`
	MaxReplicas int `json:"max_replicas"`
	Public bool `json:"public"`
	GPU bool `json:"gpu"`
	StorageGB int `json:"storage_gb"`
	Labels map[string]string `json:"labels,omitempty"`
}

type CapabilitySet struct { Autoscaling, ManagedIngress, PersistentStorage, GPU, ZeroDowntime bool }

type Plan struct {
	Provider string
	Runtime string
	InstanceClass string
	Capabilities CapabilitySet
	Warnings []string
	EstimatedCost float64
}
