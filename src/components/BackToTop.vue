<template>
  <div v-if="visible" class="back-to-top" @click="scrollToTop">
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-arrow-up">
      <line x1="12" y1="19" x2="12" y2="5"></line>
      <polyline points="5 12 12 5 19 12"></polyline>
    </svg>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const visible = ref(false)

const toggleVisibility = () => {
  visible.value = document.body.scrollTop > 800
}

const scrollToTop = () => {
  document.body.scrollTo({
    top: 0,
    behavior: 'smooth'
  });
}

onMounted(() => {
  document.body.addEventListener('scroll', toggleVisibility)
})

onUnmounted(() => {
  document.body.removeEventListener('scroll', toggleVisibility)
})
</script>

<style scoped>
.back-to-top {
  position: fixed;
  bottom: 2.5rem;
  right: 2.5rem;
  background-color: rgba(66, 184, 131, 0.7);
  color: white;
  padding: 0.8rem;
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
  transition: background-color 0.3s ease, transform 0.3s ease;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

.back-to-top:hover {
  background-color: rgba(66, 184, 131, 0.9);
  transform: scale(1.1);
}

.back-to-top svg {
  width: 20px;
  height: 20px;
  stroke: white;
}
</style> 