export default defineNuxtPlugin(async () => {
  // 認証状態の初期化
  const { checkAuthStatus } = useAuth();

  // ページ読み込み時に認証状態を確認
  await checkAuthStatus();

  console.log('認証状態の初期化が完了しました');
});
