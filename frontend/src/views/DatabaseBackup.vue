<template>
  <div class="database-backup">
    <div class="header">
      <h1>数据库备份</h1>
      <p class="subtitle">MySQL 数据库导入导出工具</p>
    </div>

    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        class="tab-btn"
        :class="{ active: activeTab === tab.id }"
        @click="activeTab = tab.id"
      >
        {{ tab.name }}
      </button>
    </div>

    <div class="tab-content">
      <Connections v-if="activeTab === 'connections'" />
      <Export v-if="activeTab === 'export'" />
      <Import v-if="activeTab === 'import'" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Connections from './database/Connections.vue'
import Export from './database/Export.vue'
import Import from './database/Import.vue'

const tabs = [
  { id: 'connections', name: '连接管理' },
  { id: 'export', name: '导出数据' },
  { id: 'import', name: '导入数据' }
]

const activeTab = ref('connections')
</script>

<style scoped>
.database-backup {
  padding: 2rem;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  margin-bottom: 1.5rem;
}

.header h1 {
  color: #00ff00;
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #888;
}

.tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid #333;
  padding-bottom: 0.5rem;
}

.tab-btn {
  background: transparent;
  border: none;
  padding: 0.75rem 1.5rem;
  color: #888;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.3s;
  border-radius: 6px 6px 0 0;
}

.tab-btn:hover {
  color: #00ff00;
  background: rgba(0, 255, 0, 0.05);
}

.tab-btn.active {
  color: #00ff00;
  background: rgba(0, 255, 0, 0.1);
  border-bottom: 2px solid #00ff00;
}

.tab-content {
  flex: 1;
  overflow: auto;
}
</style>