<template>
  <div class="config">
    <van-nav-bar title="Configuration" />
    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="config.dbPath"
          name="dbPath"
          label="Database Path"
          placeholder="Enter database path"
        />
        <van-field
          v-model="config.jwtSecret"
          name="jwtSecret"
          label="JWT Secret"
          placeholder="Enter JWT secret"
        />
        <van-field
          v-model.number="config.keyExpireDays"
          type="number"
          name="keyExpireDays"
          label="Key Expire Days"
          placeholder="Enter expire days"
        />
      </van-cell-group>
      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          Save Configuration
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { configAPI } from '../api';
import { showToast } from 'vant';

const config = ref({
  dbPath: '',
  jwtSecret: '',
  keyExpireDays: 0,
});
const loading = ref(false);

const loadConfig = async () => {
  try {
    const response = await configAPI.get();
    config.value = response.data;
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to load config' });
  }
};

const onSubmit = async () => {
  loading.value = true;
  try {
    await configAPI.update(config.value);
    showToast({ type: 'success', message: 'Config saved' });
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to save config' });
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadConfig();
});
</script>

<style scoped>
.config {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
