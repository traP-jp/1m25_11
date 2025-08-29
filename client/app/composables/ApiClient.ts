import createClient from 'openapi-fetch/dist/index.cjs';
import type { paths } from '~~/shared/types/generated';

export const apiClient = createClient<paths>({
  baseUrl: 'https://1m25-11.trap.show/api/v1/',
});
