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
    });
  }

  return apiClientInstance;
};
