import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home,
    meta: { title: '首页' }
  },
  {
    path: '/process-manager',
    name: 'process-manager',
    component: () => import('../views/ProcessManager.vue'),
    meta: { title: '进程管理器' }
  },
  {
    path: '/database-backup',
    name: 'database-backup',
    component: () => import('../views/DatabaseBackup.vue'),
    meta: { title: '数据库备份' },
    children: [
      {
        path: 'connections',
        name: 'database-connections',
        component: () => import('../views/database/Connections.vue'),
        meta: { title: '连接管理' }
      },
      {
        path: 'export',
        name: 'database-export',
        component: () => import('../views/database/Export.vue'),
        meta: { title: '导出数据' }
      },
      {
        path: 'import',
        name: 'database-import',
        component: () => import('../views/database/Import.vue'),
        meta: { title: '导入数据' }
      }
    ]
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('../views/Settings.vue'),
    meta: { title: '设置' }
  },
  {
    path: '/export-progress',
    name: 'export-progress',
    component: () => import('../views/ExportProgress.vue'),
    meta: { title: '导出进度' }
  },
  {
    path: '/import-conflicts',
    name: 'import-conflicts',
    component: () => import('../views/ImportConflicts.vue'),
    meta: { title: '导入冲突' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../views/NotFound.vue'),
    meta: { title: '页面未找到' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局前置守卫 - 设置页面标题
router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - GoTools` : 'GoTools'
  next()
})

export default router