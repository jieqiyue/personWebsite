<template>
  <article class="article-card" @click="navigateToArticle">
    <div class="article-image" v-if="article.coverImage">
      <img :src="article.coverImage" :alt="article.title" loading="lazy" />
    </div>
    <div class="article-content">
      <h3 class="article-title">{{ article.title }}</h3>
      <p class="article-excerpt">{{ article.excerpt }}</p>
      <div class="article-meta">
        <div class="meta-left">
          <span class="date">{{ formatDate(article.date) }}</span>
          <span class="reading-time" v-if="article.readingTime">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
            {{ article.readingTime }}
          </span>
        </div>
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
import { formatReadingTime } from '../utils/article'

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
  line-height: 1.4;
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
  flex-wrap: wrap;
  gap: 0.8rem;
}

.meta-left {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.reading-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--text-light, #888);
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
  color: var(--primary, #2c6e49);
  transition: color 0.2s ease, transform 0.2s ease;
}

.tag:hover {
  color: var(--secondary, #1a4d33);
  transform: translateY(-1px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .article-content {
    padding: 1.2rem;
  }
  
  .article-title {
    font-size: 1.2rem;
  }
  
  .article-excerpt {
    font-size: 0.95rem;
    -webkit-line-clamp: 3;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
}

@media (max-width: 480px) {
  .article-image {
    height: 180px;
  }
  
  .article-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.6rem;
  }
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