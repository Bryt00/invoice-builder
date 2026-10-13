<template>
  <router-view />
  <GlobalToast />
  <GlobalConfirmModal />
  <PwaStatusNotification />
</template>

<script setup lang="ts">
import GlobalToast from './components/common/GlobalToast.vue'
import GlobalConfirmModal from './components/common/GlobalConfirmModal.vue'
import PwaStatusNotification from './components/common/PwaStatusNotification.vue'
</script>

