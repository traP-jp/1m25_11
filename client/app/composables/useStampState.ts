// drawerで閲覧しているスタンプのid
export const useSelectedStampId = () => useState<string | null>('selected-stamp-id', () => null);

// drawerが開いているかどうか
export const useIsStampDrawerOpen = () => {
  const selectedStampId = useSelectedStampId();
  return computed(() => selectedStampId.value !== null);
};
