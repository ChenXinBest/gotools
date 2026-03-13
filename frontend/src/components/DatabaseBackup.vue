<script setup>
import { ref, onMounted, computed } from "vue";
import { AddDatabaseConnection, UpdateDatabaseConnection, GetDatabaseConnections, ListDatabases, ListTables } from "../../wailsjs/go/main/App";
import { FaSave } from "vue-icons-plus/fa";

// 连接配置
const connections = ref([]);
const currentConnection = ref("");
const newConnection = ref({
    id: "",
    name: "",
    host: "localhost",
    port: 3306,
    user: "root",
    password: "",
    defaultSchema: "",
});
const editingConnection = ref(null);

// 数据库和表
const databases = ref([]);
const selectedDatabases = ref(new Set());
const tables = ref([]);
const selectedTables = ref(new Set());

// 导出配置
const exportPath = ref("");
const exportTool = ref("mysql-shell");
const loading = ref(false);
const error = ref("");
const statusMessage = ref("");

// 模态框
const showConnectionModal = ref(false);
const showExportModal = ref(false);

// 计算属性
const hasSelectedDatabases = computed(() => selectedDatabases.value.size > 0);
const hasSelectedTables = computed(() => selectedTables.value.size > 0);
const canExport = computed(() => {
    return (
        exportPath.value &&
        (hasSelectedDatabases.value || hasSelectedTables.value)
    );
});

// 加载连接配置
async function loadConnections() {
    try {
        const connList = await GetDatabaseConnections();
        // 转换后端返回的数据结构以匹配前端使用的结构
        connections.value = connList.map(conn => ({
            id: conn.id,
            name: conn.name,
            host: conn.host,
            port: conn.port,
            user: conn.user,
            password: conn.password,
            defaultSchema: conn.database
        }));
    } catch (err) {
        console.error("加载连接配置失败:", err);
        // 加载失败时使用模拟数据
        connections.value = [
            {
                id: "1",
                name: "Local MySQL",
                host: "localhost",
                port: 3306,
                user: "root",
                password: "password",
                defaultSchema: "test",
            },
            {
                id: "2",
                name: "Production MySQL",
                host: "192.168.1.100",
                port: 3306,
                user: "admin",
                password: "securepassword",
                defaultSchema: "production",
            },
        ];
    }
}

// 编辑连接
function editConnection() {
    if (!currentConnection.value) return;

    const conn = connections.value.find(
        (c) => c.name === currentConnection.value,
    );
    if (conn) {
        editingConnection.value = conn.name;
        newConnection.value = { ...conn };
        showConnectionModal.value = true;
    }
}

// 刷新连接
async function refreshConnection() {
    if (!currentConnection.value) return;

    // 清空当前选择的数据库和表
    selectedDatabases.value.clear();
    selectedTables.value.clear();
    tables.value = [];

    try {
        // 获取当前连接的详细信息
        const conn = connections.value.find(c => c.name === currentConnection.value);
        if (!conn) {
            error.value = "找不到连接信息";
            return;
        }

        if (exportTool.value === "mysql-shell") {
            // 调用后端 API 获取数据库列表
            const connData = {
                id: conn.id,
                name: conn.name,
                host: conn.host,
                port: conn.port,
                user: conn.user,
                password: conn.password,
                database: conn.defaultSchema
            };

            databases.value = await ListDatabases(connData);
            showStatus(`已刷新 ${currentConnection.value} 的数据库信息`);
        } else {
            // 使用模拟数据
            databases.value = [
                "information_schema",
                "mysql",
                "performance_schema",
                "sys",
                "test",
            ];
            showStatus(`已刷新 ${currentConnection.value} 的数据库信息`);
        }
    } catch (err) {
        console.error("刷新数据库信息失败:", err);
        error.value = "刷新失败: " + (err.message || err);
        // 失败时使用模拟数据
        databases.value = [
            "information_schema",
            "mysql",
            "performance_schema",
            "sys",
            "test",
        ];
    }
}

// 保存连接（支持添加和编辑）
async function saveConnection() {
    if (!newConnection.value.name) {
        error.value = "请输入连接名称";
        return;
    }

    // 转换为后端期望的数据结构
    const connData = {
        id: newConnection.value.id || "",
        name: newConnection.value.name,
        host: newConnection.value.host,
        port: newConnection.value.port,
        user: newConnection.value.user,
        password: newConnection.value.password,
        database: newConnection.value.defaultSchema
    };

    try {
        if (editingConnection.value) {
            // 编辑现有连接
            await UpdateDatabaseConnection(connData);
            const index = connections.value.findIndex(
                (c) => c.name === editingConnection.value,
            );
            if (index !== -1) {
                connections.value[index] = { ...newConnection.value };
            }
            showStatus("连接更新成功");
            editingConnection.value = null;
        } else {
            // 添加新连接
            await AddDatabaseConnection(connData);
            // 重新加载连接列表以获取新的ID
            await loadConnections();
            showStatus("连接保存成功");
        }

        showConnectionModal.value = false;
        newConnection.value = {
            id: "",
            name: "",
            host: "localhost",
            port: 3306,
            user: "root",
            password: "",
            defaultSchema: "",
        };
    } catch (err) {
        console.error("保存连接配置失败:", err);
        error.value = "保存连接失败: " + (err.message || err);
    }
}

// 选择连接
function selectConnection(connName) {
    currentConnection.value = connName;
    // 这里应该根据选择的连接获取数据库列表
    // 暂时使用模拟数据
    databases.value = [
        "information_schema",
        "mysql",
        "performance_schema",
        "sys",
        "test",
    ];
    selectedDatabases.value.clear();
    selectedTables.value.clear();
    tables.value = [];
    showStatus(`已连接到 ${connName}`);
}

// 选择数据库
async function toggleDatabase(db) {
    if (selectedDatabases.value.has(db)) {
        selectedDatabases.value.delete(db);
    } else {
        selectedDatabases.value.add(db);
    }
    selectedDatabases.value = new Set(selectedDatabases.value);

    // 当选择数据库时，获取该数据库的表
    if (selectedDatabases.value.has(db)) {
        try {
            if (exportTool.value === "mysql-shell") {
                // 获取当前连接的详细信息
                const conn = connections.value.find(c => c.name === currentConnection.value);
                if (conn) {
                    // 调用后端 API 获取表列表
                    const connData = {
                        id: conn.id,
                        name: conn.name,
                        host: conn.host,
                        port: conn.port,
                        user: conn.user,
                        password: conn.password,
                        database: db
                    };

                    tables.value = await ListTables(connData);
                }
            } else {
                // 使用模拟数据
                tables.value = [
                    "users",
                    "orders",
                    "products",
                    "categories",
                    "transactions",
                ];
            }
        } catch (err) {
            console.error("获取表列表失败:", err);
            error.value = "获取表列表失败: " + (err.message || err);
            // 失败时使用模拟数据
            tables.value = [
                "users",
                "orders",
                "products",
                "categories",
                "transactions",
            ];
        }
    }
}

// 选择表
function toggleTable(table) {
    if (selectedTables.value.has(table)) {
        selectedTables.value.delete(table);
    } else {
        selectedTables.value.add(table);
    }
    selectedTables.value = new Set(selectedTables.value);
}

// 导出数据
async function exportData() {
    if (!canExport.value) {
        error.value = "请选择导出路径和要导出的数据库/表";
        return;
    }

    loading.value = true;
    error.value = "";

    try {
        // 这里应该调用mysql-shell进行导出
        // 暂时模拟导出过程
        await new Promise((resolve) => setTimeout(resolve, 2000));

        let exportType = "";
        if (hasSelectedTables.value) {
            exportType = `表(${selectedTables.value.size}个)`;
        } else if (hasSelectedDatabases.value) {
            exportType = `数据库(${selectedDatabases.value.size}个)`;
        } else {
            exportType = "全部数据库";
        }

        showStatus(`成功导出 ${exportType} 到 ${exportPath.value}`);
        showExportModal.value = false;
    } catch (e) {
        error.value = "导出失败: " + (e.message || e);
    } finally {
        loading.value = false;
    }
}

// 显示状态信息
function showStatus(msg) {
    statusMessage.value = msg;
    setTimeout(() => {
        statusMessage.value = "";
    }, 3000);
}

// 打开文件选择器
function selectExportPath() {
    // 这里应该打开文件选择器
    // 暂时模拟
    exportPath.value = "D:\\backup";
}

onMounted(async () => {
    await loadConnections();
});
</script>

<template>
    <div class="database-backup">
        <div class="db-header">
            <div class="db-title">
                <span class="title-icon">
                    <component :is="FaSave" />
                </span>
                导出
            </div>
            <div class="db-stats">
                <span class="stat-item">{{ connections.length }} 连接</span>
                <span class="stat-divider">|</span>
                <span class="stat-item"
                    >{{ selectedDatabases.size }} 已选数据库</span
                >
                <span class="stat-divider">|</span>
                <span class="stat-item">{{ selectedTables.size }} 已选表</span>
            </div>
        </div>

        <div class="db-toolbar">
            <div class="connection-selector">
                <select v-model="currentConnection" class="select-input">
                    <option value="">选择连接</option>
                    <option
                        v-for="conn in connections"
                        :key="conn.name"
                        :value="conn.name"
                    >
                        {{ conn.name }}
                    </option>
                </select>
                <button @click="showConnectionModal = true" class="btn">
                    添加连接
                </button>
                <button
                    @click="editConnection()"
                    class="btn"
                    :disabled="!currentConnection"
                >
                    编辑连接
                </button>
                <button
                    @click="refreshConnection()"
                    class="btn"
                    :disabled="!currentConnection"
                >
                    刷新
                </button>
            </div>
            <div class="export-settings">
                <div class="setting-item">
                    <label>导出工具:</label>
                    <select v-model="exportTool" class="select-input">
                        <option value="mysql-shell">MySQL Shell</option>
                        <option value="mysqldump">mysqldump</option>
                    </select>
                </div>
                <div class="setting-item">
                    <label>导出路径:</label>
                    <div class="path-selector">
                        <input
                            v-model="exportPath"
                            type="text"
                            class="path-input"
                            readonly
                        />
                        <button @click="selectExportPath" class="btn">
                            选择
                        </button>
                    </div>
                </div>
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

        <div class="db-content">
            <div class="db-section">
                <div class="section-header">
                    <h3>数据库</h3>
                    <span class="section-info"
                        >{{ databases.length }} 个数据库</span
                    >
                </div>
                <div class="db-list">
                    <div v-if="databases.length === 0" class="empty-state">
                        请先选择连接
                    </div>
                    <div
                        v-for="db in databases"
                        :key="db"
                        class="db-item"
                        :class="{ selected: selectedDatabases.has(db) }"
                        @click="toggleDatabase(db)"
                    >
                        <input
                            type="checkbox"
                            :checked="selectedDatabases.has(db)"
                            @change="toggleDatabase(db)"
                        />
                        <span class="db-name">{{ db }}</span>
                    </div>
                </div>
            </div>

            <div class="db-section">
                <div class="section-header">
                    <h3>表</h3>
                    <span class="section-info">{{ tables.length }} 个表</span>
                </div>
                <div class="table-list">
                    <div v-if="tables.length === 0" class="empty-state">
                        请先选择数据库
                    </div>
                    <div
                        v-for="table in tables"
                        :key="table"
                        class="table-item"
                        :class="{ selected: selectedTables.has(table) }"
                        @click="toggleTable(table)"
                    >
                        <input
                            type="checkbox"
                            :checked="selectedTables.has(table)"
                            @change="toggleTable(table)"
                        />
                        <span class="table-name">{{ table }}</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="db-footer">
            <div class="export-info">
                <span v-if="hasSelectedTables">
                    已选择 {{ selectedTables.size }} 个表
                </span>
                <span v-else-if="hasSelectedDatabases">
                    已选择 {{ selectedDatabases.size }} 个数据库
                </span>
                <span v-else> 请选择要导出的数据库或表 </span>
            </div>
            <button
                @click="showExportModal = true"
                class="btn danger"
                :disabled="!canExport || loading"
            >
                {{ loading ? "导出中..." : "开始导出" }}
            </button>
        </div>

        <!-- 连接配置模态框 -->
        <div
            v-if="showConnectionModal"
            class="modal-mask"
            @click.self="showConnectionModal = false"
        >
            <div class="modal-box">
                <div class="modal-head">
                    <span class="modal-icon">⚙</span>
                    {{ editingConnection ? "编辑连接" : "添加连接" }}
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label>连接名称:</label>
                        <input
                            v-model="newConnection.name"
                            type="text"
                            class="form-input"
                        />
                    </div>
                    <div class="form-group">
                        <label>主机:</label>
                        <input
                            v-model="newConnection.host"
                            type="text"
                            class="form-input"
                        />
                    </div>
                    <div class="form-group">
                        <label>端口:</label>
                        <input
                            v-model.number="newConnection.port"
                            type="number"
                            class="form-input"
                        />
                    </div>
                    <div class="form-group">
                        <label>用户名:</label>
                        <input
                            v-model="newConnection.user"
                            type="text"
                            class="form-input"
                        />
                    </div>
                    <div class="form-group">
                        <label>密码:</label>
                        <input
                            v-model="newConnection.password"
                            type="password"
                            class="form-input"
                        />
                    </div>
                    <div class="form-group">
                        <label>默认数据库:</label>
                        <input
                            v-model="newConnection.defaultSchema"
                            type="text"
                            class="form-input"
                        />
                    </div>
                </div>
                <div class="modal-foot">
                    <button @click="showConnectionModal = false" class="btn">
                        取消
                    </button>
                    <button @click="saveConnection" class="btn danger">
                        保存
                    </button>
                </div>
            </div>
        </div>

        <!-- 导出确认模态框 -->
        <div
            v-if="showExportModal"
            class="modal-mask"
            @click.self="showExportModal = false"
        >
            <div class="modal-box">
                <div class="modal-head">
                    <span class="modal-icon warn">⚠</span>
                    确认导出
                </div>
                <div class="modal-body">
                    <p>确定要导出以下内容吗？</p>
                    <div class="export-summary">
                        <div v-if="hasSelectedTables">
                            <h4>已选择表:</h4>
                            <div class="selected-items">
                                <span
                                    v-for="table in selectedTables"
                                    :key="table"
                                    class="item-tag"
                                >
                                    {{ table }}
                                </span>
                            </div>
                        </div>
                        <div v-else-if="hasSelectedDatabases">
                            <h4>已选择数据库:</h4>
                            <div class="selected-items">
                                <span
                                    v-for="db in selectedDatabases"
                                    :key="db"
                                    class="item-tag"
                                >
                                    {{ db }}
                                </span>
                            </div>
                        </div>
                        <div class="export-path">
                            <h4>导出路径:</h4>
                            <span>{{ exportPath }}</span>
                        </div>
                    </div>
                </div>
                <div class="modal-foot">
                    <button @click="showExportModal = false" class="btn">
                        取消
                    </button>
                    <button @click="exportData" class="btn danger">导出</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.database-backup {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: #0a0a0a;
    color: #00ff00;
    font-family: "Consolas", "Monaco", monospace;
    overflow: hidden;
}

.db-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 20px;
    background: #0d0d0d;
    border-bottom: 1px solid #1a1a1a;
}

.db-title {
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
    display: flex;
    align-items: center;
}

.title-icon svg {
    width: 16px;
    height: 16px;
}

.db-stats {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 12px;
    color: #666;
}

.stat-divider {
    color: #333;
}

.db-toolbar {
    display: flex;
    align-items: center;
    gap: 15px;
    padding: 12px 20px;
    background: #0d0d0d;
    border-bottom: 1px solid #1a1a1a;
    flex-wrap: wrap;
}

.connection-selector {
    display: flex;
    align-items: center;
    gap: 8px;
}

.select-input {
    background: #0a0a0a;
    border: 1px solid #1a1a1a;
    color: #00ff00;
    padding: 6px 10px;
    font-family: inherit;
    font-size: 12px;
    border-radius: 3px;
    outline: none;
    min-width: 200px;
}

.select-input option {
    background: #0a0a0a;
    color: #00ff00;
}

.export-settings {
    display: flex;
    align-items: center;
    gap: 15px;
    flex-wrap: wrap;
}

.setting-item {
    display: flex;
    align-items: center;
    gap: 8px;
}

.setting-item label {
    font-size: 12px;
    color: #666;
    white-space: nowrap;
}

.path-selector {
    display: flex;
    align-items: center;
    gap: 8px;
}

.path-input {
    background: #0a0a0a;
    border: 1px solid #1a1a1a;
    color: #00ff00;
    padding: 6px 10px;
    font-family: inherit;
    font-size: 12px;
    border-radius: 3px;
    outline: none;
    width: 250px;
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

.db-content {
    flex: 1;
    display: flex;
    gap: 20px;
    padding: 20px;
    overflow: hidden;
}

.db-section {
    flex: 1;
    background: #0d0d0d;
    border: 1px solid #1a1a1a;
    border-radius: 4px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 15px;
    background: #1a1a1a;
    border-bottom: 1px solid #333;
}

.section-header h3 {
    margin: 0;
    font-size: 13px;
    font-weight: bold;
    letter-spacing: 1px;
    color: #00ff00;
}

.section-info {
    font-size: 11px;
    color: #666;
}

.db-list,
.table-list {
    flex: 1;
    overflow-y: auto;
    padding: 10px;
}

.db-item,
.table-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    border-bottom: 1px solid #1a1a1a;
    cursor: pointer;
    transition: background 0.1s;
}

.db-item:hover,
.table-item:hover {
    background: rgba(0, 255, 0, 0.05);
}

.db-item.selected,
.table-item.selected {
    background: rgba(0, 255, 0, 0.12);
    border-left: 2px solid #00ff00;
    margin-left: -2px;
}

.db-item input,
.table-item input {
    accent-color: #00ff00;
}

.db-name,
.table-name {
    font-size: 12px;
    color: #ccc;
}

.empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px;
    color: #444;
    font-size: 12px;
    letter-spacing: 1px;
}

.db-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 20px;
    background: #0d0d0d;
    border-top: 1px solid #1a1a1a;
    margin-top: 15px;
}

.export-info {
    font-size: 12px;
    color: #666;
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
    border: 1px solid #00ff00;
    border-radius: 4px;
    min-width: 400px;
    box-shadow: 0 0 30px rgba(0, 255, 0, 0.2);
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

.modal-icon {
    color: #00ff00;
}

.modal-icon.warn {
    color: #ffbd2e;
}

.modal-body {
    padding: 15px;
    font-size: 12px;
    color: #888;
}

.form-group {
    margin-bottom: 12px;
}

.form-group label {
    display: block;
    margin-bottom: 4px;
    color: #666;
    font-size: 11px;
    letter-spacing: 1px;
}

.form-input {
    width: 100%;
    background: #0a0a0a;
    border: 1px solid #1a1a1a;
    color: #00ff00;
    padding: 8px 10px;
    font-family: inherit;
    font-size: 12px;
    border-radius: 3px;
    outline: none;
    box-sizing: border-box;
}

.modal-foot {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding: 12px 15px;
    border-top: 1px solid #1a1a1a;
}

.export-summary {
    margin-top: 15px;
}

.export-summary h4 {
    margin: 0 0 8px 0;
    font-size: 11px;
    color: #666;
    letter-spacing: 1px;
}

.selected-items {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 12px;
}

.item-tag {
    background: #1a1a1a;
    padding: 3px 8px;
    border-radius: 3px;
    font-size: 11px;
    color: #00ff00;
}

.export-path {
    margin-top: 12px;
}

.export-path span {
    font-size: 11px;
    color: #ccc;
    background: #1a1a1a;
    padding: 4px 8px;
    border-radius: 3px;
    display: inline-block;
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

@media (max-width: 768px) {
    .db-content {
        flex-direction: column;
    }

    .db-toolbar {
        flex-direction: column;
        align-items: flex-start;
    }

    .export-settings {
        width: 100%;
    }

    .path-input {
        width: 100%;
    }
}
</style>
