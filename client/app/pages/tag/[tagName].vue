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
          {{ formattedDate }}
        </dd>
      </div>
      <div class="flex justify-start">
        <dt class="font-medium">
          作った人
        </dt>
        <!-- 作成者のアイコンが表示されるようにしたい -->
        <UAvatar
          :src="`https://q.trap.jp/api/v3/public/icon/${creator_id}`"
          class="text-3xl ml-4"
        />
        <dd class="font-medium ml-4">
          {{ creator_id }}
        </dd>
      </div>
      <div class="flex justify-start">
        <dt class="font-medium">
          スタンプの数
        </dt>
        <dd class="font-medium ml-4">
          {{ count }}個
        </dd>
      </div>
    </dl>
    <p class="text-xl font-medium">
      {{ route.params.tagName }} スタンプ一覧
    </p>
    <StampGrid />
  </UContainer>
</template>

<script setup lang="ts">
const route = useRoute();

const tagTestData = {
  id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
  creator_id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
  created_at: '2025-08-27T10:39:39.815Z',
  count: 1,
};

const { created_at, creator_id, count } = tagTestData;

const formattedDate = new Intl.DateTimeFormat('ja-JP', {
  dateStyle: 'long', // "yyyy年mm月dd日"
  timeStyle: 'short', // "hh:mm"
}).format(new Date(created_at));
</script>
