# stamPedia

25年 1-Monthon 11 班「**stamPedia**」

ハッカソンの開発メンバー

- 

## 使用技術

- クライアント
  - [Nuxt4](https://nuxt.com/)
  - [Vue3](https://vuejs.org/)
  - [NuxtUI](https://ui.nuxt.com/)
  - [openapi-typescript](https://openapi-ts.dev/ja/introduction)
  - [openapi-fetch](https://openapi-ts.dev/ja/openapi-fetch/)
- サーバー
  - Go
  - Echo

### 機能

traQで使用可能な9000個以上のスタンプについて
- タグ付け
- 説明文の投稿

を行うことで、

## 開発環境の準備


### クライアント

クライアントディレクトリに移動し、依存関係をインストールする。

```bash
cd client
npm install
```

環境変数は `.env.example` を参考に `.env` に設定する。

```
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
# サーバーのAPIのベースとなるURL
```

設定が完了したら、開発サーバーを起動する。

```bash
npm run dev
```

> [!NOTE]
> OpenAPIの `/api/swagger.yml` をもとに、TypeScriptの型定義を自動生成している。
> APIの仕様を変更した際は、次のコマンドを実行する。
> ```bash
> npm run generate:openapi
> npm run lint:fix
> ```

### サーバー

https://github.com/ras0q/go-backend-template/ を使用させていただいた

環境変数は `.env.example` を参考に `.env` に設定する。

```bash
BOT_TOKEN_KEY=bot_token_key_here # traQ Bot のトークン (traQ API を利用するため)
CLIENT_ID=client_id_here # traQ OAuth2 Client ID
TOP_PAGE_URL=top_page_url_here # トップページのURL (OAuth2 認証後にリダイレクトする先)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080　# CORSで許可するオリジン (カンマ区切り)
```

開発時は `docker compose watch` でホットリロード付きのサーバーを起動できる。

- app
  - Go製のサーバーアプリケーション
- db
  - MariaDB
  - 本番環境では [NeoShowcase](https://github.com/traPtitech/NeoShowcase/) のデータベースを使用する。開発環境では `docker-compose.yml` で立ち上げた MariaDB を使用する。
- adminer
  - データベース管理画面