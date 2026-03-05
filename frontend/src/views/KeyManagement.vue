<template>
  <div class="key-management">
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- 操作按钮区 -->
      <van-cell-group inset style="margin: 12px 16px;">
        <van-cell center title="批量上传" icon="upgrade" is-link @click="showUpload = true" />
        <van-cell center>
          <template #title>
            <span>仅显示未使用</span>
          </template>
          <template #right-icon>
            <van-switch v-model="onlyUnused" size="20" @change="reloadKeys" />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 统计信息 -->
      <div class="stats-bar">
        <span>共 {{ total }} 个秘钥</span>
        <van-button
          type="primary"
          size="small"
          icon="records-o"
          @click="getAndCopyKey"
          :disabled="total === 0"
        >
          获取并复制
        </van-button>
      </div>

      <!-- 秘钥列表 -->
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

    <!-- 批量上传弹窗 -->
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
            :rules="[{ required: true, message: '请输入至少一个秘钥' }]"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit" :loading="uploading">
            上传
          </van-button>
        </div>
      </van-form>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { keyAPI, configAPI } from '../api';
import { showToast, showConfirmDialog } from 'vant';

const route = useRoute();
const specId = Number(route.params.specId);

const keys = ref([]);
const total = ref(0);
const loading = ref(false);
const finished = ref(false);
const refreshing = ref(false);
const onlyUnused = ref(false);
const offset = ref(0);
const limit = 20;

const showUpload = ref(false);
const uploadText = ref('');
const uploading = ref(false);

const maskKey = (val) => {
  if (!val) return '';
  if (val.length <= 8) return val;
  return val.slice(0, 4) + '***' + val.slice(-4);
};

const formatTime = (key) => {
  if (key.is_used && key.used_at) {
    return '使用于: ' + new Date(key.used_at).toLocaleString();
  }
  return '创建于: ' + new Date(key.created_at).toLocaleString();
};

const loadKeys = async (append = false) => {
  try {
    const response = await keyAPI.list(specId, {
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
  } catch (error) {
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

const getAndCopyKey = async () => {
  try {
    const response = await keyAPI.getAvailable(specId);
    const key = response.data;

    // 获取模板
    let template = '{{key}}';
    try {
      const tplRes = await configAPI.getTemplate();
      template = tplRes.data.template || '{{key}}';
    } catch {
      // 使用默认模板
    }

    const text = template.replace('{{key}}', key.key_value);
    await navigator.clipboard.writeText(text);

    // 标记为已使用
    await keyAPI.markUsed(key.id);

    showToast({ type: 'success', message: '已复制到剪贴板' });
    reloadKeys();
  } catch (error) {
    const msg = error.response?.data?.error || '没有可用的秘钥';
    showToast({ type: 'fail', message: msg });
  }
};

const copyKey = async (key) => {
  try {
    let template = '{{key}}';
    try {
      const tplRes = await configAPI.getTemplate();
      template = tplRes.data.template || '{{key}}';
    } catch {
      // 使用默认模板
    }

    const text = template.replace('{{key}}', key.key_value);
    await navigator.clipboard.writeText(text);

    if (!key.is_used) {
      await keyAPI.markUsed(key.id);
      key.is_used = true;
    }

    showToast({ type: 'success', message: '已复制到剪贴板' });
  } catch (error) {
    showToast({ type: 'fail', message: '复制失败' });
  }
};

const batchUpload = async () => {
  const lines = uploadText.value.split('\n').map(l => l.trim()).filter(l => l);
  if (lines.length === 0) {
    showToast({ type: 'fail', message: '请输入至少一个秘钥' });
    return;
  }

  uploading.value = true;
  try {
    await keyAPI.batchUpload(specId, lines);
    showToast({ type: 'success', message: `成功上传 ${lines.length} 个秘钥` });
    uploadText.value = '';
    showUpload.value = false;
    reloadKeys();
  } catch (error) {
    showToast({ type: 'fail', message: '上传失败' });
  } finally {
    uploading.value = false;
  }
};

const confirmDelete = async (id) => {
  try {
    await showConfirmDialog({ title: '确认删除', message: '确定要删除这个秘钥吗？' });
    await keyAPI.delete(id);
    showToast({ type: 'success', message: '删除成功' });
    reloadKeys();
  } catch {
    // 用户取消
  }
};

onMounted(() => {
  loadKeys();
});
</script>

<style scoped>
.key-management {
  min-height: calc(100vh - 96px);
  background-color: #f7f8fa;
}

.stats-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
