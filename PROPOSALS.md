# Proposals: nsproxy

Written by loopers. Each proposal: what, why, links. Jody promotes accepted ones into ROADMAP.md (via an agent);
the looper then marks it `promoted` or `declined` here with the date.

## Open
- **golang.org/x/net v0.62.0** (2026-10-28, #7): govulncheck finds GO-2026-6866 and GO-2026-6860 (html package)
  in golang.org/x/net v0.60.0, fixed in v0.62.0. nsproxy does not call them (Module Results only, not Symbol
  Results), so it is not a bug today; bumping keeps the next scan clean.
- **govulncheck in the gates** (2026-10-28, #7): #7 was found by the nightly scan, not by CI; running
  `govulncheck ./...` in `make check` would fail a push that reaches a known vulnerability. It would also fail
  CI when a new advisory lands without any change, which is the trade-off to decide.

## Promoted / declined
