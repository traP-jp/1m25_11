<template>
  <UContainer>
    <!-- 検索欄 -->
    <UInput
      v-model="searchValue"
      icon="i-lucide-search"
      placeholder="スタンプの名前、タグ、説明文で検索..."
      class="mb-5 w-full"
      size="xl"
    />
    <!-- 検索条件を設定するボタンと検索ボタンのコンテナ -->
    <div
      class="flex justify-around gap-5"
    >
      <!-- 検索条件を設定を押すと開くメニュー -->
      <UDrawer
        v-model:open="isSearchModalOpen"
        :handle="false"
        title="条件を設定"
      >
        <!-- input要素の下にある検索条件を設定するボタン(メニューが開く) -->
        <UButton
          label="条件を設定"
          :block="true"
          color="neutral"
          class="hover:cursor-pointer"
        />
        <!-- UDrawer コンポーネントの body slot に入れる内容 -->
        <template #body>
          <SearchInputModal />
        </template>
      </UDrawer>
      <!-- input 要素のすぐ下にある検索ボタン -->
      <UButton
        label="検索"
        :block="true"
        class="hover:cursor-pointer"
        icon="i-lucide-arrow-right"
      />
    </div>
  </UContainer>
</template>

<script setup lang="ts">
const searchValue = ref('');
const isSearchModalOpen = ref(false);
</script>
