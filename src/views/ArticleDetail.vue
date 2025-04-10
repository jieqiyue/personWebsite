<template>
  <div class="article-detail" v-if="article">
    <div class="article-header">
      <h1>{{ article.title }}</h1>
      <div class="article-meta">
        <span class="date">{{ formatDate(article.date) }}</span>
        <div class="tags">
          <span v-for="tag in article.tags" :key="tag" class="tag">
            #{{ tag }}
          </span>
        </div>
      </div>
    </div>
    <div class="article-cover" v-if="article.coverImage">
      <img :src="article.coverImage" :alt="article.title">
    </div>
    <div class="article-content">
      <MarkdownRenderer :content="markdownContent" :loading="loading" />
    </div>
  </div>
  <div v-else-if="loading" class="loading">
    正在加载文章...
  </div>
  <div v-else-if="error" class="error">
    {{ error }}
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'

const route = useRoute()
const article = ref(null)
const loading = ref(false)
const error = ref(null)
const markdownContent = ref('')

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 加载文章内容和元数据
const loadArticle = async (id) => {
  try {
    loading.value = true
    error.value = null
    
    // 获取文章元数据
    const indexResponse = await fetch('/markdown/articles/index.json')
    const indexData = await indexResponse.json()
    const articleMeta = indexData.articles.find(a => a.id === id)
    
    if (!articleMeta) {
      throw new Error('文章不存在')
    }

    // 获取文章内容
    const contentResponse = await fetch(`/markdown/articles/${id}.md`)
    const content = await contentResponse.text()
    
    // 提取文章内容，设置元数据
    article.value = {
      ...articleMeta
    }
    
    // 如果内容以---开头，说明存在YAML前置元数据
    let processedContent = content
    if (content.trim().startsWith('---')) {
      // 寻找第二个---分隔符，YAML前置元数据在两个---之间
      const secondSeparatorIndex = content.indexOf('---', 3)
      if (secondSeparatorIndex !== -1) {
        // 提取---之后的实际内容
        processedContent = content.substring(secondSeparatorIndex + 3).trim()
      }
    }
    
    // 保存处理后的Markdown内容
    markdownContent.value = processedContent
  } catch (err) {
    error.value = '获取文章内容失败，请稍后重试'
    console.error('加载文章失败:', err)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const articleId = route.params.id
  if (articleId) {
    await loadArticle(articleId)
  }
})

// 监听路由参数变化
watch(() => route.params.id, (newId) => {
  if (newId) {
    loadArticle(newId)
  }
})
</script>

<style scoped>
.article-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 0;
}

.article-header {
  margin-bottom: 2rem;
}

.article-meta {
  display: flex;
  gap: 1rem;
  color: #666;
  margin-top: 1rem;
}

.tags {
  display: flex;
  gap: 0.5rem;
}

.tag {
  color: #42b883;
}

.article-cover {
  margin-bottom: 2rem;
  border-radius: 12px;
  overflow: hidden;
}

.article-cover img {
  width: 100%;
  height: auto;
}

.article-content {
  line-height: 1.8;
}

.article-content :deep(h1),
.article-content :deep(h2) {
  margin-top: 2rem;
  margin-bottom: 1rem;
}

.article-content :deep(p) {
  margin-bottom: 1rem;
}

.article-content :deep(ul),
.article-content :deep(ol) {
  margin-bottom: 1rem;
  padding-left: 2rem;
}

.article-content :deep(li) {
  margin-bottom: 0.5rem;
}
</style> 