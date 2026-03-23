import { defineStore } from "pinia";
import { ref, computed, toRaw } from "vue";
import {
  GetDatabaseConnections,
  GetDatabaseConnection,
  AddDatabaseConnection,
  UpdateDatabaseConnection,
  DeleteDatabaseConnection,
  ConnectDatabase,
  ListDatabases,
  ListTables,
  GetExportSettings,
  SaveExportSettings,
  ExportDatabases,
  ExportTables,
  ImportDatabases,
  ImportTables,
  CheckImportConflicts,
  DropConflictingTables,
} from "../../wailsjs/go/main/App";

export const useDatabaseStore = defineStore("database", () => {
  // 连接管理状态
  const connections = ref([]);
  const currentConnection = ref(null);
  const connectionsLoading = ref(false);
  const connectionsError = ref(null);

  // 数据库/表状态
  const databases = ref([]);
  const tables = ref([]);
  const selectedDatabase = ref("");
  const selectedTables = ref([]);
  const databasesLoading = ref(false);
  const tablesLoading = ref(false);

  // 导出状态
  const exportSettings = ref({
    export_tool: "mysql-shell",
    export_path: "",
    mysql_shell: {
      threads: 4,
      compression: "gzip",
      chunk_size: "64M",
      skip_definer: true,
      skip_binlog: false,
      overwrite: true,
    },
    mysql_dump: {
      compression: "gzip",
      single_transaction: true,
      routines: true,
      events: true,
      overwrite: true,
    },
  });
  const exportLoading = ref(false);
  const exportProgress = ref(0);
  const exportStatus = ref("");

  // 导入状态
  const importLoading = ref(false);
  const importProgress = ref(0);
  const importStatus = ref("");
  const importConflicts = ref([]);

  // 计算属性
  const hasConnections = computed(() => connections.value.length > 0);
  const hasDatabases = computed(() => databases.value.length > 0);
  const hasTables = computed(() => tables.value.length > 0);
  const hasSelectedTables = computed(() => selectedTables.value.length > 0);

  const connectionOptions = computed(() => {
    return connections.value.map((conn) => ({
      value: conn.id,
      label: conn.name,
      description: `${conn.host}:${conn.port}`,
    }));
  });

  // 连接管理操作
  async function fetchConnections() {
    connectionsLoading.value = true;
    connectionsError.value = null;

    try {
      const result = await GetDatabaseConnections();
      connections.value = result || [];
    } catch (err) {
      connectionsError.value = err.message || "获取连接信息失败";
      console.error("Failed to fetch connections:", err);
    } finally {
      connectionsLoading.value = false;
    }
  }

  async function fetchConnection(id) {
    try {
      const result = await GetDatabaseConnection(id);
      currentConnection.value = result;
      return result;
    } catch (err) {
      console.error("Failed to fetch connection:", err);
      return null;
    }
  }

  async function addConnection(connection) {
    try {
      await AddDatabaseConnection(connection);
      await fetchConnections();
      return true;
    } catch (err) {
      connectionsError.value = err.message || "添加连接失败";
      console.error("Failed to add connection:", err);
      return false;
    }
  }

  async function updateConnection(connection) {
    try {
      await UpdateDatabaseConnection(connection);
      await fetchConnections();
      return true;
    } catch (err) {
      connectionsError.value = err.message || "更新连接失败";
      console.error("Failed to update connection:", err);
      return false;
    }
  }

  async function deleteConnection(id) {
    try {
      await DeleteDatabaseConnection(id);
      await fetchConnections();

      // 如果删除的是当前连接，清除当前连接
      if (currentConnection.value?.id === id) {
        currentConnection.value = null;
      }

      return true;
    } catch (err) {
      connectionsError.value = err.message || "删除连接失败";
      console.error("Failed to delete connection:", err);
      return false;
    }
  }

  async function testConnection(connection) {
    try {
      await ConnectDatabase(connection);
      return true;
    } catch (err) {
      console.error("Failed to test connection:", err);
      return false;
    }
  }

  function selectConnection(id) {
    const conn = connections.value.find((c) => c.id === id);
    currentConnection.value = conn || null;
  }

  // 数据库/表操作
  async function fetchDatabases(connection) {
    if (!connection) return [];

    databasesLoading.value = true;
    try {
      const result = await ListDatabases(connection);
      databases.value = result || [];
      return result;
    } catch (err) {
      console.error("Failed to fetch databases:", err);
      return [];
    } finally {
      databasesLoading.value = false;
    }
  }

  async function fetchTables(connection, database) {
    if (!connection || !database) return [];

    tablesLoading.value = true;
    try {
      const conn = { ...connection, Database: database };
      const result = await ListTables(conn);
      tables.value = result || [];
      return result;
    } catch (err) {
      console.error("Failed to fetch tables:", err);
      return [];
    } finally {
      tablesLoading.value = false;
    }
  }

  function selectDatabase(database) {
    selectedDatabase.value = database;
    selectedTables.value = [];
  }

  function selectTable(table, selected) {
    if (selected) {
      if (!selectedTables.value.includes(table)) {
        selectedTables.value.push(table);
      }
    } else {
      const index = selectedTables.value.indexOf(table);
      if (index > -1) {
        selectedTables.value.splice(index, 1);
      }
    }
  }

  function selectAllTables(selected) {
    if (selected) {
      selectedTables.value = [...tables.value];
    } else {
      selectedTables.value = [];
    }
  }

  // 导出设置操作
  async function fetchExportSettings() {
    try {
      const result = await GetExportSettings();
      exportSettings.value = result || exportSettings.value;
      return result;
    } catch (err) {
      console.error("Failed to fetch export settings:", err);
      return null;
    }
  }

  async function saveExportSettings(settings) {
    try {
      await SaveExportSettings(settings);
      exportSettings.value = settings;
      return true;
    } catch (err) {
      console.error("Failed to save export settings:", err);
      return false;
    }
  }

  // 导出操作
  async function exportDatabases(connection, databases, config) {
    exportLoading.value = true;
    exportProgress.value = 0;
    exportStatus.value = "准备导出...";

    try {
      const request = {
        connection_id: connection.id,
        databases: databases,
        output_dir: config.output_dir,
        threads: config.threads || exportSettings.value.Threads,
        compression: config.compression || exportSettings.value.Compression,
        overwrite: config.overwrite ?? exportSettings.value.Overwrite,
        skip_definer: config.skip_definer ?? exportSettings.value.SkipDefiner,
        skip_binlog: config.skip_binlog ?? exportSettings.value.SkipBinlog,
      };

      exportStatus.value = "正在导出...";
      const result = await ExportDatabases(request);

      if (result.success) {
        exportStatus.value = "导出完成";
        exportProgress.value = 100;
        return result;
      } else {
        throw new Error(result.message);
      }
    } catch (err) {
      exportStatus.value = "导出失败";
      throw err;
    } finally {
      exportLoading.value = false;
    }
  }

  async function exportTables(connection, database, tables, config) {
    exportLoading.value = true;
    exportProgress.value = 0;
    exportStatus.value = "准备导出...";

    try {
      const request = {
        connection_id: connection.id,
        database: database,
        tables: tables,
        output_dir: config.output_dir,
        threads: config.threads || exportSettings.value.Threads,
        compression: config.compression || exportSettings.value.Compression,
        overwrite: config.overwrite ?? exportSettings.value.Overwrite,
        skip_definer: config.skip_definer ?? exportSettings.value.SkipDefiner,
        skip_binlog: config.skip_binlog ?? exportSettings.value.SkipBinlog,
      };

      exportStatus.value = "正在导出...";
      const result = await ExportTables(request);

      if (result.success) {
        exportStatus.value = "导出完成";
        exportProgress.value = 100;
        return result;
      } else {
        throw new Error(result.message);
      }
    } catch (err) {
      exportStatus.value = "导出失败";
      throw err;
    } finally {
      exportLoading.value = false;
    }
  }

  // 导入操作
  async function importDatabases(connection, config) {
    importLoading.value = true;
    importProgress.value = 0;
    importStatus.value = "准备导入...";

    try {
      const request = {
        connection_id: connection.id,
        input_dir: config.input_dir,
        threads: config.threads || exportSettings.value.Threads,
        schema: config.schema,
        reset_progress: config.reset_progress,
        wait_timeout: config.wait_timeout,
      };

      importStatus.value = "正在导入...";
      const result = await ImportDatabases(request);

      if (result.success) {
        importStatus.value = "导入完成";
        importProgress.value = 100;
        return result;
      } else {
        throw new Error(result.message);
      }
    } catch (err) {
      importStatus.value = "导入失败";
      throw err;
    } finally {
      importLoading.value = false;
    }
  }

  async function importTables(connection, database, config) {
    importLoading.value = true;
    importProgress.value = 0;
    importStatus.value = "准备导入...";

    try {
      const request = {
        connection_id: connection.id,
        database: database,
        input_dir: config.input_dir,
        threads: config.threads || exportSettings.value.Threads,
        reset_progress: config.reset_progress,
        wait_timeout: config.wait_timeout,
      };

      importStatus.value = "正在导入...";
      const result = await ImportTables(request);

      if (result.success) {
        importStatus.value = "导入完成";
        importProgress.value = 100;
        return result;
      } else {
        throw new Error(result.message);
      }
    } catch (err) {
      importStatus.value = "导入失败";
      throw err;
    } finally {
      importLoading.value = false;
    }
  }

  // 冲突检测操作
  async function checkImportConflicts(connection, inputDir) {
    try {
      const request = {
        connection_id: connection.id,
        input_dir: inputDir,
      };

      const result = await CheckImportConflicts(request);
      importConflicts.value = result.conflicts || [];
      return result;
    } catch (err) {
      console.error("Failed to check import conflicts:", err);
      throw err;
    }
  }

  async function dropConflictingTables(connectionId, conflicts) {
    try {
      const request = {
        connection_id: connectionId,
        conflicts: conflicts.map(c => ({
          schema: c.schema,
          tables: Array.from(c.tables || []),
          views: Array.from(c.views || []),
          events: Array.from(c.events || []),
          functions: Array.from(c.functions || []),
          procedures: Array.from(c.procedures || []),
        })),
      };

      await DropConflictingTables(request);
      importConflicts.value = [];
      return true;
    } catch (err) {
      console.error("Failed to drop conflicting tables:", err);
      throw err;
    }
  }

  // 重置状态
  function resetExportStatus() {
    exportLoading.value = false;
    exportProgress.value = 0;
    exportStatus.value = "";
  }

  function resetImportStatus() {
    importLoading.value = false;
    importProgress.value = 0;
    importStatus.value = "";
    importConflicts.value = [];
  }

  function clearError() {
    connectionsError.value = null;
  }

  return {
    // 连接管理状态
    connections,
    currentConnection,
    connectionsLoading,
    connectionsError,
    // 数据库/表状态
    databases,
    tables,
    selectedDatabase,
    selectedTables,
    databasesLoading,
    tablesLoading,
    // 导出状态
    exportSettings,
    exportLoading,
    exportProgress,
    exportStatus,
    // 导入状态
    importLoading,
    importProgress,
    importStatus,
    importConflicts,
    // 计算属性
    hasConnections,
    hasDatabases,
    hasTables,
    hasSelectedTables,
    connectionOptions,
    // 连接管理操作
    fetchConnections,
    fetchConnection,
    addConnection,
    updateConnection,
    deleteConnection,
    testConnection,
    selectConnection,
    // 数据库/表操作
    fetchDatabases,
    fetchTables,
    selectDatabase,
    selectTable,
    selectAllTables,
    // 导出设置操作
    fetchExportSettings,
    saveExportSettings,
    // 导出操作
    exportDatabases,
    exportTables,
    // 导入操作
    importDatabases,
    importTables,
    // 冲突检测操作
    checkImportConflicts,
    dropConflictingTables,
    // 重置状态
    resetExportStatus,
    resetImportStatus,
    clearError,
  };
});
