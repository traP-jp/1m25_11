<template>
  <UContainer>
    <h1 class="text-3xl font-bold mb-4">
      タグ詳細
    </h1>

    <h2 class="text-2xl font-semibold mb-4">
      #{{ route.params.tagName }}
    </h2>
    <dl class="text-xl space-y-2 mb-6">
      <div class="flex justify-start">
        <dt class="font-medium">
          作られた日
        </dt>
        <dd class="font-medium ml-4">
          333
        </dd>
      </div>
      <div class="flex justify-start">
        <dt class="font-medium">
          作った人
        </dt>
        <dd class="font-medium ml-4">
          333
        </dd>
      </div>
      <div class="flex justify-start">
        <dt class="font-medium">
          スタンプの数
        </dt>
        <dd class="font-medium ml-4">
          333
        </dd>
      </div>
      <StampGrid />
    </dl>
  </UContainer>
</template>

<script setup lang="ts">
const route = useRoute();
</script>
