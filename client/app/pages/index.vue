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
const { stampsList } = useStamps();

// ランダムに選ばれた10個のスタンプを保持するref
const randomStamps = ref<Schemas['StampSummary'][]>([]);

// ランダムに10個のスタンプを選ぶ関数
const generateRandomStamps = () => {
  const allStamps = stampsList.value;
  if (allStamps.length === 0) {
    randomStamps.value = [];
    return;
  }

  // 配列をシャッフルして最初の9個を取得
  const shuffled = [...allStamps].sort(() => Math.random() - 0.5);
  randomStamps.value = shuffled.slice(0, Math.min(9, shuffled.length));
};
</script>
