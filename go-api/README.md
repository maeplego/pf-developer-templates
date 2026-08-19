# {{PROJECT}}

`pf-dev new` が生成した HTTP API です。ヘルス、商品カタログ、開発用ヘッダ認証、OIDC userinfo のスタブまでです。本番の IdP や店舗ではありません。

```powershell
copy .env.example .env
$env:DEV_AUTH="true"
go test ./...
go run ./cmd/server
```

http://localhost:{{HTTP_PORT}}/health と `/v1/products` です。Compose は `deploy/compose.yaml` です。

| 変数 | 意味 |
| --- | --- |
| `HTTP_PORT` | 待ち受け（既定 {{HTTP_PORT}}） |
| `DEV_AUTH` | `X-Dev-User-Sub` を使う |
| `OIDC_ISSUER` など | 開発認証オフのときの userinfo |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | 観測（無くても起動する） |
