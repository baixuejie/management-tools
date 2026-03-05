<template>
  <div class="config">
    <van-nav-bar title="Configuration" left-arrow @click-left="goBack" />

    <van-cell-group inset style="margin: 16px;">
      <van-cell title="Copy Template Settings" />
    </van-cell-group>

    <van-form @submit="onSubmit">
      <van-cell-group inset style="margin: 16px;">
        <van-field
          v-model="copyTemplate"
          name="copyTemplate"
          label="Template"
          type="textarea"
          placeholder="Enter copy template (use {{key}} as placeholder)"
          rows="4"
          :rules="[{ required: true, message: 'Template is required' }]"
        />
      </van-cell-group>

      <!-- Preview Section -->
      <van-cell-group inset style="margin: 16px;">
        <van-cell title="Preview" />
        <van-cell>
          <div class="preview-box">
            <div class="preview-label">Example output:</div>
            <div class="preview-content">{{ previewText }}</div>
          </div>
        </van-cell>
      </van-cell-group>

      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          Save Template
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { showToast } from 'vant';

const router = useRouter();

const copyTemplate = ref('{{key}}');
const loading = ref(false);

const previewText = computed(() => {
  return copyTemplate.value.replace('{{key}}', 'example-key-12345');
});

const loadTemplate = () => {
  const savedTemplate = localStorage.getItem('copyTemplate');
  if (savedTemplate) {
    copyTemplate.value = savedTemplate;
  }
};

const onSubmit = async () => {
  if (!copyTemplate.value.includes('{{key}}')) {
    showToast({ type: 'fail', message: 'Template must include {{key}} placeholder' });
    return;
  }

  loading.value = true;
  try {
    localStorage.setItem('copyTemplate', copyTemplate.value);
    showToast({ type: 'success', message: 'Template saved successfully' });
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to save template' });
  } finally {
    loading.value = false;
  }
};

const goBack = () => {
  router.back();
};

onMounted(() => {
  loadTemplate();
});
</script>

<style scoped>
.config {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 20px;
}

.preview-box {
  width: 100%;
  padding: 12px;
  background-color: #f7f8fa;
  border-radius: 4px;
}

.preview-label {
  font-size: 12px;
  color: #969799;
  margin-bottom: 8px;
}

.preview-content {
  font-size: 14px;
  color: #323233;
  word-break: break-all;
  white-space: pre-wrap;
  font-family: monospace;
  background-color: #fff;
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #ebedf0;
}
</style>
