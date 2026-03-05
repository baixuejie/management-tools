<template>
  <div class="more-page">
    <van-cell-group inset title="账本管理">
      <van-cell title="成本记录" is-link to="/more/costs" />
      <van-cell title="购买人管理" is-link to="/more/customers" />
    </van-cell-group>

    <van-cell-group inset title="秘钥管理" style="margin-top: 12px;">
      <van-cell title="规格管理" is-link to="/more/specs" />
      <van-cell title="复制模板" is-link to="/more/template" />
    </van-cell-group>

    <van-cell-group inset title="系统设置" style="margin-top: 12px;">
      <van-cell title="个人信息" is-link to="/more/profile" />
    </van-cell-group>
  </div>
</template>

<script setup>
</script>

<style scoped>
.more-page {
  min-height: calc(100vh - 96px);
  background: #f7f8fa;
  padding-top: 12px;
}
</style>
