# 参考 : https://content.nuxt.com/docs/deploy/docker#with-nodejs-image

# Build Stage 1
FROM node:22-alpine AS build
WORKDIR /app

# package.jsonとlockファイルをコピー
COPY package.json package-lock.json ./

# 依存関係をインストール
RUN npm install

# コピー
COPY . ./

# ビルドする
RUN npm run build

# Build Stage 2
FROM node:22-alpine
WORKDIR /app

# Only `.output` folder is needed from the build stage
COPY --from=build /app/.output/ ./

# Change the port and host
ENV PORT=80
ENV HOST=0.0.0.0

EXPOSE 80

CMD ["node", "/app/server/index.mjs"]
