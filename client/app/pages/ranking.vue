<script setup lang="ts">
import { getPaginationRowModel } from '@tanstack/vue-table';
import type { TableColumn, TabsItem } from '@nuxt/ui';
import type { Row } from '@tanstack/table-core';

const tableBody = useTemplateRef('tableBody');
const tableReaction = useTemplateRef('tableReaction');

interface StampRankingData {
  stamp_id: string;
  stamp_name: string;
  body_count: number;
  reaction_count: number;
  rank: number;
  count?: number; // body_count / reaction_count を注入して利用
}

const rankingTestData = ref<StampRankingData[]>([
  { stamp_id: 'e1e4a295', stamp_name: 'yatta-nya-1', body_count: 20, reaction_count: 10, rank: 0 },
  { stamp_id: 'e1e4a295', stamp_name: 'yatta-nya-2', body_count: 10, reaction_count: 15, rank: 0 },
  { stamp_id: 'e1e4a295', stamp_name: 'yatta-nya-3', body_count: 30, reaction_count: 5, rank: 0 },
  { stamp_id: 'e1e4a295', stamp_name: 'yatta-nya-4', body_count: 20, reaction_count: 8, rank: 0 },
  { stamp_id: 'e1e4a295', stamp_name: 'yatta-nya-5', body_count: 40, reaction_count: 12, rank: 0 },
  { stamp_id: 'e1e4a295', stamp_name: 'yatta-nya-6', body_count: 60, reaction_count: 20, rank: 0 },
]);

const medalMap: Record<number, string> = {
  1: '👑 1',
  2: '🥈 2',
  3: '🥉 3',
};

const columns: TableColumn<StampRankingData>[] = [
  {
    accessorKey: 'rank',
    header: '順位',
    cell: ({ row }: { row: Row<StampRankingData> }) =>
      medalMap[row.original.rank] || `#${row.original.rank}`,
  },
  {
    accessorKey: 'stamp_name',
    header: 'スタンプ名',
  },
  {
    accessorKey: 'count',
    header: '使用回数',
    cell: ({ row }) => row.original.count?.toLocaleString() ?? '-',
  },
];

// body_count並び替え
const sortedBodyCount = computed<StampRankingData[]>(() => {
  return [...rankingTestData.value]
    .sort((a, b) => b.body_count - a.body_count)
    .map((item, index) => ({
      ...item,
      rank: index + 1,
      count: item.body_count,
    }));
});

// reaction_count並び替え
const sortedReactionCount = computed<StampRankingData[]>(() => {
  return [...rankingTestData.value]
    .sort((a, b) => b.reaction_count - a.reaction_count)
    .map((item, index) => ({
      ...item,
      rank: index + 1,
      count: item.reaction_count,
    }));
});

const items = ref<TabsItem[]>([
  { label: '本文内', icon: 'i-lucide-user', slot: 'body_count' },
  { label: 'リアクション', icon: 'i-lucide-lock', slot: 'reaction_count' },
]);

const paginationBody = ref({ pageIndex: 0, pageSize: 20 });
const paginationReaction = ref({ pageIndex: 0, pageSize: 20 });
</script>

<template>
  <div class="w-full max-w-5xl mx-auto px-2 sm:px-4 md:px-6">
    <UTabs
      :items="items"
      variant="pill"
      class="w-full"
    >
      <!-- 本文ランキング -->
      <template #body_count>
        <UTable
          ref="tableBody"
          v-model:pagination="paginationBody"
          :data="sortedBodyCount"
          :columns="columns"
          :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
          class="w-full"
        />

        <div class="flex justify-center border-t border-default pt-4">
          <UPagination
            :default-page="(tableBody?.tableApi?.getState().pagination.pageIndex || 0) + 1"
            :items-per-page="tableBody?.tableApi?.getState().pagination.pageSize"
            :total="tableBody?.tableApi?.getFilteredRowModel().rows.length"
            @update:page="(p) => tableBody?.tableApi?.setPageIndex(p - 1)"
          />
        </div>
      </template>

      <!-- リアクションランキング -->
      <template #reaction_count>
        <UTable
          ref="tableReaction"
          v-model:pagination="paginationReaction"
          :data="sortedReactionCount"
          :columns="columns"
          :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
          class="w-full"
        />

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
