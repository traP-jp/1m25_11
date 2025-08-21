<template>
  <UContainer>
    <div class="grid gap-x-4 gap-y-6 grid-cols-[repeat(auto-fill,_96px)] justify-between content-start items-start">
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
  {
    id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
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
  {
    id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
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
  {
    id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
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
  {
    id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
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
  {
    id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
];
</script>
