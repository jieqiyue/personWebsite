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

.portfolio h1 {
  color: var(--text);
  margin-bottom: 1.5rem;
}

.categories {
  margin: 2rem 0;
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

button {
  padding: 0.7rem 1.8rem;
  border: 1px solid var(--border);
  border-radius: 25px;
  background: var(--surface);
  cursor: pointer;
  transition: all 0.3s ease;
  color: var(--text);
  font-weight: 500;
  font-size: 0.95rem;
  position: relative;
  
  /* 毛玻璃效果 */
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  box-shadow: 0 4px 15px var(--shadow);
  
  /* 确保元素层级正确 */
  z-index: 1;
  overflow: hidden;
}

/* 按钮背景发光效果 */
button::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    to right,
    transparent,
    rgba(255, 255, 255, 0.1),
    transparent
  );
  transform: translateX(-100%);
  transition: transform 0.6s ease;
  z-index: -1;
}

button:hover {
  transform: translateY(-3px);
  border-color: var(--primary);
  box-shadow: 0 6px 20px var(--shadow);
}

button:hover::before {
  transform: translateX(100%);
}

button.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
  transform: translateY(-2px);
}

/* 活跃按钮的扩散效果 */
button.active::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 5px;
  height: 5px;
  background: rgba(255, 255, 255, 0.5);
  opacity: 0;
  border-radius: 100%;
  transform: scale(1);
  animation: ripple 1s ease-out;
}

@keyframes ripple {
  0% {
    transform: scale(1);
    opacity: 0.4;
  }
  100% {
    transform: scale(20);
    opacity: 0;
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .categories {
    justify-content: center;
  }
  
  button {
    padding: 0.6rem 1.5rem;
    font-size: 0.9rem;
  }
}
</style> 