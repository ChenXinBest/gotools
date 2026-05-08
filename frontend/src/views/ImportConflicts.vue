<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Events } from '@wailsio/runtime'

const conflicts = ref([])
const loading = ref(true)

let cleanup = null

onMounted(() => {
  cleanup = Events.On('import-conflicts', (data) => {
    if (Array.isArray(data)) {
      conflicts.value = data
    } else if (data && data.conflicts) {
      conflicts.value = data.conflicts
    }
    loading.value = false
  })
})

onUnmounted(() => {
  if (cleanup) cleanup()
})

function resolveConflicts() {
  window._wails?.invoke?.('dropConflictingTables')?.then(() => {
    conflicts.value = []
  })
}
</script>

<template>
  <div class="conflict-window">
    <div class="container">
      <h3>导入冲突检测</h3>

      <div v-if="loading" class="loading">
        <p>正在检测冲突...</p>
      </div>

      <div v-else-if="conflicts.length === 0" class="no-conflicts">
        <p class="success">未检测到冲突，可以安全导入</p>
      </div>

      <div v-else class="conflict-list">
        <p class="summary">检测到 {{ conflicts.length }} 个冲突对象</p>

        <div v-for="(conflict, idx) in conflicts" :key="idx" class="conflict-item">
          <span class="conflict-name">{{ conflict.schema || conflict.database }}.{{ conflict.table || conflict.name }}</span>
          <span class="conflict-type">{{ conflict.type || '表' }}</span>
        </div>

        <div class="actions">
          <button class="btn-danger" @click="resolveConflicts">删除冲突对象并继续导入</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conflict-window {
  height: 100vh;
  display: flex;
  background: var(--bg-primary, #1a1a2e);
  color: var(--text-primary, #e0e0e0);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  padding: 20px;
}

.container {
  width: 100%;
}

h3 {
  margin-bottom: 20px;
  font-size: 18px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary, #a0a0a0);
}

.no-conflicts {
  text-align: center;
  padding: 40px;
}

.success {
  color: #4caf50;
  font-size: 16px;
}

.summary {
  color: #ff9800;
  margin-bottom: 16px;
  font-size: 14px;
}

.conflict-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  margin-bottom: 4px;
  background: var(--border-color, #2a2a4a);
  border-radius: 4px;
  font-size: 13px;
}

.conflict-name {
  color: var(--text-primary, #e0e0e0);
}

.conflict-type {
  color: #ff9800;
  font-size: 12px;
  padding: 2px 8px;
  border: 1px solid #ff9800;
  border-radius: 3px;
}

.actions {
  margin-top: 20px;
  text-align: center;
}

.btn-danger {
  padding: 10px 24px;
  background: #f44336;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.btn-danger:hover {
  opacity: 0.9;
}
</style>
