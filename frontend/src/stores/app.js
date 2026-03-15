import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 状态
  const currentTool = ref('process-manager')
  const sidebarCollapsed = ref(false)
  const theme = ref('dark')
  const locale = ref('zh-CN')
  const loading = ref(false)
  const notifications = ref([])

  // 计算属性
  const isDarkTheme = computed(() => theme.value === 'dark')
  const hasNotifications = computed(() => notifications.value.length > 0)

  // 操作
  function setCurrentTool(tool) {
    currentTool.value = tool
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTheme(newTheme) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
  }

  function setLocale(newLocale) {
    locale.value = newLocale
    localStorage.setItem('locale', newLocale)
  }

  function setLoading(value) {
    loading.value = value
  }

  function addNotification(notification) {
    const id = Date.now()
    notifications.value.push({
      id,
      ...notification,
      timestamp: new Date()
    })

    // 自动移除通知
    if (notification.duration !== 0) {
      setTimeout(() => {
        removeNotification(id)
      }, notification.duration || 5000)
    }
  }

  function removeNotification(id) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  function clearNotifications() {
    notifications.value = []
  }

  // 初始化
  function init() {
    // 从localStorage恢复设置
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
      theme.value = savedTheme
    }

    const savedLocale = localStorage.getItem('locale')
    if (savedLocale) {
      locale.value = savedLocale
    }
  }

  return {
    // 状态
    currentTool,
    sidebarCollapsed,
    theme,
    locale,
    loading,
    notifications,
    // 计算属性
    isDarkTheme,
    hasNotifications,
    // 操作
    setCurrentTool,
    toggleSidebar,
    setTheme,
    setLocale,
    setLoading,
    addNotification,
    removeNotification,
    clearNotifications,
    init
  }
})