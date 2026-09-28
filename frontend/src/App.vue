<template>
  <router-view />
  <GlobalToast />
  <GlobalConfirmModal />
</template>

<script setup lang="ts">
import GlobalToast from './components/common/GlobalToast.vue'
import GlobalConfirmModal from './components/common/GlobalConfirmModal.vue'
</script>
