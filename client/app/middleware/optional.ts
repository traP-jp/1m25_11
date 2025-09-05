/**
 * 認証オプショナルページ用ミドルウェア
 * 認証状態をチェックするだけで、リダイレクトは行いません
 * ページ内で認証状態に応じて表示を切り替える場合に使用します
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // サーバーサイドでは認証チェックしない
  if (import.meta.server) return;

  console.log(`認証オプショナルページチェック開始: ${to.path}`);

  const { isLoggedIn, checkAuthStatus, isLoading } = useAuth();

  // まだ認証状態が確認されていない場合のみチェック
  if (!isLoggedIn.value && !isLoading.value) {
    console.log('認証状態未確認のため、/me API を呼び出し');
    await checkAuthStatus();
  }

  const status = isLoggedIn.value ? '認証済み' : '未認証';
  console.log(`認証オプショナルページアクセス (${status}): ${to.path}`);
});
