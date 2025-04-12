<template>
  <div class="skeleton-loader" :class="{ 'with-image': withImage }">
    <div v-if="withImage" class="skeleton-image pulse"></div>
    <div class="skeleton-content">
      <div class="skeleton-title pulse"></div>
      <div class="skeleton-text pulse"></div>
      <div class="skeleton-text pulse"></div>
      <div class="skeleton-meta pulse"></div>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';

defineProps({
  withImage: {
    type: Boolean,
    default: true
  }
});
</script>

<style scoped>
.skeleton-loader {
  background: var(--surface, #fff);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.skeleton-image {
  width: 100%;
  height: 200px;
  background-color: var(--skeleton-base, #e9e9e9);
}

.skeleton-content {
  padding: 1.5rem;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.skeleton-title {
  height: 1.5rem;
  width: 80%;
  margin-bottom: 1rem;
  background-color: var(--skeleton-base, #e9e9e9);
  border-radius: 4px;
}

.skeleton-text {
  height: 1rem;
  margin-bottom: 0.8rem;
  background-color: var(--skeleton-base, #e9e9e9);
  border-radius: 4px;
  width: 100%;
}

.skeleton-text:last-of-type {
  width: 60%;
}

.skeleton-meta {
  height: 0.8rem;
  width: 50%;
  background-color: var(--skeleton-base, #e9e9e9);
  border-radius: 4px;
  margin-top: auto;
}

.pulse {
  animation: pulse 1.5s ease-in-out 0.5s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 适配深色模式 */
@media (prefers-color-scheme: dark) {
  .skeleton-loader {
    background: var(--surface, #1e1e1e);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
  }
  
  .skeleton-image,
  .skeleton-title,
  .skeleton-text,
  .skeleton-meta {
    background-color: var(--skeleton-base, #333);
  }
}
</style> 