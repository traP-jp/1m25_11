<template>
  <div>
    <SearchInput />
    <div class="flex justify-center my-4">
      <UButton @click="generateRandomStamps">
        ランダムに9個表示
      </UButton>
    </div>
    <StampGrid :stamps="randomStamps" />
  </div>
</template>

<script setup lang="ts">
// ランダムに選ばれた10個のスタンプを保持するref
const randomStamps = ref<Schemas['StampSummary'][]>([]);
const apiClient = useApiClient();

// ランダムに10個のスタンプを選ぶ関数
const generateRandomStamps = async () => {
  console.log(1);
  // const allStamps = stampsList.value;

  const { data: allStamps, error } = await apiClient.GET('/stamps');
  if (error || !allStamps) {
  // if (allStamps.length === 0) {
    randomStamps.value = [];
    console.log(2);
    return;
  }

  console.log(allStamps);
  // 配列をシャッフルして最初の9個を取得
  const shuffled = [...allStamps].sort(() => Math.random() - 0.5);
  randomStamps.value = shuffled.slice(0, Math.min(9, shuffled.length));
};
</script>
