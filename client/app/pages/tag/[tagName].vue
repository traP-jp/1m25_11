<template>
  <UContainer>
    <h1 class="text-3xl font-bold mb-4">
      タグ詳細
    </h1>

    <h2 class="text-xl font-semibold mb-2">
      #{{ route.params.tagName }}
    </h2>
    <dl class="space-y-2 mb-6">
      <div class="flex justify-between">
        <dt class="font-medium">
          作られた日
        </dt>
        <dd />
      </div>
      <div class="flex justify-between">
        <dt class="font-medium">
          作った人
        </dt>
        <dd />
      </div>
      <div class="flex justify-between">
        <dt class="font-medium">
          スタンプの数
        </dt>
      </div>
      <StampGrid />
    </dl>
  </UContainer>
</template>

<script setup lang="ts">
const route = useRoute();
</script>
