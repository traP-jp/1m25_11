// drawerで閲覧しているスタンプのid
export const useSelectedStampId = () => useState<string | undefined>('selected-stamp-id', () => undefined);
