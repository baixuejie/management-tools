<script setup>
import { ref, computed } from 'vue';
import { useAuthStore } from './stores/auth';
import { useRouter, useRoute } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

const showLayout = computed(() => authStore.isAuthenticated && route.path !== '/login');

const activeTab = computed(() => {
  if (route.path.startsWith('/keys') || route.path === '/key-specs' || route.path === '/') return 0;
  if (route.path === '/config') return 1;
  return 0;
});

const onTabChange = (index) => {
  if (index === 0) router.push('/key-specs');
  if (index === 1) router.push('/config');
};

const pageTitle = computed(() => {
  if (route.path === '/key-specs') return '秘钥管理';
  if (route.path.startsWith('/keys/')) return '秘钥列表';
  if (route.path === '/config') return '设置';
  return '工具箱';
});

const showBack = computed(() => {
  return route.path.startsWith('/keys/');
});

const goBack = () => {
  router.back();
};

const logout = () => {
  authStore.logout();
  router.push('/login');
};
</script>

<template>
  <div id="app-container">
    <!-- 顶部 NavBar -->
    <van-nav-bar
      v-if="showLayout"
      :title="pageTitle"
      :left-arrow="showBack"
      @click-left="goBack"
    >
      <template #right>
        <van-icon name="cross" size="18" @click="logout" />
      </template>
    </van-nav-bar>

    <!-- 内容区 -->
    <div :class="{ 'page-content': showLayout }">
      <router-view />
    </div>

    <!-- 底部 Tabbar -->
    <van-tabbar
      v-if="showLayout"
      :model-value="activeTab"
      @change="onTabChange"
    >
      <van-tabbar-item icon="apps-o">工具</van-tabbar-item>
      <van-tabbar-item icon="setting-o">设置</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<style>
#app-container {
  min-height: 100vh;
}

.page-content {
  padding-bottom: 66px; /* tabbar 高度 + 留白 */
}
</style>
