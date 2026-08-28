<template>
  <div class="flex flex-col md:flex-row min-h-screen relative w-full bg-surface">
    <!-- Admin Mobile TopBar (Mobile Only) -->
    <header class="md:hidden sticky top-0 z-40 bg-surface/95 backdrop-blur-xl border-b border-amber-500/30 px-4 py-3 flex items-center justify-between shadow-sm">
      <router-link to="/user/admin/dashboard" class="flex items-center gap-2">
        <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-7 h-7 rounded-lg object-contain ring-2 ring-amber-500/40">
        <span class="font-headline text-base font-bold text-on-surface">Teks-Invoice</span>
        <span class="px-1.5 py-0.5 rounded text-[9px] font-black uppercase bg-amber-500/20 text-amber-600 border border-amber-500/30">ADMIN</span>
      </router-link>
      <button type="button" @click="mobileNavOpen = !mobileNavOpen" class="text-on-surface-variant hover:text-amber-600 p-1.5 rounded-lg cursor-pointer">
        <span class="material-symbols-outlined text-[22px]">menu</span>
      </button>
    </header>

    <!-- Mobile Dropdown Navigation Drawer -->
    <div v-if="mobileNavOpen" class="md:hidden fixed inset-x-0 top-[57px] z-40 border-b border-outline-variant/50 bg-surface/95 backdrop-blur-xl px-4 py-3 space-y-1.5 font-label text-sm font-medium shadow-xl">
      <router-link 
        v-for="link in navLinks" 
        :key="link.path"
        :to="link.path"
        @click="mobileNavOpen = false"
        :class="['flex items-center gap-3 px-3 py-2 rounded-xl transition-colors', isCurrentPath(link.path) ? 'text-amber-600 bg-amber-500/10 font-bold' : 'text-on-surface-variant hover:text-amber-600 hover:bg-surface-container-low']"
      >
        <span class="material-symbols-outlined text-[18px]">{{ link.icon }}</span> {{ link.name }}
      </router-link>
      
      <div class="pt-2 border-t border-outline-variant/40 flex flex-col gap-1">
        <router-link to="/user/dashboard" class="flex items-center gap-3 px-3 py-2 rounded-xl text-primary font-semibold hover:bg-surface-container-low transition-colors">
          <span class="material-symbols-outlined text-[18px]">swap_horiz</span> User Workspace
        </router-link>
        <button type="button" @click="handleLogout" class="w-full flex items-center gap-3 px-3 py-2 rounded-xl text-error hover:bg-error-container/30 transition-colors cursor-pointer">
          <span class="material-symbols-outlined text-[18px]">logout</span> Sign Out
        </button>
      </div>
    </div>

    <!-- Desktop SideNavBar Component -->
    <nav class="hidden md:flex h-screen w-64 fixed left-0 top-0 bg-surface-container-lowest/90 backdrop-blur-xl flex-col py-6 justify-between border-r border-amber-500/20 z-30 shadow-sm">
      <div>
        <!-- Header Logo -->
        <router-link to="/user/admin/dashboard" class="px-6 mb-8 flex items-center gap-3 group">
          <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-10 h-10 rounded-xl object-contain shadow-md ring-2 ring-amber-500/40 group-hover:scale-105 transition-transform">
          <div>
            <h1 class="brand-title font-headline text-lg font-bold">Teks-Invoice</h1>
            <span class="px-2 py-0.5 rounded-md text-[10px] font-black uppercase bg-amber-500/20 text-amber-600 dark:text-amber-400 border border-amber-500/30">ADMIN CONSOLE</span>
          </div>
        </router-link>

        <!-- Primary Navigation List -->
        <ul class="flex flex-col gap-1 px-3 font-label text-sm font-medium">
          <li v-for="link in navLinks" :key="link.path">
            <router-link 
              :to="link.path"
              :class="['px-4 py-2.5 rounded-xl flex items-center gap-3 transition-all', isCurrentPath(link.path) ? 'text-amber-600 bg-amber-500/10 font-bold' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-container-low']"
            >
              <span class="material-symbols-outlined text-[20px]">{{ link.icon }}</span>
              {{ link.name }}
            </router-link>
          </li>
        </ul>
      </div>

      <!-- Footer Workspace Switcher & Sign Out -->
      <div class="px-3 font-label text-sm font-medium space-y-2">
        <div class="px-1">
          <router-link to="/user/dashboard" class="w-full py-2.5 px-4 rounded-xl bg-amber-500/10 text-amber-700 dark:text-amber-400 hover:bg-amber-500/20 border border-amber-500/30 font-label text-xs font-semibold flex items-center justify-center gap-2 transition-colors">
            <span class="material-symbols-outlined text-[16px]">swap_horiz</span>
            User Workspace
          </router-link>
        </div>
        <div class="pt-2 border-t border-outline-variant/40">
          <button type="button" @click="handleLogout" class="w-full text-on-surface-variant hover:text-error hover:bg-error-container/30 px-4 py-2.5 rounded-xl flex items-center gap-3 transition-all text-left cursor-pointer">
            <span class="material-symbols-outlined text-[20px]">logout</span>
            Sign Out
          </button>
        </div>
      </div>
    </nav>

    <!-- Main Content Area -->
    <main class="flex-1 w-full bg-transparent md:ml-64 px-4 sm:px-6 lg:px-8 py-6 md:py-8">
      <FlashAlert />
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import FlashAlert from '../components/common/FlashAlert.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const mobileNavOpen = ref(false)

const navLinks = [
  { name: 'Overview', path: '/user/admin/dashboard', icon: 'grid_view' },
  { name: 'Users', path: '/user/admin/users', icon: 'group' },
  { name: 'Invoices', path: '/user/admin/invoices', icon: 'receipt_long' },
  { name: 'Packages', path: '/user/admin/packages', icon: 'package_2' },
  { name: 'Payments', path: '/user/admin/payments', icon: 'payments' },
  { name: 'Credits', path: '/user/admin/credits', icon: 'bolt' },
  { name: 'Webhooks', path: '/user/admin/webhooks', icon: 'webhook' },
  { name: 'Audit Logs', path: '/user/admin/audit-logs', icon: 'history' },
  { name: 'Settings', path: '/user/admin/settings', icon: 'settings' }
]

function isCurrentPath(path) {
  return route.path === path
}

async function handleLogout() {
  await authStore.logout()
  router.push('/user/login')
}
</script>
