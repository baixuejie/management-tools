<template>
  <div class="ledger-page">
    <div class="stats-card">
      <div class="stats-grid">
        <div class="stat-item">
          <p>总成本</p>
          <h3>{{ money(ledger.statistics.total_cost) }}</h3>
        </div>
        <div class="stat-item">
          <p>总收入</p>
          <h3>{{ money(ledger.statistics.total_revenue) }}</h3>
        </div>
        <div class="stat-item">
          <p>净利润</p>
          <h3>{{ money(ledger.statistics.net_profit) }}</h3>
        </div>
        <div class="stat-item">
          <p>新客数</p>
          <h3>{{ ledger.statistics.new_customers }}</h3>
        </div>
        <div class="stat-item">
          <p>续费数</p>
          <h3>{{ ledger.statistics.renewal_customers }}</h3>
        </div>
        <div class="stat-item">
          <p>手续费</p>
          <h3>{{ money(ledger.statistics.total_commission) }}</h3>
        </div>
      </div>
    </div>

    <van-cell-group inset style="margin: 0 16px 12px;">
      <van-tabs v-model:active="activeType">
        <van-tab title="新客" name="new" />
        <van-tab title="续费" name="renewal" />
      </van-tabs>
      <template v-if="activeType === 'new'">
        <van-field v-model="customerName" label="购买人" placeholder="请输入购买人姓名" />
      </template>
      <template v-else>
        <van-field
          :model-value="selectedCustomerName"
          label="购买人"
          readonly
          is-link
          placeholder="请选择购买人"
          @click="showCustomerPicker = true"
        />
      </template>
      <van-field v-model="amount" type="number" label="金额" placeholder="请输入金额" />
      <van-field label="渠道">
        <template #input>
          <van-radio-group v-model="channel" direction="horizontal">
            <van-radio name="xianyu">闲鱼</van-radio>
            <van-radio name="wechat">微信</van-radio>
          </van-radio-group>
        </template>
      </van-field>
      <div style="margin: 12px;">
        <van-button block type="primary" :loading="submitting" @click="submitTransaction">
          提交记账
        </van-button>
      </div>
    </van-cell-group>

    <van-cell-group inset style="margin: 0 16px 12px;">
      <van-cell title="最近交易" :label="`共 ${ledger.transactionTotal} 条`" />
    </van-cell-group>

    <van-list>
      <van-cell
        v-for="item in ledger.transactions"
        :key="item.id"
        :title="`${item.customer_name}（${item.is_new_customer ? '新客' : '续费'}）`"
        :label="`${item.recorder_name} · ${time(item.created_at)}`"
      >
        <template #value>
          <div class="tx-value">
            <span class="money">{{ money(item.amount) }}</span>
            <small>{{ item.channel === 'xianyu' ? '闲鱼' : '微信' }} · 手续费 {{ money(item.commission_amount) }}</small>
          </div>
        </template>
      </van-cell>
    </van-list>

    <div class="fab" @click="showCostDialog = true">+成本</div>

    <van-popup v-model:show="showCostDialog" position="bottom" round>
      <div class="popup-header"><h3>新增成本</h3></div>
      <van-cell-group inset>
        <van-field v-model="costAmount" type="number" label="金额" placeholder="请输入成本金额" />
        <van-field v-model="costNote" type="textarea" label="备注" rows="2" placeholder="可选" />
      </van-cell-group>
      <div style="margin: 16px;">
        <van-button block type="primary" :loading="addingCost" @click="submitCost">保存成本</van-button>
      </div>
    </van-popup>

    <van-popup v-model:show="showCustomerPicker" position="bottom" round>
      <van-search v-model="customerSearch" placeholder="搜索购买人" />
      <div class="picker-list">
        <van-cell
          v-for="item in filteredCustomers"
          :key="item.id"
          :title="item.name"
          @click="selectCustomer(item)"
        />
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { showToast } from 'vant';
import { useLedgerStore } from '../stores/ledger';

const ledger = useLedgerStore();

const activeType = ref('new');
const customerName = ref('');
const selectedCustomerId = ref(null);
const amount = ref('');
const channel = ref('xianyu');
const submitting = ref(false);

const showCostDialog = ref(false);
const costAmount = ref('');
const costNote = ref('');
const addingCost = ref(false);

const showCustomerPicker = ref(false);
const customerSearch = ref('');

const filteredCustomers = computed(() => {
  const key = customerSearch.value.trim().toLowerCase();
  if (!key) return ledger.customers;
  return ledger.customers.filter((item) => item.name.toLowerCase().includes(key));
});

const selectedCustomerName = computed(() => {
  const customer = ledger.customers.find((item) => item.id === selectedCustomerId.value);
  return customer ? customer.name : '';
});

const money = (value) => `￥${Number(value || 0).toFixed(2)}`;
const time = (value) => new Date(value).toLocaleString();

const loadData = async () => {
  await Promise.all([
    ledger.fetchStatistics(),
    ledger.fetchTransactions({ page: 1, limit: 20 }),
    ledger.fetchCustomers(),
  ]);
};

const selectCustomer = (customer) => {
  selectedCustomerId.value = customer.id;
  showCustomerPicker.value = false;
};

const submitTransaction = async () => {
  const numericAmount = Number(amount.value);
  if (!numericAmount || numericAmount <= 0) {
    showToast({ type: 'fail', message: '金额必须大于 0' });
    return;
  }

  const payload = {
    amount: numericAmount,
    channel: channel.value,
    is_new_customer: activeType.value === 'new',
  };

  if (activeType.value === 'new') {
    if (!customerName.value.trim()) {
      showToast({ type: 'fail', message: '请输入购买人姓名' });
      return;
    }
    payload.customer_name = customerName.value.trim();
  } else {
    if (!selectedCustomerId.value) {
      showToast({ type: 'fail', message: '请选择购买人' });
      return;
    }
    payload.customer_id = selectedCustomerId.value;
  }

  submitting.value = true;
  try {
    await ledger.createTransaction(payload);
    showToast({ type: 'success', message: '记账成功' });
    amount.value = '';
    if (activeType.value === 'new') {
      customerName.value = '';
    }
    await loadData();
  } catch (error) {
    const msg = error.response?.data?.error || '记账失败';
    showToast({ type: 'fail', message: msg });
  } finally {
    submitting.value = false;
  }
};

const submitCost = async () => {
  const numericAmount = Number(costAmount.value);
  if (!numericAmount || numericAmount <= 0) {
    showToast({ type: 'fail', message: '成本金额必须大于 0' });
    return;
  }

  addingCost.value = true;
  try {
    await ledger.createCost({ amount: numericAmount, note: costNote.value.trim() });
    showToast({ type: 'success', message: '成本添加成功' });
    costAmount.value = '';
    costNote.value = '';
    showCostDialog.value = false;
    await loadData();
  } catch (error) {
    const msg = error.response?.data?.error || '添加成本失败';
    showToast({ type: 'fail', message: msg });
  } finally {
    addingCost.value = false;
  }
};

onMounted(loadData);
</script>

<style scoped>
.ledger-page {
  padding-bottom: 70px;
  background: #f7f8fa;
  min-height: calc(100vh - 96px);
}

.stats-card {
  margin: 12px 16px;
  padding: 14px;
  border-radius: 12px;
  color: #fff;
  background: linear-gradient(135deg, #4facfe, #00c4cc);
  box-shadow: 0 6px 16px rgba(79, 172, 254, 0.25);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.stat-item p {
  font-size: 12px;
  opacity: 0.9;
}

.stat-item h3 {
  margin-top: 4px;
  font-size: 16px;
}

.tx-value {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.tx-value .money {
  color: #07c160;
  font-weight: 700;
}

.tx-value small {
  color: #969799;
}

.fab {
  position: fixed;
  right: 16px;
  bottom: 84px;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: linear-gradient(135deg, #ff976a, #ff6a88);
  box-shadow: 0 6px 16px rgba(255, 106, 136, 0.35);
  font-weight: 700;
}

.popup-header {
  padding: 16px;
  text-align: center;
}

.picker-list {
  max-height: 280px;
  overflow: auto;
}
</style>
