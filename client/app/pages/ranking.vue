<template>
  <UContainer>
    <div class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6">
      <!-- 初期状態：ランキング読み込みボタン -->
      <div
        v-if="!initialized && !loading"
        class="text-center py-12"
      >
        <div class="space-y-4">
          <h2 class="text-2xl font-bold text-gray-800">
            スタンプ使用回数ランキング
          </h2>
          <p class="text-gray-600">
            スタンプの人気ランキングを表示します
          </p>
          <UButton
            size="lg"
            class="px-8 py-3"
            @click="loadRanking"
          >
            <UIcon
              name="material-symbols:leaderboard"
              class="mr-2"
            />
            ランキングを表示
          </UButton>
        </div>
      </div>

      <!-- ローディング状態 -->
      <div
        v-else-if="loading"
        class="flex flex-col justify-center items-center py-12"
      >
        <div class="text-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
          <div class="space-y-2">
            <p class="text-lg font-medium">
              ランキングデータを読み込み中...
            </p>
            <p class="text-sm text-gray-600">
              スタンプデータと使用回数を取得しています
            </p>
          </div>
        </div>
      </div>

      <!-- エラー状態 -->
      <div
        v-else-if="error"
        class="text-center py-12"
      >
        <div class="space-y-4">
          <div class="text-red-500 space-y-2">
            <UIcon
              name="material-symbols:error"
              class="text-6xl"
            />
            <h3 class="text-xl font-bold">
              エラーが発生しました
            </h3>
            <p class="text-sm">
              {{ error.message }}
            </p>
          </div>
          <div class="space-x-2">
            <UButton
              color="primary"
              variant="soft"
              @click="refresh"
            >
              <UIcon
                name="material-symbols:refresh"
                class="mr-2"
              />
              再試行
            </UButton>
            <UButton
              variant="ghost"
              @click="resetToInitial"
            >
              最初に戻る
            </UButton>
          </div>
        </div>
      </div>

      <!-- データ表示：ランキングテーブル -->
      <div
        v-else-if="ranking && ranking.length > 0"
        class="space-y-4"
      >
        <!-- ヘッダー部分 -->
        <div class="flex justify-between items-center">
          <h2 class="text-2xl font-bold">
            スタンプランキング
          </h2>
          <UButton
            variant="ghost"
            size="sm"
            :loading="loading"
            @click="refresh"
          >
            <UIcon
              name="material-symbols:refresh"
              class="mr-1"
            />
            更新
          </UButton>
        </div>

        <!-- タブ表示 -->
        <UTabs
          :items="items"
          variant="pill"
          class="w-full"
        >
          <!-- 総合ランキング -->
          <template #total_count>
            <UTable
              ref="tableTotal"
              v-model:pagination="paginationTotal"
              :data="totalRanking"
              :columns="columns"
              :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
              class="w-full"
            >
              <template #stamp_name-cell="{ row }">
                <div class="flex items-center gap-3">
                  <NuxtImg
                    :src="getFileUrl(row.original.stamp_id)"
                    class="m-auto w-12 h-12"
                  />
                  <p>{{ row.original.stamp_name }}</p>
                </div>
              </template>
            </UTable>
            <div class="flex justify-center border-t border-default pt-4">
              <UPagination
                v-model:page="pageTotal"
                :items-per-page="tableTotal?.tableApi?.getState().pagination.pageSize"
                :total="tableTotal?.tableApi?.getFilteredRowModel().rows.length"
              />
            </div>
          </template>

          <!-- 1か月ランキング -->
          <template #monthly_count>
            <UTable
              ref="tableMonthly"
              v-model:pagination="paginationMonthly"
              :data="monthlyRanking"
              :columns="columns"
              :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
              class="w-full"
            >
              <template #stamp_name-cell="{ row }">
                <div class="flex items-center gap-3">
                  <NuxtImg
                    :src="getFileUrl(row.original.stamp_id)"
                    class="mx-auto w-12 h-12"
                  />
                  <p>{{ row.original.stamp_name }}</p>
                </div>
              </template>
            </UTable>

            <div class="flex justify-center border-t border-default pt-4">
              <UPagination
                v-model:page="pageMonthly"
                :items-per-page="tableMonthly?.tableApi?.getState().pagination.pageSize"
                :total="tableMonthly?.tableApi?.getFilteredRowModel().rows.length"
              />
            </div>
          </template>
        </UTabs>
      </div>

      <!-- データが空の場合 -->
      <div
        v-else-if="initialized && (!ranking || ranking.length === 0)"
        class="text-center py-12"
      >
        <div class="space-y-4">
          <UIcon
            name="material-symbols:inbox"
            class="text-6xl text-gray-400"
          />
          <div class="space-y-2">
            <h3 class="text-xl font-medium text-gray-600">
              ランキングデータがありません
            </h3>
            <p class="text-sm text-gray-500">
              スタンプの使用データが見つかりませんでした
            </p>
          </div>
          <UButton
            variant="soft"
            @click="refresh"
          >
            <UIcon
              name="material-symbols:refresh"
              class="mr-2"
            />
            再読み込み
          </UButton>
        </div>
      </div>
    </div>
  </UContainer>
</template>

<script setup lang="ts">
import { getPaginationRowModel } from '@tanstack/vue-table';
import type { TableColumn, TabsItem } from '@nuxt/ui';
import type { Row } from '@tanstack/table-core';
import { getFileUrl } from '#imports';
import type { ProcessedRankingItem } from '~/composables/useRanking';

// useRanking composable を使用（ボタン方式）
const { ranking, loading, error, initialized, loadRanking, refresh } = useRanking();

// 初期状態にリセットする関数
const resetToInitial = () => {
  // composable内部の状態をリセット（実装上は新しいインスタンスの作成が必要）
  // ここでは簡易的にページリロード
  window.location.reload();
};

// 総合ランキング（total_count順でソート）
const totalRanking = computed(() => {
  if (!ranking.value || ranking.value.length === 0) return [];

  return [...ranking.value]
    .map(item => ({ ...item, count: item.total_count }))
    .sort((a, b) => b.total_count - a.total_count)
    .map((item, index) => ({ ...item, rank: index + 1 }));
});

// 月間ランキング（monthly_count順でソート）
const monthlyRanking = computed(() => {
  if (!ranking.value || ranking.value.length === 0) return [];

  return [...ranking.value]
    .map(item => ({ ...item, count: item.monthly_count }))
    .sort((a, b) => b.monthly_count - a.monthly_count)
    .map((item, index) => ({ ...item, rank: index + 1 }));
});

// テーブルのカラム定義
const columns: TableColumn<ProcessedRankingItem>[] = [
  {
    accessorKey: 'rank',
    header: '順位',
    cell: ({ row }: { row: Row<ProcessedRankingItem> }) =>
      medalMap[row.original.rank] || `#${row.original.rank}`,
    meta: {
      class: {
        td: 'font-bold',
      },
    },
  },
  {
    accessorKey: 'stamp_name',
    header: 'スタンプ名',
    meta: {
      class: {
        th: 'text-center',
        td: 'flex justify-start font-bold',
      },
    },
  },
  {
    accessorKey: 'count',
    header: '使用回数',
    cell: ({ row }) =>
      row.original.count !== undefined
        ? `${row.original.count.toLocaleString()} 回`
        : '-',
    meta: {
      class: {
        td: 'font-bold',
      },
    },
  },
];

// タブの設定
const items = ref<TabsItem[]>([
  { label: '総合', slot: 'total_count' },
  { label: '1か月以内', slot: 'monthly_count' },
]);

// ページネーション設定
const paginationTotal = ref({ pageIndex: 0, pageSize: 20 });
const paginationMonthly = ref({ pageIndex: 0, pageSize: 20 });

// テーブル参照
const tableTotal = useTemplateRef('tableTotal');
const tableMonthly = useTemplateRef('tableMonthly');

// ページ番号
const pageTotal = ref(1);
const pageMonthly = ref(1);

// ページネーションの監視
watch(pageTotal, (p) => {
  tableTotal.value?.tableApi.setPageIndex(p - 1);
});

watch(pageMonthly, (p) => {
  tableMonthly.value?.tableApi.setPageIndex(p - 1);
});

// メダル表示のマップ
const medalMap: Record<number, string> = {
  1: '👑 1',
  2: '🥈 2',
  3: '🥉 3',
};
</script>
