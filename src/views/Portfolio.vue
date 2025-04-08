<template>
  <div class="portfolio">
    <h1>摄影作品集</h1>
    <div class="categories">
      <button 
        v-for="category in categories" 
        :key="category.id"
        :class="{ active: currentCategory === category.id }"
        @click="currentCategory = category.id"
      >
        {{ category.name }}
      </button>
    </div>
    <PhotoGrid :photos="filteredPhotos" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import PhotoGrid from '../components/PhotoGrid.vue'

const currentCategory = ref('all')

const categories = [
  { id: 'all', name: '全部' },
  { id: 'landscape', name: '风景' },
  { id: 'portrait', name: '人像' },
  { id: 'street', name: '街拍' },
  { id: 'nature', name: '自然' }
]

// 示例照片数据
const photos = [
  {
    id: 1,
    title: '山间晨雾',
    url: 'https://images.unsplash.com/photo-1506744038136-46273834b3fb',
    category: 'landscape'
  },
  {
    id: 2,
    title: '城市剪影',
    url: 'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05',
    category: 'street'
  },
  {
    id: 3,
    title: '森林小径',
    url: 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e',
    category: 'nature'
  },
  {
    id: 4,
    title: '街头故事',
    url: 'https://images.unsplash.com/photo-1492684223066-81342ee5ff30',
    category: 'street'
  },
  {
    id: 5,
    title: '人像写真',
    url: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330',
    category: 'portrait'
  }
]

const filteredPhotos = computed(() => {
  if (currentCategory.value === 'all') {
    return photos
  }
  return photos.filter(photo => photo.category === currentCategory.value)
})
</script>

<style scoped>
/* 限制内容宽度 */
.portfolio {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem 0;
}

.categories {
  margin: 2rem 0;
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

button {
  padding: 0.5rem 1.5rem;
  border: 1px solid #ddd;
  border-radius: 20px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
}

button.active {
  background: #42b883;
  color: white;
  border-color: #42b883;
}

button:hover {
  border-color: #42b883;
}
</style> 