import createClient from 'openapi-fetch';
import type { paths } from '~~/shared/types/generated';

// runtimeConfig から baseUrl を取得
const { public: publicConfig } = useRuntimeConfig();
const base = publicConfig.apiBase.replace(/\/$/, '');

export const apiClient = createClient<paths>({
  baseUrl: base,
  // openapi-fetch の fetch オーバーライドは Request を受け取るシグネチャ
  fetch: (req: Request) => {
    return fetch(req, { credentials: 'include' });
  },
});
