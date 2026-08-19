# pf-developer-templates

[pf-developer-cli](https://github.com/maeplego/pf-developer-cli)（`pf-dev new`）がコピーするツリーです。ワークスペース API とコマース API の小さなスライスを元にしており、空のスタブではありません。

| テンプレート | 内容 |
| --- | --- |
| `go-api` | Go HTTP。`/health`、`/ready`、商品カタログ、OIDC の userinfo スタブ |
| `go-next` | 同じ API に加え、Next.js（ヘルスと PKCE ログインのスタブ） |

各ディレクトリに `template.json` があります。CLI が `{{MODULE}}` などを置換します。コピー後のシェルは走りません。

OpenAPI と、PR で breaking change を落とす GitHub Actions 例も入っています。スキャナやポータル本体はこのリポジトリにはありません。
