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
    @update:open="handleDrawerClose"
  >
    <template #header>
      <div class="flex items-center justify-between w-full">
        <span class="font-semibold">スタンプの詳細</span>
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
      </div>
    </template>
    <template #body>
      <StampDetail :stamp-id="selectedStampId" />
    </template>
  </UDrawer>
</template>
