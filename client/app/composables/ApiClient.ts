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

    // 認証エラー時のハンドリング強化
    apiClientInstance.use({
      onResponse({ response }) {
        console.log('📡 API レスポンス:', {
          url: response.url,
          status: response.status,
          statusText: response.statusText,
        });

        // 401エラーの場合は認証状態をリセット
        if (response.status === 401) {
          console.warn('🚨 認証エラー (401): 認証状態をリセットします');

          // グローバル認証状態をリセット（クライアントサイドでのみ）
          if (!import.meta.server) {
            try {
              const { resetAuthState } = useAuth();
              resetAuthState();
            }
            catch (error) {
              console.warn('認証状態のリセットに失敗しました:', error);
            }
          }
        }

        // その他のエラーステータスのログ
        if (response.status >= 400) {
          console.warn(`⚠️ API エラー (${response.status}):`, {
            url: response.url,
            status: response.status,
            statusText: response.statusText,
          });
        }
      },

      onRequest({ request }) {
        console.log('📤 API リクエスト:', {
          url: request.url,
          method: request.method,
        });
      },

      onRequestError({ request, error }) {
        console.error('🚨 API リクエストエラー:', {
          url: request.url,
          method: request.method,
          error: error.message,
        });
      },

      onResponseError({ response, error }) {
        console.error('🚨 API レスポンスエラー:', {
          url: response.url,
          status: response.status,
          statusText: response.statusText,
          error: error.message,
        });
      },
    });
  }

  return apiClientInstance;
};
