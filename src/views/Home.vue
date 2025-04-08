<template>
  <div class="home">
    <!-- Markdown渲染区域 -->
    <div class="markdown-container">
      <MarkdownRenderer :source="markdownSource" />
    </div>

    <!-- 分隔线 -->
    <div class="section-divider"></div>

    <!-- 精选作品部分 -->
    <section class="featured-works">
      <h2>精选作品</h2>
      <PhotoGrid :photos="featuredPhotos" />
    </section>

    <!-- 最新文章部分 -->
    <section class="latest-articles">
      <h2>最新文章</h2>
      <div class="articles-grid">
        <ArticleCard 
          v-for="article in latestArticles" 
          :key="article.id" 
          :article="article" 
        />
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import PhotoGrid from '../components/PhotoGrid.vue'
import ArticleCard from '../components/ArticleCard.vue'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'

// Markdown文件路径
const markdownSource = ref('/src/assets/markdown/Home.md')
const latestArticles = ref([])
const loading = ref(false)

// 示例数据 - 这部分也可以改为本地JSON文件
const featuredPhotos = [
  {
    id: 1,
    title: '山间晨雾',
    url: '/src/assets/images/photos/mountain.jpg'
  },
  {
    id: 2,
    title: '城市夜景',
    url: '/src/assets/images/photos/city.jpg'
  },
  {
    id: 3,
    title: '森林小径',
    url: '/src/assets/images/photos/forest.jpg'
  }
]

// 加载最新文章
const loadLatestArticles = async () => {
  try {
    loading.value = true
    const response = await fetch('/src/assets/markdown/articles/index.json')
    const data = await response.json()
    // 获取最新的2篇文章
    latestArticles.value = data.articles
      .sort((a, b) => new Date(b.date) - new Date(a.date))
      .slice(0, 2)
  } catch (error) {
    console.error('加载文章失败:', error)
    latestArticles.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLatestArticles()
})
</script>

<style scoped>
.home {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
}

.markdown-container {
  margin-bottom: 2rem;
}

.section-divider {
  height: 1px;
  background: #daebd9;
  margin: 3rem 0;
}

.featured-works,
.latest-articles {
  margin-bottom: 4rem;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}
</style> 