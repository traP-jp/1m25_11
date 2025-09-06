<template>
  <UApp>
    <!-- 認証状態確認中 -->
    <div
      v-if="isLoading"
      class="min-h-screen flex items-center justify-center"
    >
      <div class="text-center">
        <UIcon
          name="i-heroicons-arrow-path"
          class="animate-spin text-4xl mb-4 text-primary"
        />
        <h2 class="text-xl font-semibold mb-2">
          認証状態を確認中
        </h2>
        <p class="text-gray-600">
          しばらくお待ちください...
        </p>
      </div>
    </div>

    <!-- 未認証: ログインページリダイレクト通知 -->
    <div
      v-else-if="!isLoggedIn"
      class="min-h-screen flex items-center justify-center"
    >
      <div class="text-center">
        <UIcon
          name="i-heroicons-arrow-right-circle"
          class="text-4xl mb-4 text-primary"
        />
        <h2 class="text-xl font-semibold mb-2">
          ログインページに移動中
        </h2>
        <p class="text-gray-600">
          認証が完了していないため、ログインページに転送します...
        </p>
      </div>
    </div>

    <!-- 認証済み: 通常のアプリケーションレイアウト -->
    <div v-else>
      <NuxtRouteAnnouncer />
      <AppHeader />
      <NuxtPage />
      <StampDrawer />
    </div>
  </UApp>
</template>

<script setup>
const { isLoggedIn, isLoading } = useAuth();
const config = useRuntimeConfig();

// 未認証時のリダイレクト処理
watch([isLoggedIn, isLoading], ([loggedIn, loading]) => {
  if (!loading && !loggedIn) {
    console.log('未認証のため、ログインページにリダイレクト');
    window.location.href = `${config.public.apiBase}/login`;
  }
}, { immediate: true });
</script>
