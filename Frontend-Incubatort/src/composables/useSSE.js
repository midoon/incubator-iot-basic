import { onMounted, onUnmounted } from 'vue'
import { useIncubatorStore } from '../stores/incubatorStore.js'

/**
 * Composable yang mengelola koneksi SSE
 * sesuai lifecycle komponen Vue.
 */
export function useSSE() {
  const store = useIncubatorStore()

  onMounted(() => {
    store.startRealtime()
  })

  onUnmounted(() => {
    store.stopRealtime()
  })
}
