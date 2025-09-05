export default defineNuxtPlugin(async (_) => {
  // クライアントサイドでのみ実行
  if (import.meta.server) return;

  const userListState = useState<Schemas['UserProfile'][]>('user-list');
  const userMapState = useState<Map<string, Schemas['UserProfile']>>('user-map');

  // /user-list APIを叩き、ユーザーのすべてのリストを取得する
  const { data } = await useApiClient().GET('/users-list');

  if (data) {
    userListState.value = data;
    // Map形式にも変換して保持（user_id をキーとして使用）
    userMapState.value = new Map(data.map(user => [user.user_id, user]));
  }
});
