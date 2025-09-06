<template>
  <div>
    <!-- 認証状態確認中（初期状態・ローディング） -->
    <ClientOnly>
      <div
        v-if="isLoading"
        class="text-center py-8"
      >
        <UCard class="max-w-md mx-auto">
          <template #header>
            <h2 class="text-xl font-semibold text-center">
              読み込み中
            </h2>
          </template>
          <div class="text-center py-4">
            <UIcon
              name="i-heroicons-arrow-path"
              class="animate-spin text-2xl mb-2"
            />
            <p class="text-gray-600">
              認証状態を確認しています...
            </p>
          </div>
        </UCard>
      </div>

      <!-- 未認証時（ログイン誘導） -->
      <div
        v-else-if="!isLoggedIn"
        class="text-center py-8"
      >
        <UCard class="max-w-md mx-auto">
          <template #header>
            <h2 class="text-xl font-semibold text-center">
              ログインが必要です
            </h2>
          </template>
          <p class="text-gray-600 mb-4 text-center">
            このアプリケーションを利用するには、traQアカウントでのログインが必要です。
          </p>
          <div class="text-center">
            <UButton
              color="primary"
              size="lg"
              @click="handleLogin"
            >
              traQでログイン
            </UButton>
          </div>
        </UCard>
      </div>

      <!-- 認証済みユーザー向けコンテンツ -->
      <div v-else>
        <SearchInput />
        <div class="flex justify-center my-4">
          <UButton @click="generateRandomStamps">
            ランダムに9個表示
          </UButton>
        </div>
        <StampGrid :stamps="randomStamps" />
      </div>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
const { isLoggedIn, isLoading, login } = useAuth();
const { stampsList } = useStamps();

// ランダムに選ばれた9個のスタンプを保持するref
const randomStamps = ref<Schemas['StampSummary'][]>([]);

const handleLogin = () => {
  login();
};

// ランダムに9個のスタンプを選ぶ関数
const generateRandomStamps = () => {
  const allStamps = stampsList.value;
  if (allStamps.length === 0) {
    randomStamps.value = [];
    console.log('スタンプが存在しません');
    return;
  }

  // 配列をシャッフルして最初の9個を取得
  const shuffled = [...allStamps].sort(() => Math.random() - 0.5);
  randomStamps.value = shuffled.slice(0, Math.min(9, shuffled.length));
};

console.log(generateRandomStamps);
</script>
