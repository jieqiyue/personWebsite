<template>
  <div class="comments-section">
    <div class="comments-header">
      <h2 class="comments-title">评论区</h2>
      <div v-if="commentCount !== null" class="comment-count">{{ commentCount }} 条讨论</div>
    </div>
    
    <div class="comments-container">
      <!-- 加载状态 -->
      <div v-if="loading" class="comments-loading">
        <div class="spinner"></div>
        <p>正在加载评论...</p>
      </div>
      
      <!-- 评论内容区 -->
      <div ref="giscusContainer" class="giscus-wrapper"></div>
      
      <!-- 无评论时的提示 -->
      <div v-if="noComments && !loading" class="no-comments">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-message-circle">
          <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path>
        </svg>
        <p>暂无评论，来说点什么吧~</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue';
import { useRoute } from 'vue-router'; // Optional: if theme needs to react to system changes

const props = defineProps({
  pageTerm: {
    type: String,
    required: true,
  },
});

const giscusContainer = ref(null);
const route = useRoute(); // Optional: for theme reactivity
const loading = ref(true);
const noComments = ref(false);
const commentCount = ref(null);

// 监听来自giscus的消息，获取评论数量
const handleGiscusMessage = (event) => {
  if (event.origin !== 'https://giscus.app') return;
  
  if (event.data?.giscus?.discussion) {
    loading.value = false;
    const discussion = event.data.giscus.discussion;
    commentCount.value = discussion.totalCommentCount || 0;
    noComments.value = commentCount.value === 0;
  }
};

// Function to load the Giscus script
const loadGiscus = (term) => {
  if (!giscusContainer.value) return;
  
  // 重置状态
  loading.value = true;
  noComments.value = false;
  commentCount.value = null;

  // 添加消息监听
  window.addEventListener('message', handleGiscusMessage);

  // Clear previous Giscus instance if exists
  while (giscusContainer.value.firstChild) {
    giscusContainer.value.removeChild(giscusContainer.value.firstChild);
  }

  const script = document.createElement('script');
  script.src = 'https://giscus.app/client.js';
  script.async = true;
  script.crossOrigin = 'anonymous';

  // --- YOUR GISCUS CONFIGURATION BELOW ---
  script.setAttribute('data-repo', 'jieqiyue/personWebsite'); // <-- Replace with your value if needed
  script.setAttribute('data-repo-id', 'R_kgDOOV5N5Q');      // <-- Replace with your value if needed
  script.setAttribute('data-category', 'Announcements');     // <-- Replace with your value if needed
  script.setAttribute('data-category-id', 'DIC_kwDOOV5N5c4Co_3Z'); // <-- Replace with your value if needed
  script.setAttribute('data-mapping', 'pathname');          // <-- Replace with your value if needed
  script.setAttribute('data-term', term); // Use the dynamic term passed via prop
  script.setAttribute('data-strict', '0');               // <-- Replace with your value if needed
  script.setAttribute('data-reactions-enabled', '1');  // <-- Replace with your value if needed
  script.setAttribute('data-emit-metadata', '1');        // <-- Replace with your value if needed
  script.setAttribute('data-input-position', 'top');    // <-- Replace with your value if needed
  script.setAttribute('data-theme', 'preferred_color_scheme'); // <-- Replace with your value if needed
  script.setAttribute('data-lang', 'zh-CN');             // <-- Replace with your value if needed
  // --- END OF YOUR GISCUS CONFIGURATION ---

  // 设置加载超时检测
  setTimeout(() => {
    if (loading.value) {
      loading.value = false;
      noComments.value = true;
    }
  }, 5000);

  giscusContainer.value.appendChild(script);
};

// Load Giscus when the component mounts
onMounted(() => {
  loadGiscus(props.pageTerm);
});

// Reload Giscus when the pageTerm prop changes (e.g., navigating between articles)
watch(() => props.pageTerm, (newTerm) => {
  // Use nextTick to ensure the container is ready after potential DOM updates
  nextTick(() => {
    loadGiscus(newTerm);
  });
});

// 组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('message', handleGiscusMessage);
});

// Optional: Watch for system theme changes if using preferred_color_scheme
// This might require more complex setup depending on how theme is managed globally
// watch(() => route.meta.theme, (newTheme) => { // Example: if theme is in route meta
//   if (giscusContainer.value && giscusContainer.value.querySelector('iframe.giscus-frame')) {
//      const iframe = giscusContainer.value.querySelector('iframe.giscus-frame');
//      iframe.contentWindow.postMessage({ giscus: { setConfig: { theme: newTheme } } }, 'https://giscus.app');
//   }
// });

</script>

<style scoped>
.comments-section {
  margin-top: 3rem;
  padding: 1.5rem;
  background-color: var(--surface, #ffffff);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.comments-section:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
}

.comments-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border, #eee);
}

.comments-title {
  font-size: 1.5em;
  font-weight: 600;
  color: var(--text, #333);
  margin: 0;
}

.comment-count {
  font-size: 0.95rem;
  color: var(--text-light, #666);
  background-color: var(--accent, #f5f5f5);
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-weight: 500;
}

.comments-container {
  position: relative;
  min-height: 150px;
}

.comments-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2.5rem 0;
}

.spinner {
  width: 30px;
  height: 30px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--primary, #2c6e49);
  border-radius: 50%;
  animation: spin 1s ease-in-out infinite;
  margin-bottom: 0.75rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.comments-loading p {
  color: var(--text-light, #666);
  font-size: 0.95rem;
  margin: 0;
}

.no-comments {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  color: var(--text-light, #999);
  text-align: center;
}

.no-comments svg {
  color: var(--text-light, #999);
  margin-bottom: 1rem;
  opacity: 0.5;
}

.no-comments p {
  font-size: 1rem;
  margin: 0;
}

.giscus-wrapper {
  width: 100%;
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .comments-section {
    background-color: var(--surface, #222);
  }
  
  .comment-count {
    background-color: var(--accent, #333);
    color: var(--text-light, #aaa);
  }
  
  .spinner {
    border-color: rgba(255, 255, 255, 0.1);
    border-top-color: var(--primary, #4c956c);
  }
  
  .no-comments,
  .no-comments svg,
  .comments-loading p {
    color: var(--text-light, #777);
  }
}
</style> 