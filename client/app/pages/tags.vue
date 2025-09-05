<template>
  <UContainer>
    <h1 class="font-bold text-2xl">
      タグ一覧
    </h1>
    <div class="pt-4 flex flex-wrap gap-2">
      <NuxtLink
        v-for="tagItem in tagTestData"
        :key="tagItem.tag_id"
        :to="`/tag/${encodeURIComponent(tagItem.tag_name)}`"
        class="text-xl font-medium inline-block px-3 py-1 border-0 bg-gray-100 text-primary cursor-pointer hover:bg-gray-200 transition-colors"
      >
        #{{ tagItem.tag_name }}
      </NuxtLink>
    </div>
  </UContainer>
</template>

<script setup lang="ts">
interface tagSummary {
  tag_id: string;
  tag_name: string;
}

const tagTestData: tagSummary[] = [
  {
    tag_id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
    tag_name: '動物',
  },
  {
    tag_id: '7fa85f64-5717-4562-b3fc-2c963f66afa6',
    tag_name: 'スポーツ',
  },
  {
    tag_id: '9fa85f64-5717-4562-b3fc-2c963f66afa6',
    tag_name: '音楽',
  },
];
</script>
