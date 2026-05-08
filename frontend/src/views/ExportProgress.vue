<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Events } from '@wailsio/runtime'

const status = ref('准备中...')
const progress = ref(0)
const isComplete = ref(false)
const hasError = ref(false)
const errorMessage = ref('')

let cleanup = null

onMounted(() => {
  cleanup = Events.On('export-progress', (data) => {
    if (typeof data === 'object') {
      if (data.status) status.value = data.status
      if (data.progress !== undefined) progress.value = data.progress
      if (data.complete) isComplete.value = true
      if (data.error) {
        hasError.value = true
        errorMessage.value = data.error
      }
    }
  })
})

onUnmounted(() => {
  if (cleanup) cleanup()
})
</script>

<template>
  <div class="progress-window">
    <div class="container">
      <h3>导出进度</h3>

      <div class="status-bar">
        <div class="progress-track">
          <div
            class="progress-fill"
            :class="{ complete: isComplete, error: hasError }"
            :style="{ width: progress + '%' }"
          />
        </div>
        <span class="progress-text">{{ progress }}%</span>
      </div>

      <p class="status-text" :class="{ error: hasError }">
        {{ hasError ? errorMessage : status }}
      </p>

      <div v-if="isComplete" class="actions">
        <p class="success-text">导出完成</p>
        <button class="btn-close" @click="window.close()">关闭</button>
      </div>

      <div v-if="hasError" class="actions">
        <button class="btn-close" @click="window.close()">关闭</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.progress-window {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary, #1a1a2e);
  color: var(--text-primary, #e0e0e0);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.container {
  width: 90%;
  max-width: 500px;
  text-align: center;
}

h3 {
  margin-bottom: 24px;
  font-size: 18px;
}

.status-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.progress-track {
  flex: 1;
  height: 20px;
  background: var(--border-color, #2a2a4a);
  border-radius: 10px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--accent-color, #4a9eff);
  border-radius: 10px;
  transition: width 0.3s ease;
}

.progress-fill.complete {
  background: #4caf50;
}

.progress-fill.error {
  background: #f44336;
}

.progress-text {
  font-size: 14px;
  min-width: 40px;
}

.status-text {
  margin-bottom: 16px;
  font-size: 14px;
  color: var(--text-secondary, #a0a0a0);
}

.status-text.error {
  color: #f44336;
}

.success-text {
  color: #4caf50;
  margin-bottom: 12px;
  font-size: 16px;
  font-weight: bold;
}

.actions {
  margin-top: 8px;
}

.btn-close {
  padding: 8px 24px;
  background: var(--accent-color, #4a9eff);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.btn-close:hover {
  opacity: 0.9;
}
</style>
