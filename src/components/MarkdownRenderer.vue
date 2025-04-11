<template>
  <div class="markdown-body" v-html="renderedContent" ref="markdownBody"></div>
</template>

<script>
import { computed, ref, onMounted, watch, onUnmounted } from 'vue'
import md from '../plugins/markdown.js'
import 'highlight.js/styles/github.css'
import 'github-markdown-css/github-markdown.css'
import './markdown-custom.css'
import ClipboardJS from 'clipboard'

export default {
  name: 'MarkdownRenderer',
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
    
    // 监听内容变化，重新设置复制按钮
    watch(() => props.content, () => {
      // 等待DOM更新后设置按钮
      setTimeout(setupCodeCopyButtons, 0);
    });
    
    // 组件挂载后设置复制按钮
    onMounted(() => {
      setupCodeCopyButtons();
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
      markdownBody
    }
  }
}
</script>

<style>
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
</style> 