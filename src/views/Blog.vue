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
      </div>
    </div>
    <div class="articles-grid">
      <ArticleCard 
        v-for="article in filteredArticles" 
        :key="article.id" 
        :article="article"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import ArticleCard from '../components/ArticleCard.vue'

const searchQuery = ref('')
const selectedTags = ref([])
const articles = ref([])
const loading = ref(false)
const error = ref(null)

// 加载文章列表
const loadArticles = async () => {
  try {
    loading.value = true
    error.value = null
    
    const response = await fetch('/src/assets/markdown/articles/index.json')
    const data = await response.json()
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
    article.tags.forEach(tag => tags.add(tag))
  })
  return Array.from(tags)
})

const filteredArticles = computed(() => {
  return articles.value.filter(article => {
    const matchesSearch = article.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         article.excerpt.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesTags = selectedTags.value.length === 0 ||
                       article.tags.some(tag => selectedTags.value.includes(tag))
    
    return matchesSearch && matchesTags
  })
})

const toggleTag = (tag) => {
  const index = selectedTags.value.indexOf(tag)
  if (index === -1) {
    selectedTags.value.push(tag)
  } else {
    selectedTags.value.splice(index, 1)
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<style scoped>
/* 限制内容宽度 */
.blog {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem 0;
}

.blog-filters {
  margin: 2rem 0;
}

.search input {
  width: 100%;
  max-width: 400px;
  padding: 0.8rem 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.tags {
  display: flex;
  gap: 0.8rem;
  flex-wrap: wrap;
}

.tags span {
  padding: 0.4rem 1rem;
  background: #f5f5f5;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tags span.active {
  background: #42b883;
  color: white;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}
</style> 