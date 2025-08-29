const stampsList = () => useState<StampSummary[]>('stamps-list', () => []);
const stampsMap = () => useState<Map<string, StampSummary>>('stamps-map', () => new Map());

/**
 * 全スタンプ情報を管理するcomposable
 */
export const useStamps = () => {
  const list = readonly(stampsList());
  const map = readonly(stampsMap());

  // スタンプIDからスタンプ情報を取得
  const getStampById = (stampId: string) => map.value.get(stampId);

  return {
    // コンポーネントからはこれらを参照する
    stampsList: list,
    stampsMap: map,
    getStampById,
  };
};
