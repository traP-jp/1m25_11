// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ['@nuxt/ui', '@nuxt/eslint', '@nuxt/image'],
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  ui: {
    colorMode: false,
  },
  runtimeConfig: {
    // private keys (server only) を追加する場合はここに記述
    public: {
      userName: 'traP',
      // APIベースURL（環境変数で上書き可能）
      apiBase: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1',
      // ログイン開始URL
      loginUrl: process.env.NUXT_PUBLIC_LOGIN_URL || 'http://localhost:8080/api/v1/login',
      // フロントトップページURL（callback後の遷移先と整合させる）
      topPageUrl: process.env.NUXT_PUBLIC_TOP_PAGE_URL || 'http://localhost:3000',
    },
  },
  compatibilityDate: '2025-07-15',
  vite: {
    server: {
      allowedHosts: [
        '.ngrok-free.app',
      ],
    },
  },
  eslint: {
    config: {
      stylistic: {
        semi: true,
        indent: 2,
        quotes: 'single',
        commaDangle: 'always-multiline',
      },
    },
  },
  fonts: {
    families: [
      { name: 'Zen Kaku Gothic New', provider: 'google' },
    ],
  },
});
