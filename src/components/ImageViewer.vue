<template>
  <div 
    v-if="show" 
    class="image-viewer"
    @click="close"
  >
    <div class="viewer-content" @click.stop>
      <button class="close-btn" @click="close">&times;</button>
      <img :src="image.url" :alt="image.title">
      <div class="image-info">
        <h3>{{ image.title }}</h3>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  show: Boolean,
  image: Object
})

const emit = defineEmits(['close'])

const close = () => {
  emit('close')
}
</script>

<style scoped>
.image-viewer {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.viewer-content {
  position: relative;
  max-width: 90%;
  max-height: 90vh;
}

.viewer-content img {
  max-width: 100%;
  max-height: 90vh;
  object-fit: contain;
}

.close-btn {
  position: absolute;
  top: -40px;
  right: -40px;
  background: none;
  border: none;
  color: white;
  font-size: 2rem;
  cursor: pointer;
}

.image-info {
  position: absolute;
  bottom: -40px;
  left: 0;
  width: 100%;
  color: white;
  text-align: center;
}
</style> 