const stampsList = () => useState<Schemas['StampSummary'][]>('stamps-list', () => []);
const stampsMap = () => useState<Map<string, Schemas['StampSummary']>>('stamps-map', () => new Map());

export const useStamps = () => {
  const list = readonly(stampsList());
  const map = readonly(stampsMap());

  const getStampById = (stampId: string | undefined): Schemas['StampSummary'] | undefined => stampId ? map.value.get(stampId) : map.value.get('bc9a3814-f185-4b3d-ac1f-3c8f12ad7b52');

  const getStampByName = (stampName: string): Schemas['StampSummary'] | undefined => {
    return list.value.find(stamp => stamp.stamp_name === stampName);
  };

  return {
    stampsList: list,
    stampsMap: map,
    getStampById,
    getStampByName,
  };
};
