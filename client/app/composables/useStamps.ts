import { readonly, computed } from 'vue';

export const useStamps = () => {
  const list = useState<Schemas['StampSummary'][]>('stamps-list', () => []);

  const map = computed(() =>
    new Map(list.value.map(stamp => [stamp.stamp_id, stamp])),
  );

  /**
   * IDからスタンプを取得する関数
   * @param stampId スタンプID
   * @returns StampSummary | undefined
   */
  const getStampById = (stampId?: string): Schemas['StampSummary'] | undefined => {
    const targetId = stampId ?? 'bc9a3814-f185-4b3d-ac1f-3c8f12ad7b52';
    return map.value.get(targetId);
  };

  /**
   * 名前からスタンプを取得する関数
   * @param stampName スタンプ名
   * @returns StampSummary | undefined
   */
  const getStampByName = (stampName: string): Schemas['StampSummary'] | undefined => {
    return list.value.find(stamp => stamp.stamp_name === stampName);
  };

  return {
    stampsList: readonly(list),
    stampsMap: readonly(map),
    getStampById,
    getStampByName,
  };
};
