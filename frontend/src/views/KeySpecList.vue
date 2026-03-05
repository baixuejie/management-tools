<template>
  <div class="key-spec-list">
    <van-nav-bar title="Key Specifications" />
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="No more data"
        @load="onLoad"
      >
        <van-cell
          v-for="spec in keySpecs"
          :key="spec.id"
          :title="spec.name"
          :label="spec.description"
          is-link
          @click="goToKeys(spec.id)"
        />
      </van-list>
    </van-pull-refresh>
    <van-floating-bubble
      icon="plus"
      @click="showAddDialog = true"
    />
    <van-dialog
      v-model:show="showAddDialog"
      title="Add Key Specification"
      show-cancel-button
      @confirm="addKeySpec"
    >
      <van-form>
        <van-field v-model="newSpec.name" label="Name" placeholder="Enter name" />
        <van-field v-model="newSpec.description" label="Description" placeholder="Enter description" />
      </van-form>
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { keySpecAPI } from '../api';
import { showToast } from 'vant';

const router = useRouter();

const keySpecs = ref([]);
const loading = ref(false);
const finished = ref(false);
const refreshing = ref(false);
const showAddDialog = ref(false);
const newSpec = ref({ name: '', description: '' });

const loadKeySpecs = async () => {
  try {
    const response = await keySpecAPI.list();
    keySpecs.value = response.data;
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

const addKeySpec = async () => {
  try {
    await keySpecAPI.create(newSpec.value);
    showToast({ type: 'success', message: 'Key spec added' });
    newSpec.value = { name: '', description: '' };
    await loadKeySpecs();
  } catch (error) {
    showToast({ type: 'fail', message: 'Failed to add key spec' });
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
</style>
