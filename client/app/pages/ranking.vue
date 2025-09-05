<template>
  <UContainer>
    <div class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6">
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
            :data="sortedCountTotal"
            :columns="columns"
            :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
            class="w-full"
          >
            <template #stamp_name-cell="{ row }">
              <div class="flex items-center gap-3">
                <NuxtImg
                  :src="`https://q.trap.jp/api/1.0/public/emoji/${row.original.stamp_id}`"
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
            :data="sortedCountMonthly"
            :columns="columns"
            :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
            class="w-full"
          >
            <template #stamp_name-cell="{ row }">
              <div class="flex items-center gap-3">
                <NuxtImg
                  :src="`https://q.trap.jp/api/1.0/public/emoji/${row.original.stamp_id}`"
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
  </UContainer>
</template>

<script setup lang="ts">
import { getPaginationRowModel } from '@tanstack/vue-table';
import type { TableColumn, TabsItem } from '@nuxt/ui';
import type { Row } from '@tanstack/table-core';

interface StampRankingData {
  stamp_id: string;
  stamp_name: string;
  total_count: number;
  monthly_count: number;
  rank: number;
  count?: number;
}

const { getStampById } = useStamps();

const rankingData = ref<StampRankingData[]>([]);

onMounted(async () => {
  const { data } = await apiClient.GET('/stamps/ranking');

  // stamp_name を stamp_id から取得、rank は 0、count を初期化
  rankingData.value = (data ?? []).map(item => ({
    stamp_id: item.stamp_id,
    stamp_name: getStampById(item.stamp_id)?.stamp_name ?? '不明なスタンプ',
    total_count: item.total_count,
    monthly_count: item.monthly_count,
    rank: 0,
    count: item.total_count,
  }));
});

const columns: TableColumn<StampRankingData>[] = [
  {
    accessorKey: 'rank',
    header: '順位',
    cell: ({ row }: { row: Row<StampRankingData> }) =>
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

const items = ref<TabsItem[]>([
  { label: '総合', slot: 'total_count' },
  { label: '1か月以内', slot: 'monthly_count' },
]);

const paginationTotal = ref({ pageIndex: 0, pageSize: 20 });
const paginationMonthly = ref({ pageIndex: 0, pageSize: 20 });

const tableTotal = useTemplateRef('tableTotal');
const tableMonthly = useTemplateRef('tableMonthly');

const pageTotal = ref(1); // 総合ランキングのページ番号
const pageMonthly = ref(1); // 1か月ランキングのページ番号

watch(pageTotal, (p) => {
  tableTotal.value?.tableApi.setPageIndex(p - 1);
});

watch(pageMonthly, (p) => {
  tableMonthly.value?.tableApi.setPageIndex(p - 1);
});

const medalMap: Record<number, string> = {
  1: '👑 1',
  2: '🥈 2',
  3: '🥉 3',
};

function useSortedData(key: 'total_count' | 'monthly_count') {
  return computed(() =>
    [...rankingData.value]
      .map(item => ({ ...item, count: item[key] }))
      .sort((a, b) => b[key] - a[key])
      .map((item, index) => ({ ...item, rank: index + 1 })),
  );
}

const sortedCountTotal = useSortedData('total_count');
const sortedCountMonthly = useSortedData('monthly_count');
</script>
