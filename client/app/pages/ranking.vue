<template>
  <UContainer>
    <div class="max-w-5xl mx-auto px-2 sm:px-4 md:px-6">
      <UTabs
        :items="items"
        variant="pill"
        class="w-full"
      >
        <!-- 総合ランキング -->
        <template #count_total>
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
        <template #count_monthly>
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
  month_count: number;
  rank: number; // すべてのスタンプに初期値として0を付与
  count?: number; // count_total / count_monthly を注入して利用
}

const rankingTestData = ref<StampRankingData[]>([
  { stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f', stamp_name: 'kyapi-nya', month_count: 30, total_count: 25, rank: 0 },
  { stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f', stamp_name: 'kyapi-nya', month_count: 40, total_count: 30, rank: 0 },
  { stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc', stamp_name: 'korosu-nya', month_count: 25, total_count: 30, rank: 0 },
  { stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc', stamp_name: 'korosu-nya', month_count: 10, total_count: 40, rank: 0 },
  { stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068', stamp_name: '403_forbidden', month_count: 15, total_count: 30, rank: 0 },
  { stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068', stamp_name: '403_forbidden', month_count: 25, total_count: 10, rank: 0 },
]);

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
  { label: '1か月以内', slot: 'month_count' },
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

function useSortedData(key: 'total_count' | 'month_count') {
  return computed(() =>
    [...rankingTestData.value]
      .sort((a, b) => b[key] - a[key])
      .map((item, index) => ({
        ...item,
        rank: index + 1,
        count: item[key],
      })),
  );
}

const sortedCountTotal = useSortedData('total_count');
const sortedCountMonthly = useSortedData('month_count');
</script>
