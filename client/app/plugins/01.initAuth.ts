export default defineNuxtPlugin({
  name: 'auth-initialization',
  async setup() {
    // サーバーサイドでは何もしない（Hydration エラー回避）
    if (import.meta.server) {
      console.log('🔄 サーバーサイド: 認証初期化をスキップ');
      return;
    }

    console.log('🚀 クライアントサイド: アプリケーション初期化開始');

    try {
      // 1. 認証状態の確認
      const { checkAuthStatus, isLoggedIn } = useAuth();
      console.log('🔐 認証状態の確認を開始');
      await checkAuthStatus();

      // 2. 認証結果に応じた処理
      if (isLoggedIn.value) {
        console.log('✅ 認証成功 - データ初期化を開始');

        // 3. 並列でデータを初期化
        const initResults = await Promise.allSettled([
          initializeStamps(),
          initializeUsers(),
        ]);

        // 初期化結果のログ出力
        initResults.forEach((result, index) => {
          const dataType = index === 0 ? 'スタンプ' : 'ユーザー';
          if (result.status === 'fulfilled') {
            console.log(`✅ ${dataType}データ初期化成功`);
          }
          else {
            console.warn(`⚠️ ${dataType}データ初期化失敗:`, result.reason);
          }
        });

        console.log('📦 アプリケーションデータ初期化完了');
      }
      else {
        console.log('❌ 未認証 - ログインページへリダイレクト予定');
      }
    }
    catch (error) {
      console.error('🚨 アプリケーション初期化エラー:', error);
    }

    console.log('🎉 アプリケーション初期化完了');
  },
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
