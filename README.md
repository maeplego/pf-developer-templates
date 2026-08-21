# pf-developer-templates

[pf-developer-cli](https://github.com/maeplego/pf-developer-cli)（`pf-dev new`）がコピーするツリーです。ワークスペース API とコマース API の小さなスライスを元にしており、空のスタブではありません。

| テンプレート | 内容 |
| --- | --- |
| `go-api` | Go HTTP。`/health`、`/ready`、商品カタログ、OIDC の userinfo スタブ |
| `go-next` | 同じ API に加え、Next.js（ヘルスと PKCE ログインのスタブ） |

各ディレクトリに `template.json` があります。CLI が `{{MODULE}}` などを置換します。コピー後のシェルは走りません。

OpenAPI と、PR で breaking change を落とす GitHub Actions 例も入っています。スキャナやポータル本体はこのリポジトリにはありません。

## ライセンスと利用条件

本リポジトリは **デモ・学習・社内評価用** です。現状品質に **保証はありません**。

- 許可: クローン、ローカル実行、学習、非本番の評価
- 別契約が必要: 本番運用、有償サービスへの組込み、再販・托管の提供

詳細は [LICENSE](./LICENSE) と [licensing.md](https://github.com/maeplego/portfolio-plan/blob/master/portfolio-plan/licensing.md) を参照してください。

