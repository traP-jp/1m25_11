export default defineNuxtPlugin(async () => {
  // Map形式のstateを削除
  const listState = useState<Schemas['StampSummary'][]>('stamps-list', () => []);

  // アプリの初期化時にすでにデータがあれば何もしない（クライアント側での重複実行防止）
  if (listState.value.length > 0) {
    return;
  }

  const apiClient = useApiClient();
  const { data } = await apiClient.GET('/stamps');

  if (data) {
    listState.value = data;
  }
});
