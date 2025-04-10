<template>
  <div class="markdown-body" v-html="renderedContent"></div>
</template>

<script>
import { computed, ref, onMounted, watch } from 'vue'
import md from '../plugins/markdown.js'
import 'highlight.js/styles/github.css'
import 'github-markdown-css/github-markdown.css'
import './markdown-custom.css'

export default {
  name: 'MarkdownRenderer',
  props: {
    content: {
      type: String,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const renderedContent = computed(() => {
      if (props.loading) {
        return '<div class="loading">正在加载...</div>'
      }
      
      if (!props.content) {
        return '<div class="empty">暂无内容</div>'
      }
      
      try {
        return md.render(props.content)
      } catch (error) {
        console.error('渲染Markdown时出错:', error)
        return `<div class="error">渲染内容时出错: ${error.message}</div>`
      }
    })

    return {
      renderedContent
    }
  }
}
</script>

<style>
.markdown-body {
  box-sizing: border-box;
  min-width: 200px;
  max-width: 980px;
  margin: 0 auto;
  padding: 45px;
}

@media (max-width: 767px) {
  .markdown-body {
    padding: 15px;
  }
}

.loading,
.empty,
.error {
  padding: 20px;
  text-align: center;
  border-radius: 4px;
  margin: 20px 0;
}

.loading {
  background-color: #f0f8ff;
  color: #1e90ff;
}

.empty {
  background-color: #f5f5f5;
  color: #666;
}

.error {
  background-color: #fff0f0;
  color: #ff3333;
  border: 1px solid #ffdddd;
}
</style> 