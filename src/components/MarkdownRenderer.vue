<template>
  <div>
    <div class="markdown-body" v-html="renderedContent" ref="markdownBody"></div>
    <ImageViewer 
      v-if="previewVisible" 
      :visible="previewVisible" 
      :imageSrc="previewImageSrc" 
      :imageAlt="previewImageAlt" 
      :images="allImages"
      :currentIndex="currentImageIndex"
      @close="closePreview"
      @navigate="handleImageNavigation"
    />
  </div>
</template>

<script>
import { computed, ref, onMounted, watch, onUnmounted } from 'vue'
import md from '../plugins/markdown.js'
import 'highlight.js/styles/github.css'
import 'github-markdown-css/github-markdown.css'
import './markdown-custom.css'
import ClipboardJS from 'clipboard'
import ImageViewer from './ImageViewer.vue'

export default {
  name: 'MarkdownRenderer',
  components: {
    ImageViewer
  },
  props: {
    content: {
      type: String,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const markdownBody = ref(null);
    let clipboardInstance = null;
    
    // 图片预览状态
    const previewVisible = ref(false);
    const previewImageSrc = ref('');
    const previewImageAlt = ref('');
    const allImages = ref([]);
    const currentImageIndex = ref(-1);
    
    const renderedContent = computed(() => {
      if (props.loading) {
        return '<div class="loading">正在加载...</div>'
      }
      
      if (!props.content) {
        return '<div class="empty">暂无内容</div>'
      }
      
      try {
        return md.render(props.content)
      } catch (error) {
        console.error('渲染Markdown时出错:', error)
        return `<div class="error">渲染内容时出错: ${error.message}</div>`
      }
    })

    // 处理代码复制功能
    const setupCodeCopyButtons = () => {
      if (!markdownBody.value) return;
      
      // 先销毁之前的实例
      if (clipboardInstance) {
        clipboardInstance.destroy();
      }
      
      // 创建新的clipboard实例，使用text回调函数
      clipboardInstance = new ClipboardJS('.code-copy-btn', {
        text: function(trigger) {
          // 获取最近的code元素
          const pre = trigger.closest('pre');
          const code = pre.querySelector('code');
          return code ? code.innerText : '';
        }
      });
      
      // 复制成功事件监听
      clipboardInstance.on('success', (e) => {
        const button = e.trigger;
        
        // 显示按钮复制状态
        button.classList.add('copied');
        
        // 创建一个固定定位的提示元素，直接添加到body
        const notification = document.createElement('div');
        notification.textContent = '复制成功!';
        notification.style.cssText = `
          position: fixed;
          top: 20px;
          left: 50%;
          transform: translateX(-50%);
          background-color: #42b883;
          color: white;
          padding: 8px 16px;
          border-radius: 4px;
          font-size: 14px;
          z-index: 10000;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
        `;
        document.body.appendChild(notification);
        
        // 使用JavaScript动画而不是CSS过渡
        let opacity = 0;
        notification.style.opacity = opacity;
        
        // 淡入动画
        const fadeIn = setInterval(() => {
          opacity += 0.1;
          notification.style.opacity = opacity;
          if (opacity >= 1) {
            clearInterval(fadeIn);
            
            // 2秒后淡出
            setTimeout(() => {
              const fadeOut = setInterval(() => {
                opacity -= 0.1;
                notification.style.opacity = opacity;
                if (opacity <= 0) {
                  clearInterval(fadeOut);
                  document.body.removeChild(notification);
                }
              }, 30);
            }, 2000);
          }
        }, 30);
        
        // 2秒后恢复按钮状态
        setTimeout(() => {
          button.classList.remove('copied');
        }, 2000);
        
        e.clearSelection();
      });
      
      // 复制失败事件监听
      clipboardInstance.on('error', (e) => {
        const button = e.trigger;
        button.setAttribute('title', '复制失败，请手动复制');
        console.error('复制失败:', e.action);
      });
    };
    
    // 设置图片点击事件
    const setupImagePreview = () => {
      if (!markdownBody.value) return;
      
      // 每次内容更新时，重新收集图片
      allImages.value = [];
      
      // 使用事件代理，确保只绑定一次
      if (!markdownBody.value.hasAttribute('data-preview-initialized')) {
        markdownBody.value.setAttribute('data-preview-initialized', 'true');
        
        // 绑定点击事件
        markdownBody.value.addEventListener('click', (event) => {
          const target = event.target;
          // 检查点击的是否是图片
          if (target.tagName === 'IMG') {
            event.preventDefault();
            
            // 查找当前图片在图片数组中的索引
            currentImageIndex.value = allImages.value.findIndex(img => img === target);
            
            if (currentImageIndex.value === -1) {
              // 如果没找到，可能是新添加的图片，直接使用当前图片
              previewImageSrc.value = target.src;
              previewImageAlt.value = target.alt || '图片预览';
            } else {
              // 更新预览图片
              updatePreviewImage();
            }
            
            previewVisible.value = true;
          }
        });
      }
      
      // 为所有图片添加样式和提示
      const images = markdownBody.value.querySelectorAll('img');
      images.forEach(img => {
        // 添加到全局图片数组，用于导航
        if (!allImages.value.includes(img)) {
          allImages.value.push(img);
        }
        
        // 如果图片已初始化，跳过
        if (img.hasAttribute('data-preview-ready')) return;
        
        // 标记为已初始化
        img.setAttribute('data-preview-ready', 'true');
        
        // 添加加载状态指示
        img.classList.add('loading-image');
        img.addEventListener('load', () => {
          img.classList.remove('loading-image');
          img.classList.add('loaded-image');
        });
        
        img.addEventListener('error', () => {
          img.classList.remove('loading-image');
          img.classList.add('error-image');
          img.alt = '图片加载失败: ' + (img.alt || '');
          img.style.padding = '20px';
          img.style.backgroundColor = '#fff0f0';
        });
        
        // 添加鼠标悬停提示
        img.title = '点击查看大图';
        img.classList.add('zoomable-image');
        
        // 只有当图片不在图片容器内时才添加容器
        if (!img.parentElement.classList.contains('image-container')) {
          // 给图片外面包一层容器，方便添加鼠标悬停效果
          const container = document.createElement('div');
          container.className = 'image-container';
          img.parentNode.insertBefore(container, img);
          container.appendChild(img);
        }
      });
    };
    
    // 关闭图片预览
    const closePreview = () => {
      previewVisible.value = false;
      // 确保滚动状态恢复
      document.body.style.overflow = '';
      document.body.classList.remove('no-scroll');
    };
    
    // 处理图片导航
    const handleImageNavigation = (direction) => {
      if (allImages.value.length <= 1) return;
      
      if (direction === 'prev' && currentImageIndex.value > 0) {
        currentImageIndex.value--;
        updatePreviewImage();
      } else if (direction === 'next' && currentImageIndex.value < allImages.value.length - 1) {
        currentImageIndex.value++;
        updatePreviewImage();
      }
    };
    
    // 更新预览图片
    const updatePreviewImage = () => {
      if (currentImageIndex.value >= 0 && currentImageIndex.value < allImages.value.length) {
        const img = allImages.value[currentImageIndex.value];
        previewImageSrc.value = img.src;
        previewImageAlt.value = img.alt || '图片预览';
      }
    };
    
    // 监听内容变化，重新设置功能
    watch(() => props.content, () => {
      // 等待DOM更新后设置
      setTimeout(() => {
        setupCodeCopyButtons();
        setupImagePreview();
      }, 100); // 增加延迟以确保DOM完全更新
    });
    
    // 组件挂载后设置功能
    onMounted(() => {
      setupCodeCopyButtons();
      
      // 延迟执行确保DOM渲染完成
      setTimeout(() => {
        setupImagePreview();
      }, 100);
    });
    
    // 组件卸载前销毁clipboard实例
    onUnmounted(() => {
      if (clipboardInstance) {
        clipboardInstance.destroy();
        clipboardInstance = null;
      }
    });

    return {
      renderedContent,
      markdownBody,
      previewVisible,
      previewImageSrc,
      previewImageAlt,
      closePreview,
      handleImageNavigation,
      updatePreviewImage,
      allImages,
      currentImageIndex
    }
  }
}
</script>

<style>
/* 图片加载状态 */
.loading-image {
  min-height: 100px;
  background-color: #f8f9fa;
  position: relative;
}

.loading-image::before {
  content: "加载中...";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #6c757d;
  font-size: 14px;
}

.loading-image::after {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 30px;
  height: 30px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.error-image {
  position: relative;
}

.error-image::before {
  content: "⚠️ 图片加载失败";
  text-align: center;
  color: #dc3545;
}

@keyframes spin {
  0% { transform: translate(-50%, -50%) rotate(0deg); }
  100% { transform: translate(-50%, -50%) rotate(360deg); }
}

.markdown-body {
  box-sizing: border-box;
  min-width: 200px;
  max-width: 980px;
  margin: 0 auto;
  padding: 45px;
}

@media (max-width: 767px) {
  .markdown-body {
    padding: 15px;
  }
}

.loading,
.empty,
.error {
  padding: 20px;
  text-align: center;
  border-radius: 4px;
  margin: 20px 0;
}

.loading {
  background-color: #f0f8ff;
  color: #1e90ff;
}

.empty {
  background-color: #f5f5f5;
  color: #666;
}

.error {
  background-color: #fff0f0;
  color: #ff3333;
  border: 1px solid #ffdddd;
}

/* 图片容器样式 */
.image-container {
  position: relative;
  display: inline-block;
  cursor: pointer;
  overflow: hidden;
  max-width: 100%;
  vertical-align: middle;
  margin: 1rem auto;
  border-radius: 4px;
}

.zoomable-image {
  display: block;
}

/* 让图片居中显示 */
.markdown-body p img,
.markdown-body p .image-container {
  display: block;
  margin: 1rem auto;
  max-width: 100%;
}
</style> 