<template>
  <div class="profile-page">
    <van-cell-group inset style="margin: 12px 16px;">
      <van-cell title="用户名" :value="auth.user?.username || '-'" />
      <van-cell title="显示名称" :value="auth.user?.display_name || '-'" />
    </van-cell-group>

    <div style="margin: 16px;">
      <van-button block type="danger" @click="logout">退出登录</van-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const auth = useAuthStore();
const router = useRouter();

const logout = () => {
  auth.logout();
  router.push('/login');
};
</script>

<style scoped>
.profile-page {
  min-height: calc(100vh - 96px);
  background: #f7f8fa;
}
</style>
