<template>
  <div class="about">
    <h1 class="page-title">关于</h1>
    
    <div class="markdown-container">
      <MarkdownRenderer :content="markdownContent" :loading="loading" />
    </div>
    
    <div class="edit-link">
      <a :href="editUrl" target="_blank" class="edit-button">
        编辑这个页面
      </a>
      <p class="edit-note">编辑后刷新页面即可看到更新的内容</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'

// 状态变量
const markdownContent = ref('')
const loading = ref(true)
const markdownPath = '/markdown/About.md'

// 编辑链接 - 在实际应用中，这可能是指向CMS或GitHub仓库的链接
const editUrl = ref(markdownPath)

// 加载Markdown内容
const loadMarkdown = async () => {
  try {
    const response = await fetch(markdownPath)
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
.about {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.page-title {
  text-align: center;
  margin-bottom: 2rem;
  font-size: 2.5rem;
  color: #2c6e49;
}

.markdown-container {
  margin-bottom: 3rem;
}

.edit-link {
  text-align: center;
  margin-top: 2rem;
  padding: 1rem;
}

.edit-button {
  display: inline-block;
  padding: 0.8rem 1.5rem;
  background-color: #2c6e49;
  color: white;
  border-radius: 30px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
}

.edit-button:hover {
  background-color: #1a4731;
  transform: translateY(-2px);
}

.edit-note {
  margin-top: 1rem;
  color: #666;
  font-size: 0.9rem;
}
</style> 