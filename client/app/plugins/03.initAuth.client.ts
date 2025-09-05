export default defineNuxtPlugin(async () => {
  // サーバーサイドでは認証チェックをスキップ
  if (import.meta.server) return;

  // クライアントサイドでのみ実行
  const { fetchUser } = useAuth();

  try {
    // アプリ起動時に認証状態を確認
    await fetchUser();
  }
  catch (error) {
    // 認証エラーは無視（ログイン画面で処理）
    console.log('Initial auth check completed');
  }
});
