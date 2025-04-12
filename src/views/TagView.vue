<template>
  <div class="tag-view">
    <div class="tag-header">
      <button class="back-button" @click="goBack">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M19 12H5"></path>
          <path d="M12 19l-7-7 7-7"></path>
        </svg>
        <span>返回</span>
      </button>
      
      <h1 class="tag-title">#{{ tag }}</h1>
    </div>
    
    <div class="articles-count">
      找到 {{ taggedArticles.length }} 篇相关文章
    </div>
    
    <div v-if="taggedArticles.length" class="articles-list">
      <div v-for="article in taggedArticles" :key="article.id" class="article-card">
        <router-link :to="'/articles/' + article.id" class="article-link">
          <div class="article-card-content">
            <h2>{{ article.title }}</h2>
            <div class="article-meta">
              <span class="date">{{ formatDate(article.date) }}</span>
              <div class="tags">
                <span v-for="articleTag in article.tags" :key="articleTag" 
                      :class="['tag', { active: articleTag === tag }]">
                  #{{ articleTag }}
                </span>
              </div>
            </div>
            <p class="excerpt">{{ article.excerpt }}</p>
          </div>
          <div v-if="article.coverImage" class="article-image">
            <img :src="article.coverImage" :alt="article.title">
          </div>
        </router-link>
      </div>
    </div>
    
    <div v-else-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>正在加载文章...</p>
    </div>
    
    <div v-else class="no-articles">
      <svg xmlns="http://www.w3.org/2000/svg" width="60" height="60" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 6L6 18M6 6l12 12"></path>
        <circle cx="12" cy="12" r="10"></circle>
      </svg>
      <p>没有找到标签为 <strong>"{{ tag }}"</strong> 的文章</p>
      <button class="browse-all-button" @click="goToBlog">浏览所有文章</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const tag = computed(() => route.params.tag || '')
const articles = ref([])
const loading = ref(false)

// 过滤出包含当前标签的文章
const taggedArticles = computed(() => {
  if (!tag.value) return []
  
  return articles.value
    .filter(article => article.tags.includes(tag.value))
    .sort((a, b) => new Date(b.date) - new Date(a.date)) // 按日期降序排序
})

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 加载所有文章
const loadArticles = async () => {
  try {
    loading.value = true
    const response = await fetch('/markdown/articles/index.json')
    const data = await response.json()
    articles.value = data.articles
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    loading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 前往博客页面
const goToBlog = () => {
  router.push('/blog')
}

onMounted(() => {
  loadArticles()
})

// 监听路由变化
watch(() => route.params.tag, (newTag) => {
  if (newTag && !articles.value.length) {
    loadArticles()
  }
})
</script>

<style scoped>
.tag-view {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.tag-header {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
  gap: 1rem;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0.5rem 1rem;
  background-color: var(--accent, #f5f5f5);
  border: none;
  border-radius: 8px;
  color: var(--text, #333);
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-button:hover {
  background-color: var(--primary, #42b883);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.tag-title {
  font-size: 2.5rem;
  color: var(--primary, #42b883);
  margin: 0;
}

.articles-count {
  color: var(--text-light, #666);
  margin-bottom: 2rem;
  font-size: 1.1rem;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.article-card {
  background: var(--surface, #fff);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s, box-shadow 0.3s;
}

.article-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.08);
}

.article-link {
  display: flex;
  color: inherit;
  text-decoration: none;
}

.article-card-content {
  flex: 1;
  padding: 1.5rem;
}

.article-card h2 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.5rem;
  color: var(--text, #2c3e50);
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.date {
  color: var(--text-light, #666);
  font-size: 0.9rem;
}

.tags {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.tag {
  font-size: 0.8rem;
  color: var(--primary, #42b883);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  background-color: rgba(66, 184, 131, 0.1);
}

.tag.active {
  color: white;
  background-color: var(--primary, #42b883);
}

.excerpt {
  color: var(--text-light, #666);
  margin: 0;
  line-height: 1.6;
}

.article-image {
  width: 200px;
  height: 150px;
  overflow: hidden;
}

.article-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.article-card:hover .article-image img {
  transform: scale(1.05);
}

.no-articles {
  text-align: center;
  color: var(--text-light, #666);
  padding: 3rem 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.no-articles svg {
  color: var(--text-light, #999);
  opacity: 0.5;
  margin-bottom: 1rem;
}

.no-articles p {
  font-size: 1.2rem;
  margin: 0;
}

.browse-all-button {
  margin-top: 1rem;
  padding: 0.8rem 1.5rem;
  background-color: var(--primary, #42b883);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.browse-all-button:hover {
  background-color: var(--secondary, #3a9f73);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* 加载状态 */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  color: var(--text-light, #666);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top-color: var(--primary, #42b883);
  animation: spin 0.8s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .article-link {
    flex-direction: column;
  }
  
  .article-image {
    width: 100%;
    height: 180px;
  }
  
  .tag-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .tag-title {
    font-size: 2rem;
  }
  
  .browse-all-button {
    width: 100%;
  }
}
</style> 