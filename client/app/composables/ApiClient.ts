import createClient from 'openapi-fetch';
import type { paths } from '~~/shared/types/generated';

export const apiClient = createClient<paths>({
  baseUrl: 'http://localhost:3000/server/api/v1',
  // Cookie送信を有効にするためのfetchオプション
  fetch: (input: RequestInfo | URL, init?: RequestInit) => {
    return fetch(input, {
      ...init,
      credentials: 'include',
    });
  },
});
