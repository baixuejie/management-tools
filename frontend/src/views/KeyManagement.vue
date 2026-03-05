<template>
  <div class="key-management">
    <van-nav-bar title="Key Management" left-arrow @click-left="goBack" />
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="No more data"
        @load="onLoad"
      >
        <van-cell
          v-for="key in keys"
          :key="key.id"
          :title="key.key_name"
          :label="`Created: ${new Date(key.created_at).toLocaleString()}`"
          is-link
          @click="viewKey(key.id)"
        />
      </van-list>
    </van-pull-refresh>
    <van-floating-bubble
      icon="plus"
      @click="showAddDialog = true"
    />
    <van-dialog
      v-model:show="showAddDialog"
      title="Generate Key"
      show-cancel-button
      @confirm="generateKey"
    >
      <van-form>
        <van-field v-model="newKey.keyName" label="Key Name" placeholder="Enter key name" />
      </van-form>
    </van-dialog>
    <van-dialog
      v-model:show="showKeyDialog"
      title="Key Details"
      :message="keyDetails"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { keyAPI } from '../api';
import { showToast } from 'vant';

const router = useRouter();
const route = useRoute();

const keys = ref([]);
const loading = ref(false);
const finished = ref(false);
const refreshing = ref(false);
const showAddDialog = ref(false);
const showKeyDialog = ref(false);
const keyDetails = ref('');
const newKey = ref({ keyName: '' });

const specId = route.params.specId;

const loadKeys = async () => {
  try {
    const response = await keyAPI.list(specId);
    keys.value = response.data;
    finished.value = true;
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to load keys' });
  }
};

const onLoad = async () => {
  await loadKeys();
  loading.value = false;
};

const onRefresh = async () => {
  await loadKeys();
  refreshing.value = false;
};

const goBack = () => {
  router.back();
};

const generateKey = async () => {
  try {
    await keyAPI.create({ spec_id: parseInt(specId), key_name: newKey.value.keyName });
    showToast({ type: 'success', message: 'Key generated' });
    newKey.value = { keyName: '' };
    await loadKeys();
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to generate key' });
  }
};

const viewKey = async (keyId) => {
  try {
    const response = await keyAPI.get(keyId);
    keyDetails.value = `Key: ${response.data.key_value}`;
    showKeyDialog.value = true;
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to load key details' });
  }
};

onMounted(() => {
  loadKeys();
});
</script>

<style scoped>
.key-management {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
