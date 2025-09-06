/**
 * グローバル認証ミドルウェア
 * 全てのページで自動実行され、認証が必要なページで認証チェックを行います
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // サーバーサイドでは認証チェックしない（Hydration エラー回避）
  if (import.meta.server) return;

  console.log(`グローバル認証チェック開始: ${to.path}`);

  // 認証不要ページ（公開ページ）の除外リスト
  const publicPages = [
    '/', // ホームページ（認証オプショナル・ログイン誘導あり）
  ];

  // 公開ページの場合は認証チェックのみ（リダイレクトしない）
  if (publicPages.includes(to.path)) {
    console.log(`公開ページアクセス: ${to.path} - 認証チェックのみ実行`);

    const { isLoggedIn, checkAuthStatus, isLoading } = useAuth();

    // ローディング中でない場合のみ認証状態確認
    if (!isLoading.value && !isLoggedIn.value) {
      console.log('認証状態未確認のため、確認処理を実行');
      await checkAuthStatus();
    }

    const authStatus = isLoggedIn.value ? '認証済み' : '未認証';
    console.log(`公開ページ認証状態: ${authStatus}`);
    return;
  }

  // その他のページは認証必須
  console.log(`認証必須ページアクセス: ${to.path}`);

  const { isLoggedIn, checkAuthStatus, isLoading } = useAuth();

  // ローディング完了まで待機
  if (isLoading.value) {
    console.log('認証状態確認の完了を待機中...');

    // isLoadingがfalseになるまで待つ
    await new Promise<void>((resolve) => {
      const unwatch = watch(isLoading, (newValue) => {
        if (!newValue) {
          unwatch();
          resolve();
        }
      }, { immediate: true });
    });
  }

  // 認証状態が未確認の場合は確認
  if (!isLoggedIn.value) {
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
