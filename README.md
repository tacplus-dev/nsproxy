# nsproxy

Creates and deletes **tenant namespaces** for the deployments on a Kubernetes cluster, so that a deployment
needs no privilege outside its own namespace.

Several deployments share one cluster. The top pipeline creates each deployment's namespace,
`tac-<deployment>`, and deploys nsproxy once per cluster. A deployment's operator, running in
`tac-<deployment>`, calls nsproxy to create or delete its tenants' namespaces, `tac-<deployment>-<tid>`
(a tid is 8 characters of `[0-9a-z]`; all-digit ids are reserved, `people` is `00000001`). nsproxy is the only
thing holding the cluster-wide permission, and it uses it narrowly.

Work arrives as issues; see github.com/tacplus-dev/looper.
