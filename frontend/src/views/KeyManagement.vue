<template>
  <div class="key-management">
    <van-nav-bar :title="specName" left-arrow @click-left="goBack" />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- Upload Section -->
      <van-collapse v-model="activeCollapse" style="margin: 16px;">
        <van-collapse-item title="Batch Upload Keys" name="1">
          <van-form @submit="batchUpload">
            <van-field
              v-model="uploadText"
              type="textarea"
              placeholder="Enter keys (one per line)"
              rows="5"
              :rules="[{ required: true, message: 'Please enter at least one key' }]"
            />
            <div style="margin-top: 12px;">
              <van-button round block type="primary" native-type="submit" :loading="uploading">
                Upload Keys
              </van-button>
            </div>
          </van-form>
        </van-collapse-item>
      </van-collapse>

      <!-- Filter and Actions -->
      <van-cell-group inset style="margin: 16px;">
        <van-cell center>
          <template #title>
            <span>Show only unused keys</span>
          </template>
          <template #right-icon>
            <van-switch v-model="showOnlyUnused" size="20" @change="filterKeys" />
          </template>
        </van-cell>
        <van-cell>
          <van-button
            type="success"
            size="small"
            block
            @click="getAvailableKey"
            :disabled="availableCount === 0"
          >
            Get Available Key ({{ availableCount }} unused)
          </van-button>
        </van-cell>
      </van-cell-group>

      <!-- Keys List -->
      <van-empty v-if="!loading && filteredKeys.length === 0" description="No keys found" />
      <van-list
        v-else
        v-model:loading="loading"
        :finished="finished"
        finished-text="No more data"
        @load="onLoad"
      >
        <van-swipe-cell v-for="key in filteredKeys" :key="key.id">
          <van-cell :title="key.key_value" :label="`Created: ${formatDate(key.created_at)}`">
            <template #right-icon>
              <van-tag :type="key.is_used ? 'default' : 'success'" style="margin-right: 8px;">
                {{ key.is_used ? 'Used' : 'Unused' }}
              </van-tag>
              <van-button
                type="primary"
                size="small"
                @click="copyKey(key)"
              >
                Copy
              </van-button>
            </template>
          </van-cell>
          <template #right>
            <van-button square type="danger" text="Delete" @click="confirmDelete(key.id)" />
          </template>
        </van-swipe-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { keyAPI, keySpecAPI } from '../api';
import { showToast, showConfirmDialog } from 'vant';

const router = useRouter();
const route = useRoute();

const specId = route.params.specId;
const specName = ref('Key Management');
const keys = ref([]);
const loading = ref(false);
const finished = ref(false);
const refreshing = ref(false);
const uploadText = ref('');
const uploading = ref(false);
const activeCollapse = ref([]);
const showOnlyUnused = ref(false);
const copyTemplate = ref('{{key}}');

const filteredKeys = computed(() => {
  if (showOnlyUnused.value) {
    return keys.value.filter(key => !key.is_used);
  }
  return keys.value;
});

const availableCount = computed(() => {
  return keys.value.filter(key => !key.is_used).length;
});

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

const loadSpec = async () => {
  try {
    const response = await keySpecAPI.list();
    const spec = response.data?.find(s => s.id === parseInt(specId));
    if (spec) {
      specName.value = spec.name;
    }
  } catch (error) {
    console.error('Failed to load spec:', error);
  }
};

const loadKeys = async () => {
  try {
    const response = await keyAPI.list(specId);
    keys.value = response.data || [];
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

const filterKeys = () => {
  // Trigger re-render by updating computed
};

const goBack = () => {
  router.back();
};

const batchUpload = async () => {
  if (!uploadText.value.trim()) {
    showToast({ type: 'fail', message: 'Please enter at least one key' });
    return;
  }

  uploading.value = true;
  const keyLines = uploadText.value.split('\n').filter(line => line.trim());
  let successCount = 0;
  let failCount = 0;

  for (const keyValue of keyLines) {
    try {
      await keyAPI.create({
        spec_id: parseInt(specId),
        key_value: keyValue.trim(),
      });
      successCount++;
    } catch (error) {
      failCount++;
    }
  }

  uploading.value = false;
  uploadText.value = '';
  activeCollapse.value = [];

  if (failCount === 0) {
    showToast({ type: 'success', message: `Successfully uploaded ${successCount} keys` });
  } else {
    showToast({ type: 'warning', message: `Uploaded ${successCount} keys, ${failCount} failed` });
  }

  await loadKeys();
};

const getAvailableKey = async () => {
  const unusedKey = keys.value.find(key => !key.is_used);
  if (!unusedKey) {
    showToast({ type: 'fail', message: 'No available keys' });
    return;
  }

  await copyKey(unusedKey);
};

const copyKey = async (key) => {
  try {
    // Load template from localStorage
    const savedTemplate = localStorage.getItem('copyTemplate') || '{{key}}';
    const textToCopy = savedTemplate.replace('{{key}}', key.key_value);

    // Copy to clipboard
    await navigator.clipboard.writeText(textToCopy);
    showToast({ type: 'success', message: 'Key copied to clipboard' });

    // Mark as used if not already
    if (!key.is_used) {
      try {
        await keyAPI.update(key.id, { is_used: true });
        key.is_used = true;
      } catch (error) {
        console.error('Failed to mark key as used:', error);
      }
    }
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to copy key' });
  }
};

const confirmDelete = async (id) => {
  try {
    await showConfirmDialog({
      title: 'Confirm Delete',
      message: 'Are you sure you want to delete this key?',
    });
    await deleteKey(id);
  } catch {
    // User cancelled
  }
};

const deleteKey = async (id) => {
  try {
    await keyAPI.delete(id);
    showToast({ type: 'success', message: 'Key deleted' });
    await loadKeys();
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to delete key' });
  }
};

onMounted(() => {
  loadSpec();
  loadKeys();
});
</script>

<style scoped>
.key-management {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 20px;
}
</style>
