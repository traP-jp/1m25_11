<template>
  <div>
    <!-- 認証済みユーザー向けコンテンツ（app.vue で認証制御済み） -->
    <SearchInput />
    <div class="flex justify-center my-4">
      <UButton
        :disabled="!isStampsReady"
        :loading="isStampsLoading"
        @click="generateRandomStamps"
      >
        ランダムに9個表示 ({{ stampsList.length }}件)
      </UButton>
    </div>
    <StampGrid :stamps="randomStamps" />
  </div>
</template>

<script setup lang="ts">
const { stampsList } = useStamps();

// ランダムに選ばれた9個のスタンプを保持するref
const randomStamps = ref<Schemas['StampSummary'][]>([]);

// スタンプデータの準備状況
const isStampsReady = computed(() => stampsList.value && stampsList.value.length > 0);
const isStampsLoading = computed(() => !isStampsReady.value);

// ランダムに9個のスタンプを選ぶ関数
const generateRandomStamps = () => {
  console.log('🎲 ランダムスタンプ生成開始');
  console.log('📊 利用可能スタンプ数:', stampsList.value?.length || 0);

  const allStamps = stampsList.value;
  if (!allStamps || allStamps.length === 0) {
    console.warn('⚠️ スタンプデータが利用できません');
    randomStamps.value = [];
    return;
  }

  // 配列をシャッフルして最初の9個を取得
  const shuffled = [...allStamps].sort(() => Math.random() - 0.5);
  randomStamps.value = shuffled.slice(0, Math.min(9, shuffled.length));
  console.log('✅ ランダムスタンプ生成完了:', randomStamps.value.length, '件');
};

// データ準備完了時に初期表示
watch(isStampsReady, (ready) => {
  if (ready) {
    console.log('📦 スタンプデータ準備完了 - 初期表示を実行');
    generateRandomStamps();
  }
}, { immediate: true });
</script>
