<template>
  <UCard
    variant="soft"
    class="bg-white hover:bg-sky-200 "
    :ui="{
      root: 'p-2 w-1/7 min-w-24',
      body: 'p-1 sm:p-1',
      footer: 'px-1 sm:px-1 py-1 mx-0',
    }"
  >
    <NuxtImg
      :src="`https://q.trap.jp/api/1.0/public/emoji/${stampItem.id}`"
      class="m-auto"
      loading="lazy"
    />
    <template #footer>
      <div class="text-center text-xs line-clamp-1 truncate overflow-hidden whitespace-nowrap ">
        {{ stampItem.name }}
      </div>
    </template>
  </UCard>
</template>

<script setup lang="ts">
defineProps<{ stampItem: stampSummary }>();
</script>
