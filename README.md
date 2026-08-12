# Porter Cloud Contract Lab

A small Go proof-of-concept for testing a **cloud-neutral application infrastructure contract** against multiple provider adapters.

> Independent engineering prototype inspired by Porter's publicly described product and backend engineering challenges. It does not use or claim knowledge of Porter's private code or infrastructure.

## The problem

A multi-cloud PaaS has to offer a stable developer experience while AWS, GCP, Azure and Kubernetes continuously evolve underneath it. The hard part is not merely calling provider APIs; it is preserving behavioral compatibility as provider capabilities, defaults, versions and costs diverge.

This repo demonstrates one way to make that problem testable:

1. Express application requirements as an `ApplicationIntent` contract.
2. Translate the contract through provider adapters.
3. Run a shared conformance suite against every plan.
4. Surface capability mismatches rather than hiding them.
5. Compare cost and drift independently from the provider adapter.

## Run

```bash
go test ./...
go run ./cmd/contractlab ./examples/api.json
```

Example output:

```text
aws        runtime=eks        conformant=true  ...
gcp        runtime=gke        conformant=true  ...
kubernetes runtime=kubernetes conformant=true ...
```

## What I would build next

- Real AWS/GCP provider API clients behind interfaces
- EKS/GKE/AKS version and API-deprecation compatibility tests
- Ephemeral cloud conformance environments in CI
- Live desired-vs-actual drift detection
- Policy-aware cost comparison
- Capability negotiation for non-portable features
- Upgrade simulation before rolling cluster/Kubernetes versions

See [`docs/idea.md`](docs/idea.md) for the full proposal.

## Why this is relevant to Porter

The role describes owning infrastructure management, maintaining cloud-agnostic APIs, integrating new hyperscaler capabilities, and preserving compatibility with services outside the platform. This project focuses directly on that control-plane boundary.

## Author

Rahul H Bhatia
