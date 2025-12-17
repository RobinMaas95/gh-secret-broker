import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [sveltekit(), tailwindcss()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:4000',
        changeOrigin: true,
        secure: false
      },
      '/auth': {
        target: 'http://localhost:4000',
        changeOrigin: true,
        secure: false
      },
      '/logout': {
        target: 'http://localhost:4000',
        changeOrigin: true,
        secure: false
      }
    }
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    // Setup files might need adjustment later
    setupFiles: ['./src/setupTests.ts'],
    environment: 'jsdom',
    restoreMocks: true
  }
});
