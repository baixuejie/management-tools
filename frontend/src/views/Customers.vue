<template>
  <div class="customers-page">
    <van-search v-model="search" placeholder="搜索购买人" @update:model-value="loadCustomers" />

    <van-list>
      <van-cell v-for="item in customers" :key="item.id" :title="item.name" :label="new Date(item.created_at).toLocaleString()">
        <template #right-icon>
          <van-button size="small" type="primary" plain @click="startEdit(item)">编辑</van-button>
        </template>
      </van-cell>
    </van-list>

    <van-popup v-model:show="showEdit" position="bottom" round>
      <div class="popup-header"><h3>编辑购买人</h3></div>
      <van-cell-group inset>
        <van-field v-model="editingName" label="名称" placeholder="请输入购买人名称" />
      </van-cell-group>
      <div style="margin: 16px;">
        <van-button block type="primary" :loading="saving" @click="saveEdit">保存</van-button>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { showToast } from 'vant';
import { useLedgerStore } from '../stores/ledger';

const ledger = useLedgerStore();

const search = ref('');
const customers = ref([]);

const showEdit = ref(false);
const editingId = ref(null);
const editingName = ref('');
const saving = ref(false);

const loadCustomers = async () => {
  try {
    customers.value = await ledger.fetchCustomers(search.value.trim());
  } catch {
    showToast({ type: 'fail', message: '加载购买人失败' });
  }
};

const startEdit = (item) => {
  editingId.value = item.id;
  editingName.value = item.name;
  showEdit.value = true;
};

const saveEdit = async () => {
  if (!editingName.value.trim()) {
    showToast({ type: 'fail', message: '请输入名称' });
    return;
  }

  saving.value = true;
  try {
    await ledger.updateCustomer(editingId.value, { name: editingName.value.trim() });
    showToast({ type: 'success', message: '更新成功' });
    showEdit.value = false;
    await loadCustomers();
  } catch (error) {
    const msg = error.response?.data?.error || '更新失败';
    showToast({ type: 'fail', message: msg });
  } finally {
    saving.value = false;
  }
};

onMounted(loadCustomers);
</script>

<style scoped>
.customers-page {
  min-height: calc(100vh - 96px);
  background: #f7f8fa;
}

.popup-header {
  padding: 16px;
  text-align: center;
}
</style>
