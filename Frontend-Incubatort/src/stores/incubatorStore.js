import { defineStore } from 'pinia'
import { ref } from 'vue'
import { setMode, setLamp } from '../services/apiService.js'
import { connectSSE, disconnectSSE } from '../services/sseService.js'

export const useIncubatorStore = defineStore('incubator', () => {
  // State data sensor dan kontrol
  const temperature = ref(0)
  const humidity = ref(0)
  const mode = ref('auto')
  const lamp = ref(false)

  // State UI
  const connectionStatus = ref('connecting') // 'connecting' | 'connected' | 'disconnected'
  const isLoadingMode = ref(false)
  const isLoadingLamp = ref(false)
  const errorMessage = ref('')

  // Timestamp update terakhir
  const lastUpdated = ref(null)

  /**
   * Memulai koneksi SSE dan mulai menerima data realtime.
   */
  function startRealtime() {
    connectSSE(
      // Callback saat data baru diterima dari SSE
      (data) => {
        temperature.value = data.temperature ?? temperature.value
        humidity.value = data.humidity ?? humidity.value
        mode.value = data.mode ?? mode.value
        lamp.value = data.lamp ?? lamp.value
        lastUpdated.value = new Date()
        errorMessage.value = ''
      },
      // Callback saat status koneksi berubah
      (status) => {
        connectionStatus.value = status
      }
    )
  }

  /**
   * Menutup koneksi SSE.
   */
  function stopRealtime() {
    disconnectSSE()
  }

  /**
   * Mengganti mode mesin (auto/manual).
   */
  async function changeMode(newMode) {
    if (isLoadingMode.value) return
    isLoadingMode.value = true
    errorMessage.value = ''

    try {
      await setMode(newMode)
      // State akan diupdate via SSE broadcast dari backend,
      // tapi update lokal dulu agar UI responsif.
      mode.value = newMode
      if (newMode === 'auto') {
        lamp.value = false
      }
    } catch (err) {
      errorMessage.value = err.message
    } finally {
      isLoadingMode.value = false
    }
  }

  /**
   * Mengontrol lampu pemanas (hanya saat mode manual).
   */
  async function controlLamp(on) {
    if (isLoadingLamp.value) return
    isLoadingLamp.value = true
    errorMessage.value = ''

    try {
      await setLamp(on)
      lamp.value = on
    } catch (err) {
      errorMessage.value = err.message
    } finally {
      isLoadingLamp.value = false
    }
  }

  return {
    // State
    temperature,
    humidity,
    mode,
    lamp,
    connectionStatus,
    isLoadingMode,
    isLoadingLamp,
    errorMessage,
    lastUpdated,
    // Actions
    startRealtime,
    stopRealtime,
    changeMode,
    controlLamp,
  }
})
