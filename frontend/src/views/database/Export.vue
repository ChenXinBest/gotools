<template>
  <div class="export">
    <div class="export-form">
      <h3>导出配置</h3>

      <!-- 选择连接 -->
      <div class="form-group">
        <label>数据库连接 *</label>
        <select v-model="selectedConnectionId" class="form-control">
          <option value="">选择数据库连接</option>
          <option v-for="conn in connections" :key="conn.id" :value="conn.id">
            {{ conn.name }} ({{ conn.host }}:{{ conn.port }})
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
            <button @click="toggleSelectAllDatabases" class="select-all-btn">
              {{ allDatabasesSelected ? '取消全选' : '全选' }}
            </button>
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
            <button @click="toggleSelectAllTables" class="select-all-btn">
              {{ allTablesSelected ? '取消全选' : '全选' }}
            </button>
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

      <!-- 导出选项 - MySQL Shell -->
      <div class="form-group" v-if="exportTool === 'mysql-shell'">
        <label>导出选项 (MySQL Shell)</label>
        <div class="options-section">
          <div class="options-row">
            <div class="option-item">
              <label>线程数</label>
              <input v-model.number="options.threads" type="number" min="1" max="16" class="form-control small" />
            </div>
            <div class="option-item">
              <label>压缩格式</label>
              <select v-model="options.compression" class="form-control">
                <option value="gzip">gzip</option>
                <option value="zstd">zstd</option>
                <option value="none">无压缩</option>
              </select>
            </div>
            <div class="option-item">
              <label>分块大小</label>
              <select v-model="options.chunk_size" class="form-control">
                <option value="64M">64M</option>
                <option value="128M">128M</option>
                <option value="256M">256M</option>
                <option value="512M">512M</option>
              </select>
            </div>
          </div>
          <div class="options-row checkboxes">
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.overwrite" />
              <span>覆盖已存在文件</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.skip_definer" />
              <span>跳过 Definer</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.skip_binlog" />
              <span>跳过 Binlog</span>
            </label>
          </div>
        </div>
      </div>

      <!-- 导出选项 - MySQLDump -->
      <div class="form-group" v-if="exportTool === 'mysqldump'">
        <label>导出选项 (mysqldump)</label>
        <div class="options-section">
          <div class="options-row">
            <div class="option-item">
              <label>压缩格式</label>
              <select v-model="options.compression" class="form-control">
                <option value="gzip">gzip</option>
                <option value="zstd">zstd</option>
                <option value="none">无压缩</option>
              </select>
            </div>
          </div>
          <div class="options-row checkboxes">
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.overwrite" />
              <span>覆盖已存在文件</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.single_transaction" />
              <span>单事务</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.routines" />
              <span>存储过程/函数</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.events" />
              <span>事件</span>
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
import { ref, computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useDatabaseStore } from '../../stores/database'
import {
  ExportDatabasesMySQLDump,
  ExportTablesMySQLDump
} from '../../../wailsjs/go/main/App'

const props = defineProps({
  exportTool: {
    type: String,
    default: 'mysql-shell'
  }
})

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
  fetchExportSettings,
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

// MySQL Shell 选项
const mysqlShellOptions = ref({
  threads: 4,
  compression: 'gzip',
  chunk_size: '64M',
  skip_definer: true,
  skip_binlog: false,
  overwrite: true
})

// MySQLDump 选项
const mySQLDumpOptions = ref({
  compression: 'gzip',
  single_transaction: true,
  routines: true,
  events: true,
  overwrite: true
})

// 当前使用的选项（根据工具切换）
const options = computed({
  get: () => props.exportTool === 'mysql-shell' ? mysqlShellOptions.value : mySQLDumpOptions.value,
  set: (val) => {
    if (props.exportTool === 'mysql-shell') {
      mysqlShellOptions.value = val
    } else {
      mySQLDumpOptions.value = val
    }
  }
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

// 全选状态计算属性
const allDatabasesSelected = computed(() => {
  if (filteredDatabases.value.length === 0) return false
  return filteredDatabases.value.every(db => selectedDatabases.value.includes(db))
})

const allTablesSelected = computed(() => {
  if (filteredTables.value.length === 0) return false
  return filteredTables.value.every(table => selectedTables.value.includes(table))
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

function toggleSelectAllDatabases() {
  if (allDatabasesSelected.value) {
    selectedDatabases.value = selectedDatabases.value.filter(
      db => !filteredDatabases.value.includes(db)
    )
  } else {
    filteredDatabases.value.forEach(db => {
      if (!selectedDatabases.value.includes(db)) {
        selectedDatabases.value.push(db)
      }
    })
  }
}

function toggleSelectAllTables() {
  if (allTablesSelected.value) {
    selectedTables.value = selectedTables.value.filter(
      table => !filteredTables.value.includes(table)
    )
  } else {
    filteredTables.value.forEach(table => {
      if (!selectedTables.value.includes(table)) {
        selectedTables.value.push(table)
      }
    })
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
  const connection = connections.value.find(c => c.id === selectedConnectionId.value)
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
  const connection = connections.value.find(c => c.id === selectedConnectionId.value)
  
  if (!connection) {
    exportError.value = '请选择数据库连接'
    return
  }

  try {
    const config = {
      output_dir: outputDir.value,
      ...options.value
    }

    if (props.exportTool === 'mysql-shell') {
      // 使用 MySQL Shell 导出
      if (exportType.value === 'databases') {
        await exportDatabases(connection, selectedDatabases.value, config)
      } else {
        await exportTables(connection, selectedDatabases.value[0], selectedTables.value, config)
      }
    } else {
      // 使用 mysqldump 导出
      const mysqldumpRequest = {
        connection_id: connection.id,
        databases: selectedDatabases.value,
        database: selectedDatabases.value[0],
        tables: selectedTables.value,
        output_dir: outputDir.value,
        compression: options.value.compression,
        single_transaction: options.value.single_transaction,
        routines: options.value.routines,
        events: options.value.events,
        overwrite: options.value.overwrite
      }
      
      if (exportType.value === 'databases') {
        const result = await ExportDatabasesMySQLDump(mysqldumpRequest)
        if (!result.success) {
          throw new Error(result.message)
        }
      } else {
        const result = await ExportTablesMySQLDump(mysqldumpRequest)
        if (!result.success) {
          throw new Error(result.message)
        }
      }
    }
    alert('导出成功！')
  } catch (err) {
    exportError.value = err.message || '导出失败'
  }
}

// 初始化选项配置
function initOptions() {
  const settings = exportSettings.value
  if (settings) {
    // 初始化 MySQL Shell 选项
    if (settings.mysql_shell) {
      mysqlShellOptions.value = {
        threads: settings.mysql_shell.threads || 4,
        compression: settings.mysql_shell.compression || 'gzip',
        chunk_size: settings.mysql_shell.chunk_size || '64M',
        skip_definer: settings.mysql_shell.skip_definer ?? true,
        skip_binlog: settings.mysql_shell.skip_binlog ?? false,
        overwrite: settings.mysql_shell.overwrite ?? true
      }
    }
    // 初始化 MySQLDump 选项
    if (settings.mysql_dump) {
      mySQLDumpOptions.value = {
        compression: settings.mysql_dump.compression || 'gzip',
        single_transaction: settings.mysql_dump.single_transaction ?? true,
        routines: settings.mysql_dump.routines ?? true,
        events: settings.mysql_dump.events ?? true,
        overwrite: settings.mysql_dump.overwrite ?? true
      }
    }
  }
}

onMounted(async () => {
  await fetchConnections()
  await fetchExportSettings()
  initOptions()
})

// 监听工具切换，变化时重新加载对应配置
watch(() => props.exportTool, () => {
  initOptions()
})

// 监听连接选择变化，自动加载数据库列表
watch(selectedConnectionId, async (newConnectionId) => {
  if (!newConnectionId) {
    return
  }
  const connection = connections.value.find(c => c.id === newConnectionId)
  if (connection) {
    await fetchDatabases(connection)
  }
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
  color: var(--accent-color);
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.form-control {
  width: 100%;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: var(--text-tertiary);
  font-size: 1rem;
}

.form-control:focus {
  outline: none;
  border-color: var(--accent-color);
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
  color: var(--text-tertiary);
  cursor: pointer;
}

.database-selector,
.table-selector {
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
}

.search-box {
  padding: 0.5rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.search-box .form-control {
  flex: 1;
}

.select-all-btn {
  background: transparent;
  border: 1px solid var(--accent-color);
  border-radius: 4px;
  padding: 0.25rem 0.75rem;
  color: var(--accent-color);
  cursor: pointer;
  font-size: 0.85rem;
  white-space: nowrap;
  transition: all 0.3s;
}

.select-all-btn:hover {
  background: var(--accent-subtle);
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
  background: var(--accent-subtle);
}

.database-item.selected,
.table-item.selected {
  background: var(--accent-subtle);
}

.path-input {
  display: flex;
  gap: 0.5rem;
}

.browse-btn {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: var(--accent-color);
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
}

.browse-btn:hover {
  border-color: var(--accent-color);
}

.options-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.options-row {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.options-row.checkboxes {
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
  padding-top: 0.5rem;
  border-top: 1px solid var(--border-color);
}

.options-row .option-item {
  min-width: 140px;
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
  color: var(--text-tertiary);
  cursor: pointer;
  white-space: nowrap;
}

.form-actions {
  margin-top: 2rem;
}

.export-btn {
  background: var(--accent-color);
  border: none;
  border-radius: 6px;
  padding: 0.75rem 2rem;
  color: var(--text-on-accent);
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
}

.export-btn:hover:not(:disabled) {
  background: var(--accent-hover);
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
  background: var(--bg-secondary);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--accent-color), var(--accent-hover));
  transition: width 0.3s;
}

.progress-text {
  color: var(--text-secondary);
  font-size: 0.9rem;
  text-align: center;
}

.error-message {
  background: var(--danger-subtle);
  border: 1px solid var(--danger-color);
  border-radius: 6px;
  padding: 1rem;
  margin-top: 1rem;
  color: var(--danger-color);
}
</style>
