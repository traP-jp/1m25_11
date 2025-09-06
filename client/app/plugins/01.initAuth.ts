export default defineNuxtPlugin(async () => {
  // サーバーサイドでは認証状態を変更しない（Hydration エラー回避）
  if (import.meta.server) {
    console.log('サーバーサイド: 認証初期化をスキップ');
    return;
  }

  console.log('クライアントサイド: アプリケーション初期化開始');

  try {
    // 認証状態の確認
    const { checkAuthStatus, isLoggedIn } = useAuth();
    await checkAuthStatus();

    // 認証済みの場合のみデータを初期化
    if (isLoggedIn.value) {
      console.log('認証済みユーザー: データ初期化開始');
      await Promise.all([
        initializeStamps(),
        initializeUsers(),
      ]);
      console.log('データ初期化完了');
    }
    else {
      console.log('未認証ユーザー: データ初期化をスキップ');
    }
  }
  catch (error) {
    console.error('アプリケーション初期化エラー:', error);
  }

  console.log('アプリケーション初期化完了');
});

// スタンプデータ初期化
async function initializeStamps() {
  const listState = useState<Schemas['StampSummary'][]>('stamps-list', () => []);

  // アプリの初期化時にすでにデータがあれば何もしない
  if (listState.value.length > 0) {
    return;
  }

  console.log('スタンプデータ初期化開始');
  const apiClient = useApiClient();
  const { data } = await apiClient.GET('/stamps');

  if (data) {
    listState.value = data;
    console.log(`スタンプデータ初期化完了: ${data.length}件`);
  }
}

// ユーザーデータ初期化
async function initializeUsers() {
  const userListState = useState<Schemas['UserProfile'][]>('user-list');
  const userMapState = useState<Map<string, Schemas['UserProfile']>>('user-map');

  // すでにデータがあれば何もしない
  if (userListState.value && userListState.value.length > 0) {
    return;
  }

  console.log('ユーザーデータ初期化開始');
  const { data } = await useApiClient().GET('/users-list');

  if (data) {
    userListState.value = data;
    // Map形式にも変換して保持
    userMapState.value = new Map(data.map(user => [user.user_id, user]));
    console.log(`ユーザーデータ初期化完了: ${data.length}件`);
  }
}
