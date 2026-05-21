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
      stamPedia
    </h1>

    <UAvatar
      v-if="currentUser"
      :src="`https://q.trap.jp/api/v3/public/icon/${currentUser.traq_id}`"
      class="text-5xl cursor-pointer"
    />
    <!-- <UDropdownMenu
      :items="dropdownItems"
      :ui="{
        content: 'w-48',
      }"
    >
    <div v-if="currentUser">
      <UAvatar
        :src="`https://q.trap.jp/api/v3/public/icon/${currentUser.traq_id}`"
        class="text-5xl cursor-pointer"
      />
    </div>
    </UDropdownMenu> -->
  </UContainer>
</template>

<script setup lang="ts">
import type { NavigationMenuItem, DropdownMenuItem } from '@nuxt/ui';

const navigationSlideOver = ref(false);
const { currentUser } = useCurrentUser();

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
      to: '/tags',
    },
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

// const dropdownItems = ref<DropdownMenuItem[][]>([
//   [
//     {
//       label: 'Profile',
//       icon: 'material-symbols:account-circle-outline-sharp',
//       to: '/profile',
//     },
//     {
//       label: 'Settings',
//       icon: 'material-symbols:settings-outline-sharp',
//       to: '/settings',
//     },
//   ],
// ]);
</script>
