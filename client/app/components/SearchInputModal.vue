<template>
  <UContainer>
    <p>説明文から検索</p>
    <UInput
      v-model="searchDescriptionValue"
      icon="i-lucide-search"
      placeholder="説明文から検索"
      class="mb-5 w-full"
    />
    <p>スタンプ名から検索</p>
    <UInput
      v-model="searchStampNameValue"
      icon="i-lucide-search"
      placeholder="スタンプ名から検索"
      class="mb-5 w-full"
    />
    <UFormField
      label="タグ名を指定"
      required
      icon="i-lucide-search"
    >
      <UInputTags
        v-model="searchStampTagValue"
        placeholder="タグ名を入力してください"
        class="mb-5 w-full"
      />
    </UFormField>
    <p>作成日を指定</p>
    <UPopover>
      <UButton
        color="neutral"
        variant="subtle"
        icon="i-lucide-calendar"
      >
        <template v-if="searchDateValue.start">
          <template v-if="searchDateValue.end">
            {{ df.format(searchDateValue.start.toDate(getLocalTimeZone())) }} - {{ df.format(searchDateValue.end.toDate(getLocalTimeZone())) }}
          </template>

          <template v-else>
            {{ df.format(searchDateValue.start.toDate(getLocalTimeZone())) }}
          </template>
        </template>
        <template v-else>
          日付を選択
        </template>
      </UButton>

      <template #content>
        <UCalendar
          v-model="searchDateValue"
          class="p-2"
          :number-of-months="1"
          range
        />
      </template>
    </UPopover>
    <div class="flex justify-between mt-5">
      <UButton
        label="リセット"
        color="neutral"
        variant="outline"
        class="hover:cursor-pointer w-1/4 grid place-items-center"
      />
      <UButton
        label="検索"
        class="hover:cursor-pointer w-1/4 grid place-items-center"
      />
    </div>
  </UContainer>
</template>

<script setup lang="ts">
import { DateFormatter, getLocalTimeZone, today } from '@internationalized/date';

const searchDescriptionValue = ref('');
const searchStampNameValue = ref('');
const searchStampTagValue = ref<string[]>([]);

const df = new DateFormatter('jp-JP', {
  dateStyle: 'medium',
});

const now = today('Asia/Tokyo');
const searchDateValue = shallowRef({
  start: now.add({ months: -1 }),
  end: now.add({ days: -1 }),
});
</script>
