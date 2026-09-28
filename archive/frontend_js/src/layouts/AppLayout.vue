<template>
  <div class="flex flex-col flex-1 min-h-screen relative w-full">
    <!-- Subtle Background Elements -->
    <div class="fixed inset-0 z-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-[20%] -left-[10%] w-[50%] h-[50%] bg-primary-fixed rounded-full blur-[120px] opacity-40"></div>
        <div class="absolute top-[60%] -right-[10%] w-[60%] h-[60%] bg-secondary-fixed rounded-full blur-[140px] opacity-30"></div>
    </div>

    <!-- Main Application Wrapper -->
    <div class="flex flex-col flex-1 relative z-10 w-full">
        <AppNavbar />
        
        <main class="flex-1 bg-transparent max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-6 md:py-8">
            <FlashAlert />
            <router-view />
        </main>
        
        <AppFooter />
    </div>
  </div>
</template>

<script setup>
import AppNavbar from '../components/common/AppNavbar.vue'
import AppFooter from '../components/common/AppFooter.vue'
import FlashAlert from '../components/common/FlashAlert.vue'
</script>
