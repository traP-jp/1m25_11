/**
 * 認証必須ページ用ミドルウェア
 * 未認証ユーザーをログインページにリダイレクトします
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // サーバーサイドでは認証チェックしない
  if (import.meta.server) return;

  console.log(`認証チェック開始: ${to.path}`);

  const { isLoggedIn, checkAuthStatus, isLoading } = useAuth();

  // まだ認証状態が確認されていない場合
  if (!isLoggedIn.value && !isLoading.value) {
    console.log('認証状態未確認のため、/me API を呼び出し');
    await checkAuthStatus();
  }

  // 認証されていない場合はログインページへリダイレクト
  if (!isLoggedIn.value) {
    console.log('未認証のため、ログインページにリダイレクト');
    const config = useRuntimeConfig();
    return navigateTo(`${config.public.apiBase}/login`, {
      external: true,
    });
  }

  console.log(`認証済みユーザーのページアクセス許可: ${to.path}`);
});
