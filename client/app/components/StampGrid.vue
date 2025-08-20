<template>
  <UContainer>
    <div class="grid grid-cols-3 gap-5">
      <StampItem
        v-for="stampItem in stampTestData"
        :key="stampItem.id"
        :stamp-item="stampItem"
      />
    </div>
  </UContainer>
</template>

<script setup lang="ts">
const stampTestData: stampSummary[] = [
  {
    id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    name: 'kyapi-nya',
    file_id: 'e247c3ea-077c-4817-a7cd-8e0185d576d3',
  },
  {
    id: '0197a63e-44f2-7779-b843-c805c52baacc',
    name: 'korosu-nya',
    file_id: '0197a63e-4471-7779-a0f8-8225104a20ce',
  },
];
</script>
