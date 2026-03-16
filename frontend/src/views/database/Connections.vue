<template>
  <div class="connections">
    <div class="toolbar">
      <button @click="showAddDialog = true" class="add-btn">
        + 添加连接
      </button>
      <button @click="refreshConnections" class="refresh-btn" :disabled="loading">
        {{ loading ? '加载中...' : '刷新' }}
      </button>
    </div>

    <div class="connection-list" v-if="!loading">
      <div
        v-for="conn in connections"
        :key="conn.id"
        class="connection-card"
        :class="{ active: currentConnection?.id === conn.id }"
        @click="selectConnection(conn)"
      >
        <div class="connection-header">
          <h3 class="connection-name">{{ conn.name }}</h3>
          <div class="connection-actions">
            <button @click.stop="editConnection(conn)" class="edit-btn">编辑</button>
            <button @click.stop="deleteConnection(conn.id)" class="delete-btn">删除</button>
          </div>
        </div>
        <div class="connection-info">
          <div class="info-item">
            <span class="info-label">主机:</span>
            <span class="info-value">{{ conn.host }}:{{ conn.port }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">用户:</span>
            <span class="info-value">{{ conn.user }}</span>
          </div>
          <div class="info-item" v-if="conn.database">
            <span class="info-label">数据库:</span>
            <span class="info-value">{{ conn.database }}</span>
          </div>
        </div>
        <div class="connection-footer">
          <button @click.stop="testConnection(conn)" class="test-btn" :disabled="testingConnection === conn.id">
            {{ testingConnection === conn.id ? '测试中...' : '测试连接' }}
          </button>
          <span v-if="connectionStatus[conn.id]" class="status-badge" :class="connectionStatus[conn.id].success ? 'success' : 'error'">
            {{ connectionStatus[conn.id].message }}
          </span>
        </div>
      </div>

      <div v-if="connections.length === 0" class="no-data">
        <p>暂无数据库连接</p>
        <p>点击上方按钮添加第一个连接</p>
      </div>
    </div>

    <div v-else class="loading">
      <div class="loading-spinner"></div>
      <p>加载连接信息...</p>
    </div>

    <!-- 添加/编辑连接对话框 -->
    <div v-if="showAddDialog || showEditDialog" class="dialog-overlay" @click="closeDialog">
      <div class="dialog" @click.stop>
        <div class="dialog-header">
          <h3>{{ showEditDialog ? '编辑连接' : '添加连接' }}</h3>
          <button @click="closeDialog" class="close-btn">&times;</button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>连接名称 *</label>
            <input v-model="formData.name" type="text" placeholder="输入连接名称" />
          </div>
          <div class="form-row">
            <div class="form-group flex-1">
              <label>主机地址 *</label>
              <input v-model="formData.host" type="text" placeholder="localhost" />
            </div>
            <div class="form-group" style="width: 100px;">
              <label>端口 *</label>
              <input v-model.number="formData.port" type="number" placeholder="3306" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group flex-1">
              <label>用户名 *</label>
              <input v-model="formData.user" type="text" placeholder="root" />
            </div>
            <div class="form-group flex-1">
              <label>密码</label>
              <input v-model="formData.password" type="password" placeholder="可选" />
            </div>
          </div>
          <div class="form-group">
            <label>默认数据库</label>
            <input v-model="formData.database" type="text" placeholder="可选" />
          </div>
        </div>
        <div class="dialog-footer">
          <button @click="closeDialog" class="cancel-btn">取消</button>
          <button @click="saveConnection" class="save-btn" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useDatabaseStore } from '../../stores/database'

const databaseStore = useDatabaseStore()

// 从store获取响应式refs
const {
  connections,
  currentConnection,
  connectionsLoading: loading
} = storeToRefs(databaseStore)

// 从store获取方法
const {
  fetchConnections,
  addConnection,
  updateConnection,
  deleteConnection: storeDeleteConnection,
  testConnection: storeTestConnection,
  selectConnection: storeSelectConnection
} = databaseStore

// 响应式数据
const showAddDialog = ref(false)
const showEditDialog = ref(false)
const saving = ref(false)
const testingConnection = ref(null)
const connectionStatus = ref({})

const formData = ref({
  id: '',
  name: '',
  host: 'localhost',
  port: 3306,
  user: 'root',
  password: '',
  database: ''
})

// 方法
function refreshConnections() {
  fetchConnections()
}

function selectConnection(conn) {
  storeSelectConnection(conn.id)
}

function editConnection(conn) {
  formData.value = { ...conn }
  showEditDialog.value = true
}

async function deleteConnection(id) {
  if (confirm('确定要删除这个连接吗？')) {
    await storeDeleteConnection(id)
  }
}

async function testConnection(conn) {
  testingConnection.value = conn.id
  try {
    const success = await databaseStore.testConnection(conn)
    connectionStatus.value[conn.id] = {
      success,
      message: success ? '连接成功' : '连接失败'
    }
  } catch (err) {
    connectionStatus.value[conn.id] = {
      success: false,
      message: '连接失败: ' + err.message
    }
  } finally {
    setTimeout(() => {
      testingConnection.value = null
    }, 2000)
  }
}

function closeDialog() {
  showAddDialog.value = false
  showEditDialog.value = false
  resetForm()
}

function resetForm() {
  formData.value = {
    id: '',
    name: '',
    host: 'localhost',
    port: 3306,
    user: 'root',
    password: '',
    database: ''
  }
}

async function saveConnection() {
  saving.value = true
  try {
    if (showEditDialog.value) {
      await updateConnection(formData.value)
    } else {
      await addConnection(formData.value)
    }
    closeDialog()
  } catch (err) {
    alert('保存失败: ' + err.message)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchConnections()
})
</script>

<style scoped>
.connections {
  height: 100%;
}

.toolbar {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.add-btn,
.refresh-btn {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: var(--text-tertiary);
  cursor: pointer;
  transition: all 0.3s;
}

.add-btn:hover {
  border-color: var(--accent-color);
  color: var(--accent-color);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.connection-list {
  display: grid;
  gap: 1rem;
}

.connection-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  cursor: pointer;
  transition: all 0.3s;
}

.connection-card:hover {
  border-color: var(--accent-color);
}

.connection-card.active {
  border-color: var(--accent-color);
  background: var(--accent-subtle);
}

.connection-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.connection-name {
  color: var(--text-tertiary);
  margin: 0;
}

.connection-actions {
  display: flex;
  gap: 0.5rem;
}

.edit-btn,
.delete-btn {
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 0.25rem 0.75rem;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.3s;
}

.edit-btn:hover {
  border-color: var(--info-color);
  color: var(--info-color);
}

.delete-btn:hover {
  border-color: var(--danger-color);
  color: var(--danger-color);
}

.connection-info {
  margin-bottom: 1rem;
}

.info-item {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
  font-size: 0.9rem;
}

.info-label {
  color: var(--text-secondary);
  min-width: 60px;
}

.info-value {
  color: var(--text-tertiary);
}

.connection-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.test-btn {
  background: transparent;
  border: 1px solid var(--accent-color);
  border-radius: 4px;
  padding: 0.25rem 0.75rem;
  color: var(--accent-color);
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.3s;
}

.test-btn:hover:not(:disabled) {
  background: var(--accent-subtle);
}

.test-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
}

.status-badge.success {
  background: var(--accent-subtle);
  color: var(--accent-color);
}

.status-badge.error {
  background: var(--danger-subtle);
  color: var(--danger-color);
}

.no-data {
  text-align: center;
  padding: 3rem;
  color: var(--text-secondary);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--bg-secondary);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Dialog styles */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--shadow-color);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.dialog-header h3 {
  color: var(--text-tertiary);
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  font-size: 1.5rem;
  cursor: pointer;
}

.close-btn:hover {
  color: var(--text-tertiary);
}

.dialog-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.form-group input {
  width: 100%;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: var(--text-tertiary);
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-color);
}

.form-row {
  display: flex;
  gap: 1rem;
}

.flex-1 {
  flex: 1;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1.5rem;
  border-top: 1px solid var(--border-color);
}

.cancel-btn,
.save-btn {
  padding: 0.5rem 1.5rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.3s;
}

.cancel-btn {
  background: transparent;
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.cancel-btn:hover {
  border-color: var(--text-secondary);
  color: var(--text-tertiary);
}

.save-btn {
  background: var(--accent-color);
  border: none;
  color: var(--text-on-accent);
}

.save-btn:hover:not(:disabled) {
  background: var(--accent-hover);
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>