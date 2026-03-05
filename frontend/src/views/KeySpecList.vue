<template>
  <div class="key-spec-list">
    <van-nav-bar title="Key Specifications">
      <template #right>
        <van-icon name="setting-o" size="18" @click="goToConfig" />
      </template>
    </van-nav-bar>
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-empty v-if="!loading && keySpecs.length === 0" description="No key specifications yet" />
      <van-list
        v-else
        v-model:loading="loading"
        :finished="finished"
        finished-text="No more data"
        @load="onLoad"
      >
        <van-swipe-cell v-for="spec in keySpecs" :key="spec.id">
          <van-cell
            :title="spec.name"
            :label="`${spec.description || 'No description'} • Created: ${formatDate(spec.created_at)}`"
            is-link
            @click="goToKeys(spec.id)"
          />
          <template #right>
            <van-button square type="primary" text="Edit" @click="editSpec(spec)" />
            <van-button square type="danger" text="Delete" @click="confirmDelete(spec.id)" />
          </template>
        </van-swipe-cell>
      </van-list>
    </van-pull-refresh>
    <van-floating-bubble
      icon="plus"
      @click="showAddDialog = true"
    />

    <!-- Add Dialog -->
    <van-popup v-model:show="showAddDialog" position="bottom" round :style="{ height: '40%' }">
      <div class="popup-header">
        <h3>Add Key Specification</h3>
      </div>
      <van-form @submit="addKeySpec">
        <van-cell-group inset>
          <van-field
            v-model="newSpec.name"
            label="Name"
            placeholder="Enter name"
            :rules="[{ required: true, message: 'Name is required' }]"
          />
          <van-field
            v-model="newSpec.description"
            label="Description"
            type="textarea"
            placeholder="Enter description"
            rows="2"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit">
            Add
          </van-button>
          <van-button round block plain type="default" @click="showAddDialog = false" style="margin-top: 8px;">
            Cancel
          </van-button>
        </div>
      </van-form>
    </van-popup>

    <!-- Edit Dialog -->
    <van-popup v-model:show="showEditDialog" position="bottom" round :style="{ height: '40%' }">
      <div class="popup-header">
        <h3>Edit Key Specification</h3>
      </div>
      <van-form @submit="updateSpec">
        <van-cell-group inset>
          <van-field
            v-model="editingSpec.name"
            label="Name"
            placeholder="Enter name"
            :rules="[{ required: true, message: 'Name is required' }]"
          />
          <van-field
            v-model="editingSpec.description"
            label="Description"
            type="textarea"
            placeholder="Enter description"
            rows="2"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit">
            Update
          </van-button>
          <van-button round block plain type="default" @click="showEditDialog = false" style="margin-top: 8px;">
            Cancel
          </van-button>
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

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

const loadKeySpecs = async () => {
  try {
    const response = await keySpecAPI.list();
    keySpecs.value = response.data || [];
    finished.value = true;
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to load key specs' });
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

const goToConfig = () => {
  router.push('/config');
};

const addKeySpec = async () => {
  if (!newSpec.value.name) {
    showToast({ type: 'fail', message: 'Name is required' });
    return;
  }
  try {
    await keySpecAPI.create(newSpec.value);
    showToast({ type: 'success', message: 'Key spec added' });
    newSpec.value = { name: '', description: '' };
    showAddDialog.value = false;
    await loadKeySpecs();
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to add key spec' });
  }
};

const editSpec = (spec) => {
  editingSpec.value = { ...spec };
  showEditDialog.value = true;
};

const updateSpec = async () => {
  if (!editingSpec.value.name) {
    showToast({ type: 'fail', message: 'Name is required' });
    return;
  }
  try {
    await keySpecAPI.update(editingSpec.value.id, {
      name: editingSpec.value.name,
      description: editingSpec.value.description,
    });
    showToast({ type: 'success', message: 'Key spec updated' });
    showEditDialog.value = false;
    await loadKeySpecs();
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to update key spec' });
  }
};

const confirmDelete = async (id) => {
  try {
    await showConfirmDialog({
      title: 'Confirm Delete',
      message: 'Are you sure you want to delete this key specification?',
    });
    await deleteSpec(id);
  } catch {
    // User cancelled
  }
};

const deleteSpec = async (id) => {
  try {
    await keySpecAPI.delete(id);
    showToast({ type: 'success', message: 'Key spec deleted' });
    await loadKeySpecs();
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to delete key spec' });
  }
};

onMounted(() => {
  loadKeySpecs();
});
</script>

<style scoped>
.key-spec-list {
  min-height: 100vh;
  background-color: #f7f8fa;
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
