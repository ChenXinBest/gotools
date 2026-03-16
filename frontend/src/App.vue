<script setup>
import { onMounted, watch } from 'vue'
import { RouterView } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import { useAppStore } from './stores/app'

const appStore = useAppStore()

onMounted(() => {
  appStore.init()
  // 应用初始主题
  document.body.className = `theme-${appStore.theme}`
})

// 监听主题变化
watch(() => appStore.theme, (newTheme) => {
  document.body.className = `theme-${newTheme}`
})
</script>

<template>
  <div class="app-container">
    <Sidebar />
    <div class="main-content">
      <RouterView />
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  background-color: var(--bg-primary);
  min-height: 100vh;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

#app {
  width: 100%;
  min-height: 100vh;
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--bg-primary);
}

::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

::selection {
  background: var(--accent-color);
  color: var(--bg-primary);
}
</style>

<style scoped>
.app-container {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
