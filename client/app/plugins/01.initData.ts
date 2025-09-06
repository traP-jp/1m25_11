import type { Schemas } from '#imports';

export default defineNuxtPlugin(async (nuxtApp) => {
  const userState = useState<Schemas['UserStatus'] | null>('user', () => null);

  // サーバーサイドの場合、または、クライアントサイドでのレンダリングの初回
  if (!nuxtApp.payload.serverRendered) {
    const apiClient = useApiClient();
    try {
      // ログインしているか確認
      const { data: user, error } = await apiClient.GET('/me');

      if (error || user == undefined) {
        // ログインしていないとき
        userState.value = null;
        console.log('You are not logged in. ', JSON.stringify(error));
        return;
      }

      userState.value = user;

      // サーバーからフェッチしたデータを保存しておく
      const stampsListState = useState<Schemas['StampSummary'][]>('stamps-list');
      const usersListState = useState<Schemas['UserProfile'][]>('users-list');

      // データをフェッチし変数に保存
      const [{ data: stamps }, { data: users }] = await Promise.all([
        apiClient.GET('/stamps'),
        apiClient.GET('/users-list'),
      ]);

      if (stamps) {
        stampsListState.value = stamps;
      }
      if (users) {
        usersListState.value = users;
      }
    }
    catch (e) {
      console.error('Initialization failed:', e);
      userState.value = null;
    }
  }
});
