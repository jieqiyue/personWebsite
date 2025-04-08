import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    fs: {
      // 允许服务器访问父目录中的文件，这样可以访问src下的所有文件
      allow: ['..']
    }
  }
})
