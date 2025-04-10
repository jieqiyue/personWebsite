<script setup>
import Header from './components/Header.vue'
import BackToTop from './components/BackToTop.vue'
</script>

<template>
  <div class="app">
    <Header />
    <main>
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <BackToTop />
    <footer class="site-footer">
      <p>&copy; {{ new Date().getFullYear() }} 我的摄影博客. All rights reserved.</p>
    </footer>
  </div>
</template>

<style>
/* 修改应用容器，确保全宽 */
#app {
  width: 100vw; /* 使用视口宽度 */
  max-width: 100vw; /* 防止溢出 */
  margin: 0;
  padding: 0;
  background-color: var(--background);
  overflow-x: hidden; /* 防止水平滚动 */
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 100%;
  background-color: var(--background);
  margin: 0;
  padding: 0;
}

/* 内容区域样式 */
main {
  flex: 1;
  width: 100%;
  max-width: 100%;
  margin: 0 auto;
  padding: 0;
  background-color: var(--background);
}

/* 页脚样式 */
.site-footer {
  text-align: center;
  padding: 1rem 0;
  background: var(--border);
  width: 100%;
  border-top: 1px solid var(--border);
  color: var(--primary);
}

/* 防止内部元素溢出 */
* {
  box-sizing: border-box;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

body {
  margin: 0;
  padding: 0;
  background-color: var(--background);
}
</style>
