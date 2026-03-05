<template>
  <div class="config-page">
    <van-cell-group inset style="margin: 12px 16px;">
      <van-cell title="复制模板" label="使用 {{key}} 作为秘钥占位符" />
    </van-cell-group>

    <van-form @submit="onSubmit">
      <van-cell-group inset style="margin: 12px 16px;">
        <van-field
          v-model="copyTemplate"
          name="copyTemplate"
          label="模板"
          type="textarea"
          placeholder="例如：API_KEY={{key}}"
          rows="4"
          :rules="[{ required: true, message: '请输入模板' }]"
        />
      </van-cell-group>

      <van-cell-group inset style="margin: 12px 16px;">
        <van-cell title="预览" />
        <van-cell>
          <div class="preview-box">
            <code>{{ previewText }}</code>
          </div>
        </van-cell>
      </van-cell-group>

      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          保存模板
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { configAPI } from '../api';
import { showToast } from 'vant';

const copyTemplate = ref('{{key}}');
const loading = ref(false);

const previewText = computed(() => copyTemplate.value.replace('{{key}}', 'sk-example-12345abcde'));

const loadTemplate = async () => {
  try {
    const response = await configAPI.getTemplate();
    if (response.data.template) {
      copyTemplate.value = response.data.template;
    }
  } catch {
    // 使用默认值
  }
};

const onSubmit = async () => {
  if (!copyTemplate.value.includes('{{key}}')) {
    showToast({ type: 'fail', message: '模板必须包含 {{key}}' });
    return;
  }

  loading.value = true;
  try {
    await configAPI.updateTemplate(copyTemplate.value);
    showToast({ type: 'success', message: '模板保存成功' });
  } catch (error) {
    const msg = error.response?.data?.error || '保存失败';
    showToast({ type: 'fail', message: msg });
  } finally {
    loading.value = false;
  }
};

onMounted(loadTemplate);
</script>

<style scoped>
.config-page {
  min-height: calc(100vh - 96px);
  background-color: #f7f8fa;
}

.preview-box {
  width: 100%;
  padding: 12px;
  background-color: #f0f2f5;
  border-radius: 4px;
  font-size: 13px;
  word-break: break-all;
  white-space: pre-wrap;
}

.preview-box code {
  font-family: 'Courier New', monospace;
  color: #323233;
}
</style>
