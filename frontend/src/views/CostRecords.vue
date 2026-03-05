<template>
  <div class="cost-page">
    <van-cell-group inset style="margin: 12px 16px;">
      <van-field v-model="amount" type="number" label="金额" placeholder="请输入金额" />
      <van-field v-model="note" type="textarea" rows="2" label="备注" placeholder="可选" />
      <div style="margin: 12px;">
        <van-button block type="primary" :loading="submitting" @click="createCost">新增成本</van-button>
      </div>
    </van-cell-group>

    <van-cell-group inset style="margin: 0 16px 12px;">
      <van-cell :title="`共 ${total} 条`" />
    </van-cell-group>

    <van-list>
      <van-swipe-cell v-for="item in items" :key="item.id">
        <van-cell
          :title="`￥${Number(item.amount).toFixed(2)} · ${item.recorder_name}`"
          :label="`${item.note || '无备注'} · ${new Date(item.created_at).toLocaleString()}`"
        />
        <template #right>
          <van-button square type="danger" text="删除" @click="removeCost(item.id)" />
        </template>
      </van-swipe-cell>
    </van-list>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { showConfirmDialog, showToast } from 'vant';
import { useLedgerStore } from '../stores/ledger';

const ledger = useLedgerStore();

const amount = ref('');
const note = ref('');
const submitting = ref(false);

const items = ref([]);
const total = ref(0);

const loadCosts = async () => {
  try {
    const response = await ledger.fetchCosts({ page: 1, limit: 100 });
    items.value = response.items || [];
    total.value = response.total || 0;
  } catch {
    showToast({ type: 'fail', message: '加载成本记录失败' });
  }
};

const createCost = async () => {
  const numericAmount = Number(amount.value);
  if (!numericAmount || numericAmount <= 0) {
    showToast({ type: 'fail', message: '金额必须大于 0' });
    return;
  }

  submitting.value = true;
  try {
    await ledger.createCost({ amount: numericAmount, note: note.value.trim() });
    showToast({ type: 'success', message: '新增成功' });
    amount.value = '';
    note.value = '';
    await loadCosts();
  } catch (error) {
    const msg = error.response?.data?.error || '新增失败';
    showToast({ type: 'fail', message: msg });
  } finally {
    submitting.value = false;
  }
};

const removeCost = async (id) => {
  try {
    await showConfirmDialog({ title: '确认', message: '确定删除这条成本记录吗？' });
    await ledger.deleteCost(id);
    showToast({ type: 'success', message: '删除成功' });
    await loadCosts();
  } catch {
    // 用户取消
  }
};

onMounted(loadCosts);
</script>

<style scoped>
.cost-page {
  min-height: calc(100vh - 96px);
  background: #f7f8fa;
}
</style>
