// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ['@nuxt/ui', '@nuxt/eslint'],
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  ui: {
    colorMode: false,
  },
  runtimeConfig: {
    public: {
      userName: 'traP',
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
