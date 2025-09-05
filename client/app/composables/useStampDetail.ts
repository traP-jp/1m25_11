import { type MaybeRef, unref, computed } from 'vue';

export const useStampDetail = (stampId: MaybeRef<string | undefined>) => {
  const apiClient = useApiClient();
  const { getStampById } = useStamps();

  const {
    data: fetchedStamp,
    pending: loading,
    error,
    refresh: refetch,
  } = useAsyncData<Schemas['Stamp'] | null>(

    // リアクティブなユニークキーを生成
    () => `stamp-detail-${unref(stampId) ?? 'null'}`,

    // データを取得
    async () => {
      const id = unref(stampId);
      if (!id) return null;

      const { data, error } = await apiClient.GET('/stamps/{stampId}', {
        params: { path: { stampId: id } },
      });

      if (error) {
        throw createError({ statusCode: 404, message: 'スタンプが見つかりません' });
      }

      // APIの戻り値をnullに統一
      return data ?? null;
    },

    {
      // stampIdの変更を監視するため、ゲッター関数を渡す
      watch: [() => unref(stampId)],
    },
  );

  // コンポーネントに公開する表示用のデータをcomputedで生成
  const stampDetail = computed(() => {
    // APIから取得した完全なデータがあれば、それを返す
    if (fetchedStamp.value) {
      return fetchedStamp.value;
    }

    // なければ、キャッシュ済みの概要データを返す
    return getStampById(unref(stampId));
  });

  // 完全な詳細データがロード済みかを示すフラグ
  const isDetailLoaded = computed(() => !!fetchedStamp.value);

  // コンポーネントが必要とするすべての状態と関数を返す
  return {
    stampDetail,
    isDetailLoaded,
    loading,
    error,
    refetch,
  };
};
