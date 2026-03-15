<script setup>
import { ref, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { FaCog, FaDatabase, FaMicrochip, FaSave, FaHome, FaCogs } from "vue-icons-plus/fa";
import { useAppStore } from "../stores/app";

const router = useRouter();
const route = useRoute();
const appStore = useAppStore();

const menuItems = [
    {
        category: "工具",
        icon: FaCog,
        items: [
            { id: "home", name: "首页", icon: FaHome, route: "/" },
            { id: "process-manager", name: "进程管理器", icon: FaMicrochip, route: "/process-manager" },
        ],
    },
    {
        category: "数据库",
        icon: FaDatabase,
        items: [{ id: "database-backup", name: "导入/导出", icon: FaSave, route: "/database-backup" }],
    },
    {
        category: "设置",
        icon: FaCogs,
        items: [{ id: "settings", name: "应用设置", icon: FaCog, route: "/settings" }],
    },
];

const expandedCategories = ref(["工具", "数据库", "设置"]);
const isCollapsed = ref(appStore.sidebarCollapsed);

const currentRoute = computed(() => route.path);

function toggleCategory(category) {
    const index = expandedCategories.value.indexOf(category);
    if (index > -1) {
        expandedCategories.value.splice(index, 1);
    } else {
        expandedCategories.value.push(category);
    }
}

function isExpanded(category) {
    return expandedCategories.value.includes(category);
}

function navigateTo(item) {
    router.push(item.route);
}

function toggleSidebar() {
    isCollapsed.value = !isCollapsed.value;
    appStore.toggleSidebar();
}

function isActive(item) {
    return currentRoute.value === item.route;
}
</script>

<template>
    <div class="sidebar" :class="{ collapsed: isCollapsed }">
        <div class="sidebar-header">
            <div class="logo" v-if="!isCollapsed">
                <span class="logo-icon">◈</span>
                <span class="logo-text">GOTOOLS</span>
            </div>
            <div class="logo logo-collapsed" v-else>
                <span class="logo-icon">◈</span>
            </div>
            <div class="version" v-if="!isCollapsed">v1.0.0</div>
        </div>

        <div class="sidebar-menu">
            <div
                v-for="group in menuItems"
                :key="group.category"
                class="menu-group"
            >
                <div
                    class="menu-category"
                    @click="toggleCategory(group.category)"
                >
                    <span class="category-icon">
                        <component :is="group.icon" />
                    </span>
                    <span class="category-name" v-if="!isCollapsed">{{
                        group.category
                    }}</span>
                    <span
                        class="expand-icon"
                        :class="{ expanded: isExpanded(group.category) }"
                        v-if="!isCollapsed"
                        >▶</span
                    >
                </div>

                <div v-show="isExpanded(group.category)" class="menu-items">
                    <div
                        v-for="item in group.items"
                        :key="item.id"
                        class="menu-item"
                        :class="{ active: isActive(item) }"
                        @click="navigateTo(item)"
                    >
                        <span class="item-icon">
                            <component :is="item.icon" />
                        </span>
                        <span class="item-name" v-if="!isCollapsed">{{
                            item.name
                        }}</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="sidebar-footer">
            <div class="footer-line" v-if="!isCollapsed"></div>
            <div class="footer-text" v-if="!isCollapsed">SYSTEM READY</div>
        </div>

        <button class="toggle-button" @click="toggleSidebar">
            <span :class="{ rotated: isCollapsed }">◀</span>
        </button>
    </div>
</template>

<style scoped>
.sidebar {
    width: 220px;
    min-width: 220px;
    background: #0d0d0d;
    border-right: 1px solid #1a1a1a;
    display: flex;
    flex-direction: column;
    height: 100vh;
    transition:
        width 0.3s ease,
        min-width 0.3s ease;
    position: relative;
}

.sidebar.collapsed {
    width: 60px;
    min-width: 60px;
}

.sidebar-header {
    padding: 20px 15px;
    border-bottom: 1px solid #1a1a1a;
    text-align: center;
}

.logo {
    display: flex;
    align-items: center;
    gap: 10px;
    justify-content: center;
}

.logo-collapsed {
    justify-content: center;
}

.logo-icon {
    color: #00ff00;
    font-size: 20px;
    text-shadow: 0 0 10px #00ff00;
}

.logo-text {
    color: #00ff00;
    font-size: 16px;
    font-weight: bold;
    letter-spacing: 2px;
    text-shadow: 0 0 10px #00ff00;
}

.version {
    color: #333;
    font-size: 10px;
    margin-top: 5px;
    letter-spacing: 1px;
}

.sidebar-menu {
    flex: 1;
    overflow-y: auto;
    padding: 10px 0;
}

.menu-group {
    margin-bottom: 5px;
}

.menu-category {
    display: flex;
    align-items: center;
    padding: 10px 15px;
    cursor: pointer;
    color: #666;
    font-size: 12px;
    letter-spacing: 1px;
    transition: all 0.2s;
    justify-content: center;
}

.sidebar.collapsed .menu-category {
    padding: 10px 5px;
}

.menu-category:hover {
    color: #00ff00;
    background: rgba(0, 255, 0, 0.05);
}

.category-icon {
    margin-right: 8px;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
}

.category-icon svg {
    width: 14px;
    height: 14px;
}

.sidebar.collapsed .category-icon {
    margin-right: 0;
}

.category-name {
    flex: 1;
    text-align: left;
}

.expand-icon {
    font-size: 8px;
    transition: transform 0.2s;
}

.expand-icon.expanded {
    transform: rotate(90deg);
}

.menu-items {
    padding-left: 10px;
}

.sidebar.collapsed .menu-items {
    padding-left: 0;
}

.menu-item {
    display: flex;
    align-items: center;
    padding: 8px 15px 8px 25px;
    cursor: pointer;
    color: #444;
    font-size: 12px;
    letter-spacing: 1px;
    transition: all 0.2s;
    border-left: 2px solid transparent;
    justify-content: center;
}

.sidebar.collapsed .menu-item {
    padding: 8px 5px;
    border-left: none;
}

.menu-item:hover {
    color: #00ff00;
    background: rgba(0, 255, 0, 0.05);
}

.menu-item.active {
    color: #00ff00;
    background: rgba(0, 255, 0, 0.1);
    border-left-color: #00ff00;
}

.sidebar.collapsed .menu-item.active {
    border-left: none;
    background: rgba(0, 255, 0, 0.1);
}

.item-icon {
    margin-right: 8px;
    font-size: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 16px;
}

.item-icon svg {
    width: 12px;
    height: 12px;
}

.sidebar.collapsed .item-icon {
    margin-right: 0;
}

.item-name {
    flex: 1;
    text-align: left;
}

.sidebar-footer {
    padding: 15px;
    border-top: 1px solid #1a1a1a;
    text-align: center;
}

.footer-line {
    height: 1px;
    background: linear-gradient(90deg, transparent, #00ff00, transparent);
    margin-bottom: 10px;
}

.footer-text {
    color: #333;
    font-size: 10px;
    text-align: center;
    letter-spacing: 2px;
}

.toggle-button {
    position: absolute;
    top: 50%;
    right: -10px;
    transform: translateY(-50%);
    width: 20px;
    height: 40px;
    background: #0d0d0d;
    border: 1px solid #1a1a1a;
    border-left: none;
    border-radius: 0 5px 5px 0;
    color: #666;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    z-index: 10;
}

.toggle-button:hover {
    color: #00ff00;
    background: rgba(0, 255, 0, 0.05);
}

.toggle-button span {
    transition: transform 0.3s ease;
}

.toggle-button span.rotated {
    transform: rotate(180deg);
}

::-webkit-scrollbar {
    width: 4px;
}

::-webkit-scrollbar-track {
    background: #0d0d0d;
}

::-webkit-scrollbar-thumb {
    background: #1a1a1a;
    border-radius: 2px;
}
</style>
