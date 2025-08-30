<template>
  <div>
    <SearchInput />
    <div class="flex justify-center my-4">
      <UButton @click="generateRandomTen">
        <!-- <button -->
        <!-- class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors" -->
        <!-- @click="generateRandomTen" -->
        <!-- > -->
        ランダムに9個表示
      <!-- </button> -->
      </UButton>
    </div>
    <StampGrid :stamps="randomTen" />
  </div>
</template>

<script setup lang="ts">
const { stampsList } = useStamps();

// ランダムに選ばれた10個のスタンプを保持するref
const randomTen = ref<StampSummary[]>([]);

// ランダムに10個のスタンプを選ぶ関数
const generateRandomTen = () => {
  const allStamps = stampsList.value;
  console.log('allStamps length:', allStamps.length);
  console.log('allStamps sample:', JSON.stringify(allStamps.slice(0, 2), null, 2));

  if (allStamps.length === 0) {
    randomTen.value = [];
    return;
  }

  // 配列をシャッフルして最初の9個を取得
  const shuffled = [...allStamps].sort(() => Math.random() - 0.5);
  randomTen.value = shuffled.slice(0, Math.min(9, shuffled.length));
  console.log('randomTen value:', JSON.stringify(randomTen.value.slice(0, 2), null, 2));
};

// デバッグ: スタンプデータの状態を監視
watchEffect(() => {
  console.log('Index.vue Debug:');
  console.log('- stampsList length:', stampsList.value.length);
  console.log('- randomTen length:', randomTen.value.length);
});

// onMounted(() => {
//   nextTick(() => {
//     generateRandomTen();
//   });
// });

// // stampsListにデータが入ったら実行
// watchEffect(() => {
//   if (stampsList.value.length > 0 && randomTen.value.length === 0) {
//     generateRandomTen();
//   }
// });
</script>
