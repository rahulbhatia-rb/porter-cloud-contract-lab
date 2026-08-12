# Architecture

```text
Developer intent
    |
Application Intent Contract
    |
+-------------------------------+
| Policy / Capability Evaluator |
+-------------------------------+
    |          |          |
 AWS Adapter  GCP Adapter  K8s Adapter
    |          |          |
 EKS/IAM/VPC  GKE/IAM/VPC Kubernetes API
    \          |          /
       Conformance Suite
              |
       Drift + Cost Checks
```

The prototype deliberately separates **intent** from **provider implementation**. Provider-specific capabilities may differ, but the contract remains stable and incompatibilities are surfaced before rollout.
