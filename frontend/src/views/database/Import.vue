<template>
  <div class="import">
    <div class="import-form">
      <h3>导入配置</h3>

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

      <!-- 选择导入目录 -->
      <div class="form-group">
        <label>导入目录 *</label>
        <div class="path-input">
          <input
            v-model="inputDir"
            type="text"
            placeholder="选择导出文件目录"
            class="form-control"
            readonly
          />
          <button @click="selectInputDir" class="browse-btn">浏览</button>
        </div>
      </div>

      <!-- 导入选项 -->
      <div class="form-group">
        <label>导入选项</label>
        <div class="options-grid">
          <div class="option-item">
            <label>线程数</label>
            <input v-model.number="options.Threads" type="number" min="1" max="16" class="form-control small" />
          </div>
          <div class="option-item">
            <label>目标数据库</label>
            <input v-model="options.Schema" type="text" placeholder="可选" class="form-control" />
          </div>
          <div class="option-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="options.ResetProgress" />
              <span>重置进度</span>
            </label>
          </div>
          <div class="option-item">
            <label>等待超时（秒）</label>
            <input v-model.number="options.WaitTimeout" type="number" min="0" class="form-control small" />
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button
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
              <span class="schema-name">{{ conflict.Schema }}</span>
              <span class="conflict-count">
                {{ getTotalConflictCount(conflict) }} 个冲突对象
              </span>
            </div>
            <div class="conflict-details">
              <div v-if="conflict.Tables.length > 0" class="conflict-type">
                <span class="type-label">表:</span>
                <span class="type-items">{{ conflict.Tables.join(', ') }}</span>
              </div>
              <div v-if="conflict.Views.length > 0" class="conflict-type">
                <span class="type-label">视图:</span>
                <span class="type-items">{{ conflict.Views.join(', ') }}</span>
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
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useDatabaseStore } from '../../stores/database'

const databaseStore = useDatabaseStore()

// 从store获取响应式refs
const {
  connections,
  importLoading: importing,
  importProgress,
  importStatus,
  importConflicts: conflicts
} = storeToRefs(databaseStore)

// 从store获取方法
const {
  fetchConnections,
  importDatabases,
  checkImportConflicts,
  dropConflictingTables
} = databaseStore

// 响应式数据
const selectedConnectionId = ref('')
const inputDir = ref('')
const checking = ref(false)
const importError = ref('')

const options = ref({
  Threads: 4,
  Schema: '',
  ResetProgress: false,
  WaitTimeout: 300
})

// 计算属性
const canImport = computed(() => {
  return selectedConnectionId.value && inputDir.value
})

// 方法
async function selectInputDir() {
  try {
    const path = await window.go.main.App.SelectFolder()
    if (path) {
      inputDir.value = path
    }
  } catch (err) {
    console.error('Failed to select folder:', err)
  }
}

function getTotalConflictCount(conflict) {
  return conflict.Tables.length + conflict.Views.length +
         conflict.Events.length + conflict.Functions.length + conflict.Procedures.length
}

async function checkConflicts() {
  importError.value = ''
  const connection = connections.value.find(c => c.ID === selectedConnectionId.value)
  
  if (!connection) {
    importError.value = '请选择数据库连接'
    return
  }

  checking.value = true
  try {
    await checkImportConflicts(connection, inputDir.value)
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

  const connection = connections.value.find(c => c.ID === selectedConnectionId.value)
  if (!connection) return

  try {
    await dropConflictingTables(connection, conflicts)
    alert('冲突对象已删除')
  } catch (err) {
    importError.value = err.message || '删除冲突对象失败'
  }
}

async function startImport() {
  importError.value = ''
  const connection = connections.value.find(c => c.ID === selectedConnectionId.value)
  
  if (!connection) {
    importError.value = '请选择数据库连接'
    return
  }

  try {
    await importDatabases(connection, {
      InputDir: inputDir.value,
      ...options.value
    })
    alert('导入成功！')
  } catch (err) {
    importError.value = err.message || '导入失败'
  }
}

onMounted(() => {
  fetchConnections()
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
  border: 1px solid #00ccff;
  color: #00ccff;
}

.check-btn:hover:not(:disabled) {
  background: rgba(0, 204, 255, 0.1);
}

.import-btn {
  background: #00ff00;
  border: none;
  color: #000;
}

.import-btn:hover:not(:disabled) {
  background: #00cc00;
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
  background: #1a1a1a;
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
  background: linear-gradient(90deg, #00ccff, #0099cc);
  animation: pulse 1.5s infinite;
}

.progress-fill.importing {
  background: linear-gradient(90deg, #00ff00, #00cc00);
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
  color: #888;
  font-size: 0.9rem;
  text-align: center;
}

.conflicts-section {
  margin-top: 2rem;
  background: rgba(255, 68, 68, 0.1);
  border: 1px solid #ff4444;
  border-radius: 8px;
  padding: 1.5rem;
}

.conflicts-section h4 {
  color: #ff4444;
  margin-bottom: 1rem;
}

.conflict-list {
  margin-bottom: 1rem;
}

.conflict-item {
  background: rgba(0, 0, 0, 0.2);
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
  color: #fff;
  font-weight: bold;
}

.conflict-count {
  color: #ff4444;
  font-size: 0.85rem;
}

.conflict-details {
  padding-left: 1rem;
}

.conflict-type {
  margin-bottom: 0.25rem;
}

.type-label {
  color: #888;
  margin-right: 0.5rem;
}

.type-items {
  color: #ccc;
  font-size: 0.9rem;
}

.drop-btn {
  background: #ff4444;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #fff;
  cursor: pointer;
  transition: all 0.3s;
}

.drop-btn:hover {
  background: #ff6666;
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