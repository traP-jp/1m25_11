<script setup lang="ts">
const stamps = useStamps();
const route = useRoute();
const stampName = stamps.getStampByName(route.params.stampName as string)?.file_id;
</script>

<template>
  <div v-if="stampName">
    {{ stampName }}
    <!-- <StampDetail /> -->
  </div>
  <div v-else>
    スタンプが見つかりません
  </div>
</template>
