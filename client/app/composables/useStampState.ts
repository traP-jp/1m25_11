// drawerで閲覧しているスタンプのid
export const useSelectedStampId = () => useState<string | null>('selected-stamp-id', () => null);
