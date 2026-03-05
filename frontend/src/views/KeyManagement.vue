<template>
  <div class="key-management">
    <van-cell-group inset style="margin: 12px 16px;">
      <van-cell title="当前规格" :value="currentSpecName" is-link @click="showSpecPicker = true" />
      <van-cell center title="仅显示未使用">
        <template #right-icon>
          <van-switch v-model="onlyUnused" size="20" @change="reloadKeys" />
        </template>
      </van-cell>
      <van-cell title="批量上传" icon="upgrade" is-link @click="showUpload = true" />
    </van-cell-group>

    <div class="quick-copy-box">
      <van-button
        block
        type="primary"
        icon="records-o"
        @click="getAndCopyKey"
        :disabled="!selectedSpecId || total === 0"
      >
        获取并复制可用秘钥
      </van-button>
    </div>

    <van-cell-group inset style="margin: 0 16px 12px;">
      <van-cell
        title="查看全部秘钥"
        :value="showKeyList ? '已展开' : '已折叠'"
        is-link
        @click="showKeyList = !showKeyList"
      />
    </van-cell-group>

    <div v-if="showKeyList">
      <div class="stats-bar">共 {{ total }} 个秘钥</div>
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-empty v-if="!loading && keys.length === 0" description="暂无秘钥" />
        <van-list
          v-else
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="onLoad"
        >
          <van-swipe-cell v-for="key in keys" :key="key.id">
            <van-cell :label="formatTime(key)">
              <template #title>
                <div class="key-cell-title">
                  <span class="key-value">{{ maskKey(key.key_value) }}</span>
                  <van-tag :type="key.is_used ? 'default' : 'success'" size="medium">
                    {{ key.is_used ? '已使用' : '未使用' }}
                  </van-tag>
                </div>
              </template>
              <template #right-icon>
                <van-icon name="records-o" size="20" color="#1989fa" @click.stop="copyKey(key)" />
              </template>
            </van-cell>
            <template #right>
              <van-button square type="danger" text="删除" class="swipe-btn" @click="confirmDelete(key.id)" />
            </template>
          </van-swipe-cell>
        </van-list>
      </van-pull-refresh>
    </div>

    <van-popup v-model:show="showUpload" position="bottom" round :style="{ minHeight: '45%' }">
      <div class="popup-header">
        <h3>批量上传秘钥</h3>
      </div>
      <van-form @submit="batchUpload">
        <van-cell-group inset>
          <van-field
            v-model="uploadText"
            type="textarea"
            placeholder="粘贴秘钥，每行一个"
            rows="6"
            :rules="[{ required: true, message: '至少输入一个秘钥' }]"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit" :loading="uploading">
            上传
          </van-button>
        </div>
      </van-form>
    </van-popup>

    <van-popup v-model:show="showSpecPicker" position="bottom" round>
      <van-picker
        title="选择规格"
        :columns="specColumns"
        @confirm="onSpecConfirm"
        @cancel="showSpecPicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { keyAPI, keySpecAPI, configAPI } from '../api';
import { useAuthStore } from '../stores/auth';
import { showToast, showConfirmDialog } from 'vant';

const route = useRoute();
const authStore = useAuthStore();

const specOptions = ref([]);
const selectedSpecId = ref(null);
const showSpecPicker = ref(false);

const keys = ref([]);
const total = ref(0);
const loading = ref(false);
const finished = ref(false);
const refreshing = ref(false);
const onlyUnused = ref(false);
const offset = ref(0);
const limit = 20;
const showKeyList = ref(true);

const showUpload = ref(false);
const uploadText = ref('');
const uploading = ref(false);

const lastSpecStorageKey = computed(() => {
  const userId = authStore.user?.id || authStore.user?.username || 'guest';
  return `last_selected_spec_${userId}`;
});

const currentSpecName = computed(() => {
  const current = specOptions.value.find((item) => item.value === selectedSpecId.value);
  return current ? current.text : '请选择规格';
});

const specColumns = computed(() => specOptions.value);

const maskKey = (value) => {
  if (!value) return '';
  if (value.length <= 8) return value;
  return `${value.slice(0, 4)}***${value.slice(-4)}`;
};

const formatTime = (key) => {
  if (key.is_used && key.used_at) {
    return `使用于 ${new Date(key.used_at).toLocaleString()}`;
  }
  return `创建于 ${new Date(key.created_at).toLocaleString()}`;
};

const loadSpecs = async () => {
  try {
    const response = await keySpecAPI.list();
    specOptions.value = (response.data || []).map((item) => ({ text: item.name, value: item.id }));

    const querySpecId = Number(route.query.specId || 0);
    const cachedSpecId = Number(localStorage.getItem(lastSpecStorageKey.value) || 0);

    if (querySpecId && specOptions.value.some((item) => item.value === querySpecId)) {
      selectedSpecId.value = querySpecId;
    } else if (cachedSpecId && specOptions.value.some((item) => item.value === cachedSpecId)) {
      selectedSpecId.value = cachedSpecId;
    } else if (!selectedSpecId.value && specOptions.value.length > 0) {
      selectedSpecId.value = specOptions.value[0].value;
    }
  } catch {
    showToast({ type: 'fail', message: '加载规格失败' });
  }
};

const loadKeys = async (append = false) => {
  if (!selectedSpecId.value) {
    keys.value = [];
    total.value = 0;
    finished.value = true;
    return;
  }

  try {
    const response = await keyAPI.list(selectedSpecId.value, {
      onlyUnused: onlyUnused.value,
      limit,
      offset: offset.value,
    });
    const data = response.data;
    if (append) {
      keys.value = [...keys.value, ...(data.keys || [])];
    } else {
      keys.value = data.keys || [];
    }
    total.value = data.total || 0;
    finished.value = keys.value.length >= total.value;
  } catch {
    showToast({ type: 'fail', message: '加载秘钥失败' });
  }
};

const onLoad = async () => {
  if (keys.value.length === 0) {
    await loadKeys();
  } else {
    offset.value = keys.value.length;
    await loadKeys(true);
  }
  loading.value = false;
};

const onRefresh = async () => {
  offset.value = 0;
  await loadKeys();
  refreshing.value = false;
};

const reloadKeys = () => {
  offset.value = 0;
  finished.value = false;
  loadKeys();
};

const onSpecConfirm = ({ selectedOptions }) => {
  if (selectedOptions?.length) {
    selectedSpecId.value = Number(selectedOptions[0].value);
  }
  showSpecPicker.value = false;
};

const getTemplate = async () => {
  try {
    const response = await configAPI.getTemplate();
    return response.data.template || '{{key}}';
  } catch {
    return '{{key}}';
  }
};

const getAndCopyKey = async () => {
  if (!selectedSpecId.value) {
    showToast({ type: 'fail', message: '请先选择规格' });
    return;
  }

  try {
    const response = await keyAPI.getAvailable(selectedSpecId.value);
    const key = response.data;
    const template = await getTemplate();

    await navigator.clipboard.writeText(template.replace('{{key}}', key.key_value));
    await keyAPI.markUsed(key.id);

    showToast({ type: 'success', message: '已复制到剪贴板' });
    reloadKeys();
  } catch (error) {
    const msg = error.response?.data?.error || '没有可用秘钥';
    showToast({ type: 'fail', message: msg });
  }
};

const copyKey = async (key) => {
  try {
    const template = await getTemplate();
    await navigator.clipboard.writeText(template.replace('{{key}}', key.key_value));

    if (!key.is_used) {
      await keyAPI.markUsed(key.id);
      key.is_used = true;
    }

    showToast({ type: 'success', message: '已复制到剪贴板' });
  } catch {
    showToast({ type: 'fail', message: '复制失败' });
  }
};

const batchUpload = async () => {
  if (!selectedSpecId.value) {
    showToast({ type: 'fail', message: '请先选择规格' });
    return;
  }

  const lines = uploadText.value.split('\n').map((line) => line.trim()).filter(Boolean);
  if (lines.length === 0) {
    showToast({ type: 'fail', message: '请至少输入一个秘钥' });
    return;
  }

  uploading.value = true;
  try {
    await keyAPI.batchUpload(selectedSpecId.value, lines);
    showToast({ type: 'success', message: `成功上传 ${lines.length} 个秘钥` });
    uploadText.value = '';
    showUpload.value = false;
    reloadKeys();
  } catch (error) {
    const msg = error.response?.data?.error || '上传失败';
    showToast({ type: 'fail', message: msg });
  } finally {
    uploading.value = false;
  }
};

const confirmDelete = async (id) => {
  try {
    await showConfirmDialog({ title: '确认', message: '确定删除这个秘钥吗？' });
    await keyAPI.delete(id);
    showToast({ type: 'success', message: '删除成功' });
    reloadKeys();
  } catch {
    // 用户取消
  }
};

watch(selectedSpecId, (value) => {
  if (value) {
    localStorage.setItem(lastSpecStorageKey.value, String(value));
  }
  reloadKeys();
});

onMounted(async () => {
  await loadSpecs();
  await loadKeys();
});
</script>

<style scoped>
.key-management {
  min-height: calc(100vh - 96px);
  background-color: #f7f8fa;
}

.quick-copy-box {
  margin: 0 16px 12px;
}

.stats-bar {
  padding: 8px 16px;
  font-size: 13px;
  color: #969799;
}

.key-cell-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.key-value {
  font-family: 'Courier New', monospace;
  font-size: 14px;
  word-break: break-all;
}

.swipe-btn {
  height: 100%;
}

.popup-header {
  padding: 20px 16px 10px;
  text-align: center;
  border-bottom: 1px solid #ebedf0;
}

.popup-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: #323233;
}
</style>
