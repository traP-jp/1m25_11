<template>
  <UContainer>
    <!-- 検索欄 -->
    <UInput
      v-model="searchValue"
      icon="i-lucide-search"
      placeholder="スタンプの名前、タグ、説明文で検索..."
      class="mb-5 w-full"
    />
    <!-- 検索条件を設定するボタンと検索ボタンのコンテナ -->
    <div
      class="flex justify-around gap-5"
    >
      <!-- 検索条件を設定を押すと開くメニュー -->
      <UDrawer
        v-model:open="isSearchModalOpen"
        :handle="false"
        title="Drawer with title"
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
          <Placeholder class="h-48" />
        </template>
        <!-- UDrawer コンポーネントの footer slot に入れる内容 -->
        <template #footer>
          <UButton
            label="Submit"
            color="neutral"
            class="justify-center"
          />
          <UButton
            label="Cancel"
            color="neutral"
            variant="outline"
            class="justify-center"
            @click="isSearchModalOpen = false"
          />
        </template>
      </UDrawer>
      <!-- input 要素のすぐ下にある検索ボタン -->
      <UButton
        label="検索"
        :block="true"
        class="hover:cursor-pointer"
      />
    </div>
  </UContainer>
</template>

<script setup lang="ts">
const searchValue = ref('');
const isSearchModalOpen = ref(false);
</script>
