/**
 * 認証状態管理用のcomposable
 * サーバーが完全認証必須のため、/me エンドポイントでログイン状態を確認
 */

// グローバルな認証状態（シングルトン）
const globalAuthState = {
  isLoggedIn: ref(false),
  currentUser: ref<Schemas['UserStatus'] | null>(null),
  // サーバーサイドでは常にfalse、クライアントサイドで動的に変更
  isLoading: ref(false),
};

export const useAuth = () => {
  // グローバル状態を使用
  const { isLoggedIn, currentUser, isLoading } = globalAuthState;

  /**
   * 認証状態をチェック
   * /me エンドポイントが成功すればログイン済み、401エラーなら未ログイン
   * クライアントサイドでのみ実行
   */
  const checkAuthStatus = async () => {
    // サーバーサイドでは実行しない
    if (import.meta.server) return;

    try {
      isLoading.value = true;
      console.log('認証状態チェック開始');
      const result = await useApiClient().GET('/me');
      console.log('API呼び出し完了:', result);

      if (result.data && !result.error) {
        isLoggedIn.value = true;
        // UserStatusをそのまま保存
        currentUser.value = {
          user_id: result.data.user_id,
          is_admin: result.data.is_admin || false,
          stamps_user_owned: result.data.stamps_user_owned || [],
          tags_user_created: result.data.tags_user_created || [],
          stamps_user_tagged: result.data.stamps_user_tagged || [],
          descriptions_user_created: result.data.descriptions_user_created || [],
        };
        console.log('認証成功: ユーザーID', result.data.user_id);
        console.log('グローバル状態更新 - isLoggedIn:', isLoggedIn.value);
      }
      else {
        isLoggedIn.value = false;
        currentUser.value = null;
        console.log('未認証状態');
      }
    }
    catch (err) {
      // 401エラーなどの場合は未ログイン
      isLoggedIn.value = false;
      currentUser.value = null;
      console.log('認証エラー:', err);
    }
    finally {
      console.log('認証チェック完了、ローディング解除');
      isLoading.value = false;
      console.log('最終状態 - isLoggedIn:', isLoggedIn.value, 'isLoading:', isLoading.value);
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
