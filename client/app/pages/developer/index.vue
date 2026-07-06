<template>
  <iframe
    ref="iframeRef"
    src="/developer.html"
    class="w-full border-0 block"
    :style="{ height: iframeHeight }"
    @load="onLoad"
  />
</template>

<script setup lang="ts">
definePageMeta({ title: 'Swagger Viewer' });

const iframeRef = ref<HTMLIFrameElement | null>(null);
const iframeHeight = ref('100vh');
let observer: ResizeObserver | null = null;

function onLoad() {
  const iframe = iframeRef.value;
  const doc = iframe?.contentDocument;
  if (!doc?.body) return;

  // iframe 自身のスクロールバーを無効化して、外側ページのみでスクロールさせる
  const style = doc.createElement('style');
  style.textContent = 'html, body { overflow: hidden !important; height: auto !important; }';
  doc.head.appendChild(style);

  observer?.disconnect();
  observer = new ResizeObserver(() => {
    const height = doc.body.offsetHeight;
    if (height) iframeHeight.value = `${height}px`;
  });
  observer.observe(doc.body);
}

onUnmounted(() => observer?.disconnect());
</script>
