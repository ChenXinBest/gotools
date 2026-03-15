<template>
  <div class="export">
    <div class="export-form">
      <h3>导出配置</h3>

      <!-- 选择连接 -->
      <div class="form-group">
        <label>数据库连接 *</label>
        <select v-model="selectedConnectionId" class="form-control">
          <option value="">选择数据库连接</option>
          <option v-for="conn in connections" :key="conn.ID" :value="conn.ID">
            {{ conn.Name }} ({{ conn.Host }}:{{ conn.Port }})
          </option>
        </select>
      </div>

      <!-- 选择导出类型 -->
      <div class="form-group">
        <label>导出类型 *</label>
        <div class="radio-group">
          <label class="radio-label">
            <input type="radio" v-model="exportType" value="databases" />
            <span>导出数据库</span>
          </label>
          <label class="radio-label">
            <input type="radio" v-model="exportType" value="tables" />
            <span>导出表</span>
          </label>
        </div>
      </div>

      <!-- 选择数据库 -->
      <div class="form-group" v-if="exportType === 'databases' || exportType === 'tables'">
        <label>选择数据库</label>
        <div class="database-selector">
          <div class="search-box">
            <input
              v-model="databaseSearch"
              type="text"
              placeholder="搜索数据库..."
              class="form-control"
            />
          </div>
          <div class="database-list">
            <div
              v-for="db in filteredDatabases"
              :key="db"
              class="database-item"
              :class="{ selected: selectedDatabases.includes(db) }"
              @click="toggleDatabase(db)"
            >
              <input type="checkbox" :checked="selectedDatabases.includes(db)" />
              <span>{{ db }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 选择表（当选择表导出时） -->
      <div class="form-group" v-if="exportType === 'tables' && selectedDatabases.length === 1">
        <label>选择表</label>
        <div class="table-selector">
          <div class="search-box">
            <input
              v-model="tableSearch"
              type="text"
              placeholder="搜索表..."
              class="form-control"
            />
          </div>
          <div class="table-list">
            <div
              v-for="table in filteredTables"
              :key="table"
              class="table-item"
              :class="{ selected: selectedTables.includes(table) }"
              @click="toggleTable(table)"
            >
              <input type="checkbox" :checked="selectedTables.includes(table)" />
              <span>{{ table }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 选择导出目录 -->
      <div class="form-group">
        <label>导出目录 *</label>
        <div class="path-input">
          <input
            v-model="outputDir"
            type="text"
            placeholder="选择导出目录"
            class="form-control"
            readonly
          />
          <button @click="selectOutputDir" class="browse-btn">浏览</button>
        </div>
      </div>

      <!-- 导出选项 -->
      <div class="form-group">
        <label>导出选项</label>
        <div class="options-grid">
          <div class="option-item">
            <label>线程数</label>
            <input v-model.number="options.Threads" type="number" min="1" max="16" class="form-control small" />
          </div>
          <div class="option-item">
            <label>压缩格式</label>
            <select v-model="options.Compression" class="form-control">
              <option value="gzip">gzip</option>
              <option value="zstd">zstd</option>
              <option value="none">无压缩</option>
            </select>
          </div>
          <div class="option-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.Overwrite" />
              <span>覆盖已存在文件</span>
            </label>
          </div>
          <div class="option-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.SkipDefiner" />
              <span>跳过 Definer</span>
            </label>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button
          @click="startExport"
          class="export-btn"
          :disabled="!canExport || exporting"
        >
          {{ exporting ? '导出中...' : '开始导出' }}
        </button>
      </div>

      <!-- 导出进度 -->
      <div v-if="exporting" class="export-progress">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: exportProgress + '%' }"></div>
        </div>
        <p class="progress-text">{{ exportStatus }}</p>
      </div>

      <!-- 错误信息 -->
      <div v-if="exportError" class="error-message">
        {{ exportError }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useDatabaseStore } from '../../stores/database'

const databaseStore = useDatabaseStore()

// 从store获取响应式refs
const {
  connections,
  databases,
  tables,
  exportLoading: exporting,
  exportProgress,
  exportStatus,
  exportSettings
} = storeToRefs(databaseStore)

// 从store获取方法
const {
  fetchConnections,
  fetchDatabases,
  fetchTables,
  exportDatabases,
  exportTables
} = databaseStore

// 响应式数据
const selectedConnectionId = ref('')
const exportType = ref('databases')
const selectedDatabases = ref([])
const selectedTables = ref([])
const outputDir = ref('')
const databaseSearch = ref('')
const tableSearch = ref('')
const exportError = ref('')

const options = ref({
  Threads: 4,
  Compression: 'gzip',
  Overwrite: true,
  SkipDefiner: true
})

// 计算属性
const filteredDatabases = computed(() => {
  if (!databaseSearch.value) return databases.value
  return databases.value.filter(db =>
    db.toLowerCase().includes(databaseSearch.value.toLowerCase())
  )
})

const filteredTables = computed(() => {
  if (!tableSearch.value) return tables.value
  return tables.value.filter(table =>
    table.toLowerCase().includes(tableSearch.value.toLowerCase())
  )
})

const canExport = computed(() => {
  if (!selectedConnectionId.value || !outputDir.value) return false
  if (exportType.value === 'databases' && selectedDatabases.value.length === 0) return false
  if (exportType.value === 'tables' && (selectedDatabases.value.length !== 1 || selectedTables.value.length === 0)) return false
  return true
})

// 方法
function toggleDatabase(db) {
  const index = selectedDatabases.value.indexOf(db)
  if (index > -1) {
    selectedDatabases.value.splice(index, 1)
  } else {
    if (exportType.value === 'tables') {
      selectedDatabases.value = [db]
      loadTablesForDatabase(db)
    } else {
      selectedDatabases.value.push(db)
    }
  }
}

function toggleTable(table) {
  const index = selectedTables.value.indexOf(table)
  if (index > -1) {
    selectedTables.value.splice(index, 1)
  } else {
    selectedTables.value.push(table)
  }
}

async function loadTablesForDatabase(database) {
  const connection = connections.value.find(c => c.ID === selectedConnectionId.value)
  if (connection) {
    await fetchTables(connection, database)
  }
}

async function selectOutputDir() {
  try {
    const path = await window.go.main.App.SelectFolder()
    if (path) {
      outputDir.value = path
    }
  } catch (err) {
    console.error('Failed to select folder:', err)
  }
}

async function startExport() {
  exportError.value = ''
  const connection = connections.value.find(c => c.ID === selectedConnectionId.value)
  
  if (!connection) {
    exportError.value = '请选择数据库连接'
    return
  }

  try {
    if (exportType.value === 'databases') {
      await exportDatabases(connection, selectedDatabases.value, {
        OutputDir: outputDir.value,
        ...options.value
      })
    } else {
      await exportTables(connection, selectedDatabases.value[0], selectedTables.value, {
        OutputDir: outputDir.value,
        ...options.value
      })
    }
    alert('导出成功！')
  } catch (err) {
    exportError.value = err.message || '导出失败'
  }
}

onMounted(() => {
  fetchConnections()
})
</script>

<style scoped>
.export {
  height: 100%;
}

.export-form {
  max-width: 800px;
}

.export-form h3 {
  color: #00ff00;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  color: #888;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.form-control {
  width: 100%;
  background: #0d0d0d;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #fff;
  font-size: 1rem;
}

.form-control:focus {
  outline: none;
  border-color: #00ff00;
}

.form-control.small {
  width: 80px;
}

select.form-control {
  cursor: pointer;
}

.radio-group {
  display: flex;
  gap: 2rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #ccc;
  cursor: pointer;
}

.database-selector,
.table-selector {
  border: 1px solid #333;
  border-radius: 6px;
  overflow: hidden;
}

.search-box {
  padding: 0.5rem;
  border-bottom: 1px solid #333;
}

.database-list,
.table-list {
  max-height: 200px;
  overflow-y: auto;
}

.database-item,
.table-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.database-item:hover,
.table-item:hover {
  background: rgba(0, 255, 0, 0.05);
}

.database-item.selected,
.table-item.selected {
  background: rgba(0, 255, 0, 0.1);
}

.path-input {
  display: flex;
  gap: 0.5rem;
}

.browse-btn {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #00ff00;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
}

.browse-btn:hover {
  border-color: #00ff00;
}

.options-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.option-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.option-item label {
  margin-bottom: 0;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #ccc;
  cursor: pointer;
}

.form-actions {
  margin-top: 2rem;
}

.export-btn {
  background: #00ff00;
  border: none;
  border-radius: 6px;
  padding: 0.75rem 2rem;
  color: #000;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
}

.export-btn:hover:not(:disabled) {
  background: #00cc00;
}

.export-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.export-progress {
  margin-top: 1.5rem;
}

.progress-bar {
  height: 8px;
  background: #1a1a1a;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #00ff00, #00cc00);
  transition: width 0.3s;
}

.progress-text {
  color: #888;
  font-size: 0.9rem;
  text-align: center;
}

.error-message {
  background: rgba(255, 68, 68, 0.1);
  border: 1px solid #ff4444;
  border-radius: 6px;
  padding: 1rem;
  margin-top: 1rem;
  color: #ff4444;
}
</style>