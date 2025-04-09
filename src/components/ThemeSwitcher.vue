<template>
  <div class="theme-switcher">
    <button 
      class="theme-toggle" 
      @click="toggleDropdown" 
      :title="currentThemeName"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="5"></circle>
        <path d="M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4"></path>
      </svg>
      <span class="theme-label">{{ currentThemeName }}</span>
    </button>
    <div class="theme-dropdown" v-if="isOpen">
      <div
        v-for="(theme, key) in themes"
        :key="key"
        class="theme-option"
        :class="{ active: currentTheme === key }"
        @click="switchTheme(key)"
      >
        <span class="theme-color" :style="{ background: theme.colors.primary }"></span>
        <span class="theme-name">{{ theme.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { themes, currentTheme, changeTheme } from '../utils/theme'

const isOpen = ref(false)

// 获取当前主题名称
const currentThemeName = computed(() => {
  return themes[currentTheme.value]?.name || '主题'
})

// 切换下拉菜单
const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

// 切换主题
const switchTheme = (themeName) => {
  changeTheme(themeName)
  isOpen.value = false
}

// 点击外部关闭下拉菜单
const handleClickOutside = (event) => {
  const switcher = document.querySelector('.theme-switcher')
  if (switcher && !switcher.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.theme-switcher {
  position: relative;
  z-index: 1000;
}

.theme-toggle {
  background: var(--accent);
  border: 1px solid var(--border);
  border-radius: 20px;
  min-width: 44px;
  height: 44px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--primary);
  transition: all 0.3s ease;
  box-shadow: 0 2px 5px var(--shadow);
}

.theme-toggle:hover {
  background: var(--surface);
  box-shadow: 0 3px 8px var(--shadow);
}

.theme-label {
  margin-left: 8px;
  font-size: 14px;
  display: none;
}

@media (min-width: 768px) {
  .theme-label {
    display: inline;
  }
  
  .theme-toggle {
    padding: 0 16px;
  }
}

.theme-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  width: 200px;
  background: var(--surface);
  border-radius: 12px;
  box-shadow: 0 5px 15px var(--shadow);
  margin-top: 12px;
  overflow: hidden;
  animation: fadeIn 0.2s ease;
  border: 1px solid var(--border);
}

.theme-option {
  padding: 14px 18px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: background 0.2s ease;
}

.theme-option:hover {
  background: var(--accent);
}

.theme-option.active {
  background: var(--accent);
  font-weight: 500;
}

.theme-color {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  margin-right: 12px;
  border: 2px solid var(--surface);
  box-shadow: 0 0 0 1px var(--border);
}

.theme-name {
  color: var(--text);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style> 