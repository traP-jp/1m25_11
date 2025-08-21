// drawerで閲覧しているスタンプのid
export const useSelectedStampId = () => useState<string | null>('selected-stamp-id', () => null);

// drawerが開いているかどうか
export const useIsStampDrawerOpen = () => useState<boolean>('is-stamp-drawer-open', () => false);
