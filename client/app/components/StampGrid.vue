<template>
  <UContainer>
    <div class="grid gap-1 grid-cols-[repeat(auto-fill,_88px)]  justify-between mt-5 content-start items-start">
      <StampItem
        v-for="stampItem of props.stamps"
        :key="stampItem.id"
        :stamp="stampItem"
      />
    </div>
  </UContainer>
</template>

<script setup lang="ts">
const props = defineProps<{ stamps: StampSummary[] }>();

// デバッグ用
watchEffect(() => {
  console.log('StampGrid received stamps:', props.stamps.length, 'items');
  if (props.stamps.length > 0) {
    console.log('First stamp:', JSON.stringify(props.stamps[0], null, 2));
  }
});

const stampTestData: StampSummary[] = [
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'kyapi-nya',
    file_id: 'e247c3ea-077c-4817-a7cd-8e0185d576d3',
  },
  {
    stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc',
    stamp_name: 'korosu-nya',
    file_id: '0197a63e-4471-7779-a0f8-8225104a20ce',
  },
  {
    stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    stamp_name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'kyapi-nya',
    file_id: 'e247c3ea-077c-4817-a7cd-8e0185d576d3',
  },
  {
    stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc',
    stamp_name: 'korosu-nya',
    file_id: '0197a63e-4471-7779-a0f8-8225104a20ce',
  },
  {
    stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    stamp_name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'kyapi-nya',
    file_id: 'e247c3ea-077c-4817-a7cd-8e0185d576d3',
  },
  {
    stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc',
    stamp_name: 'korosu-nya',
    file_id: '0197a63e-4471-7779-a0f8-8225104a20ce',
  },
  {
    stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    stamp_name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'kyapi-nya',
    file_id: 'e247c3ea-077c-4817-a7cd-8e0185d576d3',
  },
  {
    stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc',
    stamp_name: 'korosu-nya',
    file_id: '0197a63e-4471-7779-a0f8-8225104a20ce',
  },
  {
    stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    stamp_name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
  {
    stamp_id: 'e1e4a295-cf24-4de9-936c-e72469170d8f',
    stamp_name: 'kyapi-nya',
    file_id: 'e247c3ea-077c-4817-a7cd-8e0185d576d3',
  },
  {
    stamp_id: '0197a63e-44f2-7779-b843-c805c52baacc',
    stamp_name: 'korosu-nya',
    file_id: '0197a63e-4471-7779-a0f8-8225104a20ce',
  },
  {
    stamp_id: '0197a69d-c3ce-7822-9666-ace99bd35068',
    stamp_name: '403_forbidden',
    file_id: '0197a69d-c383-7822-84da-332d68451471',
  },
];
</script>
