import { proxyRequest, getRequestHeader } from 'h3';

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  const reqUrl = event.node.req.url ?? '/';
  const suffix = reqUrl.replace(/^\/api\/v1/, '');
  const target = `${config.backendApiUrl}${suffix}`;

  const traqUser = getRequestHeader(event, 'x-forwarded-user');
  const extraHeaders: Record<string, string> = {};
  if (traqUser) {
    extraHeaders['x-forwarded-user'] = traqUser;
  }

  return proxyRequest(event, target, { headers: extraHeaders });
});
