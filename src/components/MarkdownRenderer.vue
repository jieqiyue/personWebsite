<template>
  <div class="markdown-content" v-html="renderedContent"></div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import MarkdownIt from 'markdown-it'

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
  color: var(--text);
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
  background: var(--surface);
  border-radius: 12px;
  box-shadow: 0 2px 12px var(--shadow);
}

.markdown-content img {
  max-width: 100%;
  border-radius: 8px;
  margin: 0rem 0;
}

.markdown-content h1 {
  font-size: 2.5rem;
  border-bottom: 2px solid var(--border);
  padding-bottom: 0.5rem;
  margin-bottom: 1.5rem;
  color: var(--primary);
}

.markdown-content h2 {
  font-size: 1.8rem;
  margin-top: 2rem;
  margin-bottom: 1rem;
  color: var(--primary);
}

.markdown-content h3 {
  font-size: 1.4rem;
  margin-top: 1.5rem;
  margin-bottom: 0.8rem;
  color: var(--secondary);
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
  border-left: 4px solid var(--border);
  padding-left: 1rem;
  margin: 1rem 0;
  color: var(--text);
  opacity: 0.8;
  font-style: italic;
}

.markdown-content hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 2rem 0;
}

.markdown-content code {
  background: var(--accent);
  padding: 0.2rem 0.4rem;
  border-radius: 3px;
  font-family: monospace;
  font-size: 0.9em;
}

.markdown-content a {
  color: var(--primary);
  text-decoration: none;
  border-bottom: 1px dashed var(--primary);
  transition: border-color 0.2s ease;
}

.markdown-content a:hover {
  border-bottom: 1px solid var(--primary);
}
</style> 