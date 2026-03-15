<template>
  <div class="home">
    <div class="welcome-section">
      <h1>欢迎使用 GoTools</h1>
      <p class="subtitle">一个强大的开发工具集合</p>
    </div>

    <div class="tools-grid">
      <div
        v-for="tool in tools"
        :key="tool.id"
        class="tool-card"
        @click="navigateTo(tool.route)"
      >
        <div class="tool-icon">{{ tool.icon }}</div>
        <h3 class="tool-name">{{ tool.name }}</h3>
        <p class="tool-description">{{ tool.description }}</p>
      </div>
    </div>

    <div class="features-section">
      <h2>主要特点</h2>
      <ul class="features-list">
        <li>基于 Go 和 Vue 3 构建，性能优异</li>
        <li>支持进程管理和系统监控</li>
        <li>MySQL 数据库备份和恢复</li>
        <li>跨平台支持 (Windows/macOS/Linux)</li>
        <li>现代化的用户界面</li>
      </ul>
    </div>

    <div class="version-info">
      <p>版本: {{ version }}</p>
      <p>构建时间: {{ buildTime }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'

const router = useRouter()
const appStore = useAppStore()

const version = ref('1.0.0')
const buildTime = ref(new Date().toLocaleDateString())

const tools = ref([
  {
    id: 'process-manager',
    name: '进程管理器',
    icon: '🔍',
    description: '实时监控系统进程，查看CPU、内存使用情况',
    route: '/process-manager'
  },
  {
    id: 'database-backup',
    name: '数据库备份',
    icon: '💾',
    description: 'MySQL 数据库备份和恢复工具',
    route: '/database-backup'
  }
])

function navigateTo(route) {
  router.push(route)
}

onMounted(() => {
  appStore.init()
})
</script>

<style scoped>
.home {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-section {
  text-align: center;
  margin-bottom: 3rem;
}

.welcome-section h1 {
  font-size: 2.5rem;
  color: #00ff00;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #888;
  font-size: 1.2rem;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.tool-card {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
  padding: 2rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tool-card:hover {
  border-color: #00ff00;
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 255, 0, 0.1);
}

.tool-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.tool-name {
  color: #fff;
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.tool-description {
  color: #888;
  line-height: 1.5;
}

.features-section {
  background: #1a1a1a;
  border-radius: 12px;
  padding: 2rem;
  margin-bottom: 2rem;
}

.features-section h2 {
  color: #00ff00;
  margin-bottom: 1.5rem;
}

.features-list {
  list-style: none;
  padding: 0;
}

.features-list li {
  color: #ccc;
  padding: 0.5rem 0;
  padding-left: 1.5rem;
  position: relative;
}

.features-list li::before {
  content: '✓';
  color: #00ff00;
  position: absolute;
  left: 0;
}

.version-info {
  text-align: center;
  color: #666;
  font-size: 0.9rem;
}

.version-info p {
  margin: 0.25rem 0;
}
</style>