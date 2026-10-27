# nsproxy

Creates and deletes **tenant namespaces** for the deployments on a Kubernetes cluster, so that a deployment
needs no privilege outside its own namespace.

Several deployments share one cluster. The top pipeline creates each deployment's namespace,
`tac-<deployment>`, and deploys nsproxy once per cluster. A deployment's operator, running in
`tac-<deployment>`, asks nsproxy to create or delete its tenants' namespaces, `tac-<deployment>-<tid>`
(a tid is 8 characters of `[0-9a-z]`; all-digit ids are reserved, `people` is `00000001`, `restore` is
`00000002`). nsproxy is the only
thing holding the cluster-wide permission, and it uses it narrowly.

Work arrives as issues; see github.com/tacplus-dev/looper.

## How it works

nsproxy is a Kubernetes controller for one custom resource, `TenantNamespace` (`nsproxy.tacplus.dev/v1`). It
has no network listener. The operator creates one in its own namespace:

```yaml
apiVersion: nsproxy.tacplus.dev/v1
kind: TenantNamespace
metadata:
  name: acme
  namespace: tac-prod
spec:
  tid: a1b2c3d4   # 8 of [0-9a-z], immutable
```

and waits for `status.conditions` `Ready=True`; `status.namespace` is then `tac-prod-a1b2c3d4`. When it is
refused, `Ready=False` and the reason says why (`NotDeploymentNamespace`, `InTenantNamespace`,
`InvalidDeployment`, `ReservedTid`, `NamespaceConflict`). Deleting the object deletes the namespace, at once.

- **Who may ask** is Kubernetes RBAC: the Role in [config/deployment-namespace](config/deployment-namespace),
  applied by the pipeline in each deployment namespace, lets that namespace's `operator` service account
  create `TenantNamespace` objects there and nowhere else.
- **The deployment** comes from the parent namespace's `tacplus.dev/deployment` label (set by the pipeline),
  never from its name. Objects in an unlabelled namespace, or in a tenant namespace, are refused.
- **The tenant namespace** carries `nsproxy.tacplus.dev/deployment` and `nsproxy.tacplus.dev/tid` labels and
  holds: a RoleBinding making `operator` of the deployment namespace `admin` there; NetworkPolicies denying
  everything except traffic inside the namespace, DNS, egress to the deployment namespace, and traffic with
  the deployment's shared tenants; a ResourceQuota from the ConfigMap, always with no load balancers and no
  node ports. Anything changed or deleted by hand, the namespace included, is put back.
- **Shared tenants** (`sharedTids` in the ConfigMap; shipped: `restore`, `00000002`) are reserved tenants every
  tenant of the same deployment may call: each tenant namespace of a deployment gets egress to
  `tac-<deployment>-<shared tid>`, and that namespace accepts ingress from them. Both sides are chosen by
  nsproxy's own labels, so nothing crosses deployments; every other tenant namespace is reachable from no
  other. `people` is not shared: the deployment namespace calls it.
- **Only what it owns:** nsproxy changes or deletes a namespace only if it carries its owner label for the same
  deployment; an existing namespace of the same name without it is refused, never adopted. It never touches a
  deployment namespace. Deleting a deployment namespace deletes all its tenant namespaces.

## Manifests

| | |
|---|---|
| `config/crd` | the `TenantNamespace` CRD (schema, `tid` pattern, CEL immutability rule) |
| `config/rbac` | nsproxy's service account, ClusterRole and binding |
| `config/manager` | its namespace, Deployment and ConfigMap (reserved and shared tids, quota default and per-deployment overrides) |
| `config/deployment-namespace` | the Role the pipeline applies in each deployment namespace: `NAMESPACE=tac-<deployment> envsubst '$NAMESPACE' < config/deployment-namespace/operator.yaml \| kubectl apply -f -` |

## Build and test

`make check` runs the gates: gofmt, `go vet`, build, and the tests. The controller tests run a real API server
([envtest](https://book.kubebuilder.io/reference/envtest), fetched by `setup-envtest`) with the CRD, and run the
controller as nsproxy's service account under the RBAC manifests above. `make image` builds the container image.
