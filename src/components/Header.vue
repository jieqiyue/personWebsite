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
        <ThemeSwitcher />
      </div>
    </div>
  </header>
</template>

<script setup>
import ThemeSwitcher from './ThemeSwitcher.vue'
</script>

<style scoped>
.site-header {
  padding: 1rem 2rem;
  background: var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 10px var(--shadow);
  width: 100%;
  border-bottom: 1px solid var(--border);
}

.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1280px;
  margin: 0 auto;
}

nav {
  display: flex;
  gap: 2rem;
}

.header-right {
  display: flex;
  align-items: center;
  margin-left: 1rem;
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
</style> 