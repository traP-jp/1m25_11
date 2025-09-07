<template>
  <div>
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
        onclick="refetch()"
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
        @click="navigateTo('/')"
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
        @click="navigateTo('/')"
      >
        トップページに戻る
      </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { paths } from '~~/shared/types/generated';

// URLクエリパラメータの型定義
interface SearchPageQuery {
  q?: string;
}

// ページ設定
definePageMeta({
  title: '検索結果',
});

// ルートとクエリパラメータ
const route = useRoute();

// 検索クエリを取得
const searchQuery = computed(() => route.query.q as string || '');

// openapi-typescript の型を活用した検索パラメータ
type SearchParams = paths['/stamps/search']['get']['parameters']['query'];

const searchParams = computed<SearchParams>(() => {
  const query = searchQuery.value.trim();
  return query ? { q: query } : {};
});

// 検索実行
const { searchResult, stamps, loading, error, refetch } = useStampSearch(searchParams);

// デバッグ用のログ出力（[stampName].vueと同様）
watchEffect(() => {
  console.log('Debug - Search query:', searchQuery.value);
  console.log('Debug - Loading:', loading.value);
  console.log('Debug - Error:', error.value);
  console.log('Debug - Stamps length:', stamps.value.length);
});

// [stampName].vueと同様の状態管理
const isLoading = computed(() => {
  // 検索クエリがある場合のみローディング表示
  return searchQuery.value && loading.value;
});

const hasError = computed(() => {
  // 検索クエリがある場合のみエラー表示
  return searchQuery.value && error.value;
});
</script>
