import { useState } from '#app';
import { apiClient } from './ApiClient';

// openapi-fetchで生成された型を使用
type User = Schemas['UserProfile'];

export const useAuth = () => {
  const user = useState<User | null>('auth.user', () => null);
  const isAuthenticated = computed(() => user.value !== null);

  // アプリケーション読み込み時にGoサーバーに自身の情報を問い合わせる
  const fetchUser = async () => {
    // 既にユーザー情報がある場合は処理をスキップ
    if (user.value) return user.value;

    try {
      // APIクライアントを使用してサーバーの認証状態を確認
      // Cookieはブラウザが自動で送信される (credentials: 'include' 設定済み)
      const { data, error } = await apiClient.GET('/me');

      if (error) {
        throw new Error(`API Error: ${error}`);
      }

      if (data) {
        user.value = data;
        return data;
      }

      throw new Error('No user data received');
    }
    catch (error) {
      // 401 Unauthorizedなどが返ってきた場合は未ログイン状態
      user.value = null;
      console.warn('Authentication check failed:', error);
      return null;
    }
  };

  // ログアウト処理
  const logout = async () => {
    try {
      // APIクライアントを使用してログアウト（今後実装予定）
      // await apiClient.POST('/logout');
    }
    catch (error) {
      console.warn('Logout request failed:', error);
    }
    finally {
      // ローカルの認証状態をクリア
      user.value = null;
      // ログインページにリダイレクト
      await navigateTo('/login');
    }
  };

  // ログイン開始（traQにリダイレクト）
  const login = () => {
    const config = useRuntimeConfig();
    window.location.href = config.public.loginUrl;
  };

  return {
    user: readonly(user),
    isAuthenticated,
    fetchUser,
    logout,
    login,
  };
};
