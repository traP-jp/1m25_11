<template>
  <div>
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

      <div class="space-y-2">
        <p><strong>作成日:</strong> {{ new Date(stampDetail.created_at).toLocaleDateString() }}</p>
        <p><strong>月間使用回数:</strong> {{ stampDetail.count_monthly }}回</p>
        <p><strong>総使用回数:</strong> {{ stampDetail.count_total }}回</p>

        <div v-if="sampleTag && sampleTag.length > 0">
          <strong>タグ:</strong>
          <span
            v-for="tag in sampleTag"
            :key="tag.tag_id"
            class="inline-block bg-gray-200 rounded px-2 py-1 text-sm mr-1 mt-1 "
          >
            {{ tag.tag_name }}
          </span>
        </div>

        <div v-if="stampDetail.descriptions && stampDetail.descriptions.length > 0">
          <strong>説明:</strong>
          <ul class="list-disc list-inside">
            <li
              v-for="desc in stampDetail.descriptions"
              :key="desc.creator_id + props.stampId"
            >
              {{ desc.description }}
            </li>
          </ul>
        </div>
      </div>
    </div>

    <div
      v-else
      class="text-center text-gray-500"
    >
      <p>スタンプの詳細情報がありません</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ stampId: string | undefined }>();
const { getStampById } = useStamps();

const { stampDetail, loading, error } = useStampDetail(props.stampId);

const isCopyingName = ref(false);
const isCopyingWithColon = ref(false);
const toast = useToast();

const copySelectedStampName = () => {
  if (!props.stampId) return;
  isCopyingName.value = true;
  navigator.clipboard.writeText(`${getStampById(props.stampId)?.stamp_name}`);
  toast.add({
    title: 'copied',
  });
  setTimeout(() => {
    isCopyingName.value = false;
  }, 500);
};

const copySelectedStampWithColon = () => {
  if (!props.stampId) return;
  isCopyingWithColon.value = true;
  navigator.clipboard.writeText(`${getStampById(props.stampId)?.stamp_name}`);
  toast.add({
    title: 'copied with colon',
  });
  setTimeout(() => {
    isCopyingWithColon.value = false;
  }, 500);
};

const share = () => {
  if (!props.stampId) return;
  const stamp = getStampById(props.stampId);
  if (!stamp) return;
  const shareData = {
    title: stamp.stamp_name,
    text: `Check out this stamp: ${stamp.stamp_name}`,
    url: `https://stampedia.trap.show/stamps/` + getStampById(props.stampId)?.stamp_name,
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
  }];
</script>
