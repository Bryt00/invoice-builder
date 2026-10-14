<template>
  <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-3 max-w-sm sm:max-w-md w-[calc(100vw-2rem)] pointer-events-none">
    <!-- 1. Offline Alert Pill -->
    <Transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-4 opacity-0 scale-95"
      enter-to-class="translate-y-0 opacity-100 scale-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="isOffline"
        class="pointer-events-auto p-3.5 bg-surface-container-high/95 dark:bg-surface-container/95 backdrop-blur-md border border-amber-500/30 rounded-2xl shadow-xl flex items-center justify-between gap-3 text-on-surface"
        role="status"
        aria-live="polite"
      >
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-amber-500/15 text-amber-600 dark:text-amber-400 flex items-center justify-center shrink-0">
            <span class="material-symbols-outlined text-[18px]">wifi_off</span>
          </div>
          <div class="text-xs">
            <p class="font-bold text-on-surface">Offline Mode</p>
            <p class="text-on-surface-variant text-[11px] leading-tight">Working from offline cache. Changes will sync once connected.</p>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 2. Reconnected Notice (temporary 4s toast) -->
    <Transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-4 opacity-0 scale-95"
      enter-to-class="translate-y-0 opacity-100 scale-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="showBackOnline"
        class="pointer-events-auto p-3 bg-surface-container-high/95 backdrop-blur-md border border-emerald-500/30 rounded-2xl shadow-xl flex items-center gap-3 text-on-surface"
      >
        <div class="w-7 h-7 rounded-lg bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 flex items-center justify-center shrink-0">
          <span class="material-symbols-outlined text-[16px]">wifi</span>
        </div>
        <p class="font-semibold text-xs text-on-surface">Internet connection restored</p>
      </div>
    </Transition>

    <!-- 3. Update Available Banner -->
    <Transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-4 opacity-0 scale-95"
      enter-to-class="translate-y-0 opacity-100 scale-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="needRefresh"
        class="pointer-events-auto p-4 bg-surface-container-highest/95 dark:bg-surface-container-high/95 backdrop-blur-md border border-primary/30 rounded-2xl shadow-2xl flex flex-col gap-3 text-on-surface"
        role="alert"
      >
        <div class="flex items-start gap-3">
          <div class="w-9 h-9 rounded-xl bg-primary/15 text-primary flex items-center justify-center shrink-0 mt-0.5">
            <span class="material-symbols-outlined text-[20px]">system_update</span>
          </div>
          <div class="flex-1">
            <h4 class="font-bold text-xs sm:text-sm text-on-surface">Update Available</h4>
            <p class="text-xs text-on-surface-variant leading-relaxed mt-0.5">
              A newer version of Teks-Invoice is ready with improvements.
            </p>
          </div>
          <button
            @click="dismissRefresh"
            type="button"
            class="text-on-surface-variant hover:text-on-surface p-1 rounded-lg transition-colors cursor-pointer"
            title="Dismiss"
          >
            <span class="material-symbols-outlined text-[18px]">close</span>
          </button>
        </div>

        <div class="flex items-center justify-end gap-2 pt-1 border-t border-outline-variant/20">
          <button
            type="button"
            @click="dismissRefresh"
            class="px-3 py-1.5 text-xs font-semibold text-on-surface-variant hover:text-on-surface rounded-xl hover:bg-surface-variant/30 transition-colors cursor-pointer"
          >
            Later
          </button>
          <button
            type="button"
            @click="updateApp"
            class="px-3.5 py-1.5 text-xs font-bold bg-primary text-on-primary rounded-xl shadow-xs hover:bg-primary/90 transition-all flex items-center gap-1.5 cursor-pointer"
          >
            <span class="material-symbols-outlined text-[16px]">refresh</span>
            <span>Update Now</span>
          </button>
        </div>
      </div>
    </Transition>

    <!-- 4. Install App Floating Prompt -->
    <Transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-4 opacity-0 scale-95"
      enter-to-class="translate-y-0 opacity-100 scale-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="canInstall && !isInstalled && !bannerDismissed"
        class="pointer-events-auto p-4 bg-surface-container-high/95 backdrop-blur-md border border-outline-variant/40 rounded-2xl shadow-2xl flex flex-col gap-3 text-on-surface"
      >
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-xl bg-surface-container-lowest border border-outline-variant/30 flex items-center justify-center p-1.5 shrink-0 shadow-xs">
            <img :src="'/pwa-192x192.png'" alt="Teks-Invoice Icon" class="w-full h-full object-contain rounded-lg" />
          </div>
          <div class="flex-1">
            <h4 class="font-bold text-xs sm:text-sm text-on-surface">Install Teks-Invoice</h4>
            <p class="text-xs text-on-surface-variant leading-relaxed mt-0.5">
              Install as a desktop or mobile application for instant offline access and faster loading.
            </p>
          </div>
          <button
            @click="dismissBanner"
            type="button"
            class="text-on-surface-variant hover:text-on-surface p-1 rounded-lg transition-colors cursor-pointer"
            title="Dismiss"
          >
            <span class="material-symbols-outlined text-[18px]">close</span>
          </button>
        </div>

        <div class="flex items-center justify-end gap-2 pt-1 border-t border-outline-variant/20">
          <button
            type="button"
            @click="dismissBanner"
            class="px-3 py-1.5 text-xs font-semibold text-on-surface-variant hover:text-on-surface rounded-xl hover:bg-surface-variant/30 transition-colors cursor-pointer"
          >
            Not now
          </button>
          <button
            type="button"
            @click="handleInstallClick"
            class="px-3.5 py-1.5 text-xs font-bold bg-primary text-on-primary rounded-xl shadow-xs hover:bg-primary/90 transition-all flex items-center gap-1.5 cursor-pointer"
          >
            <span class="material-symbols-outlined text-[16px]">download</span>
            <span>Install App</span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { usePwa } from '../../composables/usePwa'

const { needRefresh, isOffline, canInstall, isInstalled, promptInstall, updateApp, dismissRefresh } = usePwa()

const showBackOnline = ref(false)
const bannerDismissed = ref(sessionStorage.getItem('teks_pwa_banner_dismissed') === 'true')

// Flash "Back online" when transition from offline to online occurs
watch(isOffline, (newVal, oldVal) => {
  if (oldVal === true && newVal === false) {
    showBackOnline.value = true
    setTimeout(() => {
      showBackOnline.value = false
    }, 4000)
  }
})

const handleInstallClick = async () => {
  const installed = await promptInstall()
  if (installed) {
    bannerDismissed.value = true
  }
}

const dismissBanner = () => {
  bannerDismissed.value = true
  sessionStorage.setItem('teks_pwa_banner_dismissed', 'true')
}
</script>
