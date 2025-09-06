import createClient from 'openapi-fetch';
import type { paths } from '~~/shared/types/generated';

type ApiClient = ReturnType<typeof createClient<paths>>;

// モジュールスコープでクライアントのインスタンスを保持する変数を用意
let apiClientInstance: ApiClient | null = null;

export const useApiClient = () => {
  // インスタンスがまだ作成されていなければ作成する
  if (!apiClientInstance) {
    const config = useRuntimeConfig();

    console.log(config.public.apiBase);

    apiClientInstance = createClient<paths>({
      baseUrl: config.public.apiBase,
      credentials: 'include', // クロスオリジンでもCookieを送信
    });

    // 認証エラー時のハンドリング
    apiClientInstance.use({
      onResponse({ response }) {
        // 401エラーの場合は認証状態をリセット
        if (response.status === 401) {
          console.warn('認証エラー (401): 認証状態をリセットします');

          // グローバル認証状態をリセット（クライアントサイドでのみ）
          if (!import.meta.server) {
            const { logout } = useAuth();
            logout();
          }
        }
      },
    });
  }

  return apiClientInstance;
};
