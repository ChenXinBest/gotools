import { createPinia } from 'pinia'

export const pinia = createPinia()

// 导出所有stores
export { useAppStore } from './app'
export { useProcessStore } from './process'
export { useDatabaseStore } from './database'