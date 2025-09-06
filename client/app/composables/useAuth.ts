/**
 * 認証状態管理用のcomposable
 * サーバーが完全認証必須のため、/me エンドポイントでログイン状態を確認
 */

// グローバルな認証状態（シングルトン）
const globalAuthState = {
  isLoggedIn: ref(false),
  currentUser: ref<Schemas['UserStatus'] | null>(null),
  // 初期状態はローディング中（Hydration エラー回避のため）
  isLoading: ref(true),
  // エラー状態を追加
  error: ref<string | null>(null),
  // 最後にチェックした時刻
  lastChecked: ref<Date | null>(null),
};

export const useAuth = () => {
  // グローバル状態を使用
  const { isLoggedIn, currentUser, isLoading, error, lastChecked } = globalAuthState;

  /**
   * 認証状態をチェック
   * /me エンドポイントが成功すればログイン済み、401エラーなら未ログイン
   * クライアントサイドでのみ実行
   */
  const checkAuthStatus = async (): Promise<void> => {
    // サーバーサイドでは実行しない
    if (import.meta.server) return;

    try {
      isLoading.value = true;
      error.value = null;

      console.log('🔐 認証状態チェック開始');
      const apiClient = useApiClient();
      const result = await apiClient.GET('/me');
      console.log('📡 /me API レスポンス:', {
        hasData: !!result.data,
        hasError: !!result.error,
        errorStatus: result.error?.status,
      });

      if (result.data && !result.error) {
        // 認証成功
        isLoggedIn.value = true;
        currentUser.value = {
          user_id: result.data.user_id,
          is_admin: result.data.is_admin || false,
          stamps_user_owned: result.data.stamps_user_owned || [],
          tags_user_created: result.data.tags_user_created || [],
          stamps_user_tagged: result.data.stamps_user_tagged || [],
          descriptions_user_created: result.data.descriptions_user_created || [],
        };
        lastChecked.value = new Date();

        console.log('✅ 認証成功:', {
          userId: result.data.user_id,
          isAdmin: result.data.is_admin,
          ownedStamps: result.data.stamps_user_owned?.length || 0,
        });
      }
      else {
        // 認証失敗またはエラーレスポンス
        isLoggedIn.value = false;
        currentUser.value = null;

        if (result.error) {
          error.value = `認証エラー: ${result.error.status} ${result.error.statusText || ''}`;
          console.log('❌ 認証失敗:', {
            status: result.error.status,
            statusText: result.error.statusText,
          });
        }
        else {
          error.value = '認証レスポンスが無効です';
          console.log('❌ 認証失敗: 無効なレスポンス');
        }
      }
    }
    catch (err) {
      // ネットワークエラーなどの例外
      isLoggedIn.value = false;
      currentUser.value = null;
      error.value = err instanceof Error ? err.message : '認証チェック中にエラーが発生しました';

      console.error('🚨 認証チェック例外:', {
        error: err,
        message: error.value,
      });
    }
    finally {
      isLoading.value = false;
      console.log('🏁 認証チェック完了:', {
        isLoggedIn: isLoggedIn.value,
        isLoading: isLoading.value,
        hasError: !!error.value,
        timestamp: new Date().toISOString(),
      });
    }
  };

  /**
   * 認証状態をリセット（ログアウト時など）
   */
  const resetAuthState = (): void => {
    console.log('🔄 認証状態をリセット');
    isLoggedIn.value = false;
    currentUser.value = null;
    error.value = null;
    lastChecked.value = null;
    // isLoading は変更しない（他の処理で制御される）
  };

  /**
   * ログイン処理（サーバーの /login エンドポイントにリダイレクト）
   */
  const login = (): void => {
    console.log('🚀 ログインページにリダイレクト');
    const config = useRuntimeConfig();
    window.location.href = `${config.public.apiBase}/login`;
  };

  /**
   * ログアウト処理（認証状態のリセット）
   */
  const logout = (): void => {
    console.log('👋 ログアウト処理実行');
    resetAuthState();
    // TODO: サーバー側にログアウトエンドポイントがある場合は実装
  };

  /**
   * 認証の有効性をチェック（一定時間経過後の再チェックなど）
   */
  const isAuthValid = (): boolean => {
    if (!isLoggedIn.value || !lastChecked.value) return false;

    // 5分以内にチェックしていれば有効とみなす
    const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000);
    return lastChecked.value > fiveMinutesAgo;
  };

  return {
    // 読み取り専用として公開
    isLoggedIn: readonly(isLoggedIn),
    currentUser: readonly(currentUser),
    isLoading: readonly(isLoading),
    error: readonly(error),
    lastChecked: readonly(lastChecked),

    // メソッド
    checkAuthStatus,
    resetAuthState,
    login,
    logout,
    isAuthValid,
  };
};
