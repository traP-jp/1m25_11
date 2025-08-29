<template>
  <!-- <UCard
    variant="soft"
    class="bg-white hover:bg-sky-200 h-32"
    :ui="{
      root: 'p-0 w-1/7 min-w-24',
      body: 'p-0 sm:p-1',
      footer: 'px-0 sm:px-0 py-0 mx-0 my-auto',
    }"
  > -->
  <div class="flex flex-col h-32 w-24 justify-center gap-1 hover:bg-primary-100">
    <NuxtImg
      :src="getFileUrl(stamp.id)"
      class="mx-auto w-full px-3"
      loading="lazy"
      :alt="stamp.name"
    />
    <!-- <template #footer> -->
    <div class="text-center text-xs line-clamp-1 truncate overflow-hidden whitespace-nowrap ">
      {{ stamp.name }}
    </div>
    <!-- </template> -->
  </div>
  <!-- </UCard> -->
</template>

<script setup lang="ts">
// import type { paths } from '~~/shared/types/generated';

defineProps<{
  stamp: StampSummary;
}>();

const getFileUrl = (fileId: string) => {
  // v1.0 を使用しているため、file_idではなくstamp_idを指定する必要がある
  return `https://q.trap.jp/api/1.0/public/emoji/${fileId}`;
};
</script>
