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

    <UDropdownMenu
      :items="dropdownItems"
      :ui="{
        content: 'w-48',
      }"
    >
      <!-- <UAvatar
        :src="`https://q.trap.jp/api/v3/public/icon/${userName.getUserById(user.user.value?.user_id)?.traq_id}`"
        class="text-5xl cursor-pointer"
      /> -->
      <!-- <UAvatar
        :src="`https://q.trap.jp/api/v3/public/icon/${currentUser.value?.traq_id || 'traP'}`"
        class="text-5xl cursor-pointer"
      /> -->
      <UAvatar
        v-if="currentUser"
        :src="`https://q.trap.jp/api/v3/public/icon/${currentUser.traq_id}`"
        class="text-5xl cursor-pointer"
      />
    </UDropdownMenu>

    <!-- </div> -->
  </UContainer>
</template>

<script setup lang="ts">
import type { NavigationMenuItem, DropdownMenuItem } from '@nuxt/ui';

const navigationSlideOver = ref(false);

const { data: users } = await useApiClient().GET('/users-list');
const { data: me } = await useApiClient().GET('/me');

const currentUser = computed(() => {
  if (!me?.user_id || !users) return null;
  return users.find(user => user.user_id === me.user_id);
});

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
      label: `${currentUser.value?.user_display_name}`,
      avatar: {
        src: `https://q.trap.jp/api/v3/public/icon/${currentUser.value?.traq_id}`,
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
