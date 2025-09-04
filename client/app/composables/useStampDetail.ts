export const useStampDetail = (stampId: string | undefined) => {
  const stampDetail = ref<Schemas['Stamp'] | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchStampDetail = async () => {
    if (!stampId) {
      stampDetail.value = null;
      return;
    }

    loading.value = true;
    error.value = null;

    try {
      const { data, error: apiError } = await apiClient.GET('/stamps/{stampId}', {
        params: {
          path: {
            stampId,
          },
        },
      });

      if (apiError) {
        console.error('Failed to fetch stamp detail:', apiError);
        error.value = 'スタンプの詳細を取得できませんでした';
        stampDetail.value = null;
      }
      else {
        stampDetail.value = data || null;
      }
    }
    catch (err) {
      console.error('Error fetching stamp detail:', err);
      error.value = 'スタンプの詳細を取得できませんでした';
      stampDetail.value = null;
    }
    finally {
      loading.value = false;
    }
  };

  // stampIdが変更されたときに自動で再取得
  watch(
    () => stampId,
    () => {
      fetchStampDetail();
    },
    { immediate: true },
  );

  return {
    stampDetail: readonly(stampDetail),
    loading: readonly(loading),
    error: readonly(error),
    refetch: fetchStampDetail,
  };
};
