<template>
    <div class="settings">
        <div class="header">
            <h1>{{ t("settings.title") }}</h1>
            <p class="subtitle">{{ t("settings.subtitle") }}</p>
        </div>

        <div class="settings-content">
            <!-- 常规设置 -->
            <div class="settings-section">
                <h2>{{ t("settings.general") }}</h2>
                <div class="setting-item">
                    <div class="setting-label">
                        <span class="setting-name">{{
                            t("settings.theme")
                        }}</span>
                        <span class="setting-description"
                            >Choose application theme</span
                        >
                    </div>
                    <select v-model="settings.theme" class="setting-control">
                        <option value="dark">
                            {{ t("settings.theme.dark") }}
                        </option>
                        <option value="light">
                            {{ t("settings.theme.light") }}
                        </option>
                    </select>
                </div>

                <div class="setting-item">
                    <div class="setting-label">
                        <span class="setting-name">{{
                            t("settings.language")
                        }}</span>
                        <span class="setting-description"
                            >Choose interface display language</span
                        >
                    </div>
                    <select v-model="settings.locale" class="setting-control">
                        <option value="zh-CN">简体中文</option>
                        <option value="en-US">English</option>
                    </select>
                </div>

                <div class="setting-item">
                    <div class="setting-label">
                        <span class="setting-name">{{
                            t("settings.autoRefresh")
                        }}</span>
                        <span class="setting-description"
                            >Auto refresh process list</span
                        >
                    </div>
                    <label class="switch">
                        <input
                            type="checkbox"
                            v-model="settings.autoRefreshProcesses"
                        />
                        <span class="slider"></span>
                    </label>
                </div>

                <div class="setting-item" v-if="settings.autoRefreshProcesses">
                    <div class="setting-label">
                        <span class="setting-name">{{
                            t("settings.refreshInterval")
                        }}</span>
                        <span class="setting-description"
                            >Process list refresh interval (seconds)</span
                        >
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

            <!-- 数据库导入/导出设置 -->
            <div class="settings-section">
                <h2>{{ t("settings.database") }}</h2>

                <!-- 默认工具选择 -->
                <div class="setting-item">
                    <div class="setting-label">
                        <span class="setting-name">{{
                            t("settings.defaultExportTool")
                        }}</span>
                        <span class="setting-description"
                            >Select default import/export tool</span
                        >
                    </div>
                    <select
                        v-model="exportSettings.export_tool"
                        class="setting-control"
                    >
                        <option value="mysql-shell">
                            {{ t("db.mysqlShell") }}
                        </option>
                        <option value="mysqldump">
                            {{ t("db.mysqldump") }}
                        </option>
                    </select>
                </div>

                <!-- MySQL Shell 方案配置 -->
                <div class="scheme-config">
                    <h3>{{ t("settings.mysqlShell") }}</h3>
                    <div class="scheme-options">
                        <div class="option-row">
                            <div class="option-item">
                                <label>{{ t("settings.threads") }}</label>
                                <input
                                    v-model.number="
                                        exportSettings.mysql_shell.threads
                                    "
                                    type="number"
                                    min="1"
                                    max="16"
                                    class="setting-control small"
                                />
                            </div>
                            <div class="option-item">
                                <label>{{ t("settings.compression") }}</label>
                                <select
                                    v-model="
                                        exportSettings.mysql_shell.compression
                                    "
                                    class="setting-control"
                                >
                                    <option value="gzip">gzip</option>
                                    <option value="zstd">zstd</option>
                                    <option value="none">None</option>
                                </select>
                            </div>
                            <div class="option-item">
                                <label>{{ t("settings.chunkSize") }}</label>
                                <select
                                    v-model="
                                        exportSettings.mysql_shell.chunk_size
                                    "
                                    class="setting-control"
                                >
                                    <option value="64M">64M</option>
                                    <option value="128M">128M</option>
                                    <option value="256M">256M</option>
                                    <option value="512M">512M</option>
                                </select>
                            </div>
                        </div>
                        <div class="option-row">
                            <div class="option-item checkbox">
                                <label class="checkbox-label">
                                    <input
                                        type="checkbox"
                                        v-model="
                                            exportSettings.mysql_shell.overwrite
                                        "
                                    />
                                    <span>{{ t("settings.overwrite") }}</span>
                                </label>
                            </div>
                            <div class="option-item checkbox">
                                <label class="checkbox-label">
                                    <input
                                        type="checkbox"
                                        v-model="
                                            exportSettings.mysql_shell
                                                .skip_definer
                                        "
                                    />
                                    <span>{{ t("settings.skipDefiner") }}</span>
                                </label>
                            </div>
                            <div class="option-item checkbox">
                                <label class="checkbox-label">
                                    <input
                                        type="checkbox"
                                        v-model="
                                            exportSettings.mysql_shell
                                                .skip_binlog
                                        "
                                    />
                                    <span>{{ t("settings.skipBinlog") }}</span>
                                </label>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- MySQLDump 方案配置 -->
                <div class="scheme-config">
                    <h3>{{ t("settings.mysqldump") }}</h3>
                    <div class="scheme-options">
                        <div class="option-row">
                            <div class="option-item">
                                <label>{{ t("settings.compression") }}</label>
                                <select
                                    v-model="
                                        exportSettings.mysql_dump.compression
                                    "
                                    class="setting-control"
                                >
                                    <option value="gzip">gzip</option>
                                    <option value="zstd">zstd</option>
                                    <option value="none">None</option>
                                </select>
                            </div>
                            <div class="option-item checkbox">
                                <label class="checkbox-label">
                                    <input
                                        type="checkbox"
                                        v-model="
                                            exportSettings.mysql_dump.overwrite
                                        "
                                    />
                                    <span>{{ t("settings.overwrite") }}</span>
                                </label>
                            </div>
                            <div class="option-item checkbox">
                                <label class="checkbox-label">
                                    <input
                                        type="checkbox"
                                        v-model="
                                            exportSettings.mysql_dump
                                                .single_transaction
                                        "
                                    />
                                    <span>{{
                                        t("settings.singleTransaction")
                                    }}</span>
                                </label>
                            </div>
                        </div>
                        <div class="option-row">
                            <div class="option-item checkbox">
                                <label class="checkbox-label">
                                    <input
                                        type="checkbox"
                                        v-model="
                                            exportSettings.mysql_dump.routines
                                        "
                                    />
                                    <span>{{ t("settings.routines") }}</span>
                                </label>
                            </div>
                            <div class="option-item checkbox">
                                <label class="checkbox-label">
                                    <input
                                        type="checkbox"
                                        v-model="
                                            exportSettings.mysql_dump.events
                                        "
                                    />
                                    <span>{{ t("settings.events") }}</span>
                                </label>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 关于 -->
            <div class="settings-section">
                <h2>{{ t("settings.about") }}</h2>
                <div class="about-info">
                    <div class="info-item">
                        <span class="info-label">App Name</span>
                        <span class="info-value">GoTools</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">{{
                            t("process.version") || "Version"
                        }}</span>
                        <span class="info-value">{{
                            versionInfo.version
                        }}</span>
                    </div>
                </div>
            </div>

            <!-- 保存按钮 -->
            <div class="settings-actions">
                <button
                    @click="saveSettings"
                    class="save-btn"
                    :disabled="saving"
                >
                    {{ saving ? "Saving..." : t("settings.saveSettings") }}
                </button>
                <button @click="resetSettings" class="reset-btn">
                    {{ t("settings.resetDefault") }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { useAppStore } from "../stores/app";
import {
    GetExportSettings,
    SaveExportSettings,
} from "../../wailsjs/go/main/App";
import { t, setLocale } from "../i18n";

const appStore = useAppStore();

const saving = ref(false);

const settings = ref({
    theme: "dark",
    locale: "zh-CN",
    autoRefreshProcesses: false,
    refreshInterval: 5,
});

const exportSettings = ref({
    export_tool: "mysql-shell",
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

const versionInfo = ref({
    version: "1.0.0",
    build_time: "未知",
    git_commit: "未知",
    go_version: "未知",
    platform: "未知",
});

async function loadExportSettings() {
    try {
        const result = await GetExportSettings();
        if (result) {
            exportSettings.value = result;
        }
    } catch (err) {
        console.error("Failed to load export settings:", err);
    }
}

async function saveExportSettingsToBackend() {
    try {
        await SaveExportSettings(exportSettings.value);
    } catch (err) {
        console.error("Failed to save export settings:", err);
        throw err;
    }
}

async function saveSettings() {
    saving.value = true;
    try {
        // 保存到localStorage
        localStorage.setItem(
            "gotools-settings",
            JSON.stringify(settings.value),
        );

        // 应用设置到store
        appStore.setTheme(settings.value.theme);
        setLocale(settings.value.locale);

        // 保存导出设置到后端
        await saveExportSettingsToBackend();

        alert(t("settings.saved"));
    } catch (err) {
        alert("保存失败: " + err.message);
    } finally {
        saving.value = false;
    }
}

function resetSettings() {
    if (confirm("确定要重置所有设置为默认值吗？")) {
        settings.value = {
            theme: "dark",
            locale: "zh-CN",
            autoRefreshProcesses: false,
            refreshInterval: 5,
        };
        exportSettings.value = {
            export_tool: "mysql-shell",
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
        };
        localStorage.removeItem("gotools-settings");
        saveExportSettingsToBackend();
    }
}

async function loadVersionInfo() {
    try {
        const info = await window.go.main.App.GetVersion();
        versionInfo.value = info;
    } catch (err) {
        console.error("Failed to load version info:", err);
    }
}

onMounted(async () => {
    // 加载保存的设置
    const savedSettings = localStorage.getItem("gotools-settings");
    if (savedSettings) {
        try {
            const parsed = JSON.parse(savedSettings);
            settings.value = { ...settings.value, ...parsed };
        } catch (err) {
            console.error("Failed to load settings:", err);
        }
    }

    // 加载导出设置
    await loadExportSettings();
    await loadVersionInfo();
});
</script>

<style scoped>
.settings {
    padding: 2rem;
    max-width: 900px;
    margin: 0 auto;
}

.header {
    margin-bottom: 2.5rem;
    text-align: left;
}

.header h1 {
    color: var(--accent-color);
    font-size: 2rem;
    margin-bottom: 0.5rem;
}

.subtitle {
    color: var(--text-secondary);
    font-size: 0.95rem;
}

.settings-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

.settings-section {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 10px;
    padding: 1.75rem;
}

.settings-section h2 {
    color: var(--accent-color);
    font-size: 1.25rem;
    margin-bottom: 1.5rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--border-color);
}

.setting-item {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 1.5rem;
    align-items: center;
    padding: 1rem 0;
    border-bottom: 1px solid var(--border-subtle);
}

.setting-item:last-child {
    border-bottom: none;
}

.setting-label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    justify-content: center;
}

.setting-name {
    color: var(--text-tertiary);
    font-weight: 500;
}

.setting-description {
    color: var(--text-secondary);
    font-size: 0.85rem;
    opacity: 0.8;
}

.setting-control {
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    padding: 0.5rem 1rem;
    color: var(--text-tertiary);
    min-width: 160px;
    transition:
        border-color 0.2s,
        box-shadow 0.2s;
}

.setting-control:focus {
    outline: none;
    border-color: var(--accent-color);
    box-shadow: 0 0 0 2px var(--accent-subtle);
}

.setting-control.small {
    width: 100px;
    min-width: 100px;
    text-align: center;
}

/* 方案配置样式 */
.scheme-config {
    margin-top: 2rem;
    padding-top: 1.75rem;
    border-top: 1px solid var(--border-color);
}

.scheme-config h3 {
    color: var(--accent-color);
    font-size: 1.1rem;
    margin-bottom: 1.25rem;
    font-weight: 500;
}

.scheme-options {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
}

.option-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1.25rem;
    align-items: start;
}

.option-item {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    min-width: 0;
}

.option-item label {
    color: var(--text-secondary);
    font-size: 0.85rem;
    font-weight: 500;
}

.option-item.checkbox {
    flex-direction: row;
    align-items: center;
    min-width: auto;
}

.checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-tertiary);
    cursor: pointer;
    white-space: nowrap;
}

/* 开关样式 */
.switch {
    position: relative;
    display: inline-block;
    width: 48px;
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
    background-color: var(--border-color);
    transition: 0.3s;
    border-radius: 26px;
}

.slider:before {
    position: absolute;
    content: "";
    height: 20px;
    width: 20px;
    left: 3px;
    bottom: 3px;
    background-color: var(--text-secondary);
    transition: 0.3s;
    border-radius: 50%;
}

input:checked + .slider {
    background-color: var(--accent-color);
}

input:checked + .slider:before {
    background-color: var(--text-on-accent);
    transform: translateX(22px);
}

/* 关于信息 */
.about-info {
    padding: 1rem 0;
}

.info-item {
    display: grid;
    grid-template-columns: 120px 1fr;
    gap: 1rem;
    padding: 0.75rem 0;
    border-bottom: 1px solid var(--border-subtle);
}

.info-item:last-child {
    border-bottom: none;
}

.info-label {
    color: var(--text-secondary);
    font-weight: 500;
}

.info-value {
    color: var(--text-tertiary);
    font-family: monospace;
}

/* 操作按钮 */
.settings-actions {
    display: flex;
    gap: 1rem;
    margin-top: 2.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
}

.save-btn {
    background: var(--accent-color);
    border: none;
    border-radius: 8px;
    padding: 0.875rem 2.5rem;
    color: var(--text-on-accent);
    font-weight: bold;
    cursor: pointer;
    transition: all 0.3s;
    font-size: 0.95rem;
}

.save-btn:hover:not(:disabled) {
    background: var(--accent-hover);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px var(--accent-subtle);
}

.save-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
}

.reset-btn {
    background: transparent;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 0.875rem 2rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.3s;
    font-size: 0.95rem;
}

.reset-btn:hover {
    border-color: var(--danger-color);
    color: var(--danger-color);
    background: var(--danger-subtle);
}

/* 响应式设计 */
@media (max-width: 768px) {
    .settings {
        padding: 1.5rem;
    }

    .setting-item {
        grid-template-columns: 1fr;
        gap: 0.75rem;
    }

    .setting-control {
        min-width: 100%;
    }

    .option-row {
        grid-template-columns: 1fr;
    }

    .info-item {
        grid-template-columns: 1fr;
        gap: 0.25rem;
    }

    .settings-actions {
        flex-direction: column;
    }

    .save-btn,
    .reset-btn {
        width: 100%;
        text-align: center;
    }
}

@media (max-width: 480px) {
    .settings {
        padding: 1rem;
    }

    .header h1 {
        font-size: 1.5rem;
    }

    .settings-section {
        padding: 1.25rem;
    }
}
</style>
