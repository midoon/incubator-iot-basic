<template>
  <div class="min-h-screen bg-warm-50">

    <!-- Header -->
    <header class="bg-white border-b border-warm-100 sticky top-0 z-10">
      <div class="max-w-lg mx-auto px-4 py-3 flex items-center justify-between">
        <div>
          <h1 class="text-base font-semibold text-gray-800">🥚 Mesin Penetasan</h1>
          <p class="text-xs text-gray-400">Monitor & Kontrol Realtime</p>
        </div>
        <ConnectionIndicator :status="store.connectionStatus" />
      </div>
    </header>

    <!-- Main content -->
    <main class="max-w-lg mx-auto px-4 py-5 space-y-4">

      <!-- Error banner -->
      <div
        v-if="store.errorMessage"
        class="bg-red-50 border border-red-200 text-red-600 text-sm rounded-xl px-4 py-3 flex items-start gap-2"
      >
        <span>⚠️</span>
        <span>{{ store.errorMessage }}</span>
      </div>

      <!-- Sensor cards -->
      <section>
        <h2 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">
          Kondisi Ruangan
        </h2>
        <div class="grid grid-cols-2 gap-3">
          <SensorCard
            label="Suhu"
            :value="store.temperature"
            unit="°C"
            icon="🌡️"
            :decimals="1"
            :rangeMin="30"
            :rangeMax="45"
            :idealMin="37"
            :idealMax="38"
          />
          <SensorCard
            label="Kelembapan"
            :value="store.humidity"
            unit="%"
            icon="💧"
            :decimals="0"
            :rangeMin="40"
            :rangeMax="90"
            :idealMin="60"
            :idealMax="70"
          />
        </div>
      </section>

      <!-- Status section -->
      <section class="bg-white rounded-2xl p-5 shadow-sm border border-warm-100">
        <h2 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">
          Status Sistem
        </h2>
        <div class="flex flex-wrap gap-2">
          <StatusBadge :type="store.mode === 'auto' ? 'mode-auto' : 'mode-manual'" />
          <StatusBadge :type="store.lamp ? 'lamp-on' : 'lamp-off'" />
        </div>

        <!-- Last updated -->
        <p v-if="store.lastUpdated" class="text-xs text-gray-400 mt-3">
          Update terakhir: {{ formattedLastUpdated }}
        </p>
      </section>

      <!-- Control panel -->
      <section class="bg-white rounded-2xl p-5 shadow-sm border border-warm-100">
        <h2 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-4">
          Panel Kontrol
        </h2>

        <!-- Mode selector -->
        <div class="mb-5">
          <p class="text-sm font-medium text-gray-700 mb-2">Mode Operasi</p>
          <div class="grid grid-cols-2 gap-2">
            <button
              @click="handleSetMode('auto')"
              :disabled="store.isLoadingMode || store.connectionStatus !== 'connected'"
              :class="[
                'py-2.5 px-4 rounded-xl text-sm font-medium transition-all',
                store.mode === 'auto'
                  ? 'bg-blue-500 text-white shadow-sm'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200',
                (store.isLoadingMode || store.connectionStatus !== 'connected')
                  ? 'opacity-50 cursor-not-allowed'
                  : ''
              ]"
            >
              🤖 Auto
            </button>
            <button
              @click="handleSetMode('manual')"
              :disabled="store.isLoadingMode || store.connectionStatus !== 'connected'"
              :class="[
                'py-2.5 px-4 rounded-xl text-sm font-medium transition-all',
                store.mode === 'manual'
                  ? 'bg-orange-500 text-white shadow-sm'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200',
                (store.isLoadingMode || store.connectionStatus !== 'connected')
                  ? 'opacity-50 cursor-not-allowed'
                  : ''
              ]"
            >
              🖐️ Manual
            </button>
          </div>
        </div>

        <!-- Lamp control — hanya muncul saat mode manual -->
        <transition name="fade">
          <div v-if="store.mode === 'manual'">
            <p class="text-sm font-medium text-gray-700 mb-2">Kontrol Lampu Pemanas</p>
            <div class="grid grid-cols-2 gap-2">
              <button
                @click="handleSetLamp(true)"
                :disabled="store.isLoadingLamp || store.lamp"
                :class="[
                  'py-2.5 px-4 rounded-xl text-sm font-medium transition-all',
                  store.lamp
                    ? 'bg-yellow-400 text-yellow-900 shadow-sm'
                    : 'bg-gray-100 text-gray-600 hover:bg-yellow-100',
                  (store.isLoadingLamp || store.lamp)
                    ? 'opacity-60 cursor-not-allowed'
                    : ''
                ]"
              >
                💡 Nyalakan
              </button>
              <button
                @click="handleSetLamp(false)"
                :disabled="store.isLoadingLamp || !store.lamp"
                :class="[
                  'py-2.5 px-4 rounded-xl text-sm font-medium transition-all',
                  !store.lamp
                    ? 'bg-gray-300 text-gray-700 shadow-sm'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200',
                  (store.isLoadingLamp || !store.lamp)
                    ? 'opacity-60 cursor-not-allowed'
                    : ''
                ]"
              >
                🌑 Matikan
              </button>
            </div>
          </div>
        </transition>

        <!-- Hint saat mode auto -->
        <transition name="fade">
          <p v-if="store.mode === 'auto'" class="text-xs text-gray-400 mt-1">
            Beralih ke mode manual untuk mengontrol lampu secara langsung.
          </p>
        </transition>
      </section>

      <!-- Footer info -->
      <p class="text-center text-xs text-gray-300 pb-4">
        ESP32 + MQTT + Go + Vue.js
      </p>

    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useIncubatorStore } from '../stores/incubatorStore.js'
import { useSSE } from '../composables/useSSE.js'
import ConnectionIndicator from '../components/ConnectionIndicator.vue'
import SensorCard from '../components/SensorCard.vue'
import StatusBadge from '../components/StatusBadge.vue'

const store = useIncubatorStore()

// Mulai SSE sesuai lifecycle komponen
useSSE()

const formattedLastUpdated = computed(() => {
  if (!store.lastUpdated) return '-'
  return store.lastUpdated.toLocaleTimeString('id-ID')
})

async function handleSetMode(mode) {
  if (store.mode === mode) return
  await store.changeMode(mode)
}

async function handleSetLamp(on) {
  if (store.lamp === on) return
  await store.controlLamp(on)
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
