<template>
  <div class="theme-switcher">
    <button 
      class="theme-toggle-icon" 
      @click="cycleTheme" 
      title="切换主题"
    >
      <!-- Half-filled circle icon -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M12 2a10 10 0 1 0 10 10A10.011 10.011 0 0 0 12 2zm0 18v-16a8 8 0 0 1 0 16z"></path>
      </svg>
    </button>
    <!-- Dropdown removed -->
  </div>
</template>

<script setup>
import { /* ref, computed, onMounted, onUnmounted removed */ } from 'vue'
import { themes, currentTheme, changeTheme } from '../utils/theme'

// Function to cycle through themes
const cycleTheme = () => {
  const themeKeys = Object.keys(themes);
  const currentIndex = themeKeys.indexOf(currentTheme.value);
  const nextIndex = (currentIndex + 1) % themeKeys.length; // Cycle
  const nextThemeName = themeKeys[nextIndex];
  changeTheme(nextThemeName);
}

// onMounted/onUnmounted for handleClickOutside removed
</script>

<style scoped>
.theme-switcher {
  position: relative;
  /* Removed z-index as dropdown is gone, adjust if needed */
}

.theme-toggle-icon {
  background: none; /* Remove background */
  border: none; /* Remove border */
  padding: 8px; /* Adjust padding for icon size */
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text); /* Use text color for the icon */
  transition: color 0.3s ease, transform 0.2s ease;
  border-radius: 50%; /* Optional: make it circular */
  width: 36px; /* Adjust size */
  height: 36px; /* Adjust size */
}

.theme-toggle-icon:hover {
  color: var(--primary); /* Change color on hover */
  transform: scale(1.1); /* Optional: slight scale effect */
  background-color: rgba(128, 128, 128, 0.1); /* Optional: subtle background on hover */
}

.theme-toggle-icon svg {
  display: block; /* Ensure SVG behaves predictably */
}

/* Styles for dropdown, theme-label, theme-option, etc. removed */

</style> 