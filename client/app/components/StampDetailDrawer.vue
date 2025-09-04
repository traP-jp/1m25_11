<script setup lang="ts">
const selectedStampId = useSelectedStampId();
const { getStampById } = useStamps();

const drawerOpen = ref(false);

watch(selectedStampId, (newValue) => {
  drawerOpen.value = !!newValue;
});

const handleDrawerClose = (isOpen: boolean) => {
  if (!isOpen) {
    selectedStampId.value = null;
  }
};

const selectedStamp = computed(() => {
  if (!selectedStampId.value) return null;
  return getStampById(selectedStampId.value);
});

// const closeStampDrawer = () => {
//   selectedStampId.value = null;
// };

// const goToDetailPage = () => {
//   if (!selectedStampId.value) return;
//   closeStampDrawer();
//   navigateTo(`/stamps/${selectedStampId.value}`);
// };

const getStampImageUrl = (stampId: string) => {
  return `https://q.trap.jp/api/1.0/public/emoji/${stampId}`;
};
</script>

<template>
  <UDrawer
    v-model:open="drawerOpen"
    @update:open="handleDrawerClose"
  >
    <template #body>
      <!-- <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
            スタンプ詳細
          </h3>
          <div class="flex items-center gap-2">
            <UButton
              color="info"
              variant="ghost"
              icon="i-heroicons-arrows-pointing-out"
              @click="goToDetailPage"
            />
            <UButton
              color="info"
              variant="ghost"
              icon="i-heroicons-x-mark-20-solid"
              @click="closeStampDrawer"
            />
          </div>
        </div>
      </template> -->

      <div>
        <div
          v-if="selectedStamp"
          class="flex flex-col items-center space-y-4 p-4 mb-5"
        >
          <img
            :src="getStampImageUrl(selectedStamp.stamp_id)"
            alt="スタンプ画像"
            class="w-32 h-32 object-contain"
          >
          <div class="text-center">
            <p class="text-lg font-semibold">
              {{ selectedStamp.stamp_name }}
            </p>
            <p class="text-sm text-gray-500">
              ID: {{ selectedStamp.stamp_id }}
            </p>
          </div>
        </div>
      </div>
    </template>
  </UDrawer>
</template>
