export default defineNuxtPlugin(async () => {
  // クライアントサイドでのみ認証状態を初期化
  if (import.meta.server) return;

  const { checkAuthStatus } = useAuth();

  // 認証状態をチェック
  await checkAuthStatus();
});
