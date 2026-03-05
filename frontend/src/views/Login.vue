<template>
  <div class="login-container">
    <div class="login-header">
      <h1>账本与秘钥管理</h1>
      <p>请登录后继续</p>
    </div>

    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="username"
          is-link
          readonly
          name="username"
          label="用户"
          placeholder="请选择用户"
          @click="showUserPicker = true"
        />
        <van-field
          v-model="password"
          type="password"
          name="password"
          label="密码"
          placeholder="请输入密码"
          :rules="[{ required: true, message: '请输入密码' }]"
        />
      </van-cell-group>

      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          登录
        </van-button>
      </div>
    </van-form>

    <van-popup v-model:show="showUserPicker" position="bottom" round>
      <van-picker
        title="选择用户"
        :columns="userOptions"
        @confirm="onUserConfirm"
        @cancel="showUserPicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { showToast } from 'vant';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('admin');
const password = ref('');
const loading = ref(false);
const showUserPicker = ref(false);
const userOptions = [
  { text: 'admin（白了个白）', value: 'admin' },
  { text: 'fanchen（凡尘）', value: 'fanchen' },
];

const onUserConfirm = ({ selectedOptions }) => {
  if (selectedOptions?.length) {
    username.value = selectedOptions[0].value;
  }
  showUserPicker.value = false;
};

const onSubmit = async () => {
  loading.value = true;
  try {
    await authStore.login(username.value, password.value, false);
    showToast({ type: 'success', message: '登录成功' });
    router.push('/ledger');
  } catch (error) {
    const msg = error.response?.data?.error || '用户名或密码错误';
    showToast({ type: 'fail', message: msg });
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  if (authStore.isAuthenticated) {
    router.replace('/ledger');
  }
});
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px 16px;
  min-height: 100vh;
  background: linear-gradient(180deg, #f2f8ff 0%, #f7f8fa 100%);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h1 {
  font-size: 24px;
  font-weight: 700;
  color: #323233;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: #969799;
}
</style>
