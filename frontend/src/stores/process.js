import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { GetSystemProcessInfos, SearchPidByKeyWord, KillProcessByPID } from '../../wailsjs/go/main/App'

export const useProcessStore = defineStore('process', () => {
  // 状态
  const processes = ref([])
  const loading = ref(false)
  const error = ref(null)
  const searchKeyword = ref('')
  const selectedProcesses = ref([])
  const sortBy = ref('cpu')
  const sortOrder = ref('desc')

  // 计算属性
  const filteredProcesses = computed(() => {
    let result = processes.value

    // 搜索过滤
    if (searchKeyword.value) {
      const keyword = searchKeyword.value.toLowerCase()
      result = result.filter(p => 
        p.Name.toLowerCase().includes(keyword) ||
        p.Cmdline.toLowerCase().includes(keyword) ||
        p.PID.toString().includes(keyword) ||
        (p.ListenAddr && p.ListenAddr.includes(keyword)) ||
        (p.ListenPort && p.ListenPort.toString().includes(keyword))
      )
    }

    // 排序
    result = [...result].sort((a, b) => {
      let aVal = a[sortBy.value]
      let bVal = b[sortBy.value]

      if (typeof aVal === 'string') {
        aVal = aVal.toLowerCase()
        bVal = bVal.toLowerCase()
      }

      if (sortOrder.value === 'asc') {
        return aVal > bVal ? 1 : -1
      } else {
        return aVal < bVal ? 1 : -1
      }
    })

    return result
  })

  const selectedCount = computed(() => selectedProcesses.value.length)
  const totalCount = computed(() => processes.value.length)

  // 操作
  async function fetchProcesses() {
    loading.value = true
    error.value = null

    try {
      const result = await GetSystemProcessInfos()
      processes.value = result || []
    } catch (err) {
      error.value = err.message || '获取进程信息失败'
      console.error('Failed to fetch processes:', err)
    } finally {
      loading.value = false
    }
  }

  async function searchProcess(keyword) {
    if (!keyword) {
      return null
    }

    try {
      const result = await SearchPidByKeyWord(keyword)
      return result
    } catch (err) {
      console.error('Failed to search process:', err)
      return null
    }
  }

  async function killProcess(pid) {
    try {
      await KillProcessByPID(pid)
      // 从列表中移除已终止的进程
      processes.value = processes.value.filter(p => p.PID !== pid)
      // 从选中列表中移除
      selectedProcesses.value = selectedProcesses.value.filter(id => id !== pid)
      return true
    } catch (err) {
      error.value = err.message || '终止进程失败'
      console.error('Failed to kill process:', err)
      return false
    }
  }

  async function killSelectedProcesses() {
    const promises = selectedProcesses.value.map(pid => killProcess(pid))
    const results = await Promise.allSettled(promises)
    
    const failed = results.filter(r => r.status === 'rejected')
    if (failed.length > 0) {
      error.value = `终止 ${failed.length} 个进程失败`
    }

    // 清空选中列表
    selectedProcesses.value = []
  }

  function selectProcess(pid, selected) {
    if (selected) {
      if (!selectedProcesses.value.includes(pid)) {
        selectedProcesses.value.push(pid)
      }
    } else {
      const index = selectedProcesses.value.indexOf(pid)
      if (index > -1) {
        selectedProcesses.value.splice(index, 1)
      }
    }
  }

  function selectAllProcesses(selected) {
    if (selected) {
      selectedProcesses.value = filteredProcesses.value.map(p => p.PID)
    } else {
      selectedProcesses.value = []
    }
  }

  function clearSelection() {
    selectedProcesses.value = []
  }

  function setSortBy(field) {
    if (sortBy.value === field) {
      sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortBy.value = field
      sortOrder.value = 'desc'
    }
  }

  function setSearchKeyword(keyword) {
    searchKeyword.value = keyword
  }

  function clearError() {
    error.value = null
  }

  return {
    // 状态
    processes,
    loading,
    error,
    searchKeyword,
    selectedProcesses,
    sortBy,
    sortOrder,
    // 计算属性
    filteredProcesses,
    selectedCount,
    totalCount,
    // 操作
    fetchProcesses,
    searchProcess,
    killProcess,
    killSelectedProcesses,
    selectProcess,
    selectAllProcesses,
    clearSelection,
    setSortBy,
    setSearchKeyword,
    clearError
  }
})