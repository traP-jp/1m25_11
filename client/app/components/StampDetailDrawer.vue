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

const getStampImageUrl = (stampId: string | undefined) => {
  return stampId ? `https://q.trap.jp/api/1.0/public/emoji/${stampId}` : 'https://q.trap.jp/api/1.0/public/emoji/bc9a3814-f185-4b3d-ac1f-3c8f12ad7b52';
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
    :ui="{ header: 'flex items-center justify-between' }"
    title="aaaa"
    @update:open="handleDrawerClose"
  >
    <template #header>
      <div>
        <h3 class="text-highlighted font-semibold">
          :{{ selectedStamp?.stamp_name }}:
        </h3>
      </div>
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
    </template>
    <template #body>
      <div>
        <div
          class="flex flex-col items-center space-y-4 p-4 mb-5"
        >
          <NuxtImg
            :src="getStampImageUrl(selectedStamp?.stamp_id)"
            :alt="selectedStamp?.stamp_name"
            class="w-32 h-32 object-contain"
          />
          <!-- <div class="text-center">
            <p class="text-lg font-semibold">
              {{ selectedStamp.stamp_name }}
            </p>
            <p class="text-sm text-gray-500">
              ID: {{ selectedStamp.stamp_id }}
            </p>
          </div> -->
        </div>
      </div>
    </template>
  </UDrawer>
</template>
