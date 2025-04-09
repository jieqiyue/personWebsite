import { ref } from 'vue'

// 主题配置
export const themes = {
  nature: {
    name: '自然',
    colors: {
      primary: '#2c6e49',
      secondary: '#4c956c',
      accent: '#fefae0',
      background: '#f0f7f0',
      surface: '#ffffff',
      text: '#333333',
      border: '#daebd9',
      shadow: 'rgba(0, 0, 0, 0.05)'
    }
  },
  elegant: {
    name: '优雅',
    colors: {
      primary: '#54423A',
      secondary: '#9C8178',
      accent: '#F3EDE2',
      background: '#F6F1EB',
      surface: '#FFFFFF',
      text: '#403836',
      border: '#E5DFD9',
      shadow: 'rgba(84, 66, 58, 0.08)'
    }
  },
  ocean: {
    name: '海洋',
    colors: {
      primary: '#0A6E8F',
      secondary: '#0B96BF',
      accent: '#F0F9FC',
      background: '#EBF6FA',
      surface: '#FFFFFF',
      text: '#2D4654',
      border: '#CCE5EF',
      shadow: 'rgba(10, 110, 143, 0.08)'
    }
  }
}

// 当前主题
export const currentTheme = ref('nature')

// 切换主题
export function changeTheme(themeName) {
  if (themes[themeName]) {
    currentTheme.value = themeName
    applyTheme(themeName)
    // 保存到本地存储，以便下次访问时使用相同主题
    localStorage.setItem('selectedTheme', themeName)
  }
}

// 应用主题到CSS变量
export function applyTheme(themeName) {
  const theme = themes[themeName]
  if (!theme) return

  const root = document.documentElement
  
  // 设置CSS变量
  Object.entries(theme.colors).forEach(([key, value]) => {
    root.style.setProperty(`--${key}`, value)
  })
}

// 初始化主题
export function initTheme() {
  // 从本地存储获取之前选择的主题
  const savedTheme = localStorage.getItem('selectedTheme')
  
  // 如果有保存的主题且主题存在，则使用保存的主题
  if (savedTheme && themes[savedTheme]) {
    currentTheme.value = savedTheme
  }
  
  // 应用主题
  applyTheme(currentTheme.value)
} 