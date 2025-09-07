<template>
  <div
    class="flex flex-col h-32 w-24 justify-center gap-0.5 hover:bg-primary-100 cursor-pointer"
    @click="selectStamp(stamp.stamp_id)"
  >
    <NuxtImg
      :src="getFileUrl(stamp.stamp_id)"
      class="mx-auto w-full px-3"
      loading="lazy"
      :alt="stamp.stamp_name"
    />
    <div class="text-center text-xs truncate w-full px-2.5 overflow-hidden whitespace-nowrap overflow-ellipsis">
      {{ stamp.stamp_name }}
    </div>
  </div>
</template>

<script setup lang="ts">
const selectedStampId = useSelectedStampId();

const { stamp } = defineProps<{
  stamp: Schemas['StampSummary'];
}>();

console.log(stamp.stamp_id);

const selectStamp = (stampId: string) => {
  selectedStampId.value = stampId;
};
</script>
