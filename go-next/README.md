# {{PROJECT}}

`pf-dev new -t go-next` が生成した API + Next.js です。API は商品カタログとヘルス、Web はヘルスと PKCE ログインのスタブです。学習用です。

```powershell
go test ./...
# apps/api から

cd apps/web
npm test
```

起動は API で `DEV_AUTH=true` の `go run`、別ターミナルで `apps/web` の `npm run dev` です。
