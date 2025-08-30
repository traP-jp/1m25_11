const stampsList = () => useState<StampSummary[]>('stamps-list', () => []);
const stampsMap = () => useState<Map<string, StampSummary>>('stamps-map', () => new Map());

export const useStamps = () => {
  const list = readonly(stampsList());
  const map = readonly(stampsMap());

  const getStampById = (stampId: string) => map.value.get(stampId);

  return {
    stampsList: list,
    stampsMap: map,
    getStampById,
  };
};
