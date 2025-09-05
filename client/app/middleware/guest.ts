/**
 * 未認証ユーザー専用ページ用ミドルウェア
 * 認証済みユーザーをホームページにリダイレクトします
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // サーバーサイドでは認証チェックしない
  if (import.meta.server) return;

  console.log(`ゲスト専用ページチェック開始: ${to.path}`);

  const { isLoggedIn, checkAuthStatus, isLoading } = useAuth();

  // まだ認証状態が確認されていない場合
  if (!isLoggedIn.value && !isLoading.value) {
    console.log('認証状態未確認のため、/me API を呼び出し');
    await checkAuthStatus();
  }

  // すでにログイン済みの場合はホームへリダイレクト
  if (isLoggedIn.value) {
    console.log('認証済みユーザーのため、ホームページにリダイレクト');
    return navigateTo('/');
  }

  console.log(`未認証ユーザーのページアクセス許可: ${to.path}`);
});
