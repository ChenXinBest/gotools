<script setup>
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import { useAppStore } from './stores/app'

const appStore = useAppStore()

onMounted(() => {
  appStore.init()
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
  background-color: #0a0a0a;
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
  background: #0a0a0a;
}

::-webkit-scrollbar-thumb {
  background: #1a1a1a;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #333;
}

::selection {
  background: #00ff00;
  color: #0a0a0a;
}
</style>

<style scoped>
.app-container {
  display: flex;
  min-height: 100vh;
  background: #0a0a0a;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
