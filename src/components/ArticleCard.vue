<template>
  <article class="article-card" @click="navigateToArticle">
    <div class="article-image" v-if="article.coverImage">
      <img :src="article.coverImage" :alt="article.title" loading="lazy" />
    </div>
    <div class="article-content">
      <h3 class="article-title">{{ article.title }}</h3>
      <p class="article-excerpt">{{ article.excerpt }}</p>
      <div class="article-meta">
        <span class="date">{{ formatDate(article.date) }}</span>
        <span class="tags" v-if="article.tags">
          <span 
            v-for="tag in article.tags" 
            :key="tag" 
            class="tag"
            @click.stop="navigateToTag(tag)"
          >
            #{{ tag }}
          </span>
        </span>
      </div>
    </div>
  </article>
</template>

<script setup>
import { defineProps } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const navigateToArticle = () => {
  router.push(`/articles/${props.article.id}`)
}

const navigateToTag = (tag) => {
  router.push(`/tags/${tag}`)
}
</script>

<style scoped>
.article-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
}

.article-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}

.article-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.article-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-content {
  padding: 1.5rem;
}

.article-title {
  margin: 0 0 1rem;
  font-size: 1.25rem;
  color: #333;
  transition: color 0.3s ease;
}

.article-card:hover .article-title {
  color: #2c6e49;
}

.article-excerpt {
  color: #666;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
  color: #888;
}

.tags {
  display: flex;
  gap: 0.5rem;
}

.tag {
  color: #2c6e49;
}
</style> 