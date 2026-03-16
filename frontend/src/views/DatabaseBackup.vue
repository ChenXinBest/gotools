<template>
  <div class="database-backup">
    <div class="header">
      <h1>{{ t('db.title') }}</h1>
      <p class="subtitle">{{ t('db.subtitle') }}</p>
    </div>

    <div class="toolbar">
      <div class="tool-selector">
        <label>{{ t('db.tool') }}:</label>
        <select v-model="exportTool" class="tool-select">
          <option value="mysql-shell">{{ t('db.mysqlShell') }}</option>
          <option value="mysqldump">{{ t('db.mysqldump') }}</option>
        </select>
      </div>
      <button @click="saveAsDefault" class="save-default-btn" :disabled="saving">
        {{ saving ? t('db.saving') : t('db.saveAsDefault') }}
      </button>
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
      <Export v-if="activeTab === 'export'" :export-tool="exportTool" />
      <Import v-if="activeTab === 'import'" :export-tool="exportTool" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, provide } from 'vue'
import Connections from './database/Connections.vue'
import Export from './database/Export.vue'
import Import from './database/Import.vue'
import { useDatabaseStore } from '../stores/database'
import { storeToRefs } from 'pinia'
import { t } from '../i18n'

const databaseStore = useDatabaseStore()
const { exportSettings } = storeToRefs(databaseStore)

const tabs = [
  { id: 'connections', name: t('db.connections') },
  { id: 'export', name: t('db.export') },
  { id: 'import', name: t('db.import') }
]

const activeTab = ref('connections')
const exportTool = ref('mysql-shell')
const saving = ref(false)

// 提供给子组件使用
provide('exportTool', exportTool)
provide('exportSettings', exportSettings)

// 加载默认设置
onMounted(async () => {
  await databaseStore.fetchExportSettings()
  exportTool.value = exportSettings.value.export_tool || 'mysql-shell'
})

// 保存为默认设置
async function saveAsDefault() {
  saving.value = true
  try {
    // 更新 export_tool
    const settingsToSave = {
      ...exportSettings.value,
      export_tool: exportTool.value
    }
    await databaseStore.saveExportSettings(settingsToSave)
    alert('已保存为默认设置')
  } catch (err) {
    alert('保存失败: ' + err.message)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.database-backup {
  padding: 2rem;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  margin-bottom: 1rem;
}

.header h1 {
  color: var(--accent-color);
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: var(--text-secondary);
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
  padding: 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.tool-selector {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tool-selector label {
  color: var(--text-secondary);
}

.tool-select {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: var(--text-tertiary);
  font-size: 1rem;
  min-width: 150px;
}

.tool-select:focus {
  outline: none;
  border-color: var(--accent-color);
}

.save-default-btn {
  background: transparent;
  border: 1px solid var(--accent-color);
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: var(--accent-color);
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s;
  margin-left: auto;
}

.save-default-btn:hover:not(:disabled) {
  background: var(--accent-subtle);
}

.save-default-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.5rem;
}

.tab-btn {
  background: transparent;
  border: none;
  padding: 0.75rem 1.5rem;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.3s;
  border-radius: 6px 6px 0 0;
}

.tab-btn:hover {
  color: var(--accent-color);
  background: var(--accent-subtle);
}

.tab-btn.active {
  color: var(--accent-color);
  background: var(--accent-subtle);
  border-bottom: 2px solid var(--accent-color);
}

.tab-content {
  flex: 1;
  overflow: auto;
}
</style>
