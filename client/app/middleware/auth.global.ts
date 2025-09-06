export default defineNuxtRouteMiddleware((to) => {
  // ログインページやコールバックページへの遷移は許可
  if (to.path === '/login') {
    return;
  }

  const user = useState<Schemas['UserStatus'] | null>('user');

  if (!user.value) {
    // サーバーサイドでは外部URLにリダイレクト
    if (import.meta.server) {
      // const config = useRuntimeConfig();
      // return navigateTo(`${config.public.apiBase}/login`, { external: true });
      return;
    }
    // クライアントサイドではwindow.locationを使ってリダイレクト
    else {
      const config = useRuntimeConfig();
      window.location.href = `${config.public.apiBase}/login`;
      return;
    }
  }
});
