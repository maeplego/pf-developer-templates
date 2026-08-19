# pf-developer-templates

P11 scaffolding trees for `pf-dev new`. Files are adapted from **P04** `pf-workspace/apps/api` (and `apps/web` OIDC) and **P06** `pf-commerce/apps/api` (catalog + integer money + httptest). They are not empty stubs.

| Template | What you get |
| --- | --- |
| `go-api` | Go HTTP API: `/health`, `/ready`, OTel env, OIDC userinfo stub, catalog |
| `go-next` | Same API under `apps/api` plus Next.js with health/ready and PKCE login stub |

Each directory has `template.json`. The CLI copies the tree and substitutes `{{MODULE}}`, `{{PROJECT}}`, `{{HTTP_PORT}}`. There is no postInstall shell.

Scanner / portal / CI dashboard / review live in other `pf-developer-*` repos. `go-api` and `go-next` include `openapi.yaml` and `.github/workflows/openapi-breaking.yml` (oasdiff-action on PRs).
