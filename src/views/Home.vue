<template>
  <div class="home">
    <!-- Markdown渲染区域 -->
    <div class="markdown-container">
      <MarkdownRenderer :content="markdownContent" :loading="loading" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'

// 状态变量
const markdownContent = ref('')
const loading = ref(true)

// 加载Markdown内容
const loadMarkdown = async () => {
  try {
    const response = await fetch('/markdown/image-test.md')
    markdownContent.value = await response.text()
  } catch (error) {
    console.error('加载Markdown失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadMarkdown()
})
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