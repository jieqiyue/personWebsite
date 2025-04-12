<template>
  <header class="site-header">
    <div class="header-container">
      <nav>
        <router-link to="/">首页</router-link>
        <router-link to="/portfolio">作品集</router-link>
        <router-link to="/blog">博客</router-link>
        <router-link to="/about">关于</router-link>
      </nav>
      <div class="header-right">
        <QuickSearch />
        <ThemeSwitcher />
      </div>
    </div>
  </header>
</template>

<script setup>
import ThemeSwitcher from './ThemeSwitcher.vue'
import QuickSearch from './QuickSearch.vue'
</script>

<style scoped>
.site-header {
  padding: 1rem 2rem;
  background-color: var(--background-alt);
  color: var(--text);
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
  width: 100%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  max-width: 1280px;
  margin: 0 auto;
}

nav {
  display: flex;
  gap: 3rem;
  user-select: none; 
}

.header-right {
  display: flex;
  align-items: center;
}

a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
}

a:hover {
  color: var(--secondary);
}

a.router-link-active {
  color: var(--secondary);
  font-weight: 600;
}

@media (max-width: 768px) {
  .site-header {
    padding: 1rem;
  }
  
  nav {
    gap: 1rem;
  }
}

.logo a {
  font-size: 1.5rem;
  font-weight: bold;
  text-decoration: none;
  color: var(--primary);
}

nav ul {
  list-style: none;
  display: flex;
  gap: 1.5rem;
  margin: 0;
  padding: 0;
}

nav a {
  text-decoration: none;
  color: var(--text);
  font-weight: 500;
  transition: color 0.3s ease;
}

nav a:hover,
nav a.router-link-exact-active {
  color: var(--primary);
}
</style> 