<template>
  <!-- 1. SPROUTING COOKIE BANNER (Initial prompt) -->
  <Transition
    enter-active-class="transform ease-out duration-500 transition"
    enter-from-class="translate-y-8 opacity-0 scale-95"
    enter-to-class="translate-y-0 opacity-100 scale-100"
    leave-active-class="transition ease-in duration-300"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="translate-y-6 opacity-0 scale-95"
  >
    <div
      v-if="isBannerVisible"
      class="no-print fixed bottom-4 left-4 right-4 sm:right-auto sm:max-w-md z-50 pointer-events-auto"
      role="region"
      aria-label="Cookie consent banner"
    >
      <div class="p-4 sm:p-5 bg-surface/95 dark:bg-surface-container-high/95 backdrop-blur-md rounded-2xl shadow-2xl border border-outline-variant/60 flex flex-col gap-3.5 animate-sprout">
        <!-- Banner Header -->
        <div class="flex items-start justify-between gap-3">
          <div class="flex items-center gap-2.5">
            <div class="w-9 h-9 rounded-xl bg-primary/15 text-primary flex items-center justify-center shrink-0 shadow-xs">
              <span class="material-symbols-outlined text-[20px]">cookie</span>
            </div>
            <div>
              <h3 class="font-headline text-sm font-bold text-on-surface">Cookie Preferences</h3>
              <p class="font-body text-[11px] text-on-surface-variant font-medium">Privacy &amp; Data Security</p>
            </div>
          </div>
          <button
            type="button"
            @click="dismissBanner"
            class="text-on-surface-variant hover:text-on-surface p-1 rounded-lg hover:bg-surface-container transition-colors cursor-pointer"
            title="Dismiss notice"
          >
            <span class="material-symbols-outlined text-[18px]">close</span>
          </button>
        </div>

        <!-- Banner Content -->
        <p class="font-body text-xs text-on-surface-variant leading-relaxed">
          We use strictly necessary cookies to keep Teks-Invoice secure and functional, and optional cookies to analyze performance and save your editor preferences.
          <router-link to="/privacy" class="text-primary font-semibold hover:underline ml-1">Learn more in our Privacy Policy</router-link>.
        </p>

        <!-- Banner Actions -->
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 pt-1">
          <button
            type="button"
            @click="acceptAll"
            class="flex-1 px-3.5 py-2 rounded-xl bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-colors font-label text-xs font-bold text-center cursor-pointer shadow-xs flex items-center justify-center gap-1.5"
          >
            <span class="material-symbols-outlined text-[16px]">check_circle</span>
            <span>Accept All</span>
          </button>

          <button
            type="button"
            @click="acceptEssentialOnly"
            class="px-3.5 py-2 rounded-xl bg-surface-container hover:bg-surface-container-high text-on-surface transition-colors font-label text-xs font-semibold text-center cursor-pointer"
          >
            Essential Only
          </button>

          <button
            type="button"
            @click="openPreferencesModal"
            class="px-3 py-2 rounded-xl border border-outline-variant/60 hover:bg-surface-container-low text-on-surface transition-colors font-label text-xs font-semibold text-center cursor-pointer flex items-center justify-center gap-1"
          >
            <span class="material-symbols-outlined text-[16px]">tune</span>
            <span>Customize</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>

  <!-- 2. DETAILED COOKIE PREFERENCES MODAL (Granular settings dialog) -->
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="isModalVisible"
      class="no-print fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-5 bg-black/60 backdrop-blur-xs"
      @click.self="closePreferences"
      role="dialog"
      aria-modal="true"
      aria-labelledby="cookie-preferences-title"
    >
      <div class="w-full max-w-xl bg-surface dark:bg-surface-container-high rounded-2xl shadow-2xl border border-outline-variant/60 flex flex-col max-h-[92vh] overflow-hidden">
        <!-- Modal Header -->
        <div class="px-5 py-4 border-b border-outline-variant/40 flex items-center justify-between bg-surface-container-low/60 shrink-0">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-primary/15 text-primary flex items-center justify-center shadow-xs">
              <span class="material-symbols-outlined text-[20px]">shield_check</span>
            </div>
            <div>
              <h2 id="cookie-preferences-title" class="font-headline text-base sm:text-lg font-bold text-on-surface">Cookie Preferences</h2>
              <p class="font-body text-xs text-on-surface-variant">Manage how Teks-Invoice stores and accesses data on your browser</p>
            </div>
          </div>
          <button
            type="button"
            @click="closePreferences"
            class="text-on-surface-variant hover:text-on-surface p-1.5 rounded-lg hover:bg-surface-container transition-colors cursor-pointer"
            title="Close dialog"
          >
            <span class="material-symbols-outlined text-[20px]">close</span>
          </button>
        </div>

        <!-- Modal Body (Categories) -->
        <div class="p-5 overflow-y-auto space-y-4 font-body text-xs sm:text-sm">
          <p class="text-xs text-on-surface-variant leading-relaxed">
            When you visit Teks-Invoice, we use cookies and local storage to deliver our services, maintain your authenticated session, and understand how the app is utilized. You may customize your non-essential preferences below.
          </p>

          <!-- Category 1: Strictly Necessary (Locked) -->
          <div class="p-4 rounded-xl bg-surface-container-lowest/80 border border-outline-variant/50 space-y-2">
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-[18px]">lock</span>
                <span class="font-headline font-bold text-sm text-on-surface">Strictly Necessary Cookies</span>
              </div>
              <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-primary/15 text-primary border border-primary/20 shrink-0 flex items-center gap-1">
                <span class="material-symbols-outlined text-[12px]">check</span> Always Active
              </span>
            </div>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Required for basic operation, secure sign-in tokens, CSRF protection, and local offline draft recovery. These cannot be switched off.
            </p>
          </div>

          <!-- Category 2: Functional & Preferences -->
          <div class="p-4 rounded-xl bg-surface-container-lowest/80 border border-outline-variant/50 space-y-2">
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-[18px]">tune</span>
                <span class="font-headline font-bold text-sm text-on-surface">Functional &amp; Preferences</span>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="tempFunctional" class="sr-only peer" />
                <div class="w-11 h-6 bg-surface-variant peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Enables the site to remember your dark/light theme, preferred invoice template styles, default currency, and layout view settings.
            </p>
          </div>

          <!-- Category 3: Analytics & Performance -->
          <div class="p-4 rounded-xl bg-surface-container-lowest/80 border border-outline-variant/50 space-y-2">
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-[18px]">monitoring</span>
                <span class="font-headline font-bold text-sm text-on-surface">Analytics &amp; Performance</span>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="tempAnalytics" class="sr-only peer" />
                <div class="w-11 h-6 bg-surface-variant peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Helps us anonymously gauge page speeds, PDF generation latency, and feature usage to catch bugs and boost platform reliability.
            </p>
          </div>

          <!-- Category 4: Marketing & Communications -->
          <div class="p-4 rounded-xl bg-surface-container-lowest/80 border border-outline-variant/50 space-y-2">
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-[18px]">campaign</span>
                <span class="font-headline font-bold text-sm text-on-surface">Marketing &amp; Updates</span>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="tempMarketing" class="sr-only peer" />
                <div class="w-11 h-6 bg-surface-variant peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Provides relevant news on upcoming features, billing workflow enhancements, and occasional promotion announcements.
            </p>
          </div>
        </div>

        <!-- Modal Footer Actions -->
        <div class="px-5 py-4 border-t border-outline-variant/40 bg-surface-container-low/60 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-2.5 shrink-0">
          <button
            type="button"
            @click="rejectOptional"
            class="px-4 py-2 rounded-xl bg-surface-container hover:bg-surface-container-high text-on-surface font-label text-xs font-semibold transition-colors cursor-pointer text-center"
          >
            Reject All Optional
          </button>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="saveCustom"
              class="flex-1 sm:flex-initial px-4 py-2 rounded-xl border border-primary/50 text-primary hover:bg-primary-container/20 font-label text-xs font-bold transition-colors cursor-pointer text-center"
            >
              Save Preferences
            </button>
            <button
              type="button"
              @click="acceptAllModal"
              class="flex-1 sm:flex-initial px-5 py-2 rounded-xl bg-primary text-on-primary hover:bg-on-primary-fixed-variant font-label text-xs font-bold transition-colors cursor-pointer text-center shadow-xs"
            >
              Accept All
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useCookieConsent } from '@/composables/useCookieConsent'

const {
  hasConsented,
  isBannerVisible,
  isModalVisible,
  preferences,
  acceptAll,
  acceptEssentialOnly,
  saveCustomPreferences,
  openPreferences,
  closePreferences,
  dismissBanner,
} = useCookieConsent()

// Local copy for the modal form
const tempFunctional = ref(preferences.value.functional)
const tempAnalytics = ref(preferences.value.analytics)
const tempMarketing = ref(preferences.value.marketing)

// Keep temp state synced whenever preferences or modal visibility updates
watch(isModalVisible, (visible) => {
  if (visible) {
    tempFunctional.value = preferences.value.functional
    tempAnalytics.value = preferences.value.analytics
    tempMarketing.value = preferences.value.marketing
  }
})

function openPreferencesModal() {
  tempFunctional.value = preferences.value.functional
  tempAnalytics.value = preferences.value.analytics
  tempMarketing.value = preferences.value.marketing
  openPreferences()
}

function saveCustom() {
  saveCustomPreferences({
    functional: tempFunctional.value,
    analytics: tempAnalytics.value,
    marketing: tempMarketing.value,
  })
}

function rejectOptional() {
  acceptEssentialOnly()
}

function acceptAllModal() {
  acceptAll()
}
</script>

<style scoped>
@keyframes sprout {
  0% {
    opacity: 0;
    transform: translateY(24px) scale(0.95);
  }
  65% {
    transform: translateY(-4px) scale(1.02);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.animate-sprout {
  animation: sprout 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}
</style>
