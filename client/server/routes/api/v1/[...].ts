import { proxyRequest, getRequestHeader, createError } from 'h3';

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  if (!config.backendApiUrl) {
    throw createError({ statusCode: 500, statusMessage: 'NUXT_BACKEND_API_URL is not configured' });
  }
  const reqUrl = event.node.req.url ?? '/';
  const suffix = reqUrl.replace(/^\/api\/v1/, '');
  const target = `${config.backendApiUrl}${suffix}`;

  const traqUser = getRequestHeader(event, 'x-forwarded-user');
  const extraHeaders: Record<string, string> = {};
  if (traqUser) {
    extraHeaders['x-forwarded-user'] = traqUser;
  }

  const proxySecret = config.proxySecret as string | undefined;
  if (proxySecret) {
    extraHeaders['x-proxy-secret'] = proxySecret;
  }

  return proxyRequest(event, target, { headers: extraHeaders });
});
