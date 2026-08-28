<template>
  <div class="min-h-screen bg-background text-on-surface flex items-center justify-center p-4">
    <div class="max-w-md w-full glass-card rounded-3xl p-8 border border-outline-variant/60 text-center shadow-xl">
      <div class="w-16 h-16 rounded-2xl bg-error-container/40 text-error border border-error/20 flex items-center justify-center mx-auto mb-6">
        <span class="material-symbols-outlined text-[36px]">find_in_page</span>
      </div>

      <h1 class="font-headline text-5xl font-extrabold text-on-surface mb-2">404</h1>
      <h2 class="font-headline text-xl font-bold text-on-surface mb-3">Page Not Found</h2>
      <p class="font-body text-sm text-on-surface-variant mb-8 leading-relaxed">
        The page or resource you are looking for doesn't exist, has been moved, or is temporarily unavailable.
      </p>

      <div class="flex items-center justify-center gap-3">
        <button
          @click="$router.back()"
          class="px-5 py-2.5 rounded-xl border border-outline-variant/60 font-label text-sm font-semibold text-on-surface-variant hover:bg-surface-container-low transition-colors"
        >
          Go Back
        </button>
        <router-link
          to="/"
          class="px-6 py-2.5 rounded-xl font-label text-sm font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all inline-flex items-center gap-2 shadow-md"
        >
          <span class="material-symbols-outlined text-[18px]">home</span>
          <span>Home</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>
