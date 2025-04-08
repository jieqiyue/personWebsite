<template>
  <header class="site-header">
    <nav>
      <router-link to="/">首页</router-link>
      <router-link to="/portfolio">作品集</router-link>
      <router-link to="/blog">博客</router-link>
      <router-link to="/about">关于</router-link>
    </nav>
  </header>
</template>

<style scoped>
.site-header {
  padding: 1rem;
  background: #daebd9; /* 比背景深一点的绿色 */
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
  width: 100%;
  border-bottom: 1px solid #c8e6c7; /* 添加一个微妙的边框 */
}

nav {
  display: flex;
  gap: 2rem;
  justify-content: center;
  max-width: 1280px;
  margin: 0 auto;
}

a {
  color: #2c6e49; /* 更深的绿色作为链接颜色 */
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
}

a:hover {
  color: #1a4731; /* 悬停时更深的绿色 */
}

a.router-link-active {
  color: #1a4731; /* 活跃链接使用更深的绿色 */
  font-weight: 600;
}
</style> 