<template>
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
        Pick a date
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
</template>

<script setup lang="ts">
import { DateFormatter, getLocalTimeZone, today } from '@internationalized/date';

const df = new DateFormatter('jp-JP', {
  dateStyle: 'medium',
});

const now = today('Asia/Tokyo');
const searchDateValue = shallowRef({
  start: now.add({ months: -1 }),
  end: now.add({ days: -1 }),
});
</script>
