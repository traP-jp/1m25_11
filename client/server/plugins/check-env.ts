export default defineNitroPlugin(() => {
  const config = useRuntimeConfig();

  if (!config.backendApiUrl) {
    console.error('[FAIL] NUXT_BACKEND_API_URL: 未設定（Nitro プロキシが 500 を返す）');
  }
  else {
    console.log('[OK]   NUXT_BACKEND_API_URL');
  }

  const apiBase = config.public.apiBase as string | undefined;
  if (!apiBase) {
    console.error('[FAIL] NUXT_PUBLIC_API_BASE: 未設定');
  }
  else if (apiBase.startsWith('http://') || apiBase.startsWith('https://')) {
    console.warn(`[WARN] NUXT_PUBLIC_API_BASE=${apiBase}（ローカル開発設定のまま）`);
  }
  else {
    console.log(`[OK]   NUXT_PUBLIC_API_BASE=${apiBase}`);
  }
});
