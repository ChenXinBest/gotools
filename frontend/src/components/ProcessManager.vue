<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { GetSystemProcessInfos, SearchPidByKeyWord, KillProcessByPID } from '../../bindings/gotools/processservice.js'

const processes = ref([])
const searchKeyword = ref('')
const loading = ref(false)
const error = ref('')
const searchResult = ref(null)
const selectedPIDs = ref(new Set())
const sortKey = ref('MemoryMB')
const sortOrder = ref('desc')
const contextMenu = ref({ show: false, x: 0, y: 0, pid: null, type: 'process' })
const killConfirm = ref({ show: false, pids: [] })
const statusMessage = ref('')
const expandedApps = ref(new Set())
const columnWidths = ref({
  Name: '180px',
  PID: '80px',
  CPUPercent: '80px',
  MemoryMB: '90px',
  ListenAddr: '120px',
  ListenPort: '70px',
  Status: '80px'
})
const resizingColumn = ref(null)
const startX = ref(0)
const startWidth = ref(0)
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })
const dragEnd = ref({ x: 0, y: 0 })
const dragSelection = ref(null)

const processTree = computed(() => {
  let list = [...processes.value]
  
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.toLowerCase()
    list = list.filter(p => 
      p.Name?.toLowerCase().includes(keyword) ||
      p.Cmdline?.toLowerCase().includes(keyword) ||
      String(p.PID).includes(keyword) ||
      p.ListenAddr?.toLowerCase().includes(keyword) ||
      String(p.ListenPort).includes(keyword)
    )
  }
  
  const appMap = new Map()
  
  list.forEach(proc => {
    const appName = proc.Name || '未知进程'
    if (!appMap.has(appName)) {
      appMap.set(appName, {
        name: appName,
        processes: [],
        totalCPU: 0,
        maxCPU: 0,
        totalMemory: 0,
        totalProcesses: 0
      })
    }
    
    const app = appMap.get(appName)
    app.processes.push(proc)
    app.totalCPU += proc.CPUPercent
    if (proc.CPUPercent > app.maxCPU) {
      app.maxCPU = proc.CPUPercent
    }
    app.totalMemory += proc.MemoryMB
    app.totalProcesses++
  })
  
  const apps = Array.from(appMap.values())
  
  apps.sort((a, b) => {
    let aVal, bVal
    if (sortKey.value === 'Name') {
      aVal = a.name.toLowerCase()
      bVal = b.name.toLowerCase()
    } else if (sortKey.value === 'MemoryMB') {
      aVal = a.totalMemory
      bVal = b.totalMemory
    } else if (sortKey.value === 'CPUPercent') {
      aVal = a.maxCPU
      bVal = b.maxCPU
    } else if (sortKey.value === 'PID') {
      aVal = a.processes[0]?.PID || 0
      bVal = b.processes[0]?.PID || 0
    }
    
    if (aVal < bVal) return sortOrder.value === 'asc' ? -1 : 1
    if (aVal > bVal) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })
  
  apps.forEach(app => {
    app.processes.sort((a, b) => {
      let aVal, bVal
      if (sortKey.value === 'Name') {
        aVal = (a.Cmdline ? a.Cmdline.split(' ')[0] : '').toLowerCase()
        bVal = (b.Cmdline ? b.Cmdline.split(' ')[0] : '').toLowerCase()
      } else if (sortKey.value === 'MemoryMB') {
        aVal = a.MemoryMB
        bVal = b.MemoryMB
      } else if (sortKey.value === 'CPUPercent') {
        aVal = a.CPUPercent
        bVal = b.CPUPercent
      } else if (sortKey.value === 'PID') {
        aVal = a.PID
        bVal = b.PID
      }
      
      if (aVal < bVal) return sortOrder.value === 'asc' ? -1 : 1
      if (aVal > bVal) return sortOrder.value === 'asc' ? 1 : -1
      return 0
    })
  })
  
  return apps
})

const columns = [
  { key: 'Name', label: '应用名称', width: '180px' },
  { key: 'PID', label: '进程数', width: '80px' },
  { key: 'CPUPercent', label: 'CPU', width: '80px' },
  { key: 'MemoryMB', label: '内存', width: '90px' },
  { key: 'ListenAddr', label: '监听地址', width: '120px' },
  { key: 'ListenPort', label: '端口', width: '70px' },
  { key: 'Status', label: '状态', width: '80px' }
]

function handleSort(key) {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
}

function toggleApp(appName) {
  if (expandedApps.value.has(appName)) {
    expandedApps.value.delete(appName)
  } else {
    expandedApps.value.add(appName)
  }
  expandedApps.value = new Set(expandedApps.value)
}

function isAppExpanded(appName) {
  return expandedApps.value.has(appName)
}

async function loadProcesses() {
  loading.value = true
  error.value = ''
  try {
    const result = await GetSystemProcessInfos()
    processes.value = result || []
    showStatus(`已加载 ${processes.value.length} 个进程`)
  } catch (e) {
    error.value = '加载进程信息失败: ' + (e.message || e)
  } finally {
    loading.value = false
  }
}

async function searchProcess() {
  if (!searchKeyword.value.trim()) {
    searchResult.value = null
    return
  }
  loading.value = true
  error.value = ''
  try {
    const result = await SearchPidByKeyWord(searchKeyword.value)
    searchResult.value = result
  } catch (e) {
    searchResult.value = null
    error.value = '未找到匹配进程: ' + searchKeyword.value
  } finally {
    loading.value = false
  }
}

function clearSearch() {
  searchKeyword.value = ''
  searchResult.value = null
  error.value = ''
}

function showStatus(msg) {
  statusMessage.value = msg
  setTimeout(() => {
    statusMessage.value = ''
  }, 3000)
}

function handleRowClick(e, pid, type = 'process') {
  if (e.ctrlKey) {
    if (selectedPIDs.value.has(pid)) {
      selectedPIDs.value.delete(pid)
    } else {
      selectedPIDs.value.add(pid)
    }
  } else if (e.shiftKey && selectedPIDs.value.size > 0) {
    const lastPid = Array.from(selectedPIDs.value).pop()
    const allProcesses = processTree.value.flatMap(app => app.processes)
    const lastIndex = allProcesses.findIndex(p => p.PID === lastPid)
    const currentIndex = allProcesses.findIndex(p => p.PID === pid)
    if (lastIndex > -1 && currentIndex > -1) {
      const [start, end] = [Math.min(lastIndex, currentIndex), Math.max(lastIndex, currentIndex)]
      for (let i = start; i <= end; i++) {
        selectedPIDs.value.add(allProcesses[i].PID)
      }
    }
  } else {
    selectedPIDs.value.clear()
    selectedPIDs.value.add(pid)
  }
  selectedPIDs.value = new Set(selectedPIDs.value)
}

function handleAppClick(app, e) {
  if (e.target.closest('.expand-icon')) {
    toggleApp(app.name)
  }
}

function handleRightClick(e, pid, type = 'process') {
  e.preventDefault()
  if (!selectedPIDs.value.has(pid)) {
    selectedPIDs.value.clear()
    selectedPIDs.value.add(pid)
    selectedPIDs.value = new Set(selectedPIDs.value)
  }
  contextMenu.value = {
    show: true,
    x: e.clientX,
    y: e.clientY,
    pid: pid,
    type: type
  }
}

function closeContextMenu() {
  contextMenu.value.show = false
}

function showKillConfirm() {
  killConfirm.value = {
    show: true,
    pids: Array.from(selectedPIDs.value)
  }
  closeContextMenu()
}

async function confirmKill() {
  const pids = killConfirm.value.pids
  loading.value = true
  let success = 0
  let failed = 0
  
  for (const pid of pids) {
    try {
      await KillProcessByPID(pid)
      success++
    } catch (e) {
      failed++
    }
  }
  
  killConfirm.value.show = false
  selectedPIDs.value.clear()
  
  if (success > 0) {
    showStatus(`成功终止 ${success} 个进程${failed > 0 ? `，失败 ${failed} 个` : ''}`)
    await loadProcesses()
  } else {
    error.value = '终止进程失败'
  }
  
  loading.value = false
}

function cancelKill() {
  killConfirm.value.show = false
}

function formatMemory(mb) {
  if (mb === 0) return '0'
  if (mb < 1024) return mb.toFixed(0) + 'M'
  return (mb / 1024).toFixed(1) + 'G'
}

function formatCpu(cpu) {
  return cpu.toFixed(1) + '%'
}

function formatPort(port) {
  return port || '-' 
}

function getCpuClass(cpu) {
  if (cpu > 50) return 'cpu-critical'
  if (cpu > 20) return 'cpu-warning'
  return 'cpu-normal'
}

function getMemClass(mb) {
  if (mb > 1024) return 'mem-critical'
  if (mb > 512) return 'mem-warning'
  return 'mem-normal'
}

function startResize(e, columnKey) {
  resizingColumn.value = columnKey
  startX.value = e.clientX
  startWidth.value = parseInt(columnWidths.value[columnKey])
  document.addEventListener('mousemove', resizeColumn)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = 'col-resize'
}

function resizeColumn(e) {
  if (!resizingColumn.value) return
  const width = startWidth.value + (e.clientX - startX.value)
  if (width > 50) { // 最小宽度限制
    columnWidths.value[resizingColumn.value] = width + 'px'
  }
}

function stopResize() {
  resizingColumn.value = null
  document.removeEventListener('mousemove', resizeColumn)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
}

function startDrag(e) {
  if (e.target.closest('.tr')) {
    isDragging.value = true
    dragStart.value = { x: e.clientX, y: e.clientY }
    dragEnd.value = { x: e.clientX, y: e.clientY }
    dragSelection.value = {
      left: e.clientX,
      top: e.clientY,
      width: 0,
      height: 0
    }
    document.addEventListener('mousemove', onDrag)
    document.addEventListener('mouseup', endDrag)
  }
}

function onDrag(e) {
  if (!isDragging.value) return
  dragEnd.value = { x: e.clientX, y: e.clientY }
  dragSelection.value = {
    left: Math.min(dragStart.value.x, e.clientX),
    top: Math.min(dragStart.value.y, e.clientY),
    width: Math.abs(e.clientX - dragStart.value.x),
    height: Math.abs(e.clientY - dragStart.value.y)
  }
}

function endDrag() {
  if (!isDragging.value) return
  
  // 计算拖拽区域内的进程
  const selectionRect = dragSelection.value
  const processElements = document.querySelectorAll('.tr.process-row')
  
  processElements.forEach(el => {
    const rect = el.getBoundingClientRect()
    if (isElementInSelection(rect, selectionRect)) {
      const pid = parseInt(el.dataset.pid)
      if (!isNaN(pid)) {
        selectedPIDs.value.add(pid)
      }
    }
  })
  
  selectedPIDs.value = new Set(selectedPIDs.value)
  isDragging.value = false
  dragSelection.value = null
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', endDrag)
}

function isElementInSelection(elementRect, selectionRect) {
  return (
    elementRect.left < selectionRect.left + selectionRect.width &&
    elementRect.right > selectionRect.left &&
    elementRect.top < selectionRect.top + selectionRect.height &&
    elementRect.bottom > selectionRect.top
  )
}

function handleClickOutside(e) {
  if (!e.target.closest('.context-menu') && !e.target.closest('.tr')) {
    closeContextMenu()
  }
}

// 键盘快捷键处理
function handleKeydown(e) {
  // 忽略输入框中的按键
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
    return
  }

  // F5 - 刷新进程
  if (e.key === 'F5') {
    e.preventDefault()
    loadProcesses()
    showStatus('已刷新')
    return
  }

  // Escape - 关闭弹窗/清除选择
  if (e.key === 'Escape') {
    if (killConfirm.value.show) {
      cancelKill()
    } else if (contextMenu.value.show) {
      closeContextMenu()
    } else if (selectedPIDs.value.size > 0) {
      selectedPIDs.value.clear()
      selectedPIDs.value = new Set()
    }
    return
  }

  // Delete - 终止选中进程
  if (e.key === 'Delete') {
    if (selectedPIDs.value.size > 0 && !killConfirm.value.show) {
      showKillConfirm()
    }
    return
  }

  // Ctrl+A - 全选
  if (e.ctrlKey && e.key === 'a') {
    e.preventDefault()
    const allProcesses = processTree.value.flatMap(app => app.processes)
    allProcesses.forEach(p => selectedPIDs.value.add(p.PID))
    selectedPIDs.value = new Set(selectedPIDs.value)
    showStatus(`已选择 ${selectedPIDs.value.size} 个进程`)
    return
  }
}

onMounted(() => {
  loadProcesses()
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <div class="process-manager">
    <div class="pm-header">
      <div class="pm-title">
        <span class="title-icon">◉</span>
        进程管理器
      </div>
      <div class="pm-stats">
        <span class="stat-item">{{ processes.length }} 进程</span>
        <span class="stat-divider">|</span>
        <span class="stat-item">{{ selectedPIDs.size }} 已选</span>
      </div>
    </div>

    <div class="pm-toolbar">
      <div class="search-box">
        <span class="search-icon">$</span>
        <input 
          v-model="searchKeyword" 
          type="text" 
          placeholder="搜索进程..." 
          class="search-input"
          @keyup.enter="searchProcess"
        />
      </div>
      <div class="toolbar-btns">
        <button @click="searchProcess" class="btn" :disabled="loading">
          搜索
        </button>
        <button @click="loadProcesses" class="btn" :disabled="loading">
          刷新
        </button>
        <button @click="clearSearch" class="btn warn" v-if="searchKeyword">
          清除
        </button>
        <button 
          v-if="selectedPIDs.size > 0" 
          @click="showKillConfirm" 
          class="btn danger"
        >
          终止 [{{ selectedPIDs.size }}]
        </button>
      </div>
    </div>

    <div v-if="statusMessage" class="msg-bar success">
      <span class="msg-icon">✓</span>
      {{ statusMessage }}
    </div>

    <div v-if="error" class="msg-bar error">
      <span class="msg-icon">✗</span>
      {{ error }}
    </div>

    <div v-if="searchResult" class="result-panel">
      <div class="result-header">
        <span class="result-icon">▶</span>
        精确匹配结果
      </div>
      <div class="result-content">
        <div class="result-row">
          <div class="result-item">
            <span class="label">PID</span>
            <span class="value pid">{{ searchResult.PID }}</span>
          </div>
          <div class="result-item">
            <span class="label">名称</span>
            <span class="value">{{ searchResult.Name || '-' }}</span>
          </div>
          <div class="result-item">
            <span class="label">CPU</span>
            <span class="value" :class="getCpuClass(searchResult.CPUPercent)">
              {{ formatCpu(searchResult.CPUPercent) }}
            </span>
          </div>
          <div class="result-item">
            <span class="label">内存</span>
            <span class="value" :class="getMemClass(searchResult.MemoryMB)">
              {{ formatMemory(searchResult.MemoryMB) }}
            </span>
          </div>
        </div>
        <div class="result-row">
          <div class="result-item full">
            <span class="label">命令行</span>
            <span class="value mono">{{ searchResult.Cmdline || '-' }}</span>
          </div>
        </div>
        <div class="result-row">
          <div class="result-item">
            <span class="label">监听地址</span>
            <span class="value">{{ searchResult.ListenAddr || '-' }}</span>
          </div>
          <div class="result-item">
            <span class="label">端口</span>
            <span class="value">{{ formatPort(searchResult.ListenPort) }}</span>
          </div>
          <div class="result-item">
            <span class="label">状态</span>
            <span class="value">{{ searchResult.Status || '-' }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="process-table" @mousedown="startDrag">
      <div class="table-head">
        <div 
          v-for="col in columns" 
          :key="col.key"
          class="th"
          :style="{ width: columnWidths[col.key] }"
          @click="handleSort(col.key)"
        >
          {{ col.label }}
          <span v-if="sortKey === col.key" class="sort-arrow">
            {{ sortOrder === 'asc' ? '▲' : '▼' }}
          </span>
          <div 
            class="resize-handle" 
            @mousedown="startResize($event, col.key)"
          ></div>
        </div>
      </div>

      <div class="table-body">
        <div v-if="loading && processes.length === 0" class="loading-state">
          <span class="cursor-blink">█</span> 加载中...
        </div>

        <div v-else-if="processTree.length === 0" class="empty-state">
          无匹配进程
        </div>

        <template v-else>
          <div 
            v-for="app in processTree" 
            :key="app.name"
          >
            <div 
              class="tr app-row"
              @click="handleAppClick(app, $event)"
              @contextmenu="handleRightClick($event, app.processes[0]?.PID, 'app')"
            >
              <div class="td" :style="{ width: columnWidths.Name }">
                <span class="expand-icon" :class="{ expanded: isAppExpanded(app.name) }">▶</span>
                <span class="app-name">{{ app.name }}</span>
                <span class="process-count">({{ app.totalProcesses }}个进程)</span>
              </div>
              <div class="td num" :style="{ width: columnWidths.PID }">
                {{ app.totalProcesses }}
              </div>
              <div class="td num" :style="{ width: columnWidths.CPUPercent }" :class="getCpuClass(app.maxCPU)">
                {{ formatCpu(app.maxCPU) }}
              </div>
              <div class="td num" :style="{ width: columnWidths.MemoryMB }" :class="getMemClass(app.totalMemory)">
                {{ formatMemory(app.totalMemory) }}
              </div>
              <div class="td" :style="{ width: columnWidths.ListenAddr }">
                -</div>
              <div class="td num" :style="{ width: columnWidths.ListenPort }">
                -</div>
              <div class="td" :style="{ width: columnWidths.Status }">
                -</div>
            </div>
            
            <div v-if="isAppExpanded(app.name)" class="app-children">
              <div 
                v-for="proc in app.processes" 
                :key="proc.PID" 
                class="tr process-row"
                :class="{ selected: selectedPIDs.has(proc.PID) }"
                :data-pid="proc.PID"
                @click="handleRowClick($event, proc.PID)"
                @contextmenu="handleRightClick($event, proc.PID)"
              >
                <div class="td process-name" :style="{ width: columnWidths.Name }">
                  <span class="indent"></span>
                  <span class="process-exe" :title="proc.Cmdline || ''">{{ proc.Cmdline ? proc.Cmdline.split(' ')[0] : '-' }}</span>
                </div>
                <div class="td num" :style="{ width: columnWidths.PID }">
                  {{ proc.PID }}
                </div>
                <div class="td num" :style="{ width: columnWidths.CPUPercent }" :class="getCpuClass(proc.CPUPercent)">
                  {{ formatCpu(proc.CPUPercent) }}
                </div>
                <div class="td num" :style="{ width: columnWidths.MemoryMB }" :class="getMemClass(proc.MemoryMB)">
                  {{ formatMemory(proc.MemoryMB) }}
                </div>
                <div class="td" :style="{ width: columnWidths.ListenAddr }">
                  {{ proc.ListenAddr || '-' }}
                </div>
                <div class="td num" :style="{ width: columnWidths.ListenPort }">
                  {{ formatPort(proc.ListenPort) }}
                </div>
                <div class="td" :style="{ width: columnWidths.Status }">
                  <span v-if="proc.Status" class="status-dot"></span>
                  {{ proc.Status || '-' }}
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <div 
      v-if="dragSelection" 
      class="drag-selection"
      :style="{
        left: dragSelection.left + 'px',
        top: dragSelection.top + 'px',
        width: dragSelection.width + 'px',
        height: dragSelection.height + 'px'
      }"
    ></div>

    <div class="pm-footer">
      显示 {{ processTree.length }} 个应用 / {{ processes.length }} 个进程
      <span v-if="searchKeyword"> (已筛选)</span>
    </div>

    <div 
      v-if="contextMenu.show" 
      class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <div class="ctx-item" @click="showKillConfirm">
        <span class="ctx-icon">✕</span>
        终止进程
        <span v-if="selectedPIDs.size > 1">[{{ selectedPIDs.size }}个]</span>
      </div>
      <div class="ctx-divider"></div>
      <div class="ctx-item" @click="selectedPIDs.clear(); closeContextMenu()">
        <span class="ctx-icon">○</span>
        取消选择
      </div>
    </div>

    <div v-if="killConfirm.show" class="modal-mask" @click.self="cancelKill">
      <div class="modal-box">
        <div class="modal-head">
          <span class="modal-icon warn">⚠</span>
          确认终止进程
        </div>
        <div class="modal-body">
          <p>确定要终止 {{ killConfirm.pids.length }} 个进程吗？</p>
          <div class="pid-tags">
            <span v-for="pid in killConfirm.pids.slice(0, 10)" :key="pid" class="pid-tag">
              {{ pid }}
            </span>
            <span v-if="killConfirm.pids.length > 10" class="pid-tag more">
              +{{ killConfirm.pids.length - 10 }}
            </span>
          </div>
        </div>
        <div class="modal-foot">
          <button @click="cancelKill" class="btn">取消</button>
          <button @click="confirmKill" class="btn danger">终止</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.process-manager {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #0a0a0a;
  color: #00ff00;
  font-family: 'Consolas', 'Monaco', monospace;
  overflow: hidden;
}

.pm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #0d0d0d;
  border-bottom: 1px solid #1a1a1a;
}

.pm-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: bold;
  letter-spacing: 1px;
}

.title-icon {
  color: #00ff00;
  text-shadow: 0 0 10px #00ff00;
}

.pm-stats {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
  color: #666;
}

.stat-divider {
  color: #333;
}

.pm-toolbar {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 12px 20px;
  background: #0d0d0d;
  border-bottom: 1px solid #1a1a1a;
}

.search-box {
  flex: 1;
  display: flex;
  align-items: center;
  background: #0a0a0a;
  border: 1px solid #1a1a1a;
  border-radius: 4px;
  padding: 0 10px;
}

.search-icon {
  color: #00ff00;
  margin-right: 8px;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  color: #00ff00;
  font-family: inherit;
  font-size: 13px;
  padding: 8px 0;
  outline: none;
}

.search-input::placeholder {
  color: #333;
}

.toolbar-btns {
  display: flex;
  gap: 8px;
}

.btn {
  background: transparent;
  border: 1px solid #00ff00;
  color: #00ff00;
  padding: 6px 14px;
  font-family: inherit;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 3px;
}

.btn:hover:not(:disabled) {
  background: #00ff00;
  color: #0a0a0a;
}

.btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.btn.warn {
  border-color: #ffbd2e;
  color: #ffbd2e;
}

.btn.warn:hover {
  background: #ffbd2e;
  color: #0a0a0a;
}

.btn.danger {
  border-color: #ff5f56;
  color: #ff5f56;
}

.btn.danger:hover {
  background: #ff5f56;
  color: #0a0a0a;
}

.msg-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 20px;
  font-size: 12px;
}

.msg-bar.success {
  background: rgba(0, 255, 0, 0.1);
  border-bottom: 1px solid #00ff00;
}

.msg-bar.error {
  background: rgba(255, 95, 86, 0.1);
  border-bottom: 1px solid #ff5f56;
  color: #ff5f56;
}

.msg-icon {
  font-weight: bold;
}

.result-panel {
  margin: 15px 20px;
  background: #0d0d0d;
  border: 1px solid #00ff00;
  border-radius: 4px;
  overflow: hidden;
}

.result-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(0, 255, 0, 0.1);
  font-size: 11px;
  letter-spacing: 1px;
}

.result-icon {
  color: #00ff00;
}

.result-content {
  padding: 12px;
}

.result-row {
  display: flex;
  gap: 20px;
  margin-bottom: 8px;
}

.result-row:last-child {
  margin-bottom: 0;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.result-item.full {
  flex: 1;
}

.result-item .label {
  color: #444;
  font-size: 10px;
  letter-spacing: 1px;
}

.result-item .value {
  color: #00ff00;
  font-size: 12px;
}

.result-item .value.pid {
  color: #fff;
  text-shadow: 0 0 5px #00ff00;
}

.result-item .value.mono {
  font-size: 11px;
  word-break: break-all;
}

.process-table {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin: 0 20px;
  background: #0d0d0d;
  border: 1px solid #1a1a1a;
  border-radius: 4px;
  overflow: hidden;
}

.table-head {
  display: flex;
  background: #1a1a1a;
  border-bottom: 1px solid #333;
}

.th {
  padding: 10px 12px;
  font-size: 11px;
  color: #666;
  letter-spacing: 1px;
  cursor: pointer;
  user-select: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.th:hover {
  color: #00ff00;
}

.sort-arrow {
  margin-left: 4px;
  color: #00ff00;
  flex-shrink: 0;
}

.resize-handle {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  cursor: col-resize;
  background: transparent;
  transition: background 0.2s;
}

.resize-handle:hover {
  background: rgba(0, 255, 0, 0.3);
}

.table-body {
  flex: 1;
  overflow-y: auto;
}

.tr {
  display: flex;
  border-bottom: 1px solid #1a1a1a;
  cursor: pointer;
  transition: background 0.1s;
}

.tr:hover {
  background: rgba(0, 255, 0, 0.05);
}

.tr.selected {
  background: rgba(0, 255, 0, 0.12);
  border-left: 2px solid #00ff00;
  margin-left: -2px;
}

.tr.app-row {
  background: rgba(0, 255, 0, 0.02);
  font-weight: bold;
}

.tr.app-row:hover {
  background: rgba(0, 255, 0, 0.08);
}

.tr.process-row {
  border-left: 4px solid transparent;
}

.tr.process-row:hover {
  border-left-color: rgba(0, 255, 0, 0.3);
}

.td {
  padding: 8px 12px;
  font-size: 12px;
  color: #888;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: flex;
  align-items: center;
}

.td.num {
  text-align: right;
  font-variant-numeric: tabular-nums;
  display: flex;
  justify-content: flex-end;
}

.td.process-name {
  display: flex;
  align-items: center;
}

.expand-icon {
  display: inline-block;
  width: 20px;
  text-align: left;
  font-size: 8px;
  margin-right: 0;
  transition: transform 0.2s;
  cursor: pointer;
  flex-shrink: 0;
  padding-left: 4px;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.app-name {
  color: #00ff00;
  margin-right: 8px;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.process-count {
  color: #666;
  font-size: 10px;
  font-weight: normal;
  white-space: nowrap;
}

.indent {
  display: inline-block;
  width: 20px;
  flex-shrink: 0;
  padding-left: 4px;
}

.process-exe {
  color: #ccc;
  font-size: 11px;
  font-weight: normal;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.app-children {
  border-left: 2px solid #1a1a1a;
  margin-left: 10px;
}

.cpu-normal, .mem-normal { color: #888; }
.cpu-warning { color: #ffbd2e; }
.cpu-critical { color: #ff5f56; text-shadow: 0 0 5px #ff5f56; }
.mem-warning { color: #ffbd2e; }
.mem-critical { color: #ff5f56; }

.status-dot {
  display: inline-block;
  width: 5px;
  height: 5px;
  background: #00ff00;
  border-radius: 50%;
  margin-right: 4px;
  box-shadow: 0 0 5px #00ff00;
}

.loading-state, .empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #444;
  font-size: 12px;
  letter-spacing: 1px;
}

.cursor-blink {
  animation: blink 1s infinite;
  color: #00ff00;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

.pm-footer {
  padding: 10px 20px;
  font-size: 11px;
  color: #444;
  letter-spacing: 1px;
  border-top: 1px solid #1a1a1a;
  margin-top: 15px;
}

.context-menu {
  position: fixed;
  background: #0d0d0d;
  border: 1px solid #00ff00;
  border-radius: 4px;
  min-width: 160px;
  z-index: 1000;
  box-shadow: 0 0 20px rgba(0, 255, 0, 0.2);
}

.ctx-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  font-size: 12px;
  color: #00ff00;
  cursor: pointer;
  letter-spacing: 1px;
}

.ctx-item:hover {
  background: rgba(0, 255, 0, 0.1);
}

.ctx-icon {
  width: 14px;
  text-align: center;
}

.ctx-divider {
  height: 1px;
  background: #1a1a1a;
}

.modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal-box {
  background: #0d0d0d;
  border: 1px solid #ff5f56;
  border-radius: 4px;
  min-width: 350px;
  box-shadow: 0 0 30px rgba(255, 95, 86, 0.2);
}

.modal-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 15px;
  border-bottom: 1px solid #1a1a1a;
  font-size: 13px;
  letter-spacing: 1px;
}

.modal-icon.warn {
  color: #ff5f56;
}

.modal-body {
  padding: 15px;
  font-size: 12px;
  color: #888;
}

.modal-body p {
  margin: 0 0 10px 0;
}

.pid-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.pid-tag {
  background: #1a1a1a;
  padding: 3px 8px;
  border-radius: 3px;
  font-size: 11px;
  color: #00ff00;
}

.pid-tag.more {
  color: #666;
}

.modal-foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 15px;
  border-top: 1px solid #1a1a1a;
}

.drag-selection {
  position: fixed;
  background: rgba(0, 255, 0, 0.2);
  border: 1px solid #00ff00;
  pointer-events: none;
  z-index: 999;
}

::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #0a0a0a;
}

::-webkit-scrollbar-thumb {
  background: #1a1a1a;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #333;
}
</style>
