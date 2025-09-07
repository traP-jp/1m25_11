import type { paths, components } from '~~/shared/types/generated';

// APIの検索パラメータの型を直接活用
type SearchParams = paths['/stamps/search']['get']['parameters']['query'];

// APIのレスポンス型を直接活用
type SearchResult = components['schemas']['SearchResult'];

export const useStampSearch = (searchParams: MaybeRef<SearchParams>) => {
  const apiClient = useApiClient();
  const paramsRef = toRef(searchParams);

  const { data, pending, error, refresh } = useAsyncData<SearchResult | null>(
    // リアクティブなキーを生成
    `stamp-search-${JSON.stringify(toValue(paramsRef))}`,

    async () => {
      const currentParams = toValue(paramsRef);

      // パラメータが空または全てのキーが未定義の場合はAPIを呼び出さない
      if (!currentParams) {
        return { stamps: [] }; // 空の結果を返す
      }

      // openapi-fetch の型推論を活用
      const { data, error } = await apiClient.GET('/stamps/search', {
        params: { query: currentParams },
      });

      console.log(data);
      if (error) {
        throw new Error('Failed to search stamps');
      }

      return data ?? { stamps: [] };
    },
    {
      watch: [paramsRef],
    },
  );

  return {
    searchResult: data,
    stamps: computed(() => data.value?.stamps ?? []),
    loading: pending,
    error,
    refetch: refresh,
  };
};
