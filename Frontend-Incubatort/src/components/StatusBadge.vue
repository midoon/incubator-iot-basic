<template>
  <span
    class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium"
    :class="badgeClass"
  >
    <span>{{ icon }}</span>
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    required: true,
    // 'mode-auto' | 'mode-manual' | 'lamp-on' | 'lamp-off'
  },
})

const config = {
  'mode-auto':   { label: 'Mode Auto',   icon: '🤖', class: 'bg-blue-100 text-blue-700' },
  'mode-manual': { label: 'Mode Manual', icon: '🖐️', class: 'bg-orange-100 text-orange-700' },
  'lamp-on':     { label: 'Lampu Nyala', icon: '💡', class: 'bg-yellow-100 text-yellow-700' },
  'lamp-off':    { label: 'Lampu Mati',  icon: '🌑', class: 'bg-gray-100 text-gray-500' },
}

const current = computed(() => config[props.type] ?? config['lamp-off'])
const label = computed(() => current.value.label)
const icon = computed(() => current.value.icon)
const badgeClass = computed(() => current.value.class)
</script>
