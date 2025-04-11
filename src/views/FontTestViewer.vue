<template>
  <div class="font-test-page">
    <h1>字体渲染测试</h1>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-else class="markdown-content" v-html="renderedContent"></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { marked } from 'marked'; // Assuming marked is installed

const renderedContent = ref('');
const loading = ref(true);
const error = ref(null);

const loadAndRenderMarkdown = async () => {
  loading.value = true;
  error.value = null;
  try {
    // Fetch the content of font-test.md (relative to the public directory or root)
    // Adjust the path if font-test.md is not in the root or public folder
    const response = await fetch('/font-test.md'); 
    if (!response.ok) {
      throw new Error(`无法加载 font-test.md: ${response.statusText}`);
    }
    const markdown = await response.text();
    renderedContent.value = marked(markdown);
  } catch (e) {
    error.value = `渲染Markdown失败: ${e.message}`;
    console.error('Error loading/rendering font-test.md:', e);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadAndRenderMarkdown();
});
</script>

<style scoped>
.font-test-page {
  max-width: 900px;
  margin: 2rem auto;
  padding: 1rem;
}

.markdown-content :deep(h3) { /* Add styles to match article rendering if needed */
  margin-top: 2em;
  margin-bottom: 1em;
  border-bottom: 1px solid #eee;
  padding-bottom: 0.3em;
}

.markdown-content :deep(p) {
  line-height: 1.7; /* Ensure consistent line height */
  margin-bottom: 1em;
}

.markdown-content :deep(hr) {
  margin: 2rem 0;
  border: none;
  border-top: 1px dashed #ccc;
}

.loading, .error {
  text-align: center;
  margin-top: 2rem;
  font-size: 1.2em;
}

.error {
  color: red;
}
</style> 