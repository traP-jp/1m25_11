<template>
  <UContainer class="h-20 flex items-center justify-between font-bold">
    <!-- <div class="w-full h-20 flex items-center-safe justify-around font-bold text-3xl"> -->
    <USlideover
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
      <UAvatar
        :src="`https://q.trap.jp/api/v3/public/icon/${userName}`"
        class="text-5xl cursor-pointer"
      />
    </UDropdownMenu>

    <!-- </div> -->
  </UContainer>
</template>

<script setup lang="ts">
import type { NavigationMenuItem, DropdownMenuItem } from '@nuxt/ui';

const userName = useUser();
console.log(`userName: ${userName.value}`);

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
      label: 'Developer',
      icon: 'material-symbols:code',
      children: [
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
