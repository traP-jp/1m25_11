export const useStampDetail = (stampId: MaybeRef<string | undefined>) => {
  const apiClient = useApiClient();
  const idRef = toRef(stampId);

  const { data, pending, error, refresh } = useAsyncData<Schemas['Stamp'] | null>(

    // リアクティブなユニークキーを生成
    `stamp-detail-${idRef.value}`,

    // データを取得
    async () => {
      const currentStampId = toValue(idRef);
      if (!currentStampId) {
        return null;
      }

      const { data, error } = await apiClient.GET('/stamps/{stampId}', {
        params: {
          path: { stampId: currentStampId },
        },
      });

      if (error) {
        throw new Error('Failed to fetch stamp detail');
      }

      return data ?? null;
    },
    {
      // stampIdの変更を監視するため、ゲッター関数を渡す
      watch: [idRef],
    },
  );

  return {
    stampDetail: data,
    loading: pending,
    error,
    refetch: refresh,
  };
};
