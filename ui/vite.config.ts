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
    },
    conditions: ['browser'],
  },
  build: {
    outDir: "./dist/"
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts'],
    alias: {
      "$lib": path.resolve(process.cwd(), "./src/lib"),
    },
    server: {
      deps: {
        inline: ['@testing-library/svelte']
      }
    },
  }
})
