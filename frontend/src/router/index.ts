import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: {title: '媒体预览'}
  },
  {
    path: '/similar',
    name: 'Similar',
    component: () => import('../views/Similar.vue'),
    meta: {title: '相似图片分析'}
  },
  {
    path: '/classify',
    name: 'Classify',
    component: () => import('../views/Classify.vue'),
    meta: {title: '快捷分类'}
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫：更新页面标题
router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) || 'Media'
  next()
})

export default router
