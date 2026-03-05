<script setup>
import { computed } from 'vue';
import { useAuthStore } from './stores/auth';
import { useRouter, useRoute } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

const showLayout = computed(() => authStore.isAuthenticated && route.path !== '/login');

const activeTab = computed(() => {
  if (route.path.startsWith('/ledger')) return 0;
  if (route.path.startsWith('/keys')) return 1;
  if (route.path.startsWith('/more')) return 2;
  return 0;
});

const onTabChange = (index) => {
  if (index === 0) router.push('/ledger');
  if (index === 1) router.push('/keys');
  if (index === 2) router.push('/more');
};

const pageTitle = computed(() => {
  if (route.path.startsWith('/ledger')) return '记账';
  if (route.path.startsWith('/keys')) return '秘钥';
  if (route.path === '/more') return '更多';
  if (route.path === '/more/costs') return '成本记录';
  if (route.path === '/more/customers') return '购买人管理';
  if (route.path === '/more/specs') return '规格管理';
  if (route.path === '/more/template') return '复制模板';
  if (route.path === '/more/profile') return '个人信息';
  return '管理工具';
});

const showBack = computed(() => route.path.startsWith('/more/') && route.path !== '/more');

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
    <van-nav-bar
      v-if="showLayout"
      :title="pageTitle"
      :left-arrow="showBack"
      @click-left="goBack"
    >
      <template #right>
        <van-icon name="sign" size="18" @click="logout" />
      </template>
    </van-nav-bar>

    <div :class="{ 'page-content': showLayout }">
      <router-view />
    </div>

    <van-tabbar
      v-if="showLayout"
      :model-value="activeTab"
      @change="onTabChange"
    >
      <van-tabbar-item icon="balance-list-o">记账</van-tabbar-item>
      <van-tabbar-item icon="apps-o">秘钥</van-tabbar-item>
      <van-tabbar-item icon="setting-o">更多</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<style>
#app-container {
  min-height: 100vh;
}

.page-content {
  padding-bottom: 66px;
}
</style>
