<template>
  <div class="bg-white rounded-2xl p-5 shadow-sm border border-warm-100">
    <!-- Icon dan label -->
    <div class="flex items-center gap-3 mb-3">
      <span class="text-2xl">{{ icon }}</span>
      <span class="text-sm font-medium text-gray-500">{{ label }}</span>
    </div>

    <!-- Nilai utama -->
    <div class="flex items-end gap-1">
      <span class="text-4xl font-semibold text-gray-800 tabular-nums leading-none">
        {{ formattedValue }}
      </span>
      <span class="text-lg text-gray-400 mb-0.5">{{ unit }}</span>
    </div>

    <!-- Range indikator -->
    <div v-if="showRange" class="mt-3">
      <div class="flex justify-between text-xs text-gray-400 mb-1">
        <span>{{ rangeMin }}{{ unit }}</span>
        <span>{{ rangeMax }}{{ unit }}</span>
      </div>
      <div class="w-full bg-warm-100 rounded-full h-1.5">
        <div
          class="h-1.5 rounded-full transition-all duration-500"
          :class="barColor"
          :style="{ width: barWidth }"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: Number, default: 0 },
  unit: { type: String, default: '' },
  icon: { type: String, default: '📊' },
  decimals: { type: Number, default: 1 },
  rangeMin: { type: Number, default: 0 },
  rangeMax: { type: Number, default: 100 },
  // Nilai ideal untuk warna bar
  idealMin: { type: Number, default: null },
  idealMax: { type: Number, default: null },
  showRange: { type: Boolean, default: true },
})

const formattedValue = computed(() =>
  props.value.toFixed(props.decimals)
)

const barWidth = computed(() => {
  const pct = ((props.value - props.rangeMin) / (props.rangeMax - props.rangeMin)) * 100
  return `${Math.min(Math.max(pct, 0), 100)}%`
})

// Warna bar: hijau jika dalam range ideal, kuning jika mendekati, merah jika di luar
const barColor = computed(() => {
  if (props.idealMin === null || props.idealMax === null) {
    return 'bg-warm-400'
  }
  const v = props.value
  if (v >= props.idealMin && v <= props.idealMax) return 'bg-green-400'
  if (v >= props.idealMin - 2 && v <= props.idealMax + 2) return 'bg-yellow-400'
  return 'bg-red-400'
})
</script>
