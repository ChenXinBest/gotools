<template>
  <div class="settings">
    <div class="header">
      <h1>设置</h1>
      <p class="subtitle">应用配置和偏好设置</p>
    </div>

    <div class="settings-content">
      <!-- 常规设置 -->
      <div class="settings-section">
        <h2>常规设置</h2>
        <div class="setting-item">
          <div class="setting-label">
            <span class="setting-name">主题</span>
            <span class="setting-description">选择应用的主题颜色</span>
          </div>
          <select v-model="settings.theme" class="setting-control">
            <option value="dark">深色</option>
            <option value="light">浅色</option>
            <option value="auto">跟随系统</option>
          </select>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="setting-name">语言</span>
            <span class="setting-description">选择界面显示语言</span>
          </div>
          <select v-model="settings.locale" class="setting-control">
            <option value="zh-CN">简体中文</option>
            <option value="en-US">English</option>
          </select>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="setting-name">自动刷新进程</span>
            <span class="setting-description">自动刷新进程列表</span>
          </div>
          <label class="switch">
            <input type="checkbox" v-model="settings.autoRefreshProcesses" />
            <span class="slider"></span>
          </label>
        </div>

        <div class="setting-item" v-if="settings.autoRefreshProcesses">
          <div class="setting-label">
            <span class="setting-name">刷新间隔</span>
            <span class="setting-description">进程列表刷新间隔（秒）</span>
          </div>
          <input
            v-model.number="settings.refreshInterval"
            type="number"
            min="1"
            max="60"
            class="setting-control small"
          />
        </div>
      </div>

      <!-- 数据库设置 -->
      <div class="settings-section">
        <h2>数据库设置</h2>
        <div class="setting-item">
          <div class="setting-label">
            <span class="setting-name">默认导出工具</span>
            <span class="setting-description">选择默认使用的导出工具</span>
          </div>
          <select v-model="settings.defaultExportTool" class="setting-control">
            <option value="mysql-shell">MySQL Shell</option>
            <option value="mysqldump">mysqldump</option>
          </select>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="setting-name">默认导出路径</span>
            <span class="setting-description">设置默认的导出文件保存路径</span>
          </div>
          <div class="path-input">
            <input
              v-model="settings.defaultExportPath"
              type="text"
              placeholder="选择默认导出路径"
              class="setting-control"
              readonly
            />
            <button @click="selectExportPath" class="browse-btn">浏览</button>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="setting-name">默认线程数</span>
            <span class="setting-description">导出/导入操作的默认线程数</span>
          </div>
          <input
            v-model.number="settings.defaultThreads"
            type="number"
            min="1"
            max="16"
            class="setting-control small"
          />
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="setting-name">默认压缩格式</span>
            <span class="setting-description">导出文件的默认压缩格式</span>
          </div>
          <select v-model="settings.defaultCompression" class="setting-control">
            <option value="gzip">gzip</option>
            <option value="zstd">zstd</option>
            <option value="none">无压缩</option>
          </select>
        </div>
      </div>

      <!-- 关于 -->
      <div class="settings-section">
        <h2>关于</h2>
        <div class="about-info">
          <div class="info-item">
            <span class="info-label">应用名称</span>
            <span class="info-value">GoTools</span>
          </div>
          <div class="info-item">
            <span class="info-label">版本</span>
            <span class="info-value">{{ versionInfo.version }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">构建时间</span>
            <span class="info-value">{{ versionInfo.build_time }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Git 提交</span>
            <span class="info-value">{{ versionInfo.git_commit }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Go 版本</span>
            <span class="info-value">{{ versionInfo.go_version }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">平台</span>
            <span class="info-value">{{ versionInfo.platform }}</span>
          </div>
        </div>
      </div>

      <!-- 保存按钮 -->
      <div class="settings-actions">
        <button @click="saveSettings" class="save-btn" :disabled="saving">
          {{ saving ? '保存中...' : '保存设置' }}
        </button>
        <button @click="resetSettings" class="reset-btn">
          重置为默认
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAppStore } from '../stores/app'

const appStore = useAppStore()

const saving = ref(false)

const settings = ref({
  theme: 'dark',
  locale: 'zh-CN',
  autoRefreshProcesses: false,
  refreshInterval: 5,
  defaultExportTool: 'mysql-shell',
  defaultExportPath: '',
  defaultThreads: 4,
  defaultCompression: 'gzip'
})

const versionInfo = ref({
  version: '1.0.0',
  build_time: '未知',
  git_commit: '未知',
  go_version: '未知',
  platform: '未知'
})

async function selectExportPath() {
  try {
    const path = await window.go.main.App.SelectFolder()
    if (path) {
      settings.value.defaultExportPath = path
    }
  } catch (err) {
    console.error('Failed to select folder:', err)
  }
}

async function saveSettings() {
  saving.value = true
  try {
    // 保存到localStorage
    localStorage.setItem('gotools-settings', JSON.stringify(settings.value))
    
    // 应用设置到store
    appStore.setTheme(settings.value.theme)
    appStore.setLocale(settings.value.locale)
    
    alert('设置已保存')
  } catch (err) {
    alert('保存失败: ' + err.message)
  } finally {
    saving.value = false
  }
}

function resetSettings() {
  if (confirm('确定要重置所有设置为默认值吗？')) {
    settings.value = {
      theme: 'dark',
      locale: 'zh-CN',
      autoRefreshProcesses: false,
      refreshInterval: 5,
      defaultExportTool: 'mysql-shell',
      defaultExportPath: '',
      defaultThreads: 4,
      defaultCompression: 'gzip'
    }
    localStorage.removeItem('gotools-settings')
  }
}

async function loadVersionInfo() {
  try {
    const info = await window.go.main.App.GetVersion()
    versionInfo.value = info
  } catch (err) {
    console.error('Failed to load version info:', err)
  }
}

onMounted(() => {
  // 加载保存的设置
  const savedSettings = localStorage.getItem('gotools-settings')
  if (savedSettings) {
    try {
      const parsed = JSON.parse(savedSettings)
      settings.value = { ...settings.value, ...parsed }
    } catch (err) {
      console.error('Failed to load settings:', err)
    }
  }
  
  loadVersionInfo()
})
</script>

<style scoped>
.settings {
  padding: 2rem;
  max-width: 800px;
}

.header {
  margin-bottom: 2rem;
}

.header h1 {
  color: #00ff00;
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #888;
}

.settings-section {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.settings-section h2 {
  color: #fff;
  font-size: 1.2rem;
  margin-bottom: 1.5rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #333;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  border-bottom: 1px solid #222;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-label {
  flex: 1;
}

.setting-name {
  display: block;
  color: #fff;
  margin-bottom: 0.25rem;
}

.setting-description {
  display: block;
  color: #888;
  font-size: 0.85rem;
}

.setting-control {
  background: #0d0d0d;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  color: #fff;
  min-width: 150px;
}

.setting-control:focus {
  outline: none;
  border-color: #00ff00;
}

.setting-control.small {
  width: 80px;
  min-width: 80px;
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
}

.browse-btn:hover {
  border-color: #00ff00;
}

.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 26px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #333;
  transition: 0.4s;
  border-radius: 26px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 3px;
  bottom: 3px;
  background-color: #fff;
  transition: 0.4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #00ff00;
}

input:checked + .slider:before {
  transform: translateX(24px);
}

.about-info {
  padding: 1rem 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid #222;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: #888;
}

.info-value {
  color: #fff;
  font-family: monospace;
}

.settings-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

.save-btn {
  background: #00ff00;
  border: none;
  border-radius: 6px;
  padding: 0.75rem 2rem;
  color: #000;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
}

.save-btn:hover:not(:disabled) {
  background: #00cc00;
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.reset-btn {
  background: transparent;
  border: 1px solid #666;
  border-radius: 6px;
  padding: 0.75rem 2rem;
  color: #888;
  cursor: pointer;
  transition: all 0.3s;
}

.reset-btn:hover {
  border-color: #ff4444;
  color: #ff4444;
}
</style>