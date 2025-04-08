<template>
  <div>
    <div class="photo-grid">
      <div 
        v-for="photo in photos" 
        :key="photo.id" 
        class="photo-item"
        @click="navigateToPhoto(photo)"
      >
        <img 
          :src="photo.url" 
          :alt="photo.title"
          loading="lazy"
        />
        <div class="photo-info">
          <h3>{{ photo.title }}</h3>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
const props = defineProps({
  photos: {
    type: Array,
    required: true
  }
})

const navigateToPhoto = (photo) => {
  router.push(`/photo/${photo.id}`)
}
</script>

<style scoped>
.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
  padding: 1rem;
}

.photo-item {
  position: relative;
  aspect-ratio: 3/2;
  overflow: hidden;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.photo-item:hover {
  transform: scale(1.02);
}

.photo-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-info {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 1rem;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  color: white;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.photo-item:hover .photo-info {
  opacity: 1;
}

.photo-info h3 {
  margin: 0;
  font-size: 1.1rem;
}
</style> 