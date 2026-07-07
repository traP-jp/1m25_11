// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ['@nuxt/ui', '@nuxt/eslint', '@nuxt/image'],
  ssr: false,
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  ui: {
    colorMode: false,
  },
  runtimeConfig: {
    backendApiUrl: '',
    proxySecret: '',
    devUser: '',
    public: {
      userName: 'traP',
      apiBase: '',
    },
  },
  compatibilityDate: '2025-07-15',
  vite: {
    server: {
      allowedHosts: [
        '.ngrok-free.app',
      ],
    },
    optimizeDeps: {
      include: [
        '@internationalized/date',
        '@vue/devtools-core',
        '@vue/devtools-kit',
        'openapi-fetch',
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
