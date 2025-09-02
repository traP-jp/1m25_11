<script setup lang="ts">
const selectedStampId = useSelectedStampId();
const { getStampById } = useStamps();
const toast = useToast();

const drawerOpen = ref(false);
const isCopying = ref(false);

watch(selectedStampId, (newValue) => {
  drawerOpen.value = !!newValue;
});

const handleDrawerClose = (isOpen: boolean) => {
  if (!isOpen) {
    selectedStampId.value = undefined;
  }
};

const selectedStamp = computed(() => {
  if (!selectedStampId.value) return undefined;
  return getStampById(selectedStampId.value);
});

const closeStampDrawer = () => {
  selectedStampId.value = undefined;
};

const goToDetailPage = () => {
  if (!selectedStampId.value) return;
  navigateTo(`/stamp/${selectedStampId.value}`);
  closeStampDrawer();
};

const copySelectedStampName = () => {
  if (!selectedStamp.value) return;
  isCopying.value = true;
  navigator.clipboard.writeText(`:${selectedStamp.value.stamp_name}:`);
  toast.add({
    title: 'クリップボードにコピーしました。',
  });
  setTimeout(() => {
    isCopying.value = false;
  }, 500);
};
</script>

<template>
  <UDrawer
    v-model:open="drawerOpen"
    class=" max-h-[90%]"
    :ui="{ header: ' top-0 -mt-px bg-white w-full pb-3', container: 'pt-0 mt-4 pb-0 gap-0', body: 'overflow-y-scroll' }"
    @update:open="handleDrawerClose"
  >
    <template #header>
      <UContainer class="flex items-center justify-end">
        <div class="flex items-center gap-2">
          <UButton
            color="info"
            variant="ghost"
            icon="material-symbols:content-copy-outline-sharp"
            @click="copySelectedStampName"
          />
          <UButton
            color="info"
            variant="ghost"
            icon="material-symbols:expand-content"
            @click="goToDetailPage"
          />
          <UButton
            color="info"
            variant="ghost"
            icon="material-symbols:close"
            @click="closeStampDrawer"
          />
        </div>
      </UContainer>
    </template>
    <template #body>
      <UContainer>
        <StampDetail :stamp-id="selectedStampId" />
        <span>
          ああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああ

          ああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああ

          ああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああ
        </span>
      </UContainer>
    </template>
  </UDrawer>
</template>
