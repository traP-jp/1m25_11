<script setup lang="ts">
const stamps = useStamps();
const route = useRoute();

// デバッグ用のログ出力
watchEffect(() => {
  console.log('Debug - Loading:', stamps.loading.value);
  console.log('Debug - Error:', stamps.error.value);
  console.log('Debug - Stamps list length:', stamps.stampsList.value.length);
  console.log('Debug - Route param:', route.params.stampName);
});

// スタンプデータの読み込み完了を待ってからスタンプを検索
const stamp = computed(() => {
  if (stamps.loading.value || !stamps.stampsList.value.length) return null;
  return stamps.getStampByName(route.params.stampName as string) ?? null;
});

// ローディング状態とエラー状態を管理
const isLoading = computed(() => stamps.loading.value);
const hasError = computed(() => stamps.error.value);
</script>

<template>
  <div
    v-if="isLoading"
    class="text-center p-8"
  >
    <p>読み込み中...</p>
  </div>
  <div
    v-else-if="hasError"
    class="text-center p-8 text-red-500"
  >
    <p>データの読み込みに失敗しました</p>
  </div>
  <div
    v-else-if="stamp"
    class="w-full"
  >
    <StampDetail :stamp-id="stamp.stamp_id" />
  </div>
  <div
    v-else
    class="text-center p-8"
  >
    スタンプが見つかりません
  </div>
</template>
