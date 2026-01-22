import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react-swc";
import path from "node:path";
import { defineConfig } from "vite";

// Backend URL: Use 'backend' hostname in Docker, 'localhost' for local dev
const backendUrl = process.env.BACKEND_URL || "http://localhost:8080";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
  ],
  server: {
    host: true,
    port: 5173,
    proxy: {
      "/auth": backendUrl,
      "/health": backendUrl,
      "/api": backendUrl,
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
