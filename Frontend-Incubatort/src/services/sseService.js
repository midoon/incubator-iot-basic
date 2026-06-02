// const BACKEND_HOST = `${window.location.hostname}:8080`;
// const SSE_URL = `http://${BACKEND_HOST}/api/stream`;

const SSE_URL = `${window.location.origin}/api/stream`;

let eventSource = null;
let reconnectTimer = null;
let onMessageCallback = null;
let onStatusChangeCallback = null;

/**
 * Memulai koneksi SSE ke backend.
 * @param {Function} onMessage - dipanggil setiap ada data baru
 * @param {Function} onStatusChange - dipanggil saat status koneksi berubah
 */
export function connectSSE(onMessage, onStatusChange) {
  onMessageCallback = onMessage;
  onStatusChangeCallback = onStatusChange;

  connect();
}

function connect() {
  if (eventSource) {
    eventSource.close();
  }

  notifyStatus("connecting");

  eventSource = new EventSource(SSE_URL);

  eventSource.onopen = () => {
    console.log("[SSE] Connected");
    notifyStatus("connected");
    clearReconnectTimer();
  };

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (onMessageCallback) {
        onMessageCallback(data);
      }
    } catch (err) {
      console.error("[SSE] Failed to parse message:", err);
    }
  };

  eventSource.onerror = () => {
    console.warn("[SSE] Connection error, will reconnect...");
    notifyStatus("disconnected");
    eventSource.close();
    scheduleReconnect();
  };
}

function scheduleReconnect() {
  clearReconnectTimer();
  reconnectTimer = setTimeout(() => {
    console.log("[SSE] Reconnecting...");
    connect();
  }, 3000); // Reconnect setiap 3 detik
}

function clearReconnectTimer() {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
}

function notifyStatus(status) {
  if (onStatusChangeCallback) {
    onStatusChangeCallback(status);
  }
}

/**
 * Menutup koneksi SSE dan membersihkan semua timer.
 */
export function disconnectSSE() {
  clearReconnectTimer();
  if (eventSource) {
    eventSource.close();
    eventSource = null;
  }
  notifyStatus("disconnected");
}
