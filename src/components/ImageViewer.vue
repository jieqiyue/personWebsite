<template>
  <transition name="fade">
    <div 
      v-if="visible" 
      class="image-viewer" 
      @click.self="closeViewer"
      @touchstart.self="closeViewer"
    >
      <div class="image-wrapper">
        <transition name="zoom">
          <img 
            v-if="visible" 
            :src="imageSrc" 
            :alt="imageAlt" 
            class="preview-image" 
            @click.stop
          />
        </transition>
      </div>
      <button class="close-button" @click="closeViewer" aria-label="关闭预览">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  imageSrc: {
    type: String,
    required: true
  },
  imageAlt: {
    type: String,
    default: '图片预览'
  }
});

const emit = defineEmits(['close']);

const closeViewer = () => {
  emit('close');
  // 确保关闭时立即恢复滚动状态
  document.body.style.overflow = '';
  document.body.classList.remove('no-scroll');
};

// 监听ESC键关闭预览
const handleKeyDown = (e) => {
  if (e.key === 'Escape' && props.visible) {
    closeViewer();
  }
};

// 阻止背景滚动
const preventScroll = () => {
  if (props.visible) {
    document.body.style.overflow = 'hidden';
    document.body.classList.add('no-scroll');
  } else {
    document.body.style.overflow = '';
    document.body.classList.remove('no-scroll');
  }
};

// 确保在组件更新、关闭和卸载时都能正确恢复滚动
watch(() => props.visible, (newVal) => {
  preventScroll();
  // 额外检查：如果关闭了预览但滚动仍被禁用，强制恢复
  if (!newVal) {
    setTimeout(() => {
      document.body.style.overflow = '';
      document.body.classList.remove('no-scroll');
    }, 300); // 略大于动画时间，确保动画完成后恢复
  }
}, { immediate: true });

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);
  // 组件卸载时确保恢复滚动
  document.body.style.overflow = '';
  document.body.classList.remove('no-scroll');
});
</script>

<style scoped>
.image-viewer {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.9);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: zoom-out;
}

.image-wrapper {
  position: relative;
  max-width: 90%;
  max-height: 90%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-image {
  max-width: 100%;
  max-height: 90vh;
  object-fit: contain;
  cursor: default;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.5);
}

.close-button {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 40px;
  height: 40px;
  background-color: rgba(0, 0, 0, 0.5);
  border: none;
  border-radius: 50%;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 1001;
}

.close-button:hover {
  background-color: rgba(0, 0, 0, 0.8);
  transform: scale(1.1);
}

/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 缩放动画 */
.zoom-enter-active,
.zoom-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.zoom-enter-from {
  transform: scale(0.9);
  opacity: 0;
}

.zoom-leave-to {
  transform: scale(0.9);
  opacity: 0;
}

@media (max-width: 768px) {
  .close-button {
    top: 10px;
    right: 10px;
    width: 36px;
    height: 36px;
  }
  
  .preview-image {
    max-height: 85vh;
  }
}
</style> 