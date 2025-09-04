<template>
  <UContainer>
    <div
      v-if="props.stampId"
      class="flex flex-col items-center space-y-4 p-4 pb-0"
    >
      <NuxtImg
        :src="getFileUrl(props.stampId)"
        :alt="stampDetail?.stamp_name || 'スタンプ'"
        class="w-2/3 h-2/3 max-w-md px-0 object-contain pointer-events-none"
      />
    </div>

    <div
      v-if="loading"
      class="text-center"
    >
      <p>読み込み中...</p>
    </div>

    <div
      v-else-if="error"
      class="text-center text-red-500"
    >
      <p>{{ error }}</p>
    </div>

    <div v-else-if="stampDetail">
      <h2 class="text-2xl py-1 mt-1 font-bold text-center overflow-hidden">
        {{ stampDetail.stamp_name }}
      </h2>
      <div class="flex justify-center mt-4 gap-5">
        <UButton
          variant="subtle"
          @click="copySelectedStampName()"
        >
          <UIcon name="material-symbols:content-copy-outline-sharp" />
          <span class="ml-1">copy</span>
        </UButton>
        <UButton
          variant="subtle"
          @click="copySelectedStampWithColon()"
        >
          <UIcon name="material-symbols:content-copy-outline-sharp" />
          <span class="ml-1">copy with colon</span>
        </UButton>
        <UButton
          variant="subtle"
          @click="share()"
        >
          <UIcon name="material-symbols:ios-share-sharp" />
        </UButton>
      </div>
      <UCard class="my-4">
        <template #header>
          <h3 class="text-xl font-medium">
            タグ
          </h3>
        </template>
        <div class="flex w-full flex-wrap justify-start gap-2 text-xs mt-2">
          <UBadge
            v-for="tag in sampleTag"
            :key="tag.tag_id"
            icon="material-symbols:tag"
            class="bg-gray-100 text-primary cursor-pointer hover:bg-gray-200 transition-colors"
          >
            {{ tag.tag_name }}
          </ubadge>
        </div>
      </UCard>
      <UCard
        class="my-4"
      >
        <template #header>
          <h3 class="text-xl font-medium">
            概要
          </h3>
        </template>
        <dl class="divide-y divide-primary-100">
          <div class="px-2 py-2 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
            <dt class="text-sm/6 font-medium text-primary">
              保持者
            </dt>
            <dd class="mt-1 text-sm/6 text-gray-600 sm:col-span-2 sm:mt-0 gap-1 flex items-center min-w-0">
              <template v-if="stampDetail.is_unicode">
                <span>Unicode絵文字</span>
              </template>
              <template v-else>
                <UAvatar
                  :src="`https://q.trap.jp/api/v3/public/icon/${users.getUserById(stampDetail.creator_id)?.traq_id}`"
                  :alt="stampDetail.creator_id"
                  size="sm"
                  class="flex-shrink-0"
                />
                <span class="overflow-x-scroll whitespace-nowrap min-w-0 text-sm/7">{{ users.getUserById(stampDetail.creator_id)?.user_display_name }}</span>
                <span class="flex-shrink-0">(<span class="text-primary-400">@{{ users.getUserById(stampDetail.creator_id)?.traq_id }}</span><span
                  v-if="users.getUserById(stampDetail.creator_id)?.user_state == 0"
                  class="text-xs"
                >&nbsp;凍結済み</span>)</span>
                <UButton
                  icon="material-symbols:content-copy-outline-sharp"
                  size="xs"
                  color="primary"
                  variant="subtle"
                  class="cursor-pointer flex-shrink-0 h-7"
                  @click="copyTraqId"
                >
                  traQ ID
                </UButton>
              </template>
            </dd>
          </div>
          <div class="px-2 py-2 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
            <dt class="text-sm/6 font-medium text-primary">
              スタンプID
            </dt>
            <dd class="mt-1 text-sm/6 text-gray-600 sm:col-span-2 sm:mt-0 flex items-center gap-3">
              <span class="overflow-x-scroll whitespace-nowrap min-w-0 text-sm/7">{{ stampDetail.stamp_id }}</span>
              <UButton
                icon="material-symbols:content-copy-outline-sharp"
                size="xs"
                color="primary"
                variant="subtle"
                class="cursor-pointer h-7 flex-shrink-0"
                @click="copyStampId"
              >
                スタンプID
              </UButton>
            </dd>
          </div>
          <div class="px-2 py-2 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
            <dt class="text-sm/6 font-medium text-primary">
              作成日
            </dt>
            <dd class="mt-1 text-sm/6 text-gray-600 sm:col-span-2 sm:mt-0">
              {{ new Date(stampDetail.created_at).toLocaleDateString() }}
            </dd>
          </div>
          <div class="px-2 py-2 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
            <dt class="text-sm/6 font-medium text-primary">
              最終更新日
            </dt>
            <dd class="mt-1 text-sm/6 text-gray-600 sm:col-span-2 sm:mt-0">
              {{ new Date(stampDetail.updated_at).toLocaleDateString() }}
            </dd>
          </div>
        </dl>
      </UCard>
      <UCard class="my-4">
        <template #header>
          <h3 class="text-xl font-medium">
            説明
          </h3>
        </template>
        <div
          v-if="sampleDescriptions && sampleDescriptions.length > 0"
          class="space-y-3"
        >
          <div
            v-for="description in sampleDescriptions"
            :key="description.creator_id"
            class="border-l-4 border-primary-200 pl-4 py-2"
          >
            <p class="text-sm text-gray-700 mb-2">
              {{ description.description }}
            </p>
            <div class="flex items-center gap-1 text-xs text-gray-500 min-w-0">
              <template v-if="description.creator_id == '3b261ff3-f940-4e2c-a626-27387b6dd71b'">
                <span>LLMによる説明</span>
                <span class="flex-shrink-0">•</span>
                <span class="flex-shrink-0 whitespace-nowrap">{{ new Date(description.created_at).toLocaleDateString() }}</span>
              </template>
              <template v-else>
                <UAvatar
                  :src="`https://q.trap.jp/api/v3/public/icon/${users.getUserById(description.creator_id)?.traq_id}`"
                  :alt="description.creator_id"
                  size="xs"
                  class="flex-shrink-0"
                />
                <template v-if="users.getUserById(description.creator_id)">
                  <span class="overflow-hidden text-ellipsis whitespace-nowrap min-w-0">{{ users.getUserById(description.creator_id)?.user_display_name }}</span>
                  <span class="flex-shrink-0 whitespace-nowrap">(@{{ users.getUserById(description.creator_id)?.traq_id }}<template v-if="users.getUserById(description.creator_id)?.user_state == 0"><span class="text-xs">&nbsp;凍結済み</span></template>)</span>
                </template>
                <template v-else>
                  <span class="flex-shrink-0">凍結されたユーザー ({{ description.creator_id }})</span>
                </template>
                <span class="flex-shrink-0">•</span>
                <span class="flex-shrink-0 whitespace-nowrap">{{ new Date(description.created_at).toLocaleDateString() }}</span>
              </template>
            </div>
          </div>
        </div>
        <div
          v-else
          class="text-center text-gray-500 py-4"
        >
          <p>説明はまだありません</p>
        </div>
      </UCard>
    </div>

    <div
      v-else
      class="text-center text-gray-500"
    >
      <p>スタンプの詳細情報がありません</p>
    </div>
  </UContainer>
</template>

<script setup lang="ts">
const props = defineProps<{ stampId: string | undefined }>();
const { getStampById } = useStamps();

const { stampDetail, loading, error } = useStampDetail(props.stampId);

const isCopyingStampName = ref(false);
const isCopyingStampNameWithColon = ref(false);
const isCopyingStampId = ref(false);
const isCopyingTraqId = ref(false);
const toast = useToast();
const users = useUsers();

const copySelectedStampName = () => {
  if (!props.stampId) return;
  isCopyingStampName.value = true;
  navigator.clipboard.writeText(`${getStampById(props.stampId)?.stamp_name}`);
  toast.add({
    title: 'copied',
  });
  setTimeout(() => {
    isCopyingStampName.value = false;
  }, 500);
};

const copySelectedStampWithColon = () => {
  if (!props.stampId) return;
  isCopyingStampNameWithColon.value = true;
  navigator.clipboard.writeText(`:${getStampById(props.stampId)?.stamp_name}:`);
  toast.add({
    title: 'copied with colon',
  });
  setTimeout(() => {
    isCopyingStampNameWithColon.value = false;
  }, 500);
};

const copyStampId = () => {
  if (!props.stampId) return;
  isCopyingStampId.value = true;
  navigator.clipboard.writeText(props.stampId);
  toast.add({
    title: 'copied',
  });
  setTimeout(() => {
    isCopyingStampId.value = false;
  }, 500);
};

const copyTraqId = () => {
  if (!stampDetail.value) return;

  const user = users.getUserById(stampDetail.value.creator_id);
  if (!user?.traq_id) return;

  isCopyingTraqId.value = true;
  navigator.clipboard.writeText(user.traq_id);
  toast.add({
    title: 'copied',
  });
  setTimeout(() => {
    isCopyingTraqId.value = false;
  }, 500);
};

const share = () => {
  if (!props.stampId) return;
  const stamp = getStampById(props.stampId);
  if (!stamp) return;
  const shareData = {
    title: stamp.stamp_name,
    text: `https://stampedia.trap.show/stamps/${getStampById(props.stampId)?.stamp_name}`,
    url: `https://stampedia.trap.show/stamps/${getStampById(props.stampId)?.stamp_name}`,
  };
  navigator.share(shareData).catch((err) => {
    console.error('Error sharing:', err);
  });
};

const sampleTag: Schemas['TagSummary'][] = [
  {
    tag_id: 'sample-tag-id1',
    tag_name: 'タグ1',
  },
  {
    tag_id: 'sample-tag-id2',
    tag_name: 'タグ2',
  },
  {
    tag_id: 'sample-tag-id3',
    tag_name: 'サンプルタグ3',
  },
  {
    tag_id: 'sample-tag-id4',
    tag_name: 'タグtag４',
  },
];

const sampleDescriptions: Schemas['StampDescription'][] = [
  {
    description: 'スタンプ説明文説明文',
    creator_id: '3b261ff3-f940-4e2c-a626-27387b6dd71b',
    created_at: '2024-01-15T10:30:00Z',
    updated_at: '2024-01-15T10:30:00Z',
  },
  {
    description: 'あああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああ',
    creator_id: '01960efa-f1ed-7d54-bf3d-6d62fe8af5aa',
    created_at: '2024-02-20T14:45:00Z',
    updated_at: '2024-02-20T14:45:00Z',
  },
  {
    description: '-nya シリーズ\nあああｓ',
    creator_id: '01963cec-d1bb-7a6e-b8df-16c7ca978464',
    created_at: '2024-03-10T09:15:00Z',
    updated_at: '2024-03-10T09:15:00Z',
  },
  {
    description: '凍結されたユーザーのテスト',
    creator_id: '2cc1df43-d5d7-42aa-8831-00a4efe48ce4',
    created_at: '2024-03-10T09:15:00Z',
    updated_at: '2024-03-10T09:15:00Z',
  },
];
</script>
