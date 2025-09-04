<script setup lang="ts">
const selectedStampId = useSelectedStampId();
const { getStampById } = useStamps();

const drawerOpen = ref(false);

watch(selectedStampId, (newValue) => {
  drawerOpen.value = !!newValue;
});

const handleDrawerClose = (isOpen: boolean) => {
  if (!isOpen) {
    selectedStampId.value = undefined;
  }
};

const closeStampDrawer = () => {
  selectedStampId.value = undefined;
};

const goToDetailPage = () => {
  if (!selectedStampId.value) return;
  navigateTo(`/stamp/${getStampById(selectedStampId.value)?.stamp_name}`);
  closeStampDrawer();
};
</script>

<template>
  <UDrawer
    v-model:open="drawerOpen"
    class=" max-h-[90%]"
    :ui="{ header: ' top-0 -mt-px bg-white w-full pb-3', container: 'pt-0 mt-4 pb-0 gap-0', body: 'overflow-y-scroll' }"
    @update:open="handleDrawerClose"
  >
    <template #title>
      スタンプの詳細
    </template>
    <template #description>
      スタンプの詳細情報を表示します
    </template>
    <template #header>
      <UContainer class="flex items-center justify-end">
        <div class="flex items-center gap-2">
          <UButton
            color="primary"
            variant="ghost"
            icon="material-symbols:expand-content"
            @click="goToDetailPage"
          />
          <UButton
            color="primary"
            variant="ghost"
            icon="material-symbols:close"
            @click="closeStampDrawer"
          />
        </div>
      </UContainer>
    </template>
    <template #body>
      <StampDetail :stamp-id="selectedStampId" />
    </template>
  </UDrawer>
</template>
