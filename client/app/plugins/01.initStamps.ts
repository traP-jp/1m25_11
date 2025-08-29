import { apiClient } from '#imports';

export default defineNuxtPlugin(async (_) => {
  const listState = useState<StampSummary[]>('stamps-list');
  const mapState = useState<Map<string, StampSummary>>('stamps-map');

  // if (listState.value.length > 0) {
  //   return;
  // }

  // /stamps APIを叩き、スタンプのすべてのリストを取得する
  const { data } = await apiClient.GET('/stamps');

  if (data) {
    listState.value = data;
    // Map形式にも変換して保持
    mapState.value = new Map(data.map(stamp => [stamp.stamp_id, stamp]));
    console.log(`Loaded ${data.length} stamps.`);
  }
});
