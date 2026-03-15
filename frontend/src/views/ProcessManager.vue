<template>
  <div class="process-manager">
    <div class="header">
      <h1>进程管理器</h1>
      <p class="subtitle">实时监控系统进程</p>
    </div>

    <div class="toolbar">
      <div class="search-box">
        <input
          v-model="searchKeyword"
          type="text"
          placeholder="搜索进程 (名称/PID/命令行/端口)"
          class="search-input"
        />
        <button @click="refreshProcesses" class="refresh-btn" :disabled="loading">
          {{ loading ? '刷新中...' : '刷新' }}
        </button>
      </div>

      <div class="stats">
        <span class="stat-item">总进程数: {{ totalCount }}</span>
        <span class="stat-item">选中: {{ selectedCount }}</span>
      </div>
    </div>

    <div class="process-list" v-if="!loading">
      <table class="process-table">
        <thead>
          <tr>
            <th class="col-checkbox">
              <input
                type="checkbox"
                :checked="allSelected"
                @change="toggleSelectAll"
              />
            </th>
            <th class="col-pid" @click="setSortBy('PID')">PID</th>
            <th class="col-name" @click="setSortBy('Name')">进程名</th>
            <th class="col-cpu" @click="setSortBy('CPUPercent')">CPU %</th>
            <th class="col-memory" @click="setSortBy('MemoryMB')">内存 (MB)</th>
            <th class="col-port" @click="setSortBy('ListenPort')">端口</th>
            <th class="col-status" @click="setSortBy('Status')">状态</th>
            <th class="col-actions">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="process in filteredProcesses"
            :key="process.PID"
            class="process-row"
            :class="{ selected: isSelected(process.PID) }"
          >
            <td class="col-checkbox">
              <input
                type="checkbox"
                :checked="isSelected(process.PID)"
                @change="toggleSelect(process.PID, $event.target.checked)"
              />
            </td>
            <td class="col-pid">{{ process.PID }}</td>
            <td class="col-name" :title="process.Cmdline">
              {{ process.Name || '(无名称)' }}
            </td>
            <td class="col-cpu">
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: process.CPUPercent + '%' }"
                ></div>
                <span class="progress-text">{{ process.CPUPercent.toFixed(2) }}%</span>
              </div>
            </td>
            <td class="col-memory">
              <div class="progress-bar">
                <div
                  class="progress-fill memory"
                  :style="{ width: Math.min(process.MemoryMB / 100, 100) + '%' }"
                ></div>
                <span class="progress-text">{{ process.MemoryMB }} MB</span>
              </div>
            </td>
            <td class="col-port">
              {{ process.ListenPort || '-' }}
            </td>
            <td class="col-status">
              <span class="status-badge" :class="getStatusClass(process.Status)">
                {{ process.Status || '未知' }}
              </span>
            </td>
            <td class="col-actions">
              <button
                @click="killProcess(process.PID)"
                class="kill-btn"
                :disabled="killingProcesses.includes(process.PID)"
              >
                {{ killingProcesses.includes(process.PID) ? '终止中...' : '终止' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="filteredProcesses.length === 0" class="no-data">
        没有找到匹配的进程
      </div>
    </div>

    <div v-else class="loading">
      <div class="loading-spinner"></div>
      <p>加载进程信息...</p>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div class="batch-actions" v-if="selectedCount > 0">
      <button @click="killSelectedProcesses" class="batch-kill-btn">
        终止选中的 {{ selectedCount }} 个进程
      </button>
      <button @click="clearSelection" class="clear-selection-btn">
        清除选择
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useProcessStore } from '../stores/process'

const processStore = useProcessStore()

// 响应式数据
const killingProcesses = ref([])

// 从store获取响应式refs（包括searchKeyword）
const {
  filteredProcesses,
  loading,
  error,
  selectedCount,
  totalCount,
  sortBy,
  sortOrder,
  selectedProcesses,
  searchKeyword
} = storeToRefs(processStore)

// 从store获取方法
const {
  fetchProcesses,
  killProcess: storeKillProcess,
  selectProcess,
  selectAllProcesses,
  clearSelection: storeClearSelection,
  setSortBy: storeSetSortBy
} = processStore

// 计算属性
const allSelected = computed(() => {
  if (!filteredProcesses.value || !selectedProcesses.value) return false
  return filteredProcesses.value.length > 0 && 
         selectedProcesses.value.length === filteredProcesses.value.length
})

// 方法
function refreshProcesses() {
  fetchProcesses()
}

function toggleSelectAll(event) {
  selectAllProcesses(event.target.checked)
}

function toggleSelect(pid, selected) {
  selectProcess(pid, selected)
}

function isSelected(pid) {
  if (!selectedProcesses.value) return false
  return selectedProcesses.value.includes(pid)
}

async function killProcess(pid) {
  killingProcesses.value.push(pid)
  await storeKillProcess(pid)
  killingProcesses.value = killingProcesses.value.filter(id => id !== pid)
}

async function killSelectedProcesses() {
  const pids = [...selectedProcesses.value]
  for (const pid of pids) {
    await killProcess(pid)
  }
}

function clearSelection() {
  storeClearSelection()
}

function setSortBy(field) {
  storeSetSortBy(field)
}

function getStatusClass(status) {
  const statusMap = {
    'LISTEN': 'status-listen',
    'ESTABLISHED': 'status-established',
    'TIME_WAIT': 'status-time-wait',
    'CLOSE_WAIT': 'status-close-wait'
  }
  return statusMap[status] || 'status-unknown'
}

// 组件挂载时获取进程列表
onMounted(() => {
  fetchProcesses()
})
</script>

<style scoped>
.process-manager {
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

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.search-box {
  display: flex;
  gap: 0.5rem;
  flex: 1;
}

.search-input {
  flex: 1;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #fff;
  font-family: inherit;
}

.search-input:focus {
  outline: none;
  border-color: #00ff00;
}

.refresh-btn {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #00ff00;
  cursor: pointer;
  transition: all 0.3s;
}

.refresh-btn:hover:not(:disabled) {
  border-color: #00ff00;
  background: rgba(0, 255, 0, 0.1);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.stats {
  display: flex;
  gap: 1.5rem;
}

.stat-item {
  color: #888;
  font-size: 0.9rem;
}

.process-list {
  flex: 1;
  overflow: auto;
  background: #0d0d0d;
  border-radius: 8px;
  border: 1px solid #1a1a1a;
}

.process-table {
  width: 100%;
  border-collapse: collapse;
}

.process-table th {
  position: sticky;
  top: 0;
  background: #1a1a1a;
  padding: 1rem;
  text-align: left;
  color: #888;
  font-weight: normal;
  font-size: 0.85rem;
  border-bottom: 1px solid #333;
  cursor: pointer;
}

.process-table th:hover {
  color: #00ff00;
}

.process-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #1a1a1a;
  color: #ccc;
}

.process-row:hover {
  background: rgba(0, 255, 0, 0.05);
}

.process-row.selected {
  background: rgba(0, 255, 0, 0.1);
}

.col-checkbox {
  width: 40px;
  text-align: center;
}

.col-pid {
  width: 80px;
}

.col-cpu,
.col-memory {
  width: 150px;
}

.col-port {
  width: 100px;
}

.col-status {
  width: 120px;
}

.col-actions {
  width: 100px;
}

.progress-bar {
  position: relative;
  height: 20px;
  background: #1a1a1a;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, #00ff00, #00cc00);
  transition: width 0.3s;
}

.progress-fill.memory {
  background: linear-gradient(90deg, #00ccff, #0099cc);
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 0.75rem;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
}

.status-listen {
  background: rgba(0, 255, 0, 0.2);
  color: #00ff00;
}

.status-established {
  background: rgba(0, 204, 255, 0.2);
  color: #00ccff;
}

.status-time-wait {
  background: rgba(255, 193, 7, 0.2);
  color: #ffc107;
}

.status-close-wait {
  background: rgba(255, 87, 34, 0.2);
  color: #ff5722;
}

.status-unknown {
  background: rgba(136, 136, 136, 0.2);
  color: #888;
}

.kill-btn {
  background: transparent;
  border: 1px solid #ff4444;
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
  color: #ff4444;
  cursor: pointer;
  font-size: 0.8rem;
  transition: all 0.3s;
}

.kill-btn:hover:not(:disabled) {
  background: rgba(255, 68, 68, 0.1);
}

.kill-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.error-message {
  background: rgba(255, 68, 68, 0.1);
  border: 1px solid #ff4444;
  border-radius: 6px;
  padding: 1rem;
  margin-top: 1rem;
  color: #ff4444;
}

.batch-actions {
  position: fixed;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 1rem;
  background: #1a1a1a;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #333;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
}

.batch-kill-btn {
  background: #ff4444;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #fff;
  cursor: pointer;
  transition: all 0.3s;
}

.batch-kill-btn:hover {
  background: #ff6666;
}

.clear-selection-btn {
  background: transparent;
  border: 1px solid #666;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #888;
  cursor: pointer;
  transition: all 0.3s;
}

.clear-selection-btn:hover {
  border-color: #00ff00;
  color: #00ff00;
}

input[type="checkbox"] {
  appearance: none;
  width: 16px;
  height: 16px;
  border: 1px solid #333;
  border-radius: 3px;
  background: #1a1a1a;
  cursor: pointer;
}

input[type="checkbox"]:checked {
  background: #00ff00;
  border-color: #00ff00;
}
</style>