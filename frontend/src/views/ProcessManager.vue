<template>
    <div class="process-manager">
        <div class="header">
            <h1>{{ t("process.title") }}</h1>
            <p class="subtitle">{{ t("process.subtitle") }}</p>
        </div>

        <div class="toolbar">
            <div class="search-box">
                <input
                    v-model="searchKeyword"
                    type="text"
                    :placeholder="t('process.search')"
                    class="search-input"
                />
                <button
                    @click="refreshProcesses"
                    class="refresh-btn"
                    :disabled="loading"
                >
                    {{
                        loading ? t("process.refreshing") : t("process.refresh")
                    }}
                </button>
            </div>

            <div class="auto-refresh-controls">
                <label class="switch">
                    <input
                        type="checkbox"
                        v-model="autoRefresh"
                        @change="toggleAutoRefresh"
                    />
                    <span class="slider"></span>
                </label>
                <span>{{ t("process.autoRefresh") }}</span>
                <input
                    v-model.number="refreshInterval"
                    type="number"
                    min="1"
                    class="interval-input"
                    :disabled="!autoRefresh"
                    @change="updateRefreshInterval"
                />
                <span>秒</span>
            </div>

            <div class="stats">
                <span class="stat-item"
                    >{{ t("process.total") }}: {{ totalCount }}</span
                >
                <span class="stat-item"
                    >{{ t("process.selected") }}: {{ selectedCount }}</span
                >
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
                        <th class="col-pid" @click="setSortBy('PID')">
                            {{ t("process.pid") }}
                        </th>
                        <th class="col-name" @click="setSortBy('Name')">
                            {{ t("process.name") }}
                        </th>
                        <th class="col-cpu" @click="setSortBy('CPUPercent')">
                            {{ t("process.cpu") }}
                        </th>
                        <th class="col-memory" @click="setSortBy('MemoryMB')">
                            {{ t("process.memory") }}
                        </th>
                        <th class="col-port" @click="setSortBy('ListenPort')">
                            {{ t("process.port") }}
                        </th>
                        <th class="col-status" @click="setSortBy('Status')">
                            {{ t("process.status") }}
                        </th>
                        <th class="col-actions">{{ t("process.action") }}</th>
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
                                @change="
                                    toggleSelect(
                                        process.PID,
                                        $event.target.checked,
                                    )
                                "
                            />
                        </td>
                        <td class="col-pid">{{ process.PID }}</td>
                        <td class="col-name" :title="process.Cmdline">
                            {{ process.Name || "(无名称)" }}
                        </td>
                        <td class="col-cpu">
                            <div class="progress-bar">
                                <div
                                    class="progress-fill"
                                    :style="{ width: process.CPUPercent + '%' }"
                                ></div>
                                <span class="progress-text"
                                    >{{ process.CPUPercent.toFixed(2) }}%</span
                                >
                            </div>
                        </td>
                        <td class="col-memory">
                            <div class="progress-bar">
                                <div
                                    class="progress-fill memory"
                                    :style="{
                                        width:
                                            Math.min(
                                                process.MemoryMB / 100,
                                                100,
                                            ) + '%',
                                    }"
                                ></div>
                                <span class="progress-text"
                                    >{{ process.MemoryMB }} MB</span
                                >
                            </div>
                        </td>
                        <td class="col-port">
                            {{ process.ListenPort || "-" }}
                        </td>
                        <td class="col-status">
                            <span
                                class="status-badge"
                                :class="getStatusClass(process.Status)"
                            >
                                {{ process.Status || "Unknown" }}
                            </span>
                        </td>
                        <td class="col-actions">
                            <button
                                @click="killProcess(process.PID)"
                                class="kill-btn"
                                :disabled="
                                    killingProcesses.includes(process.PID)
                                "
                            >
                                {{
                                    killingProcesses.includes(process.PID)
                                        ? t("process.killing")
                                        : t("process.kill")
                                }}
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>

            <div v-if="filteredProcesses.length === 0" class="no-data">
                {{ t("process.noData") }}
            </div>
        </div>

        <div v-else class="loading">
            <div class="loading-spinner"></div>
            <p>{{ t("process.loading") }}</p>
        </div>

        <div v-if="error" class="error-message">
            {{ error }}
        </div>

        <div class="batch-actions" v-if="selectedCount > 0">
            <button @click="killSelectedProcesses" class="batch-kill-btn">
                {{ t("process.killSelected", { count: selectedCount }) }}
            </button>
            <button @click="clearSelection" class="clear-selection-btn">
                {{ t("process.clearSelection") }}
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { storeToRefs } from "pinia";
import { useProcessStore } from "../stores/process";
import { t } from "../i18n";

const processStore = useProcessStore();

// 响应式数据
const killingProcesses = ref([]);
const autoRefresh = ref(false);
const refreshInterval = ref(5);
let refreshTimer = null;

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
    searchKeyword,
} = storeToRefs(processStore);

// 从store获取方法
const {
    fetchProcesses,
    killProcess: storeKillProcess,
    selectProcess,
    selectAllProcesses,
    clearSelection: storeClearSelection,
    setSortBy: storeSetSortBy,
} = processStore;

// 计算属性
const allSelected = computed(() => {
    if (!filteredProcesses.value || !selectedProcesses.value) return false;
    return (
        filteredProcesses.value.length > 0 &&
        selectedProcesses.value.length === filteredProcesses.value.length
    );
});

// 方法
function refreshProcesses() {
    fetchProcesses();
}

function toggleSelectAll(event) {
    selectAllProcesses(event.target.checked);
}

function toggleSelect(pid, selected) {
    selectProcess(pid, selected);
}

function isSelected(pid) {
    if (!selectedProcesses.value) return false;
    return selectedProcesses.value.includes(pid);
}

async function killProcess(pid) {
    killingProcesses.value.push(pid);
    await storeKillProcess(pid);
    killingProcesses.value = killingProcesses.value.filter((id) => id !== pid);
}

async function killSelectedProcesses() {
    const pids = [...selectedProcesses.value];
    for (const pid of pids) {
        await killProcess(pid);
    }
}

function clearSelection() {
    storeClearSelection();
}

function setSortBy(field) {
    storeSetSortBy(field);
}

function getStatusClass(status) {
    const statusMap = {
        LISTEN: "status-listen",
        ESTABLISHED: "status-established",
        TIME_WAIT: "status-time-wait",
        CLOSE_WAIT: "status-close-wait",
    };
    return statusMap[status] || "status-unknown";
}

function toggleAutoRefresh() {
    if (autoRefresh.value) {
        refreshTimer = setInterval(() => {
            fetchProcesses();
        }, refreshInterval.value * 1000);
    } else {
        if (refreshTimer) {
            clearInterval(refreshTimer);
            refreshTimer = null;
        }
    }
}

function updateRefreshInterval() {
    if (autoRefresh.value && refreshTimer) {
        clearInterval(refreshTimer);
        refreshTimer = setInterval(() => {
            fetchProcesses();
        }, refreshInterval.value * 1000);
    }
}

// 组件挂载时获取进程列表
onMounted(() => {
    fetchProcesses();
});

// 组件卸载时清除定时器
onUnmounted(() => {
    if (refreshTimer) {
        clearInterval(refreshTimer);
        refreshTimer = null;
    }
});
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
    color: var(--accent-color);
    font-size: 2rem;
    margin-bottom: 0.5rem;
}

.subtitle {
    color: var(--text-secondary);
}

.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    gap: 1rem;
    flex-wrap: wrap;
}

.search-box {
    display: flex;
    gap: 0.5rem;
    flex: 1;
}

.search-input {
    flex: 1;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    padding: 0.5rem 1rem;
    color: var(--text-tertiary);
    font-family: inherit;
}

.search-input:focus {
    outline: none;
    border-color: var(--accent-color);
}

.refresh-btn {
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    padding: 0.5rem 1rem;
    color: var(--accent-color);
    cursor: pointer;
    transition: all 0.3s;
}

.refresh-btn:hover:not(:disabled) {
    border-color: var(--accent-color);
    background: var(--accent-subtle);
}

.refresh-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.auto-refresh-controls {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    color: var(--text-secondary);
    font-size: 0.9rem;
    line-height: 1;
}

.auto-refresh-controls > .switch {
    flex-shrink: 0;
}

.auto-refresh-controls > span {
    line-height: 1;
}

.interval-input {
    width: 60px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    padding: 0.25rem 0.5rem;
    color: var(--text-tertiary);
    text-align: center;
}

.interval-input:disabled {
    opacity: 0.5;
}

.interval-input:focus {
    border-color: var(--accent-color);
    outline: none;
}

.switch {
    position: relative;
    display: inline-block;
    width: 40px;
    height: 22px;
}

.switch input {
    opacity: 0;
    width: 0;
    height: 0;
}

.switch .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--border-color);
    transition: 0.3s;
    border-radius: 22px;
}

.switch .slider:before {
    position: absolute;
    content: "";
    height: 16px;
    width: 16px;
    left: 3px;
    bottom: 3px;
    background-color: var(--text-tertiary);
    transition: 0.3s;
    border-radius: 50%;
}

.switch input:checked + .slider {
    background-color: var(--accent-color);
}

.switch input:checked + .slider:before {
    transform: translateX(18px);
}

.stats {
    display: flex;
    gap: 1.5rem;
}

.stat-item {
    color: var(--text-secondary);
    font-size: 0.9rem;
}

.process-list {
    flex: 1;
    overflow: auto;
    background: var(--bg-secondary);
    border-radius: 8px;
    border: 1px solid var(--border-color);
}

.process-table {
    width: 100%;
    border-collapse: collapse;
}

.process-table th {
    position: sticky;
    top: 0;
    background: var(--bg-secondary);
    padding: 1rem;
    text-align: left;
    color: var(--text-secondary);
    font-weight: normal;
    font-size: 0.85rem;
    border-bottom: 1px solid var(--border-color);
    cursor: pointer;
}

.process-table th:hover {
    color: var(--accent-color);
}

.process-table td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--border-color);
    color: var(--text-tertiary);
}

.process-row:hover {
    background: var(--accent-subtle);
}

.process-row.selected {
    background: var(--accent-subtle);
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
    background: var(--bg-secondary);
    border-radius: 4px;
    overflow: hidden;
}

.progress-fill {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    background: linear-gradient(
        90deg,
        var(--accent-color),
        var(--accent-hover)
    );
    transition: width 0.3s;
}

.progress-fill.memory {
    background: linear-gradient(90deg var(--info-color), #0099cc);
}

.progress-text {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 0.75rem;
    color: var(--text-tertiary);
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.status-badge {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
}

.status-listen {
    background: var(--accent-subtle);
    color: var(--accent-color);
}

.status-established {
    background: var(--info-subtle);
    color: var(--info-color);
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
    color: var(--text-secondary);
}

.kill-btn {
    background: transparent;
    border: 1px solid var(--danger-color);
    border-radius: 4px;
    padding: 0.25rem 0.5rem;
    color: var(--danger-color);
    cursor: pointer;
    font-size: 0.8rem;
    transition: all 0.3s;
}

.kill-btn:hover:not(:disabled) {
    background: var(--danger-subtle);
}

.kill-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
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

.error-message {
    background: var(--danger-subtle);
    border: 1px solid var(--danger-color);
    border-radius: 6px;
    padding: 1rem;
    margin-top: 1rem;
    color: var(--danger-color);
}

.batch-actions {
    position: fixed;
    bottom: 2rem;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 1rem;
    background: var(--bg-secondary);
    padding: 1rem;
    border-radius: 8px;
    border: 1px solid var(--border-color);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
}

.batch-kill-btn {
    background: var(--danger-color);
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1rem;
    color: var(--text-tertiary);
    cursor: pointer;
    transition: all 0.3s;
}

.batch-kill-btn:hover {
    background: #ff6666;
}

.clear-selection-btn {
    background: transparent;
    border: 1px solid var(--border-color);
    border-radius: 6px;
    padding: 0.5rem 1rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.3s;
}

.clear-selection-btn:hover {
    border-color: var(--accent-color);
    color: var(--accent-color);
}

input[type="checkbox"] {
    appearance: none;
    width: 16px;
    height: 16px;
    border: 1px solid var(--border-color);
    border-radius: 3px;
    background: var(--bg-tertiary);
    cursor: pointer;
}

input[type="checkbox"]:checked {
    background: var(--accent-color);
    border-color: var(--accent-color);
}
</style>
