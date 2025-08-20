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
    <UFormField
      label="タグ名を指定"
      required
    >
      <UInputTags
        :v-model="searchStampTagValue"
        placeholder="タグ名を入力"
        class="mb-5 w-full"
      />
    </UFormField>
    <p>作成日を指定</p>
    <CarendarModal />
  </ucontainer>
</template>

<script setup lang="ts">
import CarendarModal from './CarendarModal.vue';

const searchDescriptionValue = ref('');
const searchStampNameValue = ref('');
const searchStampTagValue = ref(['Vue']);
</script>
