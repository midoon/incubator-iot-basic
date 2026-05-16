<template>
  <div class="flex items-center gap-2 text-sm">
    <!-- Dot indikator -->
    <span
      class="inline-block w-2 h-2 rounded-full"
      :class="dotClass"
    />
    <span :class="textClass">{{ label }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    default: 'connecting',
    validator: (v) => ['connecting', 'connected', 'disconnected'].includes(v),
  },
})

const dotClass = computed(() => ({
  'bg-green-500 pulse-dot': props.status === 'connected',
  'bg-yellow-400 pulse-dot': props.status === 'connecting',
  'bg-red-400': props.status === 'disconnected',
}))

const textClass = computed(() => ({
  'text-green-600': props.status === 'connected',
  'text-yellow-600': props.status === 'connecting',
  'text-red-500': props.status === 'disconnected',
}))

const label = computed(() => {
  const labels = {
    connected: 'Terhubung',
    connecting: 'Menghubungkan...',
    disconnected: 'Terputus',
  }
  return labels[props.status] ?? '-'
})
</script>
