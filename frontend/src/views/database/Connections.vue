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
        :key="conn.ID"
        class="connection-card"
        :class="{ active: currentConnection?.ID === conn.ID }"
        @click="selectConnection(conn)"
      >
        <div class="connection-header">
          <h3 class="connection-name">{{ conn.Name }}</h3>
          <div class="connection-actions">
            <button @click.stop="editConnection(conn)" class="edit-btn">编辑</button>
            <button @click.stop="deleteConnection(conn.ID)" class="delete-btn">删除</button>
          </div>
        </div>
        <div class="connection-info">
          <div class="info-item">
            <span class="info-label">主机:</span>
            <span class="info-value">{{ conn.Host }}:{{ conn.Port }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">用户:</span>
            <span class="info-value">{{ conn.User }}</span>
          </div>
          <div class="info-item" v-if="conn.Database">
            <span class="info-label">数据库:</span>
            <span class="info-value">{{ conn.Database }}</span>
          </div>
        </div>
        <div class="connection-footer">
          <button @click.stop="testConnection(conn)" class="test-btn" :disabled="testingConnection === conn.ID">
            {{ testingConnection === conn.ID ? '测试中...' : '测试连接' }}
          </button>
          <span v-if="connectionStatus[conn.ID]" class="status-badge" :class="connectionStatus[conn.ID].success ? 'success' : 'error'">
            {{ connectionStatus[conn.ID].message }}
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
            <input v-model="formData.Name" type="text" placeholder="输入连接名称" />
          </div>
          <div class="form-row">
            <div class="form-group flex-1">
              <label>主机地址 *</label>
              <input v-model="formData.Host" type="text" placeholder="localhost" />
            </div>
            <div class="form-group" style="width: 100px;">
              <label>端口 *</label>
              <input v-model.number="formData.Port" type="number" placeholder="3306" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group flex-1">
              <label>用户名 *</label>
              <input v-model="formData.User" type="text" placeholder="root" />
            </div>
            <div class="form-group flex-1">
              <label>密码</label>
              <input v-model="formData.Password" type="password" placeholder="可选" />
            </div>
          </div>
          <div class="form-group">
            <label>默认数据库</label>
            <input v-model="formData.Database" type="text" placeholder="可选" />
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
  ID: '',
  Name: '',
  Host: 'localhost',
  Port: 3306,
  User: 'root',
  Password: '',
  Database: ''
})

// 方法
function refreshConnections() {
  fetchConnections()
}

function selectConnection(conn) {
  storeSelectConnection(conn.ID)
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
  testingConnection.value = conn.ID
  try {
    const success = await databaseStore.testConnection(conn)
    connectionStatus.value[conn.ID] = {
      success,
      message: success ? '连接成功' : '连接失败'
    }
  } catch (err) {
    connectionStatus.value[conn.ID] = {
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
    ID: '',
    Name: '',
    Host: 'localhost',
    Port: 3306,
    User: 'root',
    Password: '',
    Database: ''
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
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #fff;
  cursor: pointer;
  transition: all 0.3s;
}

.add-btn:hover {
  border-color: #00ff00;
  color: #00ff00;
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
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 1.5rem;
  cursor: pointer;
  transition: all 0.3s;
}

.connection-card:hover {
  border-color: #00ff00;
}

.connection-card.active {
  border-color: #00ff00;
  background: rgba(0, 255, 0, 0.05);
}

.connection-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.connection-name {
  color: #fff;
  margin: 0;
}

.connection-actions {
  display: flex;
  gap: 0.5rem;
}

.edit-btn,
.delete-btn {
  background: transparent;
  border: 1px solid #333;
  border-radius: 4px;
  padding: 0.25rem 0.75rem;
  color: #888;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.3s;
}

.edit-btn:hover {
  border-color: #00ccff;
  color: #00ccff;
}

.delete-btn:hover {
  border-color: #ff4444;
  color: #ff4444;
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
  color: #888;
  min-width: 60px;
}

.info-value {
  color: #ccc;
}

.connection-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.test-btn {
  background: transparent;
  border: 1px solid #00ff00;
  border-radius: 4px;
  padding: 0.25rem 0.75rem;
  color: #00ff00;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.3s;
}

.test-btn:hover:not(:disabled) {
  background: rgba(0, 255, 0, 0.1);
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
  background: rgba(0, 255, 0, 0.2);
  color: #00ff00;
}

.status-badge.error {
  background: rgba(255, 68, 68, 0.2);
  color: #ff4444;
}

.no-data {
  text-align: center;
  padding: 3rem;
  color: #666;
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
  border: 3px solid #1a1a1a;
  border-top-color: #00ff00;
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
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #333;
}

.dialog-header h3 {
  color: #fff;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: #888;
  font-size: 1.5rem;
  cursor: pointer;
}

.close-btn:hover {
  color: #fff;
}

.dialog-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  color: #888;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.form-group input {
  width: 100%;
  background: #0d0d0d;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #fff;
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: #00ff00;
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
  border-top: 1px solid #333;
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
  border: 1px solid #333;
  color: #888;
}

.cancel-btn:hover {
  border-color: #666;
  color: #fff;
}

.save-btn {
  background: #00ff00;
  border: none;
  color: #000;
}

.save-btn:hover:not(:disabled) {
  background: #00cc00;
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>