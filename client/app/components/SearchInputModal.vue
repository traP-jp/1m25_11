<template>
  <UContainer>
    <p>説明文で検索</p>
    <UInput
      v-model="searchDescriptionValue"
      icon="i-lucide-search"
      placeholder="説明文で検索"
      class="mb-5 w-full"
    />
    <p>スタンプ名で検索</p>
    <UInput
      v-model="searchStampNameValue"
      icon="i-lucide-search"
      placeholder="スタンプ名で検索"
      class="mb-5 w-full"
    />
  </UContainer>
</template>

<script setup lang="ts">
const searchDescriptionValue = ref('');
const searchStampNameValue = ref('');
</script>
