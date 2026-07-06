export default defineNuxtRouteMiddleware(async () => {
  const user = useState<Schemas['UserStatus'] | null>('user');

  const { data: userData } = await useApiClient().GET('/me');

  user.value = userData ?? null;
});
