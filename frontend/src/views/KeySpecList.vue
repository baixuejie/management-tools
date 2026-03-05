<template>
  <div class="key-spec-list">
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-empty v-if="!loading && keySpecs.length === 0" description="暂无秘钥规格" />
      <van-list
        v-else
        v-model:loading="loading"
        :finished="finished"
        finished-text=""
        @load="onLoad"
      >
        <van-swipe-cell v-for="spec in keySpecs" :key="spec.id">
          <van-cell
            :title="spec.name"
            :label="spec.description || '无描述'"
            is-link
            @click="goToKeys(spec.id)"
          />
          <template #right>
            <van-button square type="primary" text="编辑" class="swipe-btn" @click="editSpec(spec)" />
            <van-button square type="danger" text="删除" class="swipe-btn" @click="confirmDelete(spec.id)" />
          </template>
        </van-swipe-cell>
      </van-list>
    </van-pull-refresh>

    <!-- 添加浮动按钮 -->
    <div class="fab" @click="showAddDialog = true">
      <van-icon name="plus" size="24" color="#fff" />
    </div>

    <!-- 添加规格弹窗 -->
    <van-popup v-model:show="showAddDialog" position="bottom" round :style="{ minHeight: '30%' }">
      <div class="popup-header">
        <h3>添加秘钥规格</h3>
      </div>
      <van-form @submit="addKeySpec">
        <van-cell-group inset>
          <van-field
            v-model="newSpec.name"
            label="名称"
            placeholder="请输入规格名称"
            :rules="[{ required: true, message: '请输入名称' }]"
          />
          <van-field
            v-model="newSpec.description"
            label="描述"
            type="textarea"
            placeholder="请输入描述（可选）"
            rows="2"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit">添加</van-button>
        </div>
      </van-form>
    </van-popup>

    <!-- 编辑规格弹窗 -->
    <van-popup v-model:show="showEditDialog" position="bottom" round :style="{ minHeight: '30%' }">
      <div class="popup-header">
        <h3>编辑秘钥规格</h3>
      </div>
      <van-form @submit="updateSpec">
        <van-cell-group inset>
          <van-field
            v-model="editingSpec.name"
            label="名称"
            placeholder="请输入规格名称"
            :rules="[{ required: true, message: '请输入名称' }]"
          />
          <van-field
            v-model="editingSpec.description"
            label="描述"
            type="textarea"
            placeholder="请输入描述（可选）"
            rows="2"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit">保存</van-button>
        </div>
      </van-form>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { keySpecAPI } from '../api';
import { showToast, showConfirmDialog } from 'vant';

const router = useRouter();

const keySpecs = ref([]);
const loading = ref(false);
const finished = ref(false);
const refreshing = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const newSpec = ref({ name: '', description: '' });
const editingSpec = ref({ id: null, name: '', description: '' });

const loadKeySpecs = async () => {
  try {
    const response = await keySpecAPI.list();
    keySpecs.value = response.data || [];
    finished.value = true;
  } catch (error) {
    showToast({ type: 'fail', message: '加载规格列表失败' });
  }
};

const onLoad = async () => {
  await loadKeySpecs();
  loading.value = false;
};

const onRefresh = async () => {
  await loadKeySpecs();
  refreshing.value = false;
};

const goToKeys = (specId) => {
  router.push(`/keys/${specId}`);
};

const addKeySpec = async () => {
  try {
    await keySpecAPI.create(newSpec.value);
    showToast({ type: 'success', message: '添加成功' });
    newSpec.value = { name: '', description: '' };
    showAddDialog.value = false;
    await loadKeySpecs();
  } catch (error) {
    showToast({ type: 'fail', message: '添加失败' });
  }
};

const editSpec = (spec) => {
  editingSpec.value = { ...spec };
  showEditDialog.value = true;
};

const updateSpec = async () => {
  try {
    await keySpecAPI.update(editingSpec.value.id, {
      name: editingSpec.value.name,
      description: editingSpec.value.description,
    });
    showToast({ type: 'success', message: '更新成功' });
    showEditDialog.value = false;
    await loadKeySpecs();
  } catch (error) {
    showToast({ type: 'fail', message: '更新失败' });
  }
};

const confirmDelete = async (id) => {
  try {
    await showConfirmDialog({ title: '确认删除', message: '确定要删除这个秘钥规格吗？' });
    await keySpecAPI.delete(id);
    showToast({ type: 'success', message: '删除成功' });
    await loadKeySpecs();
  } catch {
    // 用户取消
  }
};

onMounted(() => {
  loadKeySpecs();
});
</script>

<style scoped>
.key-spec-list {
  min-height: calc(100vh - 96px);
  background-color: #f7f8fa;
}

.swipe-btn {
  height: 100%;
}

.fab {
  position: fixed;
  right: 20px;
  bottom: 80px;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: #1989fa;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(25, 137, 250, 0.4);
  cursor: pointer;
  z-index: 100;
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
