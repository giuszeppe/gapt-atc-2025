<template>
  <div class="visualizer-wrapper">
    <div v-for="(bar, i) in bars" :key="i" class="bar" :style="{ height: bar + '%' }" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{ volume: number }>();

const BAR_COUNT = 11;
const bars = ref<number[]>(Array(BAR_COUNT).fill(5));

const barWeights = Array.from({ length: BAR_COUNT }, (_, i) => {
  const center = (BAR_COUNT - 1) / 2;
  const distance = Math.abs(i - center);
  return 1 - distance / center * 0.6; 
});

watch(() => props.volume, (vol) => {
  const minHeight = 4;
  const maxHeight = 100;

  const effectiveVolume = vol > 0.15 ? vol : 0;

  bars.value = barWeights.map(weight => {
    const scaledHeight = minHeight + (maxHeight - minHeight) * effectiveVolume * weight;
    return Math.max(minHeight, scaledHeight);
  });
});
</script>

<style scoped>
.visualizer-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
  width: 100%;
  gap: 8px;
}

.bar {
  width: 6px;
  background: yellow;
  border-radius: 3px;
  transition: height 0.12s ease;
}
</style>
