<template>
  <UCard
    variant="soft"
    class="hover:bg-primary-100"
    :ui="{
      root: 'p-2 w-1/7 min-w-24',
      footer: 'px-1 sm:px-1 py-1 mx-0',
    }"
  >
    <NuxtImg :src="`https://q.trap.jp/api/1.0/public/emoji/${stampItem.id}`" />

    <template #footer>
      <div class="text-center text-xs">
        {{ stampItem.name }}
      </div>
    </template>
  </UCard>
</template>

<script setup lang="ts">
defineProps<{ stampItem: stampSummary }>();
</script>
