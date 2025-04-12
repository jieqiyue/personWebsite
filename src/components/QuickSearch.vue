<template>
  <div class="quick-search">
    <div class="search-icon" @click="toggleSearch">
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"></circle>
        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
      </svg>
    </div>
    
    <div v-if="isOpen" class="search-popup-overlay" @click="closeSearch">
      <div class="search-popup" @click.stop>
        <div class="search-input-container">
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="搜索文章..." 
            ref="searchInput"
            @keydown.esc="closeSearch"
          />
          <button class="clear-button" @click="clearSearch" v-if="searchQuery">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        
        <div class="search-results-container">
          <!-- 加载状态 -->
          <div v-if="loading" class="search-loading">
            <div class="spinner"></div>
            <span>搜索中...</span>
          </div>
          
          <!-- 搜索结果 -->
          <div v-else-if="results.length > 0" class="search-results">
            <div 
              v-for="result in results.slice(0, 6)" 
              :key="result.id" 
              class="search-result"
              @click="navigateToArticle(result.id)"
            >
              <h3 class="result-title" v-html="highlightMatch(result.title)"></h3>
              <p class="result-excerpt" v-html="highlightMatch(result.searchExcerpt)"></p>
              <div class="result-tags">
                <span 
                  v-for="tag in result.tags.slice(0, 3)" 
                  :key="tag" 
                  class="result-tag"
                >
                  #{{ tag }}
                </span>
              </div>
            </div>
            
            <div v-if="results.length > 6" class="view-all">
              <router-link to="/blog" @click="closeSearch">查看全部 {{ results.length }} 个结果</router-link>
            </div>
          </div>
          
          <!-- 无结果状态 -->
          <div v-else-if="searchQuery && !loading" class="no-results">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
              <line x1="11" y1="8" x2="11" y2="14"></line>
              <line x1="8" y1="11" x2="14" y2="11"></line>
            </svg>
            <p>没有找到匹配的文章</p>
          </div>
          
          <!-- 初始状态：热门标签 -->
          <div v-else-if="!searchQuery" class="popular-tags">
            <h4>热门标签</h4>
            <div class="tags-container">
              <div 
                v-for="tag in popularTags" 
                :key="tag" 
                class="tag"
                @click="searchByTag(tag)"
              >
                #{{ tag }}
              </div>
            </div>
          </div>
        </div>
        
        <div class="search-footer">
          <span class="keyboard-hint">按 ESC 关闭</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import { searchArticles, buildSearchIndex } from '../utils/search';

const router = useRouter();
const isOpen = ref(false);
const searchQuery = ref('');
const searchInput = ref(null);
const results = ref([]);
const loading = ref(false);
const popularTags = ref([]);

// 打开搜索框
const toggleSearch = () => {
  isOpen.value = !isOpen.value;
  
  if (isOpen.value) {
    nextTick(() => {
      searchInput.value.focus();
      fetchPopularTags();
    });
  }
};

// 关闭搜索框
const closeSearch = () => {
  isOpen.value = false;
};

// 清空搜索内容
const clearSearch = () => {
  searchQuery.value = '';
  results.value = [];
  nextTick(() => {
    searchInput.value.focus();
  });
};

// 跳转到文章
const navigateToArticle = (id) => {
  router.push(`/articles/${id}`);
  closeSearch();
};

// 按标签搜索
const searchByTag = (tag) => {
  searchQuery.value = `#${tag}`;
  performSearch();
};

// 执行搜索
const performSearch = async () => {
  if (!searchQuery.value.trim()) {
    results.value = [];
    return;
  }
  
  loading.value = true;
  
  // 检查是否是按标签搜索 (以#开头)
  const isTagSearch = searchQuery.value.startsWith('#');
  const query = isTagSearch ? searchQuery.value.substring(1).trim() : searchQuery.value;
  const tags = isTagSearch ? [query] : [];
  
  try {
    results.value = await searchArticles(isTagSearch ? '' : query, {
      tags: tags,
      limit: 10
    });
  } catch (error) {
    console.error('搜索失败:', error);
  } finally {
    loading.value = false;
  }
};

// 高亮匹配的搜索词
const highlightMatch = (text) => {
  if (!text || !searchQuery.value.trim()) return text;
  
  // 如果是标签搜索，高亮标签而非文本
  if (searchQuery.value.startsWith('#')) {
    const tag = searchQuery.value.substring(1).trim().toLowerCase();
    if (tag) {
      return text.replace(
        new RegExp(`(#?${tag})`, 'gi'), 
        '<mark>$1</mark>'
      );
    }
    return text;
  }
  
  const searchTerms = searchQuery.value
    .toLowerCase()
    .trim()
    .split(/\s+/)
    .filter(term => term.length > 1);
  
  let highlightedText = text;
  
  searchTerms.forEach(term => {
    const escapedTerm = term.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    highlightedText = highlightedText.replace(
      new RegExp(`(${escapedTerm})`, 'gi'), 
      '<mark>$1</mark>'
    );
  });
  
  return highlightedText;
};

// 获取热门标签
const fetchPopularTags = async () => {
  try {
    const response = await fetch('/markdown/articles/index.json');
    const data = await response.json();
    
    // 统计标签出现频率
    const tagCount = {};
    data.articles.forEach(article => {
      if (article.tags) {
        article.tags.forEach(tag => {
          tagCount[tag] = (tagCount[tag] || 0) + 1;
        });
      }
    });
    
    // 按出现频率排序并取前8个
    popularTags.value = Object.entries(tagCount)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 8)
      .map(entry => entry[0]);
      
  } catch (error) {
    console.error('获取热门标签失败:', error);
  }
};

// 监听搜索查询变化，自动搜索
let searchTimeout;
watch(searchQuery, (newQuery) => {
  clearTimeout(searchTimeout);
  
  if (!newQuery.trim()) {
    results.value = [];
    return;
  }
  
  // 防抖：延迟300ms执行搜索，避免频繁请求
  searchTimeout = setTimeout(() => {
    performSearch();
  }, 300);
});

// 监听ESC键关闭搜索框
const handleKeyDown = (e) => {
  if (e.key === 'Escape' && isOpen.value) {
    closeSearch();
  } else if (e.key === '/' && !isOpen.value && 
            !(e.target instanceof HTMLInputElement) && 
            !(e.target instanceof HTMLTextAreaElement)) {
    // 按"/"键打开搜索框(当未在输入框中时)
    e.preventDefault();
    toggleSearch();
  }
};

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
  
  // 预加载搜索索引
  buildSearchIndex();
});

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeyDown);
});
</script>

<style scoped>
.quick-search {
  position: relative;
}

.search-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: var(--accent, #f5f5f5);
  margin-right: 10px;
  color: var(--text, #333);
  cursor: pointer;
  transition: all 0.3s ease;
}

.search-icon:hover {
  transform: translateY(-2px);
  background-color: var(--primary, #42b883);
  color: white;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.search-icon svg {
  transition: transform 0.3s ease;
}

.search-icon:hover svg {
  transform: rotate(-15deg);
}

/* 搜索弹窗 */
.search-popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.search-popup {
  position: relative;
  width: 90%;
  max-width: 600px;
  background-color: var(--background, #fff);
  margin-top: 80px;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  max-height: 80vh;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from { transform: translateY(-20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.search-input-container {
  position: relative;
  padding: 16px;
  border-bottom: 1px solid var(--border, #eee);
}

.search-input-container input {
  width: 100%;
  padding: 12px 40px 12px 16px;
  border: 2px solid var(--border, #ddd);
  border-radius: 8px;
  font-size: 16px;
  outline: none;
  background-color: var(--background, #fff);
  color: var(--text, #333);
}

.search-input-container input:focus {
  border-color: var(--primary, #42b883);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.1);
}

.clear-button {
  position: absolute;
  right: 28px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: var(--text-light, #999);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  border-radius: 50%;
}

.clear-button:hover {
  color: var(--text, #333);
  background-color: var(--accent, #f5f5f5);
}

.search-results-container {
  overflow-y: auto;
  padding: 16px;
  flex: 1;
  max-height: 60vh;
}

/* 搜索结果样式 */
.search-results {
  display: flex;
  flex-direction: column;
}

.search-result {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.search-result:hover {
  background-color: var(--accent, #f5f5f5);
  transform: translateY(-2px);
}

.result-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--text, #333);
}

.result-excerpt {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: var(--text-light, #666);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.result-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.result-tag {
  font-size: 12px;
  color: var(--primary, #42b883);
}

mark {
  background-color: rgba(255, 230, 0, 0.3);
  color: inherit;
  font-weight: 600;
  padding: 0 2px;
  border-radius: 2px;
}

.view-all {
  text-align: center;
  padding: 16px;
  border-top: 1px solid var(--border, #eee);
}

.view-all a {
  color: var(--primary, #42b883);
  text-decoration: none;
  font-weight: 500;
  font-size: 14px;
}

.view-all a:hover {
  text-decoration: underline;
}

/* 无结果状态 */
.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: var(--text-light, #666);
  text-align: center;
}

.no-results svg {
  color: var(--text-light, #999);
  margin-bottom: 16px;
  opacity: 0.5;
}

.no-results p {
  font-size: 16px;
  margin: 0;
}

/* 加载状态 */
.search-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: var(--text-light, #666);
}

.spinner {
  width: 30px;
  height: 30px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--primary, #42b883);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 热门标签样式 */
.popular-tags {
  padding: 8px 0;
}

.popular-tags h4 {
  font-size: 14px;
  color: var(--text-light, #666);
  margin: 0 0 16px 0;
  font-weight: 500;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.tag {
  padding: 6px 12px;
  background-color: var(--accent, #f5f5f5);
  border-radius: 16px;
  font-size: 14px;
  color: var(--text, #333);
  cursor: pointer;
  transition: all 0.2s ease;
}

.tag:hover {
  background-color: var(--primary, #42b883);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* 页脚 */
.search-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border, #eee);
  font-size: 12px;
  color: var(--text-light, #999);
  text-align: center;
}

.keyboard-hint {
  display: inline-block;
  padding: 2px 6px;
  background-color: var(--accent, #f5f5f5);
  border-radius: 4px;
  font-family: monospace;
}

@media (max-width: 768px) {
  .search-popup {
    width: 95%;
    margin-top: 60px;
  }

  .search-results-container {
    max-height: 50vh;
  }
}
</style> 