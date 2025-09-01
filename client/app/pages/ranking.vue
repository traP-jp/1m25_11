<script setup lang="ts">
import { getPaginationRowModel } from '@tanstack/vue-table';
import type { TableColumn, TabsItem } from '@nuxt/ui';

const table = useTemplateRef('table');

interface User {
  stamp_id: string;
  stamp_name: string;
  body_count: string;
  reaction_count: string;
}

const data = ref<User[]>([
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'yatta-nya-1',
    body_count: '20',
    reaction_count: '10',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'yatta-nya-2',
    body_count: '10',
    reaction_count: '10',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'yatta-nya-3',
    body_count: '30',
    reaction_count: '10',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'yatta-nya-4',
    body_count: '20',
    reaction_count: '10',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'yatta-nya-5',
    body_count: '40',
    reaction_count: '10',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'yatta-nya-6',
    body_count: '60',
    reaction_count: '10',
  },
]);

const columns: TableColumn<User>[] = [
  {
    accessorKey: 'stamp_name',
    header: 'スタンプ名',
  },
  {
    accessorKey: 'body_count',
    header: '使用回数',
  },
  {
    accessorKey: 'reaction_count',
    header: '使用回数',
  },
];

const pagination = ref({
  pageIndex: 0,
  pageSize: 5,
});

// body_count (count_monthly) の降順でソート
const sortedData = computed(() => {
  return [...data.value].sort((a, b) => {
    return Number(b.body_count) - Number(a.body_count);
  });
});

const items = ref<TabsItem[]>([
  {
    label: '本文内',
    icon: 'i-lucide-user',
  },
  {
    label: 'リアクション',
    icon: 'i-lucide-lock',
  },
]);
</script>

<template>
  <div class="w-full max-w-5xl mx-auto px-2 sm:px-4 md:px-6">
    <UTabs
      :items="items"
      variant="pill"
      class="w-full"
    />
    <UTable
      ref="table"
      v-model:pagination="pagination"
      :data="sortedData"
      :columns="columns"
      :pagination-options="{
        getPaginationRowModel: getPaginationRowModel(),
      }"
      class="w-full"
    />

    <div class="flex justify-center border-t border-default pt-4">
      <UPagination
        :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1"
        :items-per-page="table?.tableApi?.getState().pagination.pageSize"
        :total="table?.tableApi?.getFilteredRowModel().rows.length"
        @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)"
      />
    </div>
  </div>
</template>
