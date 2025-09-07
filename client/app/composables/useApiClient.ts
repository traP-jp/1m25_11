import createClient from 'openapi-fetch';
import type { paths } from '~~/shared/types/generated';

type ApiClient = ReturnType<typeof createClient<paths>>;

let apiClientInstance: ApiClient | null = null;

export const useApiClient = () => {
  if (!apiClientInstance) {
    const config = useRuntimeConfig();
    apiClientInstance = createClient<paths>({
      baseUrl: config.public.apiBase,
      credentials: 'include', // クロスオリジンに対してもCookieを送信する
    });
  };

  return apiClientInstance;
};
