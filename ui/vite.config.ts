import path from "path";
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "$lib": path.resolve(__dirname, "./src/lib"),
    },
  },
  build: {
    outDir: "./dist/"
  }
})
