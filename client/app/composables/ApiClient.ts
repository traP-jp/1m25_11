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
        // 401エラーの場合は未ログイン状態として扱う
        // 強制リダイレクトはせず、各ページで適切にハンドリング
        if (response.status === 401) {
          console.warn('認証が必要です。未ログイン状態として処理します。');
          // 認証状態をリセット（グローバル状態管理があれば）
        }
      },
    });
  }

  return apiClientInstance;
};
