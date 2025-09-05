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
        // 401エラーの場合はログインページにリダイレクト
        if (response.status === 401) {
          console.warn('認証が必要です。ログインページにリダイレクトします。');
          // サーバーのログインエンドポイントにリダイレクト
          window.location.href = `${config.public.apiBase}/login`;
        }
      },
    });
  }

  return apiClientInstance;
};
