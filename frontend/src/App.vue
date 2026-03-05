<script setup>
import { useAuthStore } from './stores/auth';
import { useRouter, useRoute } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

const onTabChange = (index) => {
  const routes = ['/key-specs', '/config'];
  router.push(routes[index]);
};

const logout = () => {
  authStore.logout();
  router.push('/login');
};
</script>

<template>
  <div id="app">
    <router-view />
    <van-tabbar v-if="authStore.isAuthenticated && route.path !== '/login'" @change="onTabChange">
      <van-tabbar-item icon="apps-o">Key Specs</van-tabbar-item>
      <van-tabbar-item icon="setting-o">Config</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<style>
#app {
  min-height: 100vh;
}
</style>
