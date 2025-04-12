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
  background: var(--surface, #fff);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
  height: 100%;
  display: flex;
  flex-direction: column;
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
  transition: transform 0.5s ease;
}

.article-card:hover .article-image img {
  transform: scale(1.05);
}

.article-content {
  padding: 1.5rem;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.article-title {
  margin: 0 0 1rem;
  font-size: 1.25rem;
  color: var(--text, #333);
  transition: color 0.3s ease;
}

.article-card:hover .article-title {
  color: var(--primary, #2c6e49);
}

.article-excerpt {
  color: var(--text-secondary, #666);
  margin-bottom: 1rem;
  line-height: 1.5;
  flex: 1;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
  color: var(--text-light, #888);
  margin-top: auto;
}

.tags {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.tag {
  color: var(--primary, #2c6e49);
  transition: color 0.2s ease, transform 0.2s ease;
}

.tag:hover {
  color: var(--secondary, #1a4d33);
  transform: translateY(-1px);
}

/* 适配深色模式 */
@media (prefers-color-scheme: dark) {
  .article-card {
    background: var(--surface, #1e1e1e);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
  }
  
  .article-card:hover {
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
  }
}
</style> 