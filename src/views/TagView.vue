<template>
  <div class="tag-view">
    <h1 class="tag-title">#{{ tag }}</h1>
    
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
    
    <div v-else class="no-articles">
      没有找到标签为 "{{ tag }}" 的文章
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
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

.tag-title {
  font-size: 2.5rem;
  color: #42b883;
  margin-bottom: 1rem;
}

.articles-count {
  color: #666;
  margin-bottom: 2rem;
  font-size: 1.1rem;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.article-card {
  background: #fff;
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
  color: #2c3e50;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.date {
  color: #666;
  font-size: 0.9rem;
}

.tags {
  display: flex;
  gap: 0.5rem;
}

.tag {
  font-size: 0.8rem;
  color: #42b883;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  background-color: rgba(66, 184, 131, 0.1);
}

.tag.active {
  color: #fff;
  background-color: #42b883;
}

.excerpt {
  color: #666;
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
}

.no-articles {
  text-align: center;
  color: #666;
  padding: 3rem 0;
  font-style: italic;
}

@media (max-width: 768px) {
  .article-link {
    flex-direction: column;
  }
  
  .article-image {
    width: 100%;
    height: 180px;
  }
}
</style> 