<template>
  <UContainer>
    <UInput
      v-model="searchValue"
      icon="i-lucide-search"
      placeholder="スタンプの名前、タグ、説明文で検索..."
      class="mb-5"
    />
    <div
      class="flex justify-around gap-5"
    >
      <UDrawer
        v-model:open="isSearchModalOpen"
        :handle="false"
        title="Drawer with title"
      >
        <UButton
          :block="true"
          color="neutral"
          class="hover:cursor-pointer"
        >
          条件を指定
        </UButton>
        <template #body>
          <Placeholder class="h-48" />
        </template>
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
      <UButton
        :block="true"
        class="hover:cursor-pointer"
      >
        検索
      </UButton>
    </div>
  </UContainer>
</template>

<script setup lang="ts">
const searchValue = ref('');
const isSearchModalOpen = ref(false);
</script>
