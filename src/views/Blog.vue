<template>
  <div class="blog">
    <h1>博客文章</h1>
    <div class="blog-filters">
      <div class="search">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="搜索文章..."
        >
      </div>
      <div class="tags">
        <span 
          v-for="tag in allTags" 
          :key="tag"
          :class="{ active: selectedTags.includes(tag) }"
          @click="toggleTag(tag)"
        >
          #{{ tag }}
        </span>
        <router-link 
          v-for="tag in allTags" 
          :key="`link-${tag}`"
          :to="`/tags/${tag}`"
          class="tag-link"
        >
          查看所有 #{{ tag }} 文章
        </router-link>
      </div>
    </div>
    
    <!-- 文章列表 -->
    <div class="articles-grid" v-if="filteredArticles.length > 0">
      <ArticleCard 
        v-for="article in visibleArticles" 
        :key="article.id" 
        :article="article"
        class="article-item"
      />
    </div>
    
    <!-- 加载中状态 -->
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>正在加载文章...</p>
    </div>
    
    <!-- 加载更多按钮 -->
    <div class="load-more-container" v-if="hasMoreArticles">
      <button 
        class="load-more-button" 
        @click="loadMore" 
        :disabled="isLoadingMore"
      >
        <span v-if="!isLoadingMore">加载更多</span>
        <span v-else class="loading-spinner"></span>
      </button>
    </div>
    
    <!-- 全部加载完毕提示 -->
    <div v-if="!hasMoreArticles && visibleArticles.length > 0" class="all-loaded">
      <div class="line"></div>
      <p>已加载全部文章</p>
      <div class="line"></div>
    </div>
    
    <!-- 无结果提示 -->
    <div v-if="filteredArticles.length === 0 && !loading" class="no-results">
      <div class="no-results-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          <line x1="11" y1="8" x2="11" y2="14"></line>
          <line x1="8" y1="11" x2="14" y2="11"></line>
        </svg>
      </div>
      <h2 class="no-results-title">没有找到匹配的文章</h2>
      <p class="no-results-message">
        {{ getNoResultsMessage() }}
      </p>
      <div class="search-tips">
        <h4>搜索技巧：</h4>
        <ul>
          <li>检查关键词拼写是否正确</li>
          <li>尝试使用不同的关键词</li>
          <li>使用更通用的关键词</li>
          <li>减少标签筛选条件</li>
        </ul>
      </div>
      <button @click="resetFilters" class="reset-button">重置所有筛选条件</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import ArticleCard from '../components/ArticleCard.vue'

const router = useRouter()
const searchQuery = ref('')
const selectedTags = ref([])
const articles = ref([])
const loading = ref(false)
const error = ref(null)

// 分页相关参数
const itemsPerPage = 6 // 每次显示的文章数量
const currentPage = ref(1) // 当前页码
const isLoadingMore = ref(false) // 是否正在加载更多

// 加载文章列表
const loadArticles = async () => {
  try {
    loading.value = true
    error.value = null
    
    const response = await fetch('/markdown/articles/index.json')
    const data = await response.json()
    
    // 直接使用API返回的真实数据
    articles.value = data.articles
  } catch (err) {
    error.value = '加载文章列表失败，请稍后重试'
    console.error('加载文章列表失败:', err)
    articles.value = []
  } finally {
    loading.value = false
  }
}

const allTags = computed(() => {
  const tags = new Set()
  articles.value.forEach(article => {
    article.tags && article.tags.forEach(tag => tags.add(tag))
  })
  return Array.from(tags)
})

const filteredArticles = computed(() => {
  return articles.value.filter(article => {
    const matchesSearch = article.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         article.excerpt.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesTags = selectedTags.value.length === 0 ||
                       (article.tags && article.tags.some(tag => selectedTags.value.includes(tag)))
    
    return matchesSearch && matchesTags
  })
})

// 当前可见的文章列表
const visibleArticles = computed(() => {
  const endIndex = currentPage.value * itemsPerPage
  return filteredArticles.value.slice(0, endIndex)
})

// 是否还有更多文章可加载
const hasMoreArticles = computed(() => {
  return visibleArticles.value.length < filteredArticles.value.length
})

// 加载更多文章
const loadMore = () => {
  if (isLoadingMore.value || !hasMoreArticles.value) return
  
  isLoadingMore.value = true
  
  // 模拟网络请求延迟
  setTimeout(() => {
    currentPage.value++
    isLoadingMore.value = false
  }, 500)
}

const toggleTag = (tag) => {
  const index = selectedTags.value.indexOf(tag)
  if (index === -1) {
    selectedTags.value.push(tag)
  } else {
    selectedTags.value.splice(index, 1)
  }
}

// 重置筛选条件
const resetFilters = () => {
  searchQuery.value = ''
  selectedTags.value = []
}

// 根据筛选条件生成无结果消息
const getNoResultsMessage = () => {
  if (searchQuery.value && selectedTags.value.length > 0) {
    return `没有找到包含"${searchQuery.value}"且标签为"${selectedTags.value.join('、')}"的文章`;
  } else if (searchQuery.value) {
    return `没有找到包含"${searchQuery.value}"的文章`;
  } else if (selectedTags.value.length > 0) {
    return `没有找到标签为"${selectedTags.value.join('、')}"的文章`;
  } else {
    return "没有找到任何文章";
  }
}

// 监听筛选条件变化，重置分页
watch([searchQuery, selectedTags], () => {
  currentPage.value = 1
}, { deep: true })

onMounted(() => {
  loadArticles()
})
</script>

<style scoped>
/* 限制内容宽度 */
.blog {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem 1rem;
  user-select: none; 
}

.blog-filters {
  margin: 2rem 0;
}

.search input {
  width: 100%;
  max-width: 400px;
  padding: 0.8rem 1rem;
  border: 1px solid var(--border, #ddd);
  border-radius: 8px;
  margin-bottom: 1rem;
  background-color: var(--surface, #fff);
  color: var(--text, #333);
  transition: all 0.3s ease;
}

.search input:focus {
  outline: none;
  border-color: var(--primary, #42b883);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.1);
}

.tags {
  display: flex;
  gap: 0.8rem;
  flex-wrap: wrap;
  margin-bottom: 1rem;
}

.tags span {
  padding: 0.4rem 1rem;
  background: var(--accent, #f5f5f5);
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: var(--text, #333);
}

.tags span:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tags span.active {
  background: var(--primary, #42b883);
  color: white;
}

.tag-link {
  display: none;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.article-item {
  animation: fadeIn 0.5s ease forwards;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 加载更多按钮 */
.load-more-container {
  text-align: center;
  margin: 3rem 0 2rem;
}

.load-more-button {
  padding: 0.8rem 2rem;
  background-color: var(--primary, #42b883);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.3s ease;
  min-width: 140px;
  min-height: 48px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.load-more-button:hover:not(:disabled) {
  background-color: var(--secondary, #3a9f73);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(66, 184, 131, 0.2);
}

.load-more-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 加载动画 */
.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
  display: inline-block;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 全部加载完毕提示 */
.all-loaded {
  display: flex;
  align-items: center;
  color: var(--text-light, #999);
  margin: 3rem 0;
  padding: 0 1rem;
}

.all-loaded .line {
  flex-grow: 1;
  height: 1px;
  background-color: var(--border, #eee);
  margin: 0 1rem;
}

.all-loaded p {
  white-space: nowrap;
  font-size: 0.9rem;
  margin: 0;
}

/* 无搜索结果状态 */
.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  margin: 2rem 0;
  color: var(--text-light, #666);
  animation: fadeIn 0.5s ease;
}

.no-results-icon {
  margin-bottom: 1.5rem;
  color: var(--text-light, #999);
  opacity: 0.6;
  animation: pulse 3s infinite ease-in-out;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 0.6; }
  50% { transform: scale(1.05); opacity: 0.8; }
}

.no-results-title {
  font-size: 1.5rem;
  margin: 0 0 1rem 0;
  color: var(--text, #333);
}

.no-results-message {
  font-size: 1.1rem;
  margin-bottom: 1.5rem;
  text-align: center;
}

.search-tips {
  background-color: var(--accent, #f5f5f5);
  border-radius: 10px;
  padding: 1.2rem 1.5rem;
  margin-bottom: 1.5rem;
  width: 100%;
  max-width: 450px;
  border-left: 4px solid var(--primary, #42b883);
}

.search-tips h4 {
  margin-top: 0;
  margin-bottom: 0.8rem;
  color: var(--text, #333);
}

.search-tips ul {
  margin: 0;
  padding-left: 1.5rem;
  color: var(--text-light, #666);
}

.search-tips li {
  margin-bottom: 0.5rem;
}

.search-tips li:last-child {
  margin-bottom: 0;
}

.reset-button {
  padding: 0.8rem 1.5rem;
  background-color: var(--primary, #42b883);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
}

.reset-button:hover {
  background-color: var(--secondary, #3a9f73);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(66, 184, 131, 0.2);
}

/* 加载中状态 */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 0;
  color: var(--text-light, #666);
  text-align: center;
}

.loading-container .spinner {
  width: 50px;
  height: 50px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top-color: var(--primary, #42b883);
  animation: spin 0.8s linear infinite;
  margin-bottom: 1.5rem;
}

.loading-container p {
  font-size: 1.1rem;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .blog {
    padding: 1rem;
  }
  
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
  
  .search input {
    max-width: 100%;
    font-size: 16px; /* 防止iOS缩放 */
  }
  
  .load-more-button {
    width: 100%;
    max-width: 300px;
  }
  
  .tags {
    gap: 0.6rem;
    margin-top: 1rem;
  }
  
  .tags span {
    padding: 0.3rem 0.8rem;
    font-size: 0.9rem;
    white-space: nowrap;
  }
  
  .all-loaded {
    flex-direction: column;
    gap: 1rem;
  }
  
  .all-loaded .line {
    width: 80%;
    margin: 0;
  }
  
  .search-tips {
    padding: 1rem;
    width: calc(100% - 2rem);
  }
  
  .reset-button {
    width: 100%;
  }
}

/* 针对超小屏幕的优化 */
@media (max-width: 480px) {
  .blog h1 {
    font-size: 1.8rem;
  }
  
  .search input {
    padding: 0.7rem;
  }
  
  .tags span {
    padding: 0.3rem 0.7rem;
    font-size: 0.85rem;
  }
  
  .load-more-button {
    font-size: 0.9rem;
    min-height: 42px;
  }
  
  .no-results-title {
    font-size: 1.3rem;
  }
  
  .no-results-message {
    font-size: 1rem;
  }
}
</style> 