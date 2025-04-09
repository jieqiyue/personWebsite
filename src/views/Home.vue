<template>
  <div class="home">
    <!-- Markdown渲染区域 -->
    <div class="markdown-container">
      <MarkdownRenderer :source="markdownSource" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'

// Markdown文件路径
const markdownSource = ref('/src/assets/markdown/Home.md')
</script>

<style scoped>
.home {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
}

.markdown-container {
  width: 100%;
}
</style> 