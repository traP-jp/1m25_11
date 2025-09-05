import { readonly, computed } from 'vue';

// 状態のソースはスタンプの「リスト（配列）」のみにする
const useStampsListState = () => useState<Schemas['StampSummary'][]>('stamps-list', () => []);

// このComposableがスタンプデータへの唯一のアクセスポイントとなる
export const useStamps = () => {
  // 唯一の状態（配列）を取得
  const list = useStampsListState();

  // 配列のStateからMapを「算出プロパティ（computed）」として生成する
  // listが変更されると、このmapも自動的に再計算される
  const map = computed(() =>
    new Map(list.value.map(stamp => [stamp.stamp_id, stamp])),
  );

  /**
   * IDからスタンプを取得する関数
   * @param stampId スタンプID
   * @returns StampSummary | undefined
   */
  const getStampById = (stampId?: string): Schemas['StampSummary'] | undefined => {
    // デフォルトIDを設定（オプショナル）
    const targetId = stampId ?? 'bc9a3814-f185-4b3d-ac1f-3c8f12ad7b52';
    return map.value.get(targetId);
  };

  /**
   * 名前からスタンプを取得する関数
   * @param stampName スタンプ名
   * @returns StampSummary | undefined
   */
  const getStampByName = (stampName: string): Schemas['StampSummary'] | undefined => {
    // こちらは低頻度な操作と仮定し、都度findする
    return list.value.find(stamp => stamp.stamp_name === stampName);
  };

  return {
    // 外部からは読み取り専用として公開する
    stampsList: readonly(list),
    stampsMap: readonly(map),
    getStampById,
    getStampByName,
  };
};
