<template>
  <div class="markdown-content" v-html="renderedContent"></div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import MarkdownIt from 'markdown-it'
import { defineProps } from 'vue'

const props = defineProps({
  source: {
    type: String,
    required: true
  }
})

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

const renderedContent = ref('')

const renderMarkdown = async () => {
  try {
    const response = await fetch(props.source)
    const text = await response.text()
    renderedContent.value = md.render(text)
  } catch (error) {
    console.error('加载Markdown文件失败:', error)
    renderedContent.value = '<p>内容加载失败</p>'
  }
}

onMounted(() => {
  renderMarkdown()
})

watch(() => props.source, () => {
  renderMarkdown()
})
</script>

<style>
.markdown-content {
  text-align: left;
  line-height: 1.6;
  color: #333;
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.markdown-content img {
  max-width: 100%;
  border-radius: 8px;
  margin: 1rem 0;
}

.markdown-content h1 {
  font-size: 2.5rem;
  border-bottom: 2px solid #f0f0f0;
  padding-bottom: 0.5rem;
  margin-bottom: 1.5rem;
  color: #2c6e49;
}

.markdown-content h2 {
  font-size: 1.8rem;
  margin-top: 2rem;
  margin-bottom: 1rem;
  color: #2c6e49;
}

.markdown-content h3 {
  font-size: 1.4rem;
  margin-top: 1.5rem;
  margin-bottom: 0.8rem;
}

.markdown-content p {
  margin-bottom: 1rem;
}

.markdown-content ul, .markdown-content ol {
  margin: 1rem 0;
  padding-left: 2rem;
}

.markdown-content li {
  margin-bottom: 0.5rem;
}

.markdown-content blockquote {
  border-left: 4px solid #daebd9;
  padding-left: 1rem;
  margin: 1rem 0;
  color: #666;
  font-style: italic;
}

.markdown-content hr {
  border: none;
  border-top: 1px solid #eaeaea;
  margin: 2rem 0;
}

.markdown-content code {
  background: #f5f5f5;
  padding: 0.2rem 0.4rem;
  border-radius: 3px;
  font-family: monospace;
  font-size: 0.9em;
}

.markdown-content a {
  color: #2c6e49;
  text-decoration: none;
}

.markdown-content a:hover {
  text-decoration: underline;
}
</style> 