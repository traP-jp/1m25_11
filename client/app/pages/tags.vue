<template>
  <UContainer>
    <h1 class="font-bold text-2xl">
      タグ一覧
    </h1>
    <div
      v-if="tagList.length === 0"
      class="pt-4 flex flex-wrap gap-2"
    >
      <NuxtLink
        v-for="tagItem in tagTestData"
        :key="tagItem.tag_id"
        :to="`/tag/${encodeURIComponent(tagItem.tag_name)}`"
        class="text-xl font-medium inline-block px-3 py-1 border-0 bg-gray-100 text-primary cursor-pointer hover:bg-gray-200 transition-colors"
      >
        #{{ tagItem.tag_name }}
      </NuxtLink>
    </div>
    <div
      v-else
      class="pt-4 flex flex-wrap gap-2"
    >
      <NuxtLink
        v-for="tagItem in tagList"
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
const tagTestData = ref<Schemas['TagSummary'][]>([
  { tag_id: '1', tag_name: 'ありがとう' },
  { tag_id: '2', tag_name: 'おめでとう' },
  { tag_id: '3', tag_name: 'よろしく' },
]);

const tagList = ref<Schemas['TagSummary'][]>([]);
onMounted(async () => {
  try {
    const res = await apiClient.GET('/tags');
    console.log('API レスポンス:', res);
    if (res.data) {
      tagList.value = res.data;
      console.log('tagList updated:', tagList.value);
    }
  }
  catch (err) {
    console.error('タグ取得エラー:', err);
  }
});
</script>
