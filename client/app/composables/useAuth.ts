/**
 * 認証状態管理用のcomposable
 * サーバーが完全認証必須のため、/me エンドポイントでログイン状態を確認
 */
export const useAuth = () => {
  const isLoggedIn = ref(false);
  const currentUser = ref<Schemas['UserStatus'] | null>(null);
  const isLoading = ref(true);

  /**
   * 認証状態をチェック
   * /me エンドポイントが成功すればログイン済み、401エラーなら未ログイン
   */
  const checkAuthStatus = async () => {
    try {
      isLoading.value = true;
      const { data, error } = await useApiClient().GET('/me');

      if (data && !error) {
        isLoggedIn.value = true;
        currentUser.value = data;
        console.log('ログイン済み:', data);
      }
      else {
        isLoggedIn.value = false;
        currentUser.value = null;
        console.log('未ログインまたはエラー');
      }
    }
    catch (err) {
      // 401エラーなどの場合は未ログイン
      isLoggedIn.value = false;
      currentUser.value = null;
      console.log('認証エラー:', err);
    }
    finally {
      isLoading.value = false;
    }
  };

  /**
   * ログイン処理（サーバーの /login エンドポイントにリダイレクト）
   */
  const login = () => {
    const config = useRuntimeConfig();
    window.location.href = `${config.public.apiBase}/login`;
  };

  /**
   * ログアウト処理（簡易版）
   * 実際のログアウトはサーバー側で実装が必要
   */
  const logout = () => {
    isLoggedIn.value = false;
    currentUser.value = null;
    // TODO: サーバー側にログアウトエンドポイントが必要な場合は実装
  };

  return {
    isLoggedIn: readonly(isLoggedIn),
    currentUser: readonly(currentUser),
    isLoading: readonly(isLoading),
    checkAuthStatus,
    login,
    logout,
  };
};
