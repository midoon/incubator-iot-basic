const BASE_URL = `http://${window.location.hostname}:8080/api`;

/**
 * Mengambil state terkini dari backend (untuk initial load).
 */
export async function fetchState() {
  const res = await fetch(`${BASE_URL}/state`);
  if (!res.ok) throw new Error("Failed to fetch state");
  return res.json();
}

/**
 * Mengirim perintah ganti mode ke backend.
 * @param {'auto'|'manual'} mode
 */
export async function setMode(mode) {
  const res = await fetch(`${BASE_URL}/mode`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ mode }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Failed to set mode");
  }
  return res.json();
}

/**
 * Mengirim perintah kontrol lampu ke backend.
 * @param {boolean} lamp - true = ON, false = OFF
 */
export async function setLamp(lamp) {
  const res = await fetch(`${BASE_URL}/lamp`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ lamp }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Failed to set lamp");
  }
  return res.json();
}
