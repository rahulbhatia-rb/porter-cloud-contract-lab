# Porter Cloud Contract Lab

A Go proof-of-concept for a **cloud-neutral infrastructure contract, provider compatibility layer, conformance framework, drift detector, and cost comparison model** for applications that need to run consistently across AWS, GCP, and Kubernetes environments.

> **Important context:** this is an independent engineering prototype inspired by Porter's publicly described product, architecture, and Senior Backend Engineer challenges. It does not use, reproduce, or claim knowledge of Porter's private code, infrastructure, APIs, customer data, or implementation details. The project exists to demonstrate how I would reason about one of the hardest problems in a multi-cloud PaaS: keeping a simple developer-facing contract stable while the infrastructure underneath it constantly changes.

---

## Table of contents

1. [Why this project exists](#why-this-project-exists)
2. [The core engineering problem](#the-core-engineering-problem)
3. [The idea: Cloud Contract Compatibility Lab](#the-idea-cloud-contract-compatibility-lab)
4. [What the prototype does today](#what-the-prototype-does-today)
5. [What the prototype intentionally does not do yet](#what-the-prototype-intentionally-does-not-do-yet)
6. [Architecture](#architecture)
7. [End-to-end execution flow](#end-to-end-execution-flow)
8. [Application Intent Contract](#application-intent-contract)
9. [Provider interface](#provider-interface)
10. [AWS adapter](#aws-adapter)
11. [GCP adapter](#gcp-adapter)
12. [Generic Kubernetes adapter](#generic-kubernetes-adapter)
13. [Capability model](#capability-model)
14. [Provider plan model](#provider-plan-model)
15. [Conformance engine](#conformance-engine)
16. [Drift detection](#drift-detection)
17. [Cost model](#cost-model)
18. [CLI](#cli)
19. [Example application intent](#example-application-intent)
20. [Tests](#tests)
21. [CI](#ci)
22. [Makefile](#makefile)
23. [Repository structure](#repository-structure)
24. [How to run locally](#how-to-run-locally)
25. [Expected output](#expected-output)
26. [Design principles](#design-principles)
27. [Important tradeoffs](#important-tradeoffs)
28. [How this maps to Porter's engineering challenges](#how-this-maps-to-porters-engineering-challenges)
29. [How I would evolve this into a production system](#how-i-would-evolve-this-into-a-production-system)
30. [Potential production architecture](#potential-production-architecture)
31. [Provider version and API compatibility](#provider-version-and-api-compatibility)
32. [Kubernetes upgrade safety](#kubernetes-upgrade-safety)
33. [Capability negotiation](#capability-negotiation)
34. [Drift reconciliation](#drift-reconciliation)
35. [Cost-aware placement](#cost-aware-placement)
36. [Security and isolation considerations](#security-and-isolation-considerations)
37. [Observability](#observability)
38. [Failure handling](#failure-handling)
39. [Testing strategy for a production implementation](#testing-strategy-for-a-production-implementation)
40. [Possible API surface](#possible-api-surface)
41. [Example product workflow](#example-product-workflow)
42. [Why Go](#why-go)
43. [Current limitations](#current-limitations)
44. [Roadmap](#roadmap)
45. [Why I built this for Porter](#why-i-built-this-for-porter)
46. [Author](#author)

---

## Why this project exists

Porter's product proposition is compelling because the end user should not need to think deeply about Kubernetes, cloud networking, managed cluster operations, autoscaling, or provider-specific infrastructure APIs just to deploy an application.

The challenge is that the underlying infrastructure is not actually uniform.

AWS, GCP, Azure, and Kubernetes distributions differ in:

- API shapes
- supported Kubernetes versions
- upgrade windows
- deprecated resources
- autoscaling behavior
- load balancer implementation
- identity and access models
- persistent storage semantics
- GPU availability
- managed database integrations
- network topology
- availability constructs
- quotas
- pricing
- regional capability
- security defaults
- lifecycle behavior

A platform can hide those differences from users only up to a point. If it hides too much, it can accidentally create false equivalence between providers. If it exposes too much, the platform stops feeling simple.

The core question behind this prototype is therefore:

> **How can a multi-cloud PaaS offer a stable, simple developer-facing infrastructure contract while explicitly testing and surfacing the places where cloud-provider behavior diverges?**

---

## The core engineering problem

A typical developer intent might be very simple:

- run this application
- give it 1 CPU
- give it 2 GB memory
- keep at least 2 replicas
- allow scaling up to 20 replicas
- expose it publicly
- attach 20 GB persistent storage
- deploy it in a target region

That request looks cloud-neutral.

The provider-specific implementation is not.

For AWS, the platform may need to reason about EKS, node groups, VPC networking, IAM, load balancers, storage classes, autoscaling, quotas, region support, and version compatibility.

For GCP, the equivalent implementation may involve GKE, IAM, VPCs, Google Cloud load balancers, persistent disks, autoscaling semantics, node pools, and different version schedules.

For a user-supplied Kubernetes cluster, Porter may not control the cloud primitives at all. A managed ingress controller may or may not exist. GPU nodes may or may not exist. The cluster could expose different storage classes and Kubernetes capabilities.

The application intent should remain stable even when these implementations differ.

This repository separates those concerns explicitly.

---

## The idea: Cloud Contract Compatibility Lab

The proposed concept is a **Cloud Contract Compatibility Lab**.

The system has five major responsibilities:

1. **Capture developer intent** in a provider-neutral contract.
2. **Translate that intent** through a provider adapter.
3. **Describe the resulting provider plan and capabilities** explicitly.
4. **Run conformance rules** to determine whether the provider plan satisfies the original intent.
5. **Evaluate drift and cost separately** so portability is not reduced to a simple yes/no provider translation.

The important architectural decision is this:

> **Intent is the source of truth. Provider implementation is an adapter. Compatibility is a testable result.**

This is different from embedding provider-specific logic throughout the application-management layer.

---

## What the prototype does today

The current repository is intentionally small, but every piece maps to a larger production concept.

Today it can:

- read an application infrastructure intent from JSON
- represent CPU, memory, replicas, autoscaling range, public exposure, GPU requirements, persistent storage, region, and labels
- pass the same intent to multiple provider adapters
- generate an abstract provider plan
- represent provider capabilities
- run a shared conformance engine against each provider plan
- detect capability mismatches
- estimate a simple provider-relative monthly cost
- compare desired and actual application state for selected drift fields
- run a CLI that evaluates AWS, GCP, and Kubernetes in one execution
- run unit tests
- run Go static validation using `go vet`
- run all validation in GitHub Actions

---

## What the prototype intentionally does not do yet

This repository is **not pretending to be a production cloud control plane**.

It currently does not:

- call real AWS APIs
- call real GCP APIs
- call Azure APIs
- call a live Kubernetes API server
- provision infrastructure
- mutate clusters
- reconcile state continuously
- authenticate to cloud accounts
- manage IAM policies
- perform real cloud pricing lookups
- discover quotas dynamically
- detect Kubernetes API removals from live clusters
- validate real EKS/GKE versions
- create load balancers
- create VPCs
- create managed databases
- perform rolling upgrades
- create autoscalers
- schedule workloads
- perform live health checks
- persist control-plane state
- manage asynchronous workflows
- implement retries or idempotency tokens
- model every possible application feature

The goal is to demonstrate **architecture and reasoning**, not to produce a fake production platform in a small repository.

---

## Architecture

The conceptual architecture is:

```text
                          Developer / API Client
                                  |
                                  v
                       +-----------------------+
                       | Application Intent    |
                       | Contract              |
                       +-----------------------+
                                  |
                                  v
                       +-----------------------+
                       | Policy / Capability   |
                       | Evaluation Boundary   |
                       +-----------------------+
                          /         |         \
                         /          |          \
                        v           v           v
                +-----------+ +-----------+ +-------------+
                | AWS       | | GCP       | | Kubernetes  |
                | Adapter   | | Adapter   | | Adapter     |
                +-----------+ +-----------+ +-------------+
                     |             |              |
                     v             v              v
                   EKS           GKE        Kubernetes API
                  IAM/VPC       IAM/VPC      Cluster-specific
                     \             |              /
                      \            |             /
                       +-------------------------+
                       | Provider Plan           |
                       +-------------------------+
                                  |
                    +-------------+-------------+
                    |                           |
                    v                           v
            +---------------+            +--------------+
            | Conformance   |            | Cost / Drift |
            | Engine        |            | Evaluation   |
            +---------------+            +--------------+
                    |                           |
                    +-------------+-------------+
                                  |
                                  v
                       Safe / Unsafe / Warning
                                  |
                                  v
                        Reconcile or Reject
```

The prototype deliberately separates **intent**, **provider translation**, **capabilities**, **conformance**, **cost**, and **drift**.

That separation is the main architectural point of the project.

---

## End-to-end execution flow

When the CLI runs, the following sequence occurs:

1. `cmd/contractlab/main.go` receives a path to an intent JSON file.
2. The file is loaded from disk.
3. JSON is unmarshaled into `contract.ApplicationIntent`.
4. The CLI constructs three provider adapters:
   - AWS
   - GCP
   - Kubernetes
5. Each adapter receives the exact same `ApplicationIntent`.
6. Each adapter returns a `contract.Plan`.
7. The generic cost package estimates monthly cost for that provider.
8. The conformance engine checks whether the returned provider capabilities satisfy the original application intent.
9. The CLI prints:
   - provider name
   - runtime
   - conformance result
   - estimated cost
   - provider warnings
   - conformance failures

Nothing in the intent object is changed to make a provider pass.

That is deliberate. If a provider cannot satisfy the contract, the platform should surface that mismatch.

---

## Application Intent Contract

The central model is defined in:

```text
internal/contract/model.go
```

The structure is:

```go
type ApplicationIntent struct {
    Name        string            `json:"name"`
    Region      string            `json:"region"`
    CPU         int               `json:"cpu_millicores"`
    MemoryMB    int               `json:"memory_mb"`
    Replicas    int               `json:"replicas"`
    MinReplicas int               `json:"min_replicas"`
    MaxReplicas int               `json:"max_replicas"`
    Public      bool              `json:"public"`
    GPU         bool              `json:"gpu"`
    StorageGB   int               `json:"storage_gb"`
    Labels      map[string]string `json:"labels,omitempty"`
}
```

### Field meanings

#### `name`

Logical application or service name.

Example:

```json
"name": "api"
```

#### `region`

Desired deployment region.

In the current prototype the value is carried as intent but is not yet validated per provider.

A production version would translate a logical region requirement into provider-specific region availability and capability checks.

#### `cpu_millicores`

Requested CPU expressed using Kubernetes-style millicores.

Example:

```json
"cpu_millicores": 1000
```

This represents approximately one CPU core.

Using a provider-neutral unit avoids exposing instance-family selection to the application contract.

#### `memory_mb`

Requested memory in megabytes.

Example:

```json
"memory_mb": 2048
```

#### `replicas`

Requested initial replica count.

This is validated by the conformance engine to ensure it is at least one.

#### `min_replicas`

Minimum replica count for autoscaling.

#### `max_replicas`

Maximum replica count for autoscaling.

If `max_replicas > min_replicas`, the intent is treated as requiring autoscaling capability.

#### `public`

Whether the application needs public ingress.

A provider plan must expose the appropriate ingress capability unless the generic Kubernetes model is being used, where ingress availability can depend on cluster configuration.

#### `gpu`

Whether the workload requires GPU support.

The conformance engine validates GPU capability explicitly.

#### `storage_gb`

Requested persistent storage capacity.

Any value greater than zero means the provider plan must support persistent storage.

#### `labels`

Optional metadata that can later be used for:

- policy
- placement
- billing allocation
- team ownership
- environment identification
- observability correlation
- deployment controls

The sample uses:

```json
"labels": {
  "tier": "production"
}
```

---

## Provider interface

The provider abstraction lives in:

```text
internal/providers/provider.go
```

The interface is deliberately small:

```go
type Provider interface {
    Name() string
    Plan(context.Context, contract.ApplicationIntent) (contract.Plan, error)
}
```

A provider therefore has two responsibilities in the prototype:

### `Name()`

Returns the provider identifier.

Examples:

- `aws`
- `gcp`
- `kubernetes`

### `Plan()`

Accepts an `ApplicationIntent` and returns an abstract provider plan.

In a production implementation, `Plan()` could become a much richer planning phase that:

- discovers provider capabilities
- reads quotas
- selects regions
- chooses node types
- selects storage classes
- validates Kubernetes versions
- computes upgrade paths
- applies policy
- resolves network requirements
- creates an execution DAG

The provider interface keeps that logic behind a stable boundary.

---

## AWS adapter

Located at:

```text
internal/providers/aws/aws.go
```

The prototype AWS adapter returns:

- provider: `aws`
- runtime: `eks`
- general-purpose instance class by default
- GPU instance class when GPU is requested

The AWS capability set currently declares support for:

- autoscaling
- managed ingress
- persistent storage
- GPU
- zero-downtime deployment capability

The implementation is intentionally abstract.

A production AWS adapter might eventually reason about:

- EKS versions
- managed node groups
- Karpenter
- EC2 families
- Graviton compatibility
- EBS volume types
- EFS
- ALB/NLB
- VPC and subnet topology
- security groups
- IAM roles for service accounts / pod identity
- AWS Load Balancer Controller
- Route 53
- ACM
- CloudWatch
- ECR
- RDS compatibility
- ElastiCache compatibility
- service quotas
- spot capacity
- region availability
- upgrade policies

---

## GCP adapter

Located at:

```text
internal/providers/gcp/gcp.go
```

The GCP adapter follows the same contract.

It currently returns:

- provider: `gcp`
- runtime: `gke`
- general-purpose class by default
- GPU class when requested

Its current declared capabilities are equivalent to the AWS adapter for the features represented by this prototype.

A real implementation could reason about:

- GKE Standard vs Autopilot
- release channels
- node pools
- machine families
- GPU node support
- persistent disks
- Filestore
- Google Cloud Load Balancing
- VPC-native clusters
- Workload Identity
- Cloud DNS
- Artifact Registry
- Cloud SQL
- Memorystore
- quotas
- regional vs zonal clusters
- upgrade schedules
- maintenance windows
- GKE deprecations

---

## Generic Kubernetes adapter

Located at:

```text
internal/providers/kubernetes/kubernetes.go
```

This adapter represents an important case: Porter may be dealing with Kubernetes itself rather than a fully controlled hyperscaler environment.

It returns:

- provider: `kubernetes`
- runtime: `kubernetes`
- instance class: `portable`

The adapter currently declares:

- autoscaling support
- persistent storage support
- GPU support
- zero-downtime capability
- **no guaranteed managed ingress**

That difference is intentional.

A generic Kubernetes cluster cannot guarantee that a cloud-managed ingress implementation exists.

The adapter also emits a warning when GPU is requested:

```text
GPU availability depends on node pool configuration
```

This demonstrates a key idea:

> A feature can be conceptually supported by Kubernetes but still depend on runtime environment configuration.

A production implementation should distinguish:

- supported
- unsupported
- supported with prerequisites
- unknown until discovery
- supported but degraded

---

## Capability model

Capabilities are represented by:

```go
type CapabilitySet struct {
    Autoscaling        bool
    ManagedIngress     bool
    PersistentStorage bool
    GPU                bool
    ZeroDowntime       bool
}
```

The current set is intentionally small.

It demonstrates the principle that provider compatibility should be explicit rather than assumed.

Possible future capabilities include:

- IPv6
- private ingress
- public ingress
- internal load balancing
- network policies
- service mesh
- GPU types
- spot capacity
- persistent block storage
- RWX storage
- workload identity
- secret manager integration
- horizontal autoscaling
- vertical autoscaling
- cluster autoscaling
- zero-downtime node upgrades
- topology spread
- multi-zone scheduling
- managed DNS
- managed certificates
- WebSocket support
- UDP support
- static egress IP
- private networking
- external managed databases

A capability model can eventually become versioned and machine-readable.

---

## Provider plan model

A provider adapter returns:

```go
type Plan struct {
    Provider      string
    Runtime       string
    InstanceClass string
    Capabilities  CapabilitySet
    Warnings      []string
    EstimatedCost float64
}
```

### `Provider`

The selected provider implementation.

### `Runtime`

The execution runtime.

Current examples:

- `eks`
- `gke`
- `kubernetes`

### `InstanceClass`

A high-level scheduling/compute class.

Current values include:

- `general-purpose`
- `gpu`
- `portable`

In production this might be expanded into an internal placement object rather than exposing raw instance types.

### `Capabilities`

The provider capabilities available for the generated plan.

### `Warnings`

Non-fatal issues or prerequisites.

Warnings are different from conformance failures.

That distinction is useful because not every provider difference should block deployment.

### `EstimatedCost`

Approximate monthly cost returned by the generic cost evaluation layer.

---

## Conformance engine

Located at:

```text
internal/conformance/check.go
```

The conformance layer receives:

- original application intent
- provider-generated plan

It does **not** know how AWS or GCP work internally.

It only determines whether the plan satisfies the contract.

Current checks are:

### Autoscaling

If:

```text
max_replicas > min_replicas
```

then the provider must expose:

```text
Autoscaling = true
```

Otherwise the failure is:

```text
autoscaling unsupported
```

### Managed ingress

If the application is public, a provider plan is expected to support managed ingress.

The generic Kubernetes provider is treated differently because ingress availability depends on cluster configuration.

### Persistent storage

If:

```text
storage_gb > 0
```

then:

```text
PersistentStorage = true
```

is required.

Otherwise:

```text
persistent storage unsupported
```

is returned.

### GPU

If the application requests GPU resources, the provider plan must expose GPU capability.

Otherwise:

```text
GPU unsupported
```

is returned.

### Replica validation

The current contract also validates:

```text
replicas >= 1
```

If not:

```text
replicas must be >= 1
```

is returned with the invalid replica count.

### Result

The engine returns:

```go
type Result struct {
    Provider string
    Passed   bool
    Failures []string
}
```

This makes compatibility testable and composable.

---

## Drift detection

Located at:

```text
internal/drift/drift.go
```

Drift detection compares a desired `ApplicationIntent` with an actual observed `ApplicationIntent`.

The current implementation checks:

- CPU
- memory
- replicas

A mismatch produces:

```go
type Finding struct {
    Field   string
    Desired string
    Actual  string
}
```

For example:

```text
Field: replicas
Desired: 3
Actual: 2
```

The current implementation is deliberately narrow.

A production drift system would need to distinguish between:

### User-controlled drift

The user changed infrastructure outside Porter.

### Provider drift

The cloud provider changed a managed resource or default.

### Controller drift

A Kubernetes controller modified state.

### Expected runtime drift

Autoscaling changed replica count legitimately.

### Dangerous drift

A change violates a platform invariant.

The reconciliation policy should not blindly revert every mismatch.

For example, HPA-managed replica changes should not be treated the same way as an unauthorized security group change.

---

## Cost model

Located at:

```text
internal/cost/score.go
```

The cost model is intentionally illustrative.

It estimates a base monthly amount using:

- requested CPU
- requested memory
- persistent storage
- GPU requirement
- replica count

The rough formula uses:

```text
CPU       -> 22 units per requested core
Memory    -> 6 units per requested GiB
Storage   -> 0.12 units per GB
GPU       -> +450 units when requested
```

The total is multiplied by replica count.

The GCP adapter currently receives a `0.97` multiplier simply to demonstrate that the same application intent can produce different cost outcomes across providers.

This is **not a cloud pricing calculator**.

A production version would need:

- live AWS pricing data
- live GCP pricing data
- Azure pricing
- region-specific rates
- cluster control plane cost
- network egress
- NAT cost
- load balancer cost
- storage IOPS
- snapshots
- managed databases
- logging/metrics ingestion
- spot/preemptible discounts
- committed-use discounts
- Savings Plans / reservations
- autoscaling utilization
- idle capacity

The architectural point is that cost evaluation is separate from provider translation.

This enables the platform to ask:

> Both providers satisfy the contract, but do they satisfy it with materially different economics?

---

## CLI

The command-line entry point is:

```text
cmd/contractlab/main.go
```

Usage:

```bash
go run ./cmd/contractlab <intent.json>
```

The CLI requires exactly one argument.

If no intent file is supplied, it prints:

```text
usage: contractlab <intent.json>
```

and exits with status code `2`.

The CLI then:

1. reads the file
2. unmarshals JSON
3. constructs AWS, GCP, and Kubernetes adapters
4. generates a plan for each
5. estimates cost
6. performs conformance validation
7. prints a result per provider

The CLI is intentionally synchronous and easy to inspect.

A real platform would likely make planning and provisioning asynchronous.

---

## Example application intent

The sample input lives at:

```text
examples/api.json
```

It contains:

```json
{
  "name": "api",
  "region": "us-east-1",
  "cpu_millicores": 1000,
  "memory_mb": 2048,
  "replicas": 2,
  "min_replicas": 2,
  "max_replicas": 20,
  "public": true,
  "gpu": false,
  "storage_gb": 20,
  "labels": {
    "tier": "production"
  }
}
```

This describes a production API requiring:

- 1 CPU per replica
- 2 GB memory per replica
- 2 initial replicas
- minimum 2 replicas
- maximum 20 replicas
- autoscaling capability
- public ingress
- 20 GB storage
- no GPU

The same intent is evaluated by all providers.

---

## Tests

The repository currently contains a focused unit test in:

```text
internal/conformance/check_test.go
```

The test verifies that a workload requiring GPU fails conformance when the provider plan declares:

```text
GPU = false
```

This validates the basic contract/conformance relationship.

The current test suite is intentionally small.

Production expansion should include:

- table-driven tests for every capability
- invalid input tests
- autoscaling boundary tests
- ingress tests
- persistent storage tests
- provider warning tests
- region compatibility tests
- pricing tests
- drift tests
- provider adapter contract tests
- fuzz tests for application intent parsing
- golden plan tests

---

## CI

GitHub Actions workflow:

```text
.github/workflows/ci.yml
```

The workflow triggers on:

- every push
- every pull request

It uses:

```text
ubuntu-latest
```

and Go:

```text
1.23.x
```

The CI pipeline performs three validations:

```bash
go test ./...
go vet ./...
go run ./cmd/contractlab ./examples/api.json
```

This means the repository checks:

- unit-test correctness
- Go static analysis
- successful execution of the example contract through every adapter

The workflow requests only:

```yaml
permissions:
  contents: read
```

which follows least-privilege principles for this CI job.

---

## Makefile

The repository includes a minimal Makefile.

### Run tests

```bash
make test
```

Equivalent to:

```bash
go test ./...
```

### Run the demo

```bash
make run
```

Equivalent to:

```bash
go run ./cmd/contractlab ./examples/api.json
```

---

## Repository structure

```text
porter-cloud-contract-lab/
|
|-- .github/
|   `-- workflows/
|       `-- ci.yml
|
|-- cmd/
|   `-- contractlab/
|       `-- main.go
|
|-- docs/
|   |-- architecture.md
|   `-- idea.md
|
|-- examples/
|   `-- api.json
|
|-- internal/
|   |-- conformance/
|   |   |-- check.go
|   |   `-- check_test.go
|   |
|   |-- contract/
|   |   `-- model.go
|   |
|   |-- cost/
|   |   `-- score.go
|   |
|   |-- drift/
|   |   `-- drift.go
|   |
|   `-- providers/
|       |-- provider.go
|       |-- aws/
|       |   `-- aws.go
|       |-- gcp/
|       |   `-- gcp.go
|       `-- kubernetes/
|           `-- kubernetes.go
|
|-- go.mod
|-- Makefile
`-- README.md
```

### File-by-file purpose

#### `.github/workflows/ci.yml`

Automated validation for tests, vetting, and example execution.

#### `cmd/contractlab/main.go`

CLI orchestration layer.

#### `docs/architecture.md`

Compact architecture diagram and intent/provider separation explanation.

#### `docs/idea.md`

Short product proposal explaining why the compatibility lab matters.

#### `examples/api.json`

Representative production application intent.

#### `internal/contract/model.go`

Core provider-neutral models.

#### `internal/providers/provider.go`

Provider interface.

#### `internal/providers/aws/aws.go`

AWS/EKS planning adapter prototype.

#### `internal/providers/gcp/gcp.go`

GCP/GKE planning adapter prototype.

#### `internal/providers/kubernetes/kubernetes.go`

Generic Kubernetes adapter with environment-dependent warnings.

#### `internal/conformance/check.go`

Contract satisfaction checks.

#### `internal/conformance/check_test.go`

Unit test demonstrating capability failure behavior.

#### `internal/cost/score.go`

Simple provider-relative cost estimation.

#### `internal/drift/drift.go`

Desired-vs-actual comparison logic.

#### `go.mod`

Go module definition.

#### `Makefile`

Local developer commands.

---

## How to run locally

### Requirements

- Go 1.23 or compatible Go toolchain
- Git

No external services are required.

No AWS credentials are required.

No GCP credentials are required.

No Kubernetes cluster is required.

### Clone

```bash
git clone https://github.com/rahulbhatia-rb/porter-cloud-contract-lab.git
cd porter-cloud-contract-lab
```

### Run tests

```bash
go test ./...
```

or:

```bash
make test
```

### Run static analysis

```bash
go vet ./...
```

### Run the demo

```bash
go run ./cmd/contractlab ./examples/api.json
```

or:

```bash
make run
```

### Try a custom contract

Create a JSON file:

```json
{
  "name": "worker",
  "region": "us-east-1",
  "cpu_millicores": 2000,
  "memory_mb": 4096,
  "replicas": 3,
  "min_replicas": 3,
  "max_replicas": 50,
  "public": false,
  "gpu": true,
  "storage_gb": 100,
  "labels": {
    "tier": "production",
    "workload": "ml"
  }
}
```

Then run:

```bash
go run ./cmd/contractlab ./worker.json
```

---

## Expected output

The exact estimated numbers depend on the input, but output follows this form:

```text
aws        runtime=eks        conformant=true  est_cost=$.../mo warnings=[] failures=[]
gcp        runtime=gke        conformant=true  est_cost=$.../mo warnings=[] failures=[]
kubernetes runtime=kubernetes conformant=true  est_cost=$.../mo warnings=[] failures=[]
```

For GPU workloads, the generic Kubernetes adapter may emit:

```text
GPU availability depends on node pool configuration
```

This is a warning, not automatically a hard failure.

---

## Design principles

### 1. Intent before implementation

The developer says what the application needs.

The provider adapter decides how that requirement could be fulfilled.

### 2. Cloud neutrality does not mean pretending clouds are identical

The system explicitly represents capability differences.

### 3. Compatibility should be testable

A provider adapter should not be considered safe simply because it compiles.

It should satisfy the same contract suite as every other adapter.

### 4. Provider-specific code stays behind provider boundaries

The rest of the control plane should not need AWS-specific conditionals everywhere.

### 5. Warnings are different from failures

A prerequisite or operational caveat should not always block deployment.

### 6. Cost is part of infrastructure correctness

A plan that works technically but becomes dramatically more expensive is still important to surface.

### 7. Drift requires semantics

Not all drift is bad.

Reconciliation must understand ownership and controller behavior.

### 8. Abstractions need escape hatches

A truly useful cloud abstraction must eventually support provider-specific capabilities without corrupting the portable core contract.

---

## Important tradeoffs

### Lowest common denominator vs provider innovation

If a platform supports only features common to every provider, it becomes portable but weak.

If it exposes every provider-specific capability directly, portability disappears.

A better model is:

```text
Portable core contract
        +
Explicit capability negotiation
        +
Provider-specific extensions
```

### Declarative intent vs imperative operations

The prototype is declarative.

Real cloud APIs frequently require multi-stage imperative workflows.

A production system therefore needs a reconciliation/execution engine behind the declarative API.

### Planning vs provisioning

This repository stops at planning and validation.

That is intentional.

A production design should ideally separate:

```text
Plan -> Validate -> Approve -> Apply -> Observe -> Reconcile
```

rather than immediately provisioning from unvalidated input.

### Generic Kubernetes vs managed Kubernetes

A generic Kubernetes adapter cannot make the same assumptions as EKS or GKE.

This should be reflected in capability discovery rather than hidden.

### Cost accuracy vs planning speed

Live provider pricing can be expensive or slow to calculate comprehensively.

A production system may need cached price catalogs and approximate planning estimates followed by more detailed forecasting.

---

## How this maps to Porter's engineering challenges

The Senior Backend Engineer description calls out several specific challenges.

### "Own our infrastructure management system"

This prototype is centered on the infrastructure-management boundary:

```text
Intent -> Provider Plan -> Conformance -> Reconciliation
```

### "Maintain internal cloud-agnostic infrastructure APIs"

`ApplicationIntent` is the simplified representation of that idea.

The provider interface prevents cloud-specific implementation from leaking into the contract.

### "Allow users to seamlessly switch clouds"

A switch is safe only if the target provider satisfies the same contract.

That is why the conformance layer exists.

### "Continually integrate newest developments released by hyperscalers"

A provider adapter can evolve independently while remaining subject to the same conformance suite.

New capabilities can be introduced without silently changing the application contract.

### "Ensure Porter-native infrastructure remains compatible with other cloud services"

Capability-based planning and explicit drift detection help model resources that Porter may not own directly.

A production version could classify resource ownership and avoid overwriting customer-managed infrastructure.

### "Infrastructure that just works and doesn't cost an arm and a leg"

The separate cost layer exists because portability without cost awareness can still produce a bad user outcome.

---

## How I would evolve this into a production system

The next major step is to turn the static adapters into **real provider discovery + planning engines**.

The architecture would likely evolve into these components:

```text
API
 |
 v
Intent Service
 |
 v
Policy Engine
 |
 v
Capability Resolver
 |
 v
Provider Planner
 |
 v
Conformance Engine
 |
 v
Execution DAG
 |
 v
Provider Controllers
 |
 v
Observed State Store
 |
 +----> Drift Engine
 |
 +----> Cost Engine
 |
 +----> Upgrade Safety Engine
```

---

## Potential production architecture

### API layer

Accept application and infrastructure intent.

Responsibilities:

- authentication
- authorization
- schema validation
- idempotency
- request versioning
- tenant isolation

### Intent service

Stores the desired state.

The desired state should be immutable/versioned so deployments can be rolled back and audited.

### Provider discovery service

Queries live infrastructure to determine:

- cloud account capabilities
- quotas
- enabled services
- Kubernetes versions
- storage classes
- GPU/node availability
- network topology

### Policy engine

Applies organization-level rules such as:

- approved regions
- allowed instance families
- budget constraints
- security requirements
- public ingress restrictions
- mandatory encryption

### Planner

Produces an execution plan without making changes.

### Conformance engine

Verifies that the plan satisfies the intent.

### Execution engine

Applies the plan using idempotent operations.

A DAG model would be useful because infrastructure has dependencies.

For example:

```text
VPC
 |
 +--> Subnets
       |
       +--> Cluster
             |
             +--> Node pools
             |
             +--> Controllers
                     |
                     +--> Application
```

### State store

Tracks:

- desired state
- observed state
- provider resource identifiers
- operation state
- failures
- revisions

### Reconciler

Periodically compares desired and observed state and decides whether to:

- do nothing
- warn
- repair
- escalate

---

## Provider version and API compatibility

One of the most valuable extensions would be a version compatibility matrix.

Example conceptual record:

```text
Provider: AWS
Runtime: EKS
Kubernetes: 1.34
API: policy/v1beta1 PodSecurityPolicy
Status: removed
Deadline: incompatible
Replacement: Pod Security Admission
```

The platform could use this data before upgrades.

The same model could track:

- EKS-supported Kubernetes versions
- GKE release channels
- AKS versions
- removed APIs
- deprecated ingress APIs
- CSI requirements
- autoscaler compatibility
- controller version requirements

---

## Kubernetes upgrade safety

A production-grade version of this project could include an **upgrade simulator**.

Before moving a cluster from one Kubernetes version to another:

1. inspect workload manifests
2. inspect installed CRDs
3. inspect admission webhooks
4. inspect controllers
5. check removed API versions
6. validate Helm chart compatibility
7. check node-image compatibility
8. check autoscaler compatibility
9. run conformance tests
10. create a rollout plan

Only then should the platform attempt the upgrade.

For customer workloads, the output could be:

```text
Upgrade blocked:
- 3 resources use a removed API
- ingress controller version is incompatible
- one CRD webhook does not support the target version
```

This turns infrastructure upgrades into a testable product workflow rather than an operator-only activity.

---

## Capability negotiation

Not every feature can be truly portable.

The platform therefore needs something richer than booleans over time.

A capability might have states such as:

```text
SUPPORTED
UNSUPPORTED
REQUIRES_CONFIGURATION
REQUIRES_UPGRADE
REGION_LIMITED
QUOTA_LIMITED
DEGRADED
UNKNOWN
```

Example:

```text
GPU:
  status: REGION_LIMITED
  supported_regions:
    - us-east-1
    - us-west-2
```

or:

```text
ManagedIngress:
  status: REQUIRES_CONFIGURATION
  prerequisite: aws-load-balancer-controller
```

This gives the user a simple experience without lying about the infrastructure.

---

## Drift reconciliation

The current drift package compares simple values.

A production drift engine should also understand **ownership**.

Possible ownership modes:

```text
PORTER_MANAGED
CUSTOMER_MANAGED
SHARED
CONTROLLER_MANAGED
```

Example:

- Deployment replicas may be HPA-managed.
- VPC may be customer-managed.
- EKS cluster may be Porter-managed.
- DNS may be shared.

The reconciler should only enforce fields it owns.

This avoids destructive automation.

---

## Cost-aware placement

A future cost engine could compare valid provider plans.

For example:

```text
Intent is conformant on:

AWS:   $1,420/month
GCP:   $1,190/month
Azure: $1,310/month
```

But raw price should not be the only score.

A useful placement score could combine:

```text
cost
+ reliability
+ latency
+ regional availability
+ capacity risk
+ operational complexity
+ existing data gravity
```

That turns the platform from a simple provisioning wrapper into an infrastructure decision engine.

---

## Security and isolation considerations

A real implementation would need strong tenant boundaries.

Important areas include:

- short-lived cloud credentials
- role assumption rather than long-lived keys
- scoped IAM permissions
- workload identity
- encrypted control-plane state
- audit logs
- tenant-level authorization
- secret isolation
- provider API rate limiting
- protection against confused-deputy problems
- secure handling of kubeconfigs
- no plaintext credentials in plans
- network isolation between management plane and customer workloads

Provider adapters should receive only the credentials and permissions required for their operation.

---

## Observability

A production control plane should expose metrics around both platform behavior and provider behavior.

Useful metrics could include:

### Planning

- plan duration
- plan failures
- unsupported capabilities
- provider API latency

### Provisioning

- operation duration
- retry counts
- provider throttling
- resource creation failure rate

### Reconciliation

- drift findings
- auto-repaired drift
- unresolved drift
- reconciliation latency

### Upgrades

- compatibility failures
- upgrade duration
- rollback frequency

### Cost

- estimated vs actual cost
- cost delta after deployment
- idle-resource detection

Tracing would also be important because a single user action could fan out across many provider calls.

---

## Failure handling

Cloud APIs fail frequently for reasons outside the platform's direct control.

A production execution engine should handle:

- rate limiting
- eventual consistency
- asynchronous provider operations
- duplicate requests
- partial failure
- timeout
- permission errors
- quota exhaustion
- unavailable capacity
- conflicting user changes

Important implementation properties would be:

### Idempotency

Retries must not create duplicate infrastructure.

### Checkpointing

Long operations should resume rather than restart from zero.

### Compensation

Partial failures need explicit rollback or repair strategies.

### Error classification

Errors should be classified as:

- retryable
- terminal
- user-action-required
- provider incident
- quota issue
- policy violation

---

## Testing strategy for a production implementation

A production provider layer needs far more than unit tests.

### Unit tests

Fast validation of planning and conformance logic.

### Contract tests

Every provider adapter must satisfy the same behavioral interface.

### Golden plan tests

Known intent should generate expected plans.

### Fake-provider tests

Simulate throttling, failures, quota exhaustion, and eventual consistency.

### Ephemeral integration environments

CI should periodically create temporary environments in:

- AWS
- GCP
- Azure

and validate actual behavior.

### Kubernetes version matrix

Run tests across supported Kubernetes minor versions.

### Upgrade tests

Create an old supported version, deploy workloads, upgrade, verify availability.

### Chaos tests

Inject provider/API failures during reconciliation.

### Drift tests

Modify resources outside the control plane and verify safe detection/reconciliation.

---

## Possible API surface

A production API could look conceptually like:

### Create application intent

```http
POST /v1/applications
```

```json
{
  "name": "api",
  "resources": {
    "cpu": "1",
    "memory": "2Gi"
  },
  "scale": {
    "min": 2,
    "max": 20
  },
  "network": {
    "public": true
  },
  "storage": {
    "size": "20Gi"
  }
}
```

### Plan against providers

```http
POST /v1/applications/api/plan
```

Response:

```json
{
  "plans": [
    {
      "provider": "aws",
      "runtime": "eks",
      "conformant": true,
      "warnings": [],
      "estimated_monthly_cost": 1234
    },
    {
      "provider": "gcp",
      "runtime": "gke",
      "conformant": true,
      "warnings": [],
      "estimated_monthly_cost": 1190
    }
  ]
}
```

### Apply plan

```http
POST /v1/plans/{plan_id}/apply
```

### Inspect drift

```http
GET /v1/applications/api/drift
```

This is illustrative only; it is not implemented in the repository.

---

## Example product workflow

From a user perspective, an eventual product workflow might be:

```text
1. Connect cloud account
2. Point Porter at repository
3. Porter detects application shape
4. User chooses resource intent
5. Porter evaluates provider capabilities
6. Porter shows plan + warnings + cost
7. User deploys
8. Porter provisions infrastructure
9. Porter continuously observes state
10. Porter detects drift
11. Porter tests future upgrades before rollout
12. Porter reconciles safely
```

The user sees a simple workflow.

The control plane handles the provider complexity.

---

## Why Go

Porter's backend is publicly described as Go-based, so Go is appropriate for the prototype.

More importantly, Go is well suited to infrastructure control planes because it offers:

- simple concurrency
- strong static typing
- efficient binaries
- mature Kubernetes libraries
- mature AWS/GCP SDKs
- good tooling
- easy deployment
- strong ecosystem support for controllers and operators

The current repository deliberately uses only the standard library and local packages, keeping the concept easy to inspect.

---

## Current limitations

To be explicit, the current implementation has several simplifications.

### Provider capabilities are static

They are hard-coded in the adapters rather than discovered dynamically.

### Region is not validated

The contract carries a region but provider adapters do not yet translate or validate it.

### Cost is illustrative

The cost function is not tied to real provider billing APIs.

### Drift coverage is partial

Only CPU, memory, and replicas are currently compared.

### No Azure adapter yet

The architecture supports additional providers, but Azure is not currently implemented.

### No persistence

There is no database or state store.

### No reconciliation loop

The system performs one-shot planning rather than continuous reconciliation.

### No live cloud credentials

This is intentional for a public demonstration repository.

### No API server

The demonstration is CLI-driven.

### Capability model is boolean

A production capability model should represent prerequisites, degradation, region limitations, and version dependencies.

### No schema versioning

A production contract should be versioned to preserve API compatibility.

### No asynchronous workflow engine

Cloud operations often take minutes and need durable orchestration.

---

## Roadmap

### Phase 1 - Strengthen the contract

- add schema version
- add environment type
- add availability requirements
- add network policy intent
- add database intent
- add service dependencies
- add secret requirements
- add provider-specific extension block

### Phase 2 - Real provider discovery

- AWS SDK integration
- GCP SDK integration
- Kubernetes discovery client
- region validation
- quota lookup
- cluster version lookup
- storage-class discovery

### Phase 3 - Real planning

- EKS/GKE planning
- node class selection
- ingress planning
- storage planning
- network planning
- IAM planning

### Phase 4 - Rich conformance

- version-aware capabilities
- prerequisite states
- degraded capability support
- policy violations
- quota constraints

### Phase 5 - Cost engine

- AWS pricing catalog
- GCP pricing catalog
- Azure pricing
- network-cost estimation
- workload utilization model
- spot/preemptible support

### Phase 6 - Drift engine

- Kubernetes observed-state adapter
- cloud resource discovery
- ownership metadata
- reconciliation policies

### Phase 7 - Upgrade simulator

- Kubernetes API deprecation scanner
- EKS version matrix
- GKE version matrix
- controller compatibility matrix
- Helm compatibility checks

### Phase 8 - Execution engine

- durable workflows
- idempotent operations
- retries
- rollback
- checkpoints
- operation audit history

### Phase 9 - Ephemeral conformance environments

- temporary AWS environment
- temporary GCP environment
- temporary Azure environment
- automated deploy/test/destroy cycle

### Phase 10 - Product integration

- API server
- UI plan visualization
- provider comparison
- warning explanations
- cost comparison
- upgrade readiness dashboard

---

## Why I built this for Porter

I built this repository because Porter's backend role is not a generic CRUD backend role.

The interesting problem is the boundary between:

```text
simple developer intent
```

and:

```text
complex, changing, provider-specific infrastructure
```

That boundary involves exactly the areas I enjoy working in:

- cloud architecture
- Kubernetes
- infrastructure APIs
- platform engineering
- Terraform and declarative infrastructure concepts
- reliability
- networking
- autoscaling
- observability
- security
- cost optimization
- production operations

Rather than only saying that I am interested in those problems, I wanted to put forward an actual technical model for one of them.

The most important idea in the repository is not the amount of code.

It is the separation of concerns:

```text
Developer Intent
      |
      v
Provider-Neutral Contract
      |
      v
Provider Adapter
      |
      v
Explicit Capability Plan
      |
      v
Shared Conformance Checks
      |
      +--> Cost
      |
      +--> Drift
      |
      +--> Upgrade Safety
```

That is the direction I would explore when building a multi-cloud infrastructure management system that needs to remain simple for developers without ignoring the reality that clouds are different.

---

## Related design documents

Additional short documents are included in:

- [`docs/architecture.md`](docs/architecture.md)
- [`docs/idea.md`](docs/idea.md)

The README is intentionally the complete explanation of the project, while those documents provide shorter focused views.

---

## Author

**Rahul H Bhatia**

Cloud / DevOps / SRE / Platform Engineering

GitHub: [rahulbhatia-rb](https://github.com/rahulbhatia-rb)

This project was created specifically as an engineering demonstration around the problems described by Porter, not as a claim about Porter's internal implementation.