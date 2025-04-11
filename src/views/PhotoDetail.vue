<template>
  <div class="photo-detail">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>正在加载照片信息...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h2>{{ error }}</h2>
      <button class="retry-button" @click="fetchPhotoData(route.params.id)">
        重试
      </button>
      <button class="back-button" @click="router.back()">
        返回上一页
      </button>
    </div>

    <!-- 数据加载成功 -->
    <template v-else-if="photo">
      <!-- 返回按钮独立放置 -->
      <div class="back-navigation">
        <button class="back-button" @click="router.back()">
          <span>&larr;</span> 返回
        </button>
      </div>
      
      <div class="content-container">
        <div class="photo-container">
          <img 
            :src="photo.url" 
            :alt="photo.title"
            @click="toggleFullscreen"
            :class="{ 'fullscreen': isFullscreen }"
            @load="handleImageLoad"
            v-show="imageLoaded"
          />
          <!-- 图片加载时的占位 -->
          <div v-if="!imageLoaded" class="image-placeholder">
            <div class="loading-spinner"></div>
          </div>
        </div>
        
        <div class="photo-info">
          <h1>{{ photo.title }}</h1>
          
          <div class="metadata">
            <div class="info-group">
              <h3>拍摄信息</h3>
              <div class="info-grid">
                <div class="info-item" v-if="photo.dateTaken">
                  <span class="label">拍摄时间</span>
                  <span class="value">{{ formatDate(photo.dateTaken) }}</span>
                </div>
                <div class="info-item" v-if="photo.camera">
                  <span class="label">相机</span>
                  <span class="value">{{ photo.camera }}</span>
                </div>
                <div class="info-item" v-if="photo.lens">
                  <span class="label">镜头</span>
                  <span class="value">{{ photo.lens }}</span>
                </div>
                <div class="info-item" v-if="photo.aperture">
                  <span class="label">光圈</span>
                  <span class="value">{{ photo.aperture }}</span>
                </div>
                <div class="info-item" v-if="photo.shutter">
                  <span class="label">快门</span>
                  <span class="value">{{ photo.shutter }}</span>
                </div>
                <div class="info-item" v-if="photo.iso">
                  <span class="label">ISO</span>
                  <span class="value">{{ photo.iso }}</span>
                </div>
              </div>
            </div>
            
            <div class="info-group" v-if="photo.categories?.length">
              <h3>分类</h3>
              <div class="categories">
                <span 
                  v-for="category in photo.categories" 
                  :key="category"
                  class="category-tag"
                >
                  {{ category }}
                </span>
              </div>
            </div>
            
            <div class="info-group" v-if="photo.description">
              <h3>作品描述</h3>
              <p class="description">{{ photo.description }}</p>
            </div>

            <div class="info-group" v-if="photo.location">
              <h3>拍摄地点</h3>
              <p class="location">{{ photo.location }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Wrap the Comments component -->
      <div class="comments-wrapper">
        <Comments :pageTerm="photoPath" />
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BackToTop from '../components/BackToTop.vue'
import Comments from '../components/Comments.vue'

const route = useRoute()
const router = useRouter()
const photo = ref(null)
const loading = ref(true)
const error = ref(null)
const isFullscreen = ref(false)
const imageLoaded = ref(false)
const photoPath = ref('')

// 模拟API调用
const fetchPhotoData = async (id) => {
  try {
    loading.value = true
    error.value = null
    imageLoaded.value = false
    photoPath.value = route.path
    
    // 模拟API延迟
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 模拟数据
    const mockData = {
      id: id,
      title: '山间晨雾',
      url: 'https://images.unsplash.com/photo-1506744038136-46273834b3fb',
      description: '这张照片拍摄于清晨，山间雾气缭绕，展现了大自然的神秘与壮美。',
      dateTaken: '2024-04-01 06:30',
      camera: 'Sony A7R IV',
      lens: 'FE 16-35mm f/2.8 GM',
      aperture: 'f/11',
      shutter: '1/60s',
      iso: '100',
      categories: ['风景', '自然', '晨景'],
      location: '黄山',
    }
    
    photo.value = mockData
  } catch (err) {
    error.value = '获取照片信息失败，请稍后重试'
    console.error('Error fetching photo:', err)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString) => {
  try {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date)
  } catch (err) {
    return dateString
  }
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  // 处理全屏时的滚动
  document.body.style.overflow = isFullscreen.value ? 'hidden' : ''
}

const handleImageLoad = () => {
  imageLoaded.value = true
}

// 组件卸载时清理
onUnmounted(() => {
  document.body.style.overflow = ''
  isFullscreen.value = false
})

// 监听路由参数变化
watch(() => route.params.id, (newId) => {
  if (newId) {
    fetchPhotoData(newId)
  }
})

onMounted(async () => {
  const photoId = route.params.id
  if (photoId) {
    await fetchPhotoData(photoId)
  }

  // 添加样式覆盖
  const style = document.createElement('style')
  style.id = 'photo-detail-style-override'
  style.textContent = `
    .site-footer {
      background: transparent !important;
      padding: 0.5rem 0 !important;
      margin: 0 !important;
    }
  `
  document.head.appendChild(style)
})

onUnmounted(() => {
  // 清理样式覆盖
  const style = document.getElementById('photo-detail-style-override')
  if (style) {
    document.head.removeChild(style)
  }
})
</script>

<style scoped>
.photo-detail {
  padding: 2rem 2rem 0 2rem;
  max-width: 1400px;
  margin: 0 auto;
  background-color: var(--background); /* 使用主题变量 */
}

/* 返回按钮独立容器 */
.back-navigation {
  margin-bottom: 1.5rem;
  position: relative;
  z-index: 10;
}

.back-button {
  padding: 0.5rem 1rem;
  background: var(--surface); /* 使用主题变量 */
  border: 1px solid var(--border); /* 添加边框增强可见性 */
  border-radius: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px var(--shadow); /* 使用主题变量 */
  color: var(--primary); /* 使用主题变量使文字颜色更加深色 */
  font-weight: 500; /* 加粗字体 */
}

.back-button:hover {
  background: var(--accent); /* 使用主题变量 */
  transform: translateX(-5px);
}

/* 内容容器 */
.content-container {
  display: grid;
  grid-template-columns: 3fr 2fr; /* 图片在左(3份)，介绍在右(2份) */
  gap: 2rem;
  margin-bottom: 0; /* 移除底部外边距 */
  background: transparent; /* 确保无背景 */
}

/* 加载状态样式 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 3px solid var(--border); /* 使用主题变量 */
  border-top: 3px solid var(--primary); /* 使用主题变量 */
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 错误状态样式 */
.error-state {
  text-align: center;
  padding: 3rem;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.retry-button {
  padding: 0.8rem 1.5rem;
  margin: 1rem;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: var(--primary); /* 使用主题变量 */
  color: white;
}

.retry-button:hover {
  background: var(--secondary); /* 使用主题变量 */
}

/* 图片容器样式 */
.photo-container {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  margin-top: 2.5rem; /* 图片下移，与标题下边框对齐 */
  margin-bottom: 0; /* 移除底部外边距 */
  background: none; /* 移除背景色 */
  box-shadow: none; /* 移除阴影 */
}

.photo-container img {
  width: 100%;
  height: auto;
  display: block;
  border-radius: 8px;
  cursor: zoom-in;
  transition: all 0.3s ease;
  margin-bottom: 0; /* 确保图片底部没有外边距 */
}

.photo-container img.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  object-fit: contain;
  background: rgba(0, 0, 0, 0.9);
  z-index: 1000;
  border-radius: 0;
  cursor: zoom-out;
}

.image-placeholder {
  aspect-ratio: 3/2;
  display: flex;
  align-items: center;
  justify-content: center;
  background: none; /* 移除背景色 */
}

/* 照片信息样式 */
.photo-info {
  background: var(--surface); /* 使用主题变量 */
  border-radius: 8px;
  box-shadow: 0 4px 15px var(--shadow); /* 使用主题变量 */
  padding: 1.5rem;
  align-self: start;
  margin-bottom: 0; /* 移除底部外边距 */
}

.photo-info h1 {
  margin: 0 0 1.5rem 0;
  font-size: 1.8rem;
  color: var(--text); /* 使用主题变量 */
  border-bottom: 2px solid var(--border); /* 使用主题变量 */
  padding-bottom: 0.8rem;
}

.info-group {
  margin-bottom: 1.5rem;
}

/* 最后一个info-group不需要底部边距 */
.info-group:last-child {
  margin-bottom: 0;
}

.info-group h3 {
  margin-bottom: 0.8rem;
  color: var(--primary); /* 使用主题变量 */
  font-size: 1.1rem;
  font-weight: 600;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 1rem;
  background: var(--accent); /* 使用主题变量 */
  padding: 1rem;
  border-radius: 6px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.info-item .label {
  color: var(--text); /* 使用主题变量 */
  opacity: 0.7; /* 降低不透明度使其变淡 */
  font-size: 0.85rem;
}

.info-item .value {
  font-weight: 500;
  color: var(--text); /* 使用主题变量 */
}

.categories {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.category-tag {
  padding: 0.3rem 1rem;
  background: var(--accent); /* 使用主题变量 */
  border-radius: 15px;
  font-size: 0.9rem;
}

.description {
  line-height: 1.6;
  color: var(--text); /* 使用主题变量 */
  background: var(--accent); /* 使用主题变量 */
  padding: 1rem;
  border-radius: 6px;
}

.location {
  color: var(--text); /* 使用主题变量 */
  opacity: 0.8; /* 降低不透明度使其变淡 */
}

/* 响应式布局 */
@media (max-width: 768px) {
  .content-container {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
  
  .photo-container {
    margin-top: 0;
  }
  
  .photo-detail {
    padding: 1rem 1rem 0 1rem;
  }
  
  .photo-info {
    padding: 1rem;
  }
}

/* Wrapper for the comments section */
.comments-wrapper {
  max-width: 800px; /* Set to 800px to match ArticleDetail */
  margin: 0 auto;  /* Center the wrapper */
  /* The top margin is handled by Comments.vue's internal .comments-section */
  padding-bottom: 2rem; /* Add some padding at the very bottom */
}

/* Ensure the site footer still behaves correctly */
/* Styles added via JS override this, but good practice */
::v-deep(.site-footer) {
  background: transparent !important;
  padding: 0.5rem 0 !important;
  margin: 0 !important;
  border-top: none !important; /* Ensure no extra borders appear */
  box-shadow: none !important; /* Ensure no extra shadows appear */
}
</style> 