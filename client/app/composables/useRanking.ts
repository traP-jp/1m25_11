// 処理済みランキングデータの型定義（generated型を拡張）
export type ProcessedRankingItem = Schemas['RankingResult'] & {
  stamp_name: string;
  rank: number;
  count: number; // 表示用（total_count または monthly_count）
};

export const useRanking = () => {
  const apiClient = useApiClient();

  // リアクティブな状態管理
  const ranking = ref<ProcessedRankingItem[]>([]);
  const loading = ref(false);
  const error = ref<Error | null>(null);
  const initialized = ref(false);

  // ランキングデータを手動で読み込む関数
  const loadRanking = async () => {
    // 既に読み込み済みの場合はスキップ
    if (initialized.value && ranking.value.length > 0) {
      return;
    }

    loading.value = true;
    error.value = null;

    try {
      console.log('ランキングデータの読み込みを開始...');

      // 1. 両方のAPIを並行して取得
      const [rankingResponse, stampsResponse] = await Promise.all([
        apiClient.GET('/stamps/ranking'),
        apiClient.GET('/stamps'),
      ]);

      console.log('API レスポンス - ranking:', rankingResponse);
      console.log('API レスポンス - stamps:', stampsResponse);

      // 2. エラーハンドリング
      if (rankingResponse.error) {
        throw new Error(`ランキングデータの取得に失敗しました: ${rankingResponse.error}`);
      }
      if (stampsResponse.error) {
        throw new Error(`スタンプデータの取得に失敗しました: ${stampsResponse.error}`);
      }

      // 3. データの準備
      const rankings = rankingResponse.data || [];
      const stamps = stampsResponse.data || [];

      console.log(`取得データ - rankings: ${rankings.length}件, stamps: ${stamps.length}件`);

      // 4. スタンプIDをキーとしたMapを作成
      const stampMap = new Map(stamps.map(stamp => [stamp.stamp_id, stamp]));

      // 5. データを結合して処理済みランキングデータを生成
      const processedData = rankings.map(item => ({
        stamp_id: item.stamp_id,
        total_count: item.total_count,
        monthly_count: item.monthly_count,
        stamp_name: stampMap.get(item.stamp_id)?.stamp_name ?? '不明なスタンプ',
        rank: 0, // ソート時に設定
        count: item.total_count, // 初期値
      }));

      ranking.value = processedData;
      initialized.value = true;

      console.log('ランキングデータの処理完了:', processedData.length, '件');
    }
    catch (err) {
      console.error('ランキングデータの読み込みエラー:', err);
      error.value = err instanceof Error ? err : new Error('不明なエラーが発生しました');
      ranking.value = [];
    }
    finally {
      loading.value = false;
    }
  };

  // データをリセットして再読み込みする関数
  const refresh = async () => {
    initialized.value = false;
    ranking.value = [];
    await loadRanking();
  };

  return {
    ranking: readonly(ranking),
    loading: readonly(loading),
    error: readonly(error),
    initialized: readonly(initialized),
    loadRanking,
    refresh,
  };
};
