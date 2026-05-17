import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";

import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools(), tailwindcss()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 5173,
    // Izinkan akses dari perangkat lain di jaringan yang sama (HP, tablet, dll)
    host: true,
    proxy: {
      // Semua request /api diteruskan ke backend Go
      // Gunakan IP server, bukan localhost, agar proxy bekerja dari semua perangkat
      "/api": {
        target: "http://10.10.20.254:8080",
        changeOrigin: true,
        // Penting untuk SSE: matikan buffering agar data mengalir realtime
        configure: (proxy) => {
          proxy.on("proxyReq", (proxyReq) => {
            proxyReq.setHeader("Accept-Encoding", "identity");
          });
        },
      },
    },
  },
});
