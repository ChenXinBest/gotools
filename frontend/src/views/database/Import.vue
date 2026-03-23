<template>
  <div class="import">
    <div class="import-form">
      <h3>导入配置</h3>

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

      <!-- 选择导入目录/文件 -->
      <div class="form-group">
        <label>{{ exportTool === 'mysqldump' ? '导入文件 *' : '导入目录 *' }}</label>
        <div class="path-input">
          <input
            v-model="inputPath"
            type="text"
            :placeholder="exportTool === 'mysqldump' ? '选择 SQL 文件' : '选择导出文件目录'"
            class="form-control"
            readonly
          />
          <button @click="selectInputPath" class="browse-btn">浏览</button>
        </div>
      </div>

      <!-- 导入选项 - MySQL Shell -->
      <div class="form-group" v-if="exportTool === 'mysql-shell'">
        <label>导入选项 (MySQL Shell)</label>
        <div class="options-section">
          <div class="options-row">
            <div class="option-item">
              <label>线程数</label>
              <input v-model.number="options.threads" type="number" min="1" max="16" class="form-control small" />
            </div>
            <div class="option-item">
              <label>目标数据库</label>
              <input v-model="options.schema" type="text" placeholder="可选" class="form-control" />
            </div>
            <div class="option-item">
              <label>等待超时（秒）</label>
              <input v-model.number="options.wait_timeout" type="number" min="0" class="form-control small" />
            </div>
          </div>
          <div class="options-row checkboxes">
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.reset_progress" />
              <span>重置进度</span>
            </label>
          </div>
        </div>
      </div>

      <!-- 导入选项 - MySQLDump -->
      <div class="form-group" v-if="exportTool === 'mysqldump'">
        <label>导入选项 (mysqldump)</label>
        <div class="options-section">
          <div class="options-row">
            <div class="option-item">
              <label>目标数据库 *</label>
              <input v-model="options.database" type="text" placeholder="数据库名" class="form-control" />
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button
          v-if="exportTool === 'mysql-shell'"
          @click="checkConflicts"
          class="check-btn"
          :disabled="!canImport || checking"
        >
          {{ checking ? '检测中...' : '检测冲突' }}
        </button>
        <button
          @click="startImport"
          class="import-btn"
          :disabled="!canImport || importing"
        >
          {{ importing ? '导入中...' : '开始导入' }}
        </button>
      </div>

      <!-- 检测进度 -->
      <div v-if="checking" class="progress-section">
        <div class="progress-bar">
          <div class="progress-fill checking"></div>
        </div>
        <p class="progress-text">正在检测导入冲突...</p>
      </div>

      <!-- 冲突列表 -->
      <div v-if="conflicts.length > 0" class="conflicts-section">
        <h4>发现冲突</h4>
        <div class="conflict-list">
          <div v-for="(conflict, index) in conflicts" :key="index" class="conflict-item">
            <div class="conflict-header">
              <span class="schema-name">{{ conflict.schema }}</span>
              <span class="conflict-count">
                {{ getTotalConflictCount(conflict) }} 个冲突对象
              </span>
            </div>
            <div class="conflict-details">
              <div v-if="conflict.tables && conflict.tables.length > 0" class="conflict-type">
                <span class="type-label">表:</span>
                <span class="type-items">{{ conflict.tables.join(', ') }}</span>
              </div>
              <div v-if="conflict.views && conflict.views.length > 0" class="conflict-type">
                <span class="type-label">视图:</span>
                <span class="type-items">{{ conflict.views.join(', ') }}</span>
              </div>
            </div>
          </div>
        </div>
        <button @click="dropConflicts" class="drop-btn">
          删除冲突对象
        </button>
      </div>

      <!-- 导入进度 -->
      <div v-if="importing" class="progress-section">
        <div class="progress-bar">
          <div class="progress-fill importing"></div>
        </div>
        <p class="progress-text">{{ importStatus }}</p>
      </div>

      <!-- 错误信息 -->
      <div v-if="importError" class="error-message">
        {{ importError }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useDatabaseStore } from '../../stores/database'

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
  importLoading: importing,
  importProgress,
  importStatus,
  importConflicts: conflicts,
  exportSettings
} = storeToRefs(databaseStore)

// 从store获取方法
const {
  fetchConnections,
  fetchExportSettings,
  importDatabases,
  importTables,
  checkImportConflicts,
  dropConflictingTables
} = databaseStore

// 响应式数据
const selectedConnectionId = ref('')
const inputPath = ref('')
const checking = ref(false)
const importError = ref('')

// MySQL Shell 选项
const mysqlShellOptions = ref({
  threads: 4,
  schema: '',
  reset_progress: false,
  wait_timeout: 300
})

// MySQLDump 选项
const mySQLDumpOptions = ref({
  database: ''
})

// 当前使用的选项
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
const canImport = computed(() => {
  if (!selectedConnectionId.value || !inputPath.value) return false
  if (props.exportTool === 'mysqldump' && !options.value.database) return false
  return true
})

// 方法
async function selectInputPath() {
  try {
    let path
    if (props.exportTool === 'mysqldump') {
      path = await window.go.main.App.SelectFile()
    } else {
      path = await window.go.main.App.SelectFolder()
    }
    if (path) {
      inputPath.value = path
    }
  } catch (err) {
    console.error('Failed to select path:', err)
  }
}

function getTotalConflictCount(conflict) {
  const tables = conflict.tables?.length || 0
  const views = conflict.views?.length || 0
  const events = conflict.events?.length || 0
  const functions = conflict.functions?.length || 0
  const procedures = conflict.procedures?.length || 0
  return tables + views + events + functions + procedures
}

async function checkConflicts() {
  importError.value = ''
  const connection = connections.value.find(c => c.id === selectedConnectionId.value)
  
  if (!connection) {
    importError.value = '请选择数据库连接'
    return
  }

  checking.value = true
  try {
    await checkImportConflicts(connection, inputPath.value)
  } catch (err) {
    importError.value = err.message || '检测冲突失败'
  } finally {
    checking.value = false
  }
}

async function dropConflicts() {
  if (!confirm('确定要删除这些冲突对象吗？此操作不可恢复！')) {
    return
  }

  const connection = connections.value.find(c => c.id === selectedConnectionId.value)
  if (!connection) return

  try {
    await dropConflictingTables(connection.id, conflicts.value)
    alert('冲突对象已删除')
  } catch (err) {
    importError.value = err.message || '删除冲突对象失败'
  }
}

async function startImport() {
  importError.value = ''
  const connection = connections.value.find(c => c.id === selectedConnectionId.value)
  
  if (!connection) {
    importError.value = '请选择数据库连接'
    return
  }

  try {
    if (props.exportTool === 'mysql-shell') {
      // 使用 MySQL Shell 导入
      await importDatabases(connection, {
        input_dir: inputPath.value,
        ...options.value
      })
    } else {
      // 使用 mysqldump 导入
      await window.go.main.App.ImportDumpMySQLDump({
        connection_id: connection.id,
        input_file: inputPath.value,
        database: options.value.database
      })
    }
    alert('导入成功！')
  } catch (err) {
    importError.value = err.message || '导入失败'
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
        schema: '',
        reset_progress: false,
        wait_timeout: 300
      }
    }
  }
}

onMounted(async () => {
  await fetchConnections()
  await fetchExportSettings()
  initOptions()
})

// 监听工具切换
watch(() => props.exportTool, () => {
  initOptions()
  inputPath.value = ''
})
</script>

<style scoped>
.import {
  height: 100%;
}

.import-form {
  max-width: 800px;
}

.import-form h3 {
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
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

.check-btn,
.import-btn {
  padding: 0.75rem 2rem;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
}

.check-btn {
  background: transparent;
  border: 1px solid var(--info-color);
  color: var(--info-color);
}

.check-btn:hover:not(:disabled) {
  background: var(--info-subtle);
}

.import-btn {
  background: var(--accent-color);
  border: none;
  color: var(--text-on-accent);
}

.import-btn:hover:not(:disabled) {
  background: var(--accent-hover);
}

.check-btn:disabled,
.import-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.progress-section {
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
  transition: width 0.3s;
}

.progress-fill.checking {
  width: 100%;
  background: linear-gradient(90deg, var(--info-color), var(--info-subtle));
  animation: pulse 1.5s infinite;
}

.progress-fill.importing {
  background: linear-gradient(90deg, var(--accent-color), var(--accent-hover));
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.progress-text {
  color: var(--text-secondary);
  font-size: 0.9rem;
  text-align: center;
}

.conflicts-section {
  margin-top: 2rem;
  background: var(--danger-subtle);
  border: 1px solid var(--danger-color);
  border-radius: 8px;
  padding: 1.5rem;
}

.conflicts-section h4 {
  color: var(--danger-color);
  margin-bottom: 1rem;
}

.conflict-list {
  margin-bottom: 1rem;
}

.conflict-item {
  background: var(--bg-primary);
  border-radius: 6px;
  padding: 1rem;
  margin-bottom: 0.5rem;
}

.conflict-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.schema-name {
  color: var(--text-tertiary);
  font-weight: bold;
}

.conflict-count {
  color: var(--danger-color);
  font-size: 0.85rem;
}

.conflict-details {
  padding-left: 1rem;
}

.conflict-type {
  margin-bottom: 0.25rem;
}

.type-label {
  color: var(--text-secondary);
  margin-right: 0.5rem;
}

.type-items {
  color: var(--text-tertiary);
  font-size: 0.9rem;
}

.drop-btn {
  background: var(--danger-color);
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: var(--text-on-accent);
  cursor: pointer;
  transition: all 0.3s;
}

.drop-btn:hover {
  background: var(--danger-hover);
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
