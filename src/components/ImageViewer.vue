<template>
  <transition name="fade">
    <div 
      v-if="visible" 
      class="image-viewer" 
      @click.self="closeViewer"
      @touchstart.self="closeViewer"
      @keydown.left="navigatePrev"
      @keydown.right="navigateNext"
      tabindex="0"
      ref="viewerContainer"
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
      
      <!-- 导航按钮 -->
      <div class="navigation-controls">
        <button 
          v-if="hasPrevious" 
          class="nav-button prev-button" 
          @click.stop="navigatePrev" 
          aria-label="上一张图片"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
        </button>
        <button 
          v-if="hasNext" 
          class="nav-button next-button" 
          @click.stop="navigateNext" 
          aria-label="下一张图片"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue';

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
  },
  // 新增属性
  images: {
    type: Array,
    default: () => []
  },
  currentIndex: {
    type: Number,
    default: -1
  }
});

const emit = defineEmits(['close', 'navigate']);
const viewerContainer = ref(null);
const hasPrevious = ref(false);
const hasNext = ref(false);

// 导航到上一张图片
const navigatePrev = () => {
  emit('navigate', 'prev');
};

// 导航到下一张图片
const navigateNext = () => {
  emit('navigate', 'next');
};

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
  } else if (e.key === 'ArrowLeft' && props.visible) {
    navigatePrev();
  } else if (e.key === 'ArrowRight' && props.visible) {
    navigateNext();
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
  
  // 当预览变为可见时，聚焦容器以启用键盘导航
  if (newVal) {
    nextTick(() => {
      if (viewerContainer.value) {
        viewerContainer.value.focus();
      }
    });
  }
  
  // 额外检查：如果关闭了预览但滚动仍被禁用，强制恢复
  if (!newVal) {
    setTimeout(() => {
      document.body.style.overflow = '';
      document.body.classList.remove('no-scroll');
    }, 300); // 略大于动画时间，确保动画完成后恢复
  }
}, { immediate: true });

// 设置是否有上一张/下一张图片的状态
watch(() => props.currentIndex, (newIndex) => {
  if (props.images && props.images.length > 0) {
    hasPrevious.value = newIndex > 0;
    hasNext.value = newIndex < props.images.length - 1;
  } else {
    hasPrevious.value = false;
    hasNext.value = false;
  }
}, { immediate: true });

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
  
  // 添加触摸滑动支持
  if (viewerContainer.value) {
    let touchStartX = 0;
    let touchEndX = 0;
    
    viewerContainer.value.addEventListener('touchstart', (e) => {
      touchStartX = e.changedTouches[0].screenX;
    });
    
    viewerContainer.value.addEventListener('touchend', (e) => {
      touchEndX = e.changedTouches[0].screenX;
      handleSwipe();
    });
    
    const handleSwipe = () => {
      const swipeThreshold = 100; // 最小滑动距离
      
      if (touchEndX - touchStartX > swipeThreshold) {
        // 向右滑动，显示上一张
        navigatePrev();
      } else if (touchStartX - touchEndX > swipeThreshold) {
        // 向左滑动，显示下一张
        navigateNext();
      }
    };
  }
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
  cursor: pointer;
  outline: none; /* 移除聚焦时的轮廓 */
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

/* 导航按钮 */
.navigation-controls {
  position: absolute;
  width: 100%;
  display: flex;
  justify-content: space-between;
  padding: 0 40px;
  pointer-events: none; /* 防止影响背景点击 */
}

.nav-button {
  width: 50px;
  height: 50px;
  background-color: rgba(0, 0, 0, 0.5);
  border: none;
  border-radius: 50%;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  pointer-events: auto; /* 允许按钮点击 */
  opacity: 0.7;
}

.nav-button:hover {
  background-color: rgba(0, 0, 0, 0.8);
  transform: scale(1.1);
  opacity: 1;
}

.prev-button {
  margin-right: auto;
}

.next-button {
  margin-left: auto;
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
  
  .navigation-controls {
    padding: 0 20px;
  }
  
  .nav-button {
    width: 40px;
    height: 40px;
  }
  
  .preview-image {
    max-height: 85vh;
  }
}
</style> 