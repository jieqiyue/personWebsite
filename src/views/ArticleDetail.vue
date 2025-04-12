<template>
  <div class="article-container" v-if="article">
    <!-- 已在App.vue中添加了进度条，此处移除 -->
    
    <div class="article-detail">
      <div class="article-header">
        <h1>{{ article.title }}</h1>
        <div class="article-meta">
          <span class="date">{{ formatDate(article.date) }}</span>
          <span class="reading-time" v-if="article.readingTime">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
            {{ article.readingTime }}
          </span>
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
      <div class="comments-wrapper">
        <Comments :pageTerm="articlePath" />
      </div>
    </div>
    
    <aside class="article-sidebar">
      <div class="related-articles">
        <h3>相关文章</h3>
        <ul v-if="relatedArticles.length > 0">
          <li v-for="relatedArticle in relatedArticles" :key="relatedArticle.id">
            <router-link :to="'/articles/' + relatedArticle.id">
              {{ relatedArticle.title }}
              <span class="related-tag">{{ getCommonTag(relatedArticle) }}</span>
            </router-link>
          </li>
        </ul>
        <div v-else class="no-related">
          暂无相关文章
        </div>
      </div>
      
      <div class="popular-tags">
        <h3>热门标签</h3>
        <div class="tag-cloud">
          <router-link 
            v-for="tag in popularTags" 
            :key="tag.name" 
            :to="'/tags/' + tag.name"
            class="tag-item"
            :style="{ fontSize: getTagSize(tag.count) }"
          >
            #{{ tag.name }} ({{ tag.count }})
          </router-link>
        </div>
      </div>
    </aside>
    <BackToTop />
  </div>
  <div v-else-if="loading" class="loading">
    正在加载文章...
  </div>
  <div v-else-if="error" class="error">
    {{ error }}
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import { marked } from 'marked'
import Prism from 'prismjs'
import 'prismjs/themes/prism-tomorrow.css'
import BackToTop from '../components/BackToTop.vue'
import Comments from '../components/Comments.vue'
import { estimateReadingTime, formatReadingTime } from '../utils/article'

const route = useRoute()
const router = useRouter()
const article = ref(null)
const loading = ref(false)
const error = ref(null)
const markdownContent = ref('')
const allArticles = ref([])
const popularTags = ref([])
const articlePath = ref('')

// 计算相关文章（共享至少一个标签的文章）
const relatedArticles = computed(() => {
  if (!article.value || !allArticles.value.length) return []
  
  return allArticles.value
    .filter(otherArticle => 
      // 排除当前文章
      otherArticle.id !== article.value.id && 
      // 至少有一个共同标签
      otherArticle.tags.some(tag => article.value.tags.includes(tag))
    )
    .sort((a, b) => new Date(b.date) - new Date(a.date)) // 按日期排序
    .slice(0, 5) // 最多显示5篇
})

// 获取两篇文章的共同标签（用于显示）
const getCommonTag = (relatedArticle) => {
  if (!article.value) return ''
  
  const commonTag = relatedArticle.tags.find(tag => 
    article.value.tags.includes(tag)
  )
  
  return commonTag || ''
}

// 计算标签大小（根据出现次数）
const getTagSize = (count) => {
  const baseSize = 0.9
  const scale = 0.1
  return `${baseSize + scale * Math.min(count, 5)}rem`
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 加载所有文章的元数据
const loadAllArticles = async () => {
  try {
    const indexResponse = await fetch('/markdown/articles/index.json')
    const indexData = await indexResponse.json()
    allArticles.value = indexData.articles
    
    // 计算标签统计
    const tagCounts = {}
    allArticles.value.forEach(article => {
      article.tags.forEach(tag => {
        tagCounts[tag] = (tagCounts[tag] || 0) + 1
      })
    })
    
    // 转换为数组并排序
    popularTags.value = Object.entries(tagCounts)
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count)
    
    // 为所有文章添加阅读时间估计
    if (allArticles.value.length > 0) {
      await addReadingTimeEstimations();
    }
  } catch (err) {
    console.error('加载文章索引失败:', err)
  }
}

// 为所有文章添加阅读时间估计
const addReadingTimeEstimations = async () => {
  await Promise.all(allArticles.value.map(async (article) => {
    try {
      // 如果文章内容已缓存，直接使用缓存计算
      if (article.cachedContent) {
        const minutes = estimateReadingTime(article.cachedContent);
        article.readingTime = formatReadingTime(minutes);
        return;
      }
      
      // 否则获取文章内容
      const response = await fetch(`/markdown/articles/${article.id}.md`);
      const content = await response.text();
      
      // 提取正文内容（去除YAML前缀）
      let processedContent = content;
      if (content.trim().startsWith('---')) {
        const secondSeparatorIndex = content.indexOf('---', 3);
        if (secondSeparatorIndex !== -1) {
          processedContent = content.substring(secondSeparatorIndex + 3).trim();
        }
      }
      
      // 缓存处理后的内容
      article.cachedContent = processedContent;
      
      // 计算阅读时间
      const minutes = estimateReadingTime(processedContent);
      article.readingTime = formatReadingTime(minutes);
    } catch (error) {
      console.error(`无法获取文章 ${article.id} 的内容:`, error);
    }
  }));
}

// 加载文章内容和元数据
const loadArticle = async (id) => {
  try {
    loading.value = true
    error.value = null
    
    // 获取文章元数据
    const articleMeta = allArticles.value.find(a => a.id === id)
    
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
    
    // 缓存处理后的内容
    article.value.cachedContent = processedContent;
    
    // 保存处理后的Markdown内容
    markdownContent.value = processedContent
    
    // 计算阅读时间（如果尚未计算）
    if (!article.value.readingTime) {
      const minutes = estimateReadingTime(processedContent);
      article.value.readingTime = formatReadingTime(minutes);
    }

    // Highlight code blocks after content is rendered
    await nextTick(() => {
      Prism.highlightAll()
    })
  } catch (err) {
    error.value = '获取文章内容失败，请稍后重试'
    console.error('加载文章失败:', err)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  // 先加载所有文章元数据
  await loadAllArticles()
  
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
.article-container {
  display: flex;
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
  gap: 2rem;
}

.article-detail {
  flex: 1;
  max-width: 800px;
}

.article-sidebar {
  width: 300px;
  position: sticky;
  top: 2rem;
  align-self: flex-start;
}

@media (max-width: 1024px) {
  .article-container {
    flex-direction: column;
  }
  
  .article-sidebar {
    width: 100%;
    position: static;
    margin-top: 2rem;
  }
}

.article-header {
  margin-bottom: 2rem;
}

.article-meta {
  display: flex;
  gap: 1rem;
  color: #666;
  margin-top: 1rem;
  flex-wrap: wrap;
  align-items: center;
}

.reading-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--text-light, #666);
}

.reading-time svg {
  color: var(--primary, #42b883);
}

.tags {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
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

/* 侧边栏样式 */
.related-articles,
.popular-tags {
  background-color: #f8f9fa;
  border-radius: 10px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.related-articles h3,
.popular-tags h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #2c3e50;
  font-size: 1.2rem;
  border-bottom: 1px solid #eee;
  padding-bottom: 0.5rem;
}

.related-articles ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.related-articles li {
  padding: 0.5rem 0;
  border-bottom: 1px dashed #eee;
}

.related-articles li:last-child {
  border-bottom: none;
}

.related-articles a {
  color: #2c3e50;
  text-decoration: none;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: color 0.2s;
}

.related-articles a:hover {
  color: #42b883;
}

.related-tag {
  font-size: 0.8rem;
  color: #42b883;
  background-color: rgba(66, 184, 131, 0.1);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}

.no-related {
  color: #999;
  font-style: italic;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag-item {
  color: #42b883;
  text-decoration: none;
  transition: all 0.2s;
}

.tag-item:hover {
  color: #3aa171;
  transform: scale(1.05);
}

/* Wrapper for the comments section */
.comments-wrapper {
  max-width: 800px; /* Align with .article-detail max-width */
  margin: 0 auto;  /* Center the wrapper */
  /* The top margin is handled by Comments.vue's internal .comments-section */
  padding-bottom: 2rem; /* Add some padding at the very bottom */
}

/* 响应式设计 */
@media (max-width: 768px) {
  .article-meta {
    gap: 0.8rem;
    margin-top: 0.8rem;
  }
  
  .article-header {
    margin-bottom: 1.5rem;
  }
  
  .article-cover {
    margin-bottom: 1.5rem;
  }
}

@media (prefers-color-scheme: dark) {
  .related-articles,
  .popular-tags {
    background-color: var(--surface, #2a2a2a);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
  
  .related-articles h3,
  .popular-tags h3 {
    color: var(--text, #e0e0e0);
    border-bottom-color: var(--border, #444);
  }
  
  .related-articles li {
    border-bottom-color: var(--border, #444);
  }
  
  .related-articles a {
    color: var(--text, #e0e0e0);
  }
  
  .related-articles a:hover {
    color: var(--primary, #42b883);
  }
  
  .no-related {
    color: var(--text-light, #888);
  }
}
</style> 