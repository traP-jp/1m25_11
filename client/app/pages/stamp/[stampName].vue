<script setup lang="ts">
const stamps = useStamps();
const route = useRoute();
const stamp = stamps.getStampByName(route.params.stampName as string);
</script>

<template>
  <div v-if="stamp">
    <StampDetail :stamp-id="stamp?.stamp_id" />
    <!-- <StampDetail /> -->
  </div>
  <div v-else>
    スタンプが見つかりません
  </div>
</template>
