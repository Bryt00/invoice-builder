<template>
  <div class="min-h-screen flex flex-col justify-center py-12 sm:px-6 lg:px-8 bg-background relative overflow-hidden">
    <!-- Subtle Gradient Background (No particles/animations) -->
    <div class="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-tertiary/5 pointer-events-none"></div>

    <router-view />
  </div>
</template>

<script setup>
</script>
