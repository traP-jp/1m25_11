<template>
  <UContainer>
    <div class="grid gap-1 grid-cols-[repeat(auto-fill,_96px)]  justify-between mt-5 content-start items-start">
      <StampItem
        v-for="stampItem of props.stamps"
        :key="stampItem.stamp_id"
        :stamp="stampItem"
      />
    </div>
  </UContainer>
</template>

<script setup lang="ts">
const props = defineProps<{ stamps: Schemas['StampSummary'][] }>();
console.log('StampGrid props:', props);
</script>
