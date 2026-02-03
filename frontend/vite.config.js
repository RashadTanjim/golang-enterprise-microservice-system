import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api/v1/users': {
        target: 'http://localhost:8081',
        changeOrigin: true
      },
      '/api/v1/auth': {
        target: 'http://localhost:8081',
        changeOrigin: true
      },
      '/api/v1/orders': {
        target: 'http://localhost:8082',
        changeOrigin: true
      },
      '/health/user': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        rewrite: () => '/health'
      },
      '/health/order': {
        target: 'http://localhost:8082',
        changeOrigin: true,
        rewrite: () => '/health'
      },
      '/metrics/user': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        rewrite: () => '/metrics'
      },
      '/metrics/order': {
        target: 'http://localhost:8082',
        changeOrigin: true,
        rewrite: () => '/metrics'
      }
    }
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/tests/setup.js'
  }
});
