import { apiClient } from '#imports';

export default defineNuxtPlugin(async (_) => {
  const listState = useState<Schemas['StampSummary'][]>('stamps-list');
  const mapState = useState<Map<string, Schemas['StampSummary']>>('stamps-map');

  // /stamps APIを叩き、スタンプのすべてのリストを取得する
  const { data } = await apiClient.GET('/stamps');

  if (data) {
    listState.value = data;
    // Map形式にも変換して保持
    mapState.value = new Map(data.map(stamp => [stamp.stamp_id, stamp]));
  }
});
