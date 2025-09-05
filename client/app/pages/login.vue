<template>
  <UContainer class="flex min-h-screen items-center justify-center">
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <h1 class="text-2xl font-bold">
            Service Name
          </h1>
          <p class="text-gray-600 mt-2">
            traQアカウントでログインしてください
          </p>
        </div>
      </template>

      <div class="space-y-4">
        <UButton
          size="lg"
          color="primary"
          block
          icon="material-symbols:login"
          @click="handleLogin"
        >
          traQでログイン
        </UButton>
      </div>

      <template #footer>
        <div class="text-center text-sm text-gray-500">
          <p>traP内部サービスです</p>
        </div>
      </template>
    </UCard>
  </UContainer>
</template>

<script setup lang="ts">
const { login, isAuthenticated } = useAuth();

// 既にログイン済みの場合はホームにリダイレクト
watchEffect(() => {
  if (isAuthenticated.value) {
    navigateTo('/');
  }
});

const handleLogin = () => {
  login();
};

// メタ設定
definePageMeta({
  layout: false, // レイアウトを使用しない
});

useHead({
  title: 'ログイン - Service Name',
});
</script>
