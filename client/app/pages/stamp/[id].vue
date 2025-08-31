<script setup lang="ts">
const stamps = useStamps();
const route = useRoute();
const stampId = stamps.getStampById(route.params.id as string)?.file_id;
</script>

<template>
  <div v-if="stampId">
    {{ stampId }}
    <StampDetailDrawer />
  </div>
  <div v-else>
    スタンプが見つかりません
  </div>
</template>
