<template>
  <div class="w-full max-w-5xl mx-auto px-2 sm:px-4 md:px-6">
    <UTabs
      :items="items"
      variant="pill"
      class="w-full"
    >
      <!-- 総合ランキング -->
      <template #count_total>
        <UTable
          ref="tableBody"
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
            :default-page="(tableBody?.tableApi?.getState().pagination.pageIndex || 0) + 1"
            :items-per-page="tableBody?.tableApi?.getState().pagination.pageSize"
            :total="tableBody?.tableApi?.getFilteredRowModel().rows.length"
            @update:page="(p) => tableBody?.tableApi?.setPageIndex(p - 1)"
          />
        </div>
      </template>

      <!-- 1か月ランキング -->
      <template #count_monthly>
        <UTable
          ref="tableReaction"
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
            :default-page="(tableReaction?.tableApi?.getState().pagination.pageIndex || 0) + 1"
            :items-per-page="tableReaction?.tableApi?.getState().pagination.pageSize"
            :total="tableReaction?.tableApi?.getFilteredRowModel().rows.length"
            @update:page="(p) => tableReaction?.tableApi?.setPageIndex(p - 1)"
          />
        </div>
      </template>
    </UTabs>
  </div>
</template>

<script setup lang="ts">
import { getPaginationRowModel } from '@tanstack/vue-table';
import type { TableColumn, TabsItem } from '@nuxt/ui';
import type { Row } from '@tanstack/table-core';

const tableBody = useTemplateRef('tableBody');
const tableReaction = useTemplateRef('tableReaction');

interface StampRankingData {
  stamp_id: string;
  stamp_name: string;
  count_total: number;
  count_monthly: number;
  rank: number; // すべてのスタンプに初期値として0を付与
  count?: number; // count_total / count_monthly を注入して利用
}

const rankingTestData = ref<StampRankingData[]>([
  { stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f', stamp_name: 'kyapi-nya', count_monthly: 30, count_total: 25, rank: 0 },
  { stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f', stamp_name: 'kyapi-nya', count_monthly: 40, count_total: 30, rank: 0 },
  { stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc', stamp_name: 'korosu-nya', count_monthly: 25, count_total: 30, rank: 0 },
  { stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc', stamp_name: 'korosu-nya', count_monthly: 10, count_total: 40, rank: 0 },
  { stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068', stamp_name: '403_forbidden', count_monthly: 15, count_total: 30, rank: 0 },
  { stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068', stamp_name: '403_forbidden', count_monthly: 25, count_total: 10, rank: 0 },
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
  { label: '総合', slot: 'count_total' },
  { label: '1か月以内', slot: 'count_monthly' },
]);

const paginationTotal = ref({ pageIndex: 0, pageSize: 20 });
const paginationMonthly = ref({ pageIndex: 0, pageSize: 20 });

const medalMap: Record<number, string> = {
  1: '👑 1',
  2: '🥈 2',
  3: '🥉 3',
};

function useSortedData(key: 'count_total' | 'count_monthly') {
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

const sortedCountTotal = useSortedData('count_total');
const sortedCountMonthly = useSortedData('count_monthly');
</script>
