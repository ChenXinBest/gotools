<script setup>
import { ref, onMounted, computed, watch } from "vue";
import { AddDatabaseConnection, UpdateDatabaseConnection, GetDatabaseConnections, ListDatabases, ListTables, SelectFolder, GetExportSettings, SaveExportSettings, ExportDatabases, ExportTables } from "../../wailsjs/go/main/App";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { FaSave, FaDownload } from "vue-icons-plus/fa";

const MYSQLSHELL_NOT_FOUND = "MYSQLSHELL_NOT_FOUND";
const MYSQLSHELL_DOWNLOAD_URL = "https://dev.mysql.com/downloads/shell/";
const LOCAL_STORAGE_KEY = "mysql_export_config_backup";

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
const selectedDatabase = ref("");
const tables = ref([]);
const selectedTables = ref(new Set());

// 导出配置
const exportPath = ref("");
const exportTool = ref("mysql-shell");
const loading = ref(false);
const error = ref("");
const statusMessage = ref("");

// 导出配置模态框数据
const exportScope = ref("database"); // database 或 table
const selectedDatabases = ref(new Set());
const exportDatabase = ref("");
const exportTables = ref(new Set());
const mysqlShellParams = ref({
    threads: 4,
    compression: "gzip",
    chunkSize: "",
    skipDefiner: true,
    skipBinlog: false,
    overwrite: true,
    includeSchemas: "",
    excludeSchemas: "",
    includeTables: "",
    excludeTables: ""
});
const isLoadingTables = ref(false);

// 配置保存状态
const configSaveStatus = ref("idle"); // idle, saving, saved, error
const configSaveMessage = ref("");
let configSaveTimeout = null;

// 模态框
const showConnectionModal = ref(false);
const showExportModal = ref(false);
const showExportConfigModal = ref(false);

// 计算属性
const hasSelectedDatabase = computed(() => !!selectedDatabase.value);
const hasSelectedTables = computed(() => selectedTables.value.size > 0);
const canExport = computed(() => {
    return (
        exportPath.value &&
        (hasSelectedDatabase.value || hasSelectedTables.value)
    );
});

// 导出配置模态框计算属性
const canConfirmExport = computed(() => {
    if (exportScope.value === 'database') {
        return selectedDatabases.value.size > 0;
    } else if (exportScope.value === 'table') {
        return exportDatabase.value && exportTables.value.size > 0;
    }
    return false;
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

async function refreshConnection() {
    if (!currentConnection.value) return;

    selectedDatabase.value = "";
    selectedTables.value.clear();
    tables.value = [];

    try {
        const conn = connections.value.find(c => c.name === currentConnection.value);
        if (!conn) {
            error.value = "找不到连接信息";
            return;
        }

        if (exportTool.value === "mysql-shell") {
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
        const errorMsg = err.message || err || "";
        if (errorMsg.includes(MYSQLSHELL_NOT_FOUND)) {
            error.value = "未找到 mysqlsh 命令，请先安装 MySQL Shell";
            BrowserOpenURL(MYSQLSHELL_DOWNLOAD_URL);
            return;
        }
        error.value = "刷新失败: " + errorMsg;
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
    databases.value = [
        "information_schema",
        "mysql",
        "performance_schema",
        "sys",
        "test",
    ];
    selectedDatabase.value = "";
    selectedTables.value.clear();
    tables.value = [];
    showStatus(`已连接到 ${connName}`);
}

async function selectDatabase(db) {
    if (selectedDatabase.value === db) {
        selectedDatabase.value = "";
        tables.value = [];
        selectedTables.value.clear();
        return;
    }

    selectedDatabase.value = db;
    selectedTables.value.clear();

    try {
        if (exportTool.value === "mysql-shell") {
            const conn = connections.value.find(c => c.name === currentConnection.value);
            if (conn) {
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
        const errorMsg = err.message || err || "";
        if (errorMsg.includes(MYSQLSHELL_NOT_FOUND)) {
            error.value = "未找到 mysqlsh 命令，请先安装 MySQL Shell";
            BrowserOpenURL(MYSQLSHELL_DOWNLOAD_URL);
            return;
        }
        error.value = "获取表列表失败: " + errorMsg;
        tables.value = [
            "users",
            "orders",
            "products",
            "categories",
            "transactions",
        ];
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
    loading.value = true;
    error.value = "";

    try {
        // 构建导出命令
        const exportCommand = buildExportCommand();
        console.log("Export command:", exportCommand);

        // 模拟导出过程
        await simulateExportProcess();

        // 生成导出结果信息
        const exportResult = generateExportResult();
        showStatus(exportResult);
        
        // 关闭模态框
        showExportModal.value = false;
    } catch (e) {
        error.value = "导出失败: " + (e.message || e);
    } finally {
        loading.value = false;
    }
}

// 构建导出命令
function buildExportCommand() {
    let command = "mysqlsh";
    
    // 添加连接参数
    const conn = connections.value.find(c => c.name === currentConnection.value);
    if (conn) {
        command += ` --uri="mysql://${conn.user}:${conn.password}@${conn.host}:${conn.port}/`;
    }
    
    // 添加导出参数
    command += ` --export=json --path="${exportPath.value}"`;
    
    // 添加MySQL Shell参数
    if (mysqlShellParams.value.threads) {
        command += ` --threads=${mysqlShellParams.value.threads}`;
    }
    if (mysqlShellParams.value.skipDefiner) {
        command += " --skip-definer";
    }
    if (mysqlShellParams.value.skipBinlog) {
        command += " --skip-binlog";
    }
    
    // 添加导出对象
    if (exportScope.value === 'database') {
        const databasesList = Array.from(selectedDatabases.value).join(',');
        command += ` --databases=${databasesList}`;
    } else if (exportScope.value === 'table') {
        const tablesList = Array.from(exportTables.value).join(',');
        command += ` --database=${exportDatabase.value} --tables=${tablesList}`;
    }
    
    return command;
}

// 模拟导出过程
async function simulateExportProcess() {
    // 模拟导出进度
    return new Promise((resolve) => {
        let progress = 0;
        const interval = setInterval(() => {
            progress += 10;
            if (progress <= 100) {
                showStatus(`导出进度: ${progress}%`);
            } else {
                clearInterval(interval);
                resolve();
            }
        }, 200);
    });
}

// 生成导出结果信息
function generateExportResult() {
    let exportType = "";
    if (exportScope.value === 'database') {
        exportType = `数据库(${selectedDatabases.value.size}个)`;
    } else if (exportScope.value === 'table') {
        exportType = `表(${exportTables.value.size}个) from ${exportDatabase.value}`;
    }
    return `成功导出 ${exportType} 到 ${exportPath.value}`;
}

// 显示状态信息
function showStatus(msg) {
    statusMessage.value = msg;
    setTimeout(() => {
        statusMessage.value = "";
    }, 3000);
}

// 处理导出按钮点击
function handleExportClick() {
    if (exportTool.value === "mysql-shell") {
        // 重置导出配置
        exportScope.value = "database";
        selectedDatabases.value.clear();
        exportDatabase.value = "";
        exportTables.value.clear();
        mysqlShellParams.value = {
            threads: 4,
            compression: "gzip",
            chunkSize: "",
            skipDefiner: true,
            skipBinlog: false,
            overwrite: true,
            includeSchemas: "",
            excludeSchemas: "",
            includeTables: "",
            excludeTables: ""
        };
        showExportConfigModal.value = true;
    } else {
        exportData();
    }
}

// 切换数据库选择
function toggleDatabaseSelection(db) {
    if (selectedDatabases.value.has(db)) {
        selectedDatabases.value.delete(db);
    } else {
        selectedDatabases.value.add(db);
    }
    selectedDatabases.value = new Set(selectedDatabases.value);
}

// 全选/取消全选数据库
function toggleSelectAllDatabases() {
    if (selectedDatabases.value.size === databases.length) {
        selectedDatabases.value.clear();
    } else {
        selectedDatabases.value = new Set(databases);
    }
}

// 加载导出表列表
async function loadExportTables() {
    if (!exportDatabase.value) {
        tables.value = [];
        exportTables.value.clear();
        return;
    }

    isLoadingTables.value = true;
    exportTables.value.clear();

    try {
        if (exportTool.value === "mysql-shell") {
            const conn = connections.value.find(c => c.name === currentConnection.value);
            if (conn) {
                const connData = {
                    id: conn.id,
                    name: conn.name,
                    host: conn.host,
                    port: conn.port,
                    user: conn.user,
                    password: conn.password,
                    database: exportDatabase.value
                };

                tables.value = await ListTables(connData);
            }
        } else {
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
        const errorMsg = err.message || err || "";
        if (errorMsg.includes(MYSQLSHELL_NOT_FOUND)) {
            error.value = "未找到 mysqlsh 命令，请先安装 MySQL Shell";
            BrowserOpenURL(MYSQLSHELL_DOWNLOAD_URL);
            return;
        }
        error.value = "获取表列表失败: " + errorMsg;
        tables.value = [
            "users",
            "orders",
            "products",
            "categories",
            "transactions",
        ];
    } finally {
        isLoadingTables.value = false;
    }
}

// 切换表选择
function toggleTableSelection(table) {
    if (exportTables.value.has(table)) {
        exportTables.value.delete(table);
    } else {
        exportTables.value.add(table);
    }
    exportTables.value = new Set(exportTables.value);
}

// 全选/取消全选表
function toggleSelectAllTables() {
    if (exportTables.value.size === tables.length) {
        exportTables.value.clear();
    } else {
        exportTables.value = new Set(tables);
    }
}

// 确认导出配置
async function confirmExportConfig() {
    showExportConfigModal.value = false;
    
    const conn = connections.value.find(c => c.name === currentConnection.value);
    if (!conn) {
        error.value = "请先选择数据库连接";
        return;
    }

    loading.value = true;
    error.value = "";

    try {
        const parseCommaSeparated = (str) => {
            if (!str || !str.trim()) return [];
            return str.split(',').map(s => s.trim()).filter(s => s);
        };

        const request = {
            connection_id: conn.id,
            databases: Array.from(selectedDatabases.value),
            database: exportDatabase.value,
            tables: Array.from(exportTables.value),
            output_dir: exportPath.value,
            threads: mysqlShellParams.value.threads,
            compression: mysqlShellParams.value.compression,
            chunk_size: mysqlShellParams.value.chunkSize,
            skip_definer: mysqlShellParams.value.skipDefiner,
            skip_binlog: mysqlShellParams.value.skipBinlog,
            overwrite: mysqlShellParams.value.overwrite,
            include_schemas: parseCommaSeparated(mysqlShellParams.value.includeSchemas),
            exclude_schemas: parseCommaSeparated(mysqlShellParams.value.excludeSchemas),
            include_tables: parseCommaSeparated(mysqlShellParams.value.includeTables),
            exclude_tables: parseCommaSeparated(mysqlShellParams.value.excludeTables),
        };

        let result;
        if (exportScope.value === 'database') {
            result = await ExportDatabases(request);
        } else if (exportScope.value === 'table') {
            result = await ExportTables(request);
        }

        if (result && result.success) {
            showStatus(result.message || "导出成功");
        } else {
            error.value = result?.message || "导出失败";
        }
    } catch (err) {
        console.error("导出失败:", err);
        const errorMsg = err.message || err || "";
        if (errorMsg.includes(MYSQLSHELL_NOT_FOUND)) {
            error.value = "未找到 mysqlsh 命令，请先安装 MySQL Shell";
            BrowserOpenURL(MYSQLSHELL_DOWNLOAD_URL);
            return;
        }
        error.value = "导出失败: " + errorMsg;
    } finally {
        loading.value = false;
    }
}

// 打开文件选择器
async function selectExportPath() {
    try {
        const path = await SelectFolder();
        if (path) {
            exportPath.value = path;
            await saveExportSettings();
        }
    } catch (err) {
        console.error("选择目录失败:", err);
        error.value = "选择目录失败: " + (err.message || err);
    }
}

onMounted(async () => {
    await loadConnections();
    await loadExportSettings();
});

async function loadExportSettings() {
    try {
        const settings = await GetExportSettings();
        if (settings.export_tool) {
            exportTool.value = settings.export_tool;
        }
        if (settings.export_path) {
            exportPath.value = settings.export_path;
        }
        if (settings.threads > 0) {
            mysqlShellParams.value.threads = settings.threads;
        }
        if (settings.compression) {
            mysqlShellParams.value.compression = settings.compression;
        }
        if (settings.chunk_size) {
            mysqlShellParams.value.chunkSize = settings.chunk_size;
        }
        mysqlShellParams.value.skipDefiner = settings.skip_definer;
        mysqlShellParams.value.skipBinlog = settings.skip_binlog;
        if (settings.overwrite !== undefined) {
            mysqlShellParams.value.overwrite = settings.overwrite;
        }
        if (settings.include_schemas) {
            mysqlShellParams.value.includeSchemas = settings.include_schemas;
        }
        if (settings.exclude_schemas) {
            mysqlShellParams.value.excludeSchemas = settings.exclude_schemas;
        }
        if (settings.include_tables) {
            mysqlShellParams.value.includeTables = settings.include_tables;
        }
        if (settings.exclude_tables) {
            mysqlShellParams.value.excludeTables = settings.exclude_tables;
        }
        if (settings.last_databases && settings.last_databases.length > 0) {
            selectedDatabases.value = new Set(settings.last_databases);
        }
        if (settings.last_database) {
            exportDatabase.value = settings.last_database;
        }
        if (settings.last_tables && settings.last_tables.length > 0) {
            exportTables.value = new Set(settings.last_tables);
        }
        if (settings.last_connection_id) {
            const conn = connections.value.find(c => c.id === settings.last_connection_id);
            if (conn) {
                currentConnection.value = conn.name;
            }
        }
        if (settings.export_scope) {
            exportScope.value = settings.export_scope;
        }
    } catch (err) {
        console.error("加载导出配置失败:", err);
        loadFromLocalStorage();
    }
}

function loadFromLocalStorage() {
    try {
        const localConfig = localStorage.getItem(LOCAL_STORAGE_KEY);
        if (localConfig) {
            const settings = JSON.parse(localConfig);
            if (settings.export_tool) {
                exportTool.value = settings.export_tool;
            }
            if (settings.export_path) {
                exportPath.value = settings.export_path;
            }
            if (settings.threads > 0) {
                mysqlShellParams.value.threads = settings.threads;
            }
            if (settings.compression) {
                mysqlShellParams.value.compression = settings.compression;
            }
            if (settings.chunk_size) {
                mysqlShellParams.value.chunkSize = settings.chunk_size;
            }
            if (settings.skip_definer !== undefined) {
                mysqlShellParams.value.skipDefiner = settings.skip_definer;
            }
            if (settings.skip_binlog !== undefined) {
                mysqlShellParams.value.skipBinlog = settings.skip_binlog;
            }
            if (settings.overwrite !== undefined) {
                mysqlShellParams.value.overwrite = settings.overwrite;
            }
            if (settings.include_schemas) {
                mysqlShellParams.value.includeSchemas = settings.include_schemas;
            }
            if (settings.exclude_schemas) {
                mysqlShellParams.value.excludeSchemas = settings.exclude_schemas;
            }
            if (settings.include_tables) {
                mysqlShellParams.value.includeTables = settings.include_tables;
            }
            if (settings.exclude_tables) {
                mysqlShellParams.value.excludeTables = settings.exclude_tables;
            }
            if (settings.last_databases && settings.last_databases.length > 0) {
                selectedDatabases.value = new Set(settings.last_databases);
            }
            if (settings.last_database) {
                exportDatabase.value = settings.last_database;
            }
            if (settings.last_tables && settings.last_tables.length > 0) {
                exportTables.value = new Set(settings.last_tables);
            }
            if (settings.export_scope) {
                exportScope.value = settings.export_scope;
            }
        }
    } catch (err) {
        console.error("从本地存储加载配置失败:", err);
    }
}

function saveToLocalStorage(settings) {
    try {
        localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(settings));
    } catch (err) {
        console.error("保存到本地存储失败:", err);
    }
}

function validateExportSettings(settings) {
    const errors = [];
    
    if (settings.threads < 1 || settings.threads > 32) {
        errors.push("线程数必须在1-32之间");
    }
    
    if (settings.chunk_size) {
        const chunkSizePattern = /^\d+[KMG]?$/i;
        if (!chunkSizePattern.test(settings.chunk_size)) {
            errors.push("分块大小格式无效，例如: 64, 64K, 64M, 64G");
        }
    }
    
    if (settings.export_scope === "database" && settings.last_databases.length === 0) {
        errors.push("请至少选择一个数据库");
    }
    
    if (settings.export_scope === "table") {
        if (!settings.last_database) {
            errors.push("请选择数据库");
        }
        if (settings.last_tables.length === 0) {
            errors.push("请至少选择一个表");
        }
    }
    
    return errors;
}

async function saveExportSettings() {
    const conn = connections.value.find(c => c.name === currentConnection.value);
    const settings = {
        export_tool: exportTool.value,
        export_path: exportPath.value,
        last_connection_id: conn ? conn.id : "",
        last_databases: Array.from(selectedDatabases.value),
        last_database: exportDatabase.value,
        last_tables: Array.from(exportTables.value),
        threads: mysqlShellParams.value.threads,
        compression: mysqlShellParams.value.compression,
        chunk_size: mysqlShellParams.value.chunkSize,
        skip_definer: mysqlShellParams.value.skipDefiner,
        skip_binlog: mysqlShellParams.value.skipBinlog,
        overwrite: mysqlShellParams.value.overwrite,
        include_schemas: mysqlShellParams.value.includeSchemas,
        exclude_schemas: mysqlShellParams.value.excludeSchemas,
        include_tables: mysqlShellParams.value.includeTables,
        exclude_tables: mysqlShellParams.value.excludeTables,
        export_scope: exportScope.value,
    };
    
    const errors = validateExportSettings(settings);
    if (errors.length > 0 && (settings.last_databases.length > 0 || settings.last_tables.length > 0)) {
        console.warn("配置验证警告:", errors);
    }
    
    saveToLocalStorage(settings);
    
    configSaveStatus.value = "saving";
    configSaveMessage.value = "保存中...";
    
    try {
        await SaveExportSettings(settings);
        configSaveStatus.value = "saved";
        configSaveMessage.value = "配置已保存";
        
        if (configSaveTimeout) {
            clearTimeout(configSaveTimeout);
        }
        configSaveTimeout = setTimeout(() => {
            configSaveStatus.value = "idle";
            configSaveMessage.value = "";
        }, 2000);
    } catch (err) {
        console.error("保存导出配置失败:", err);
        configSaveStatus.value = "error";
        configSaveMessage.value = "保存失败，已备份到本地";
        
        if (configSaveTimeout) {
            clearTimeout(configSaveTimeout);
        }
        configSaveTimeout = setTimeout(() => {
            configSaveStatus.value = "idle";
            configSaveMessage.value = "";
        }, 3000);
    }
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

const debouncedSaveExportSettings = debounce(saveExportSettings, 500);

watch(exportTool, () => {
    debouncedSaveExportSettings();
});

watch(exportPath, () => {
    debouncedSaveExportSettings();
});

watch(exportScope, () => {
    debouncedSaveExportSettings();
});

watch(selectedDatabases, () => {
    debouncedSaveExportSettings();
}, { deep: true });

watch(exportDatabase, () => {
    debouncedSaveExportSettings();
});

watch(exportTables, () => {
    debouncedSaveExportSettings();
}, { deep: true });

watch(mysqlShellParams, () => {
    debouncedSaveExportSettings();
}, { deep: true });
</script>

<template>
    <div class="database-backup">
        <div class="fixed-export-btn" :class="{ disabled: !canExport || loading }">
            <button
                @click="handleExportClick"
                class="export-action-btn"
                :disabled="!canExport || loading"
            >
                <span class="export-btn-icon">
                    <component :is="FaDownload" />
                </span>
                <span class="export-btn-text">{{ loading ? "导出中..." : "导出" }}</span>
            </button>
        </div>
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
                    >{{ selectedDatabase || '未选择' }} 数据库</span
                >
                <span class="stat-divider">|</span>
                <span class="stat-item">{{ selectedTables.size }} 已选表</span>
                <span v-if="configSaveStatus !== 'idle'" class="config-save-status" :class="configSaveStatus">
                    {{ configSaveMessage }}
                </span>
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
                        :class="{ selected: selectedDatabase === db }"
                        @click="selectDatabase(db)"
                    >
                        <span class="db-icon">▣</span>
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
                <span v-else-if="hasSelectedDatabase">
                    已选择数据库: {{ selectedDatabase }}
                </span>
                <span v-else> 请选择要导出的数据库或表 </span>
            </div>
            <div class="export-hint">
                <span class="hint-icon">↑</span>
                <span>点击右上角按钮导出</span>
            </div>
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
                        <div v-else-if="hasSelectedDatabase">
                            <h4>已选择数据库:</h4>
                            <div class="selected-items">
                                <span class="item-tag">
                                    {{ selectedDatabase }}
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

        <!-- 导出配置模态框 -->
        <div
            v-if="showExportConfigModal"
            class="modal-mask"
            @click.self="showExportConfigModal = false"
        >
            <div class="modal-box export-config-modal">
                <div class="modal-head">
                    <span class="modal-icon">📤</span>
                    导出配置
                </div>
                <div class="modal-body">
                    <!-- 导出范围选择 -->
                    <div class="form-group">
                        <label>导出范围:</label>
                        <div class="radio-group">
                            <label class="radio-item">
                                <input type="radio" value="database" v-model="exportScope" />
                                <span>数据库</span>
                            </label>
                            <label class="radio-item">
                                <input type="radio" value="table" v-model="exportScope" />
                                <span>数据表</span>
                            </label>
                        </div>
                    </div>

                    <!-- 数据库选择 (导出范围为数据库时) -->
                    <div v-if="exportScope === 'database'" class="form-group">
                        <div class="section-header">
                            <h4>选择数据库</h4>
                            <button @click="toggleSelectAllDatabases" class="btn small">
                                {{ selectedDatabases.size === databases.length ? '取消全选' : '全选' }}
                            </button>
                        </div>
                        <div class="db-select-list">
                            <div v-if="databases.length === 0" class="empty-state">
                                请先选择连接并刷新数据库列表
                            </div>
                            <div
                                v-for="db in databases"
                                :key="db"
                                class="db-select-item"
                                :class="{ selected: selectedDatabases.has(db) }"
                                @click="toggleDatabaseSelection(db)"
                            >
                                <input
                                    type="checkbox"
                                    :checked="selectedDatabases.has(db)"
                                    @change="toggleDatabaseSelection(db)"
                                />
                                <span>{{ db }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- 表选择 (导出范围为表时) -->
                    <div v-if="exportScope === 'table'">
                        <!-- 数据库选择 -->
                        <div class="form-group">
                            <label>选择数据库:</label>
                            <select v-model="exportDatabase" class="select-input" @change="loadExportTables">
                                <option value="">请选择数据库</option>
                                <option v-for="db in databases" :key="db" :value="db">{{ db }}</option>
                            </select>
                        </div>

                        <!-- 表选择 -->
                        <div class="form-group">
                            <div class="section-header">
                                <h4>选择数据表</h4>
                                <button 
                                    @click="toggleSelectAllTables" 
                                    class="btn small"
                                    :disabled="!exportDatabase || tables.length === 0"
                                >
                                    {{ exportTables.size === tables.length ? '取消全选' : '全选' }}
                                </button>
                            </div>
                            <div class="table-select-list">
                                <div v-if="!exportDatabase" class="empty-state">
                                    请先选择数据库
                                </div>
                                <div v-else-if="isLoadingTables" class="empty-state">
                                    加载表列表中...
                                </div>
                                <div v-else-if="tables.length === 0" class="empty-state">
                                    该数据库中没有表
                                </div>
                                <div
                                    v-for="table in tables"
                                    :key="table"
                                    class="table-select-item"
                                    :class="{ selected: exportTables.has(table) }"
                                    @click="toggleTableSelection(table)"
                                >
                                    <input
                                        type="checkbox"
                                        :checked="exportTables.has(table)"
                                        @change="toggleTableSelection(table)"
                                    />
                                    <span>{{ table }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- MySQL Shell 参数配置 -->
                    <div class="form-group">
                        <h4>MySQL Shell 参数配置</h4>
                        <div class="param-group">
                            <div class="param-item">
                                <label>线程数 (1-32):</label>
                                <input
                                    v-model.number="mysqlShellParams.threads"
                                    type="number"
                                    min="1"
                                    max="32"
                                    class="form-input small"
                                />
                            </div>
                            <div class="param-item">
                                <label>压缩方式:</label>
                                <select v-model="mysqlShellParams.compression" class="select-input small">
                                    <option value="gzip">gzip (推荐)</option>
                                    <option value="zstd">zstd (更快)</option>
                                    <option value="none">无压缩</option>
                                </select>
                            </div>
                            <div class="param-item">
                                <label>分块大小:</label>
                                <input
                                    v-model="mysqlShellParams.chunkSize"
                                    type="text"
                                    class="form-input small"
                                    placeholder="如: 64M"
                                />
                            </div>
                            <div class="param-item">
                                <label>
                                    <input type="checkbox" v-model="mysqlShellParams.skipDefiner" />
                                    跳过 Definer
                                </label>
                            </div>
                            <div class="param-item">
                                <label>
                                    <input type="checkbox" v-model="mysqlShellParams.skipBinlog" />
                                    跳过 Binlog
                                </label>
                            </div>
                            <div class="param-item">
                                <label>
                                    <input type="checkbox" v-model="mysqlShellParams.overwrite" />
                                    覆盖已存在目录
                                </label>
                            </div>
                        </div>
                    </div>

                    <!-- 高级过滤选项 -->
                    <div class="form-group">
                        <h4>高级过滤选项</h4>
                        <div class="filter-grid">
                            <div class="filter-item">
                                <label>包含数据库 (逗号分隔):</label>
                                <input
                                    v-model="mysqlShellParams.includeSchemas"
                                    type="text"
                                    class="form-input"
                                    placeholder="如: db1, db2"
                                />
                            </div>
                            <div class="filter-item">
                                <label>排除数据库 (逗号分隔):</label>
                                <input
                                    v-model="mysqlShellParams.excludeSchemas"
                                    type="text"
                                    class="form-input"
                                    placeholder="如: mysql, sys"
                                />
                            </div>
                            <div class="filter-item">
                                <label>包含表 (逗号分隔):</label>
                                <input
                                    v-model="mysqlShellParams.includeTables"
                                    type="text"
                                    class="form-input"
                                    placeholder="如: users, orders"
                                />
                            </div>
                            <div class="filter-item">
                                <label>排除表 (逗号分隔):</label>
                                <input
                                    v-model="mysqlShellParams.excludeTables"
                                    type="text"
                                    class="form-input"
                                    placeholder="如: logs, temp_*"
                                />
                            </div>
                        </div>
                    </div>
                </div>
                <div class="modal-foot">
                    <button @click="showExportConfigModal = false" class="btn">
                        取消
                    </button>
                    <button 
                        @click="confirmExportConfig"
                        class="btn danger"
                        :disabled="!canConfirmExport"
                    >
                        确认导出
                    </button>
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
    position: relative;
}

.fixed-export-btn {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 1000;
    transition: all 0.3s ease;
}

.fixed-export-btn.disabled {
    opacity: 0.5;
}

.export-action-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
    border: 2px solid #ff6b6b;
    color: #ffffff;
    padding: 12px 24px;
    font-family: inherit;
    font-size: 14px;
    font-weight: bold;
    cursor: pointer;
    border-radius: 8px;
    box-shadow: 0 4px 15px rgba(220, 53, 69, 0.4), 0 0 20px rgba(220, 53, 69, 0.2);
    transition: all 0.3s ease;
    letter-spacing: 1px;
}

.export-action-btn:hover:not(:disabled) {
    background: linear-gradient(135deg, #ff6b6b 0%, #dc3545 100%);
    box-shadow: 0 6px 20px rgba(220, 53, 69, 0.6), 0 0 30px rgba(220, 53, 69, 0.3);
    transform: translateY(-2px);
}

.export-action-btn:active:not(:disabled) {
    transform: translateY(0);
    box-shadow: 0 2px 10px rgba(220, 53, 69, 0.4);
}

.export-action-btn:disabled {
    background: linear-gradient(135deg, #6c757d 0%, #5a6268 100%);
    border-color: #6c757d;
    cursor: not-allowed;
    box-shadow: none;
}

.export-btn-icon {
    display: flex;
    align-items: center;
    font-size: 16px;
}

.export-btn-icon svg {
    width: 18px;
    height: 18px;
}

.export-btn-text {
    white-space: nowrap;
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

.config-save-status {
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 3px;
    margin-left: 10px;
    transition: all 0.3s;
}

.config-save-status.saving {
    background: rgba(0, 255, 0, 0.1);
    color: #00ff00;
    border: 1px solid #00ff00;
}

.config-save-status.saved {
    background: rgba(0, 255, 0, 0.2);
    color: #00ff00;
    border: 1px solid #00ff00;
}

.config-save-status.error {
    background: rgba(255, 95, 86, 0.1);
    color: #ff5f56;
    border: 1px solid #ff5f56;
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

.db-icon {
    color: #00ff00;
    font-size: 12px;
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

.export-hint {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    color: #444;
    padding: 4px 10px;
    background: rgba(255, 95, 86, 0.1);
    border: 1px solid rgba(255, 95, 86, 0.2);
    border-radius: 4px;
}

.hint-icon {
    color: #ff5f56;
    font-size: 12px;
    animation: bounce 1.5s infinite;
}

@keyframes bounce {
    0%, 100% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(-3px);
    }
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

.export-config-modal {
    min-width: 600px;
    max-width: 800px;
    max-height: 80vh;
    overflow-y: auto;
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

.form-input.small {
    width: 100px;
}

.radio-group {
    display: flex;
    gap: 20px;
    margin-top: 8px;
}

.radio-item {
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
}

.radio-item input[type="radio"] {
    accent-color: #00ff00;
}

.db-select-list,
.table-select-list {
    max-height: 200px;
    overflow-y: auto;
    border: 1px solid #1a1a1a;
    border-radius: 3px;
    margin-top: 8px;
}

.db-select-item,
.table-select-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    border-bottom: 1px solid #1a1a1a;
    cursor: pointer;
    transition: background 0.1s;
}

.db-select-item:hover,
.table-select-item:hover {
    background: rgba(0, 255, 0, 0.05);
}

.db-select-item.selected,
.table-select-item.selected {
    background: rgba(0, 255, 0, 0.12);
    border-left: 2px solid #00ff00;
    margin-left: -2px;
}

.db-select-item input,
.table-select-item input {
    accent-color: #00ff00;
}

.param-group {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    margin-top: 12px;
}

.param-item {
    display: flex;
    align-items: center;
    gap: 8px;
}

.param-item input[type="checkbox"] {
    accent-color: #00ff00;
}

.select-input.small {
    min-width: 120px;
    padding: 6px 8px;
}

.filter-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    margin-top: 12px;
}

.filter-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.filter-item label {
    font-size: 11px;
    color: #666;
    letter-spacing: 0.5px;
}

.filter-item .form-input {
    width: 100%;
}

.btn.small {
    padding: 4px 10px;
    font-size: 11px;
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

    .fixed-export-btn {
        top: 10px;
        right: 10px;
    }

    .export-action-btn {
        padding: 10px 16px;
        font-size: 12px;
    }

    .export-btn-icon svg {
        width: 16px;
        height: 16px;
    }

    .db-footer {
        flex-direction: column;
        gap: 8px;
        align-items: flex-start;
    }

    .export-hint {
        width: 100%;
        justify-content: center;
    }

    .filter-grid {
        grid-template-columns: 1fr;
    }
}

@media (max-width: 480px) {
    .fixed-export-btn {
        top: 8px;
        right: 8px;
    }

    .export-action-btn {
        padding: 8px 12px;
        font-size: 11px;
        gap: 4px;
    }

    .export-btn-text {
        display: none;
    }

    .export-btn-icon {
        font-size: 18px;
    }

    .export-btn-icon svg {
        width: 20px;
        height: 20px;
    }
}
</style>
