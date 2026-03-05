<template>
  <div class="login-container">
    <div class="login-header">
      <h1>Key Management</h1>
      <p>Please login to continue</p>
    </div>
    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="username"
          name="username"
          label="Username"
          placeholder="Enter username"
          :rules="[{ required: true, message: 'Username is required' }]"
        />
        <van-field
          v-model="password"
          type="password"
          name="password"
          label="Password"
          placeholder="Enter password"
          :rules="[{ required: true, message: 'Password is required' }]"
        />
        <van-cell center>
          <template #title>
            <van-checkbox v-model="rememberMe">Remember me</van-checkbox>
          </template>
        </van-cell>
      </van-cell-group>
      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          Login
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { showToast } from 'vant';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');
const rememberMe = ref(false);
const loading = ref(false);

onMounted(() => {
  const savedUsername = localStorage.getItem('savedUsername');
  if (savedUsername) {
    username.value = savedUsername;
    rememberMe.value = true;
  }
});

const onSubmit = async () => {
  loading.value = true;
  try {
    await authStore.login(username.value, password.value);

    if (rememberMe.value) {
      localStorage.setItem('savedUsername', username.value);
    } else {
      localStorage.removeItem('savedUsername');
    }

    showToast({ type: 'success', message: 'Login successful' });
    router.push('/key-specs');
  } catch (error) {
    showToast({ type: 'fail', message: 'Login failed' });
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
  font-size: 28px;
  font-weight: bold;
  color: #323233;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: #969799;
}
</style>
