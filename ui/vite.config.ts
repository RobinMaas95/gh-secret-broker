/// <reference types="vitest" />
import path from "path";
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

import tailwindcss from '@tailwindcss/vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(process.cwd(), "./src"),
      "$lib": path.resolve(process.cwd(), "./src/lib"),
      "$app/environment": path.resolve(process.cwd(), "./src/test/sveltekit_mocks.ts"),
      "$app/stores": path.resolve(process.cwd(), "./src/test/sveltekit_mocks.ts"),
      "$app/navigation": path.resolve(process.cwd(), "./src/test/sveltekit_mocks.ts"),
      "$app/forms": path.resolve(process.cwd(), "./src/test/sveltekit_mocks.ts"),
    },
    conditions: ['browser'],
  },
  build: {
    outDir: "./dist/"
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts', './src/test/setup_kit.ts'],
    globals: true,
    restoreMocks: true
  }
})
