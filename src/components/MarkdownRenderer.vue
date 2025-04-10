<template>
  <div class="markdown-content" v-html="renderedContent"></div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import MarkdownIt from 'markdown-it'
import { full as emoji } from 'markdown-it-emoji'

const props = defineProps({
  source: {
    type: String,
    required: true
  }
})

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true
})

// 使用emoji插件，配置自定义表情符号快捷方式
md.use(emoji, {
  shortcuts: {
    // 添加一些常用表情符号的快捷方式
    'smile': [':)', ':-)', '😊'],
    'laughing': [':D', ':-D', '😄'],
    'wink': [';)', ';-)', '😉'],
    'frown': [':(', ':-(', '😞'],
    'stuck_out_tongue': [':P', ':-P', '😛'],
    'confused': [':/', ':-/', '😕'],
    'open_mouth': [':O', ':-O', '😮'],
    'heart': ['<3', '❤️']
  }
})

const renderedContent = ref('')

const renderMarkdown = async () => {
  try {
    const response = await fetch(props.source)
    if (!response.ok) {
      throw new Error(`无法加载文件: ${response.status} ${response.statusText}`);
    }
    
    const text = await response.text()
    renderedContent.value = md.render(text)
  } catch (error) {
    console.error('加载Markdown文件失败:', error)
    renderedContent.value = `<div class="markdown-error">
      <h3>内容加载失败</h3>
      <p>无法加载文件: ${props.source}</p>
      <p>错误: ${error.message}</p>
    </div>`;
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
  margin: 0.5rem 0;
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

.markdown-content pre {
  background: #f5f5f5;
  border-radius: 6px;
  padding: 1rem;
  overflow-x: auto;
  margin: 1.5rem 0;
}

.markdown-content pre code {
  background: transparent;
  padding: 0;
}

/* Emoji样式 */
.emoji {
  height: 1.2em;
  vertical-align: middle;
}

/* 表格样式 */
.markdown-content table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5rem 0;
  overflow-x: auto;
  display: block;
}

.markdown-content table th {
  background-color: var(--surface);
  color: var(--primary);
  font-weight: bold;
  padding: 10px 15px;
  border: 1px solid var(--border);
  text-align: left;
}

.markdown-content table td {
  padding: 8px 15px;
  border: 1px solid var(--border);
}

.markdown-content table tr:nth-child(even) {
  background-color: var(--accent);
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

.markdown-error {
  padding: 1rem;
  background: #fdd;
  border-left: 4px solid #f66;
  color: #700;
  border-radius: 4px;
}
</style> 