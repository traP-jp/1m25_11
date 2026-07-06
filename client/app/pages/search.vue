<template>
  <div>
    <SearchInput />
    <div class="mt-6">
      <!-- 検索条件表示 -->
      <div
        v-if="searchQuery"
        class="mb-6 p-4 bg-gray-50 rounded-lg"
      >
        <h2 class="text-lg font-semibold mb-2">
          検索結果
        </h2>
        <p class="text-gray-600">
          「{{ searchQuery }}」の検索結果
          <span
            v-if="!isLoading && stamps.length > 0"
            class="ml-2"
          >
            ({{ stamps.length }}件)
          </span>
        </p>
      </div>

      <!-- ローディング状態 -->
      <div
        v-if="isLoading"
        class="text-center p-8"
      >
        <p>検索中...</p>
      </div>

      <!-- エラー状態 -->
      <div
        v-else-if="hasError"
        class="text-center p-8 text-red-500"
      >
        <p>検索に失敗しました</p>
        <UButton
          class="mt-4"
          variant="outline"
          @click="() => refetch()"
        >
          再試行
        </UButton>
      </div>

      <div v-else-if="searchQuery && stamps.length > 0">
        <StampGrid :stamps="stamps" />
      </div>

      <div
        v-else-if="searchQuery && stamps.length === 0"
        class="text-center p-8 text-gray-500"
      >
        <p>「{{ searchQuery }}」に一致するスタンプが見つかりませんでした</p>
        <UButton
          class="mt-4"
          variant="outline"
          to="/"
        >
          トップページに戻る
        </UButton>
      </div>

      <!-- 検索条件が未入力 -->
      <div
        v-else
        class="text-center p-8 text-gray-500"
      >
        <p>検索キーワードを入力してください</p>
        <UButton
          class="mt-4"
          variant="outline"
          to="/"
        >
          トップページに戻る
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { paths } from '~~/shared/types/generated';

definePageMeta({
  title: '検索結果',
});

const route = useRoute();

const searchQuery = computed(() => route.query.q as string || '');

type SearchParams = paths['/stamps/search']['get']['parameters']['query'];

const searchParams = computed<SearchParams>(() => {
  const query = searchQuery.value.trim();
  return query ? { q: query } : {};
});

const { stamps, loading, error, refetch } = useStampSearch(searchParams);

const isLoading = computed(() => searchQuery.value && loading.value);
const hasError = computed(() => searchQuery.value && error.value);
</script>
