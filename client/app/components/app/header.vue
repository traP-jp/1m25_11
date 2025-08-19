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
    <UAvatar
      v-if="username"
      :src="`https://q.trap.jp/api/v3/public/icon/${username}`"
      class="text-5xl"
    />
    <!-- </div> -->
  </UContainer>
</template>

<script setup lang="ts">
const username = useForwardedUser();
</script>
