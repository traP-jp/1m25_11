<script setup lang="ts">
const items = [
  'https://placehold.jp/3d4070/ffffff/800x600.png?text=Carousel%E2%86%92',
  'https://placehold.jp/07a0a5/ffffff/800x600.png?text=%E2%86%90Carousel',
];

const { data, error, status } = useFetch('/api/v1/ping');
if (error.value) {
  console.error('GETリクエストエラー:', error.value);
}
else {
  console.log('GETリクエスト結果:', data.value);
}
</script>

<template>
  <UCarousel
    v-slot="{ item }"
    :items="items"
    class="w-full max-w-xs mx-auto"
  >
    <img
      :src="item"
      width="320"
      height="320"
      class="rounded-lg"
    >
  </UCarousel>
</template>
