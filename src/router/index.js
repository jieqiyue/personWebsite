import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Portfolio from '../views/Portfolio.vue'
import Blog from '../views/Blog.vue'
import ArticleDetail from '../views/ArticleDetail.vue'
import About from '../views/About.vue'
import PhotoDetail from '../views/PhotoDetail.vue'
import TagView from '../views/TagView.vue'
import FontTestViewer from '../views/FontTestViewer.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/portfolio',
    name: 'Portfolio',
    component: Portfolio
  },
  {
    path: '/photo/:id',
    name: 'PhotoDetail',
    component: PhotoDetail
  },
  {
    path: '/blog',
    name: 'Blog',
    component: Blog
  },
  {
    path: '/articles/:id',
    name: 'ArticleDetail',
    component: ArticleDetail
  },
  {
    path: '/tags/:tag',
    name: 'TagView',
    component: TagView
  },
  {
    path: '/about',
    name: 'About',
    component: About
  },
  {
    path: '/font-test',
    name: 'FontTest',
    component: FontTestViewer
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // 页面切换时滚动到顶部
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

export default router 