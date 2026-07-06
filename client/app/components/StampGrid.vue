<template>
  <div class="grid gap-1 grid-cols-[repeat(auto-fill,96px)]  justify-between mt-5 content-start items-start">
    <StampItem
      v-for="stampItem of props.stamps"
      :key="stampItem.stamp_id"
      :stamp="stampItem"
    />
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ stamps: Schemas['StampSummary'][] }>();
console.log('StampGrid props:', props);
</script>
