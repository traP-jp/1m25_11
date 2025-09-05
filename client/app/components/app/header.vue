<template>
  <UContainer class="h-20 flex items-center justify-between font-bold">
    <!-- <div class="w-full h-20 flex items-center-safe justify-around font-bold text-3xl"> -->
    <USlideover
      v-model:open="navigationSlideOver"
      title="メニュー"
      side="left"
      :close="{
        color: 'primary',
        class: 'rounded-full',
      }"
    >
      <UIcon
        name="material-symbols:dehaze"
        class="cursor-pointer text-4xl text-primary"
      />
      <template #body>
        <UNavigationMenu
          orientation="vertical"
          :items="navigationItems"
          @click="navigationSlideOver=false"
        />
      </template>
    </USlideover>
    <h1 class="text-3xl">
      Service Name
    </h1>

    <!-- 認証状態に応じた表示（ClientOnlyで包む） -->
    <ClientOnly>
      <div v-if="isLoading">
        <!-- 認証状態確認中の表示 -->
        <UIcon
          name="material-symbols:progress-activity"
          class="text-2xl animate-spin text-primary"
        />
      </div>
      <div v-else-if="isLoggedIn">
        <UDropdownMenu
          :items="dropdownItems"
          :ui="{
            content: 'w-48',
          }"
        >
          <UAvatar
            :src="`https://q.trap.jp/api/v3/public/icon/${userName}`"
            class="text-5xl cursor-pointer"
          />
        </UDropdownMenu>
      </div>
      <div v-else>
        <UButton
          color="primary"
          size="lg"
          @click="handleLogin"
        >
          ログイン
        </UButton>
      </div>

      <!-- サーバーサイドでのフォールバック表示 -->
      <template #fallback>
        <UButton
          color="primary"
          size="lg"
          @click="handleLogin"
        >
          ログイン
        </UButton>
      </template>
    </ClientOnly>

    <!-- </div> -->
  </UContainer>
</template>

<script setup lang="ts">
import type { NavigationMenuItem, DropdownMenuItem } from '@nuxt/ui';

const { isLoggedIn, isLoading, login } = useAuth();
const userName = useUser();

// リアクティブな値の変化を監視
watch([isLoggedIn, isLoading, userName], ([newLoggedIn, newLoading, newUserName]) => {
  console.log(`状態変化 - userName: ${newUserName}, isLoggedIn: ${newLoggedIn}, isLoading: ${newLoading}`);
});

// 初期値をログ出力
console.log(`初期値 - userName: ${userName.value}, isLoggedIn: ${isLoggedIn.value}, isLoading: ${isLoading.value}`);

const handleLogin = () => {
  login();
};

const navigationSlideOver = ref(false);

const navigationItems = ref<NavigationMenuItem[][]>([
  [
    {
      label: 'Home',
      icon: 'material-symbols:home',
      to: '/',
    },
    {
      label: 'Search',
      icon: 'material-symbols:search',
      to: 'search',
    },
    {
      label: 'Ranking',
      icon: 'material-symbols:leaderboard',
      to: '/ranking',
    },
    {
      label: 'Tags',
      icon: 'material-symbols:tag',
      to: '/tag',
    },
    {
      label: 'Swagger viewer',
      to: '/developer',
      icon: 'material-symbols:info',
    },
    {
      label: 'GitHub',
      to: 'https://github.com/traP-jp/1m25_11',
      icon: 'material-symbols:arrow-outward',
    },
  ],
]);

const dropdownItems = ref<DropdownMenuItem[][]>([
  [
    {
      label: `${userName.value}`,
      avatar: {
        src: `https://q.trap.jp/api/v3/public/icon/${userName.value}`,
      },
      type: 'label',
    },
  ],
  [
    {
      label: 'Profile',
      icon: 'material-symbols:account-circle-outline-sharp',
      to: '/profile',
    },
    {
      label: 'Settings',
      icon: 'material-symbols:settings-outline-sharp',
      to: '/settings',
    },
  ],
]);
</script>
