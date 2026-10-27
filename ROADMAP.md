# Roadmap: nsproxy

Owned by Jody; changed only by an agent at Jody's request. Loopers propose in PROPOSALS.md.

## Now
nsproxy, one per cluster: a Kubernetes controller for `TenantNamespace` objects that creates and deletes a deployment's tenant namespaces (`tac-<deployment>-<tid>`) with baseline isolation, so no deployment needs cluster-wide rights.
Latest release v0.1.2 (2026-10-27), line 0.1.

## Next
- No open issues.

## Later
- Nothing recorded yet.

## Not doing
- Deployment namespaces (`tac-<deployment>`): the top pipeline creates them; nsproxy handles tenant namespaces only (decisions/deployments-per-cluster-2026-10-26.md)
- Keeping a tenant's data on deletion: deletion is immediate; keeping data is the caller's responsibility (decisions/deployments-per-cluster-2026-10-26.md)
