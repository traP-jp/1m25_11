<template>
  <UContainer class="h-20 flex items-center justify-between font-bold">
    <!-- <div class="w-full h-20 flex items-center-safe justify-around font-bold text-3xl"> -->
    <UModal
      title="メニュー"
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
        <ul>
          <li><ULink to="/">Top Page</ULink></li>
          <li><ULink to="/">Top Page</ULink></li>
        </ul>
      </template>
    </UModal>
    <h1 class="text-3xl">
      Service Name
    </h1>

    <UDropdownMenu
      :items="items"
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
// const runtimeConfig = useRuntimeConfig();
// const username = useForwardedUser();
import type { DropdownMenuItem } from '@nuxt/ui';

const userName = useUser();
// const userAvaterImage = runtimeConfig
console.log(`useName: ${userName.value}`);

const items = ref<DropdownMenuItem[][]>([
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
