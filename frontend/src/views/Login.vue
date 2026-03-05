<template>
  <div class="login-container">
    <div class="login-header">
      <h1>秘钥管理工具</h1>
      <p>请登录后使用</p>
    </div>
    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="username"
          name="username"
          label="用户名"
          placeholder="请输入用户名"
          :rules="[{ required: true, message: '请输入用户名' }]"
        />
        <van-field
          v-model="password"
          type="password"
          name="password"
          label="密码"
          placeholder="请输入密码"
          :rules="[{ required: true, message: '请输入密码' }]"
        />
        <van-cell center>
          <template #title>
            <van-checkbox v-model="rememberMe">记住我</van-checkbox>
          </template>
        </van-cell>
      </van-cell-group>
      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          登录
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { showToast } from 'vant';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');
const rememberMe = ref(false);
const loading = ref(false);

const onSubmit = async () => {
  loading.value = true;
  try {
    await authStore.login(username.value, password.value, rememberMe.value);
    showToast({ type: 'success', message: '登录成功' });
    router.push('/key-specs');
  } catch (error) {
    showToast({ type: 'fail', message: '用户名或密码错误' });
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px 16px;
  min-height: 100vh;
  background-color: #f7f8fa;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h1 {
  font-size: 24px;
  font-weight: bold;
  color: #323233;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: #969799;
}
</style>
