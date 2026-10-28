# State: nsproxy

Overwritten to stay current (loopers and agents). Last updated: 2026-10-28 by looper (issue #8).

- **Version / line:** v0.1.4 (https://github.com/tacplus-dev/nsproxy/releases/tag/v0.1.4), line 0.1; Go floor go 1.27.3 (go.mod); no tacplus.yaml
- **Works:** the `TenantNamespace` controller: tenant namespaces `tac-<deployment>-<tid>` with RoleBinding, NetworkPolicies, ResourceQuota and shared-tenant (`restore`) traffic; gates `make check` (envtest) and `make image`; govulncheck reports no vulnerabilities (source and binary)
- **Released / deployed:** image ghcr.io/tacplus-dev/nsproxy:0.1.4 and line tag 0.1 (linux/amd64, linux/arm64)
- **Known gaps:** none known
- **Blocked on:** nothing
