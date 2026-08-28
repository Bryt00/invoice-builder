<template>
  <header class="bg-[#FDFBF7] w-full sticky top-0 z-50 border-b border-outline-variant/30 shadow-sm">
      <div class="max-w-screen-2xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex justify-between items-center relative">
          <!-- Left: Logo & Brand -->
          <div class="flex items-center gap-3 w-1/4">
              <router-link to="/user/dashboard" class="flex items-center gap-3 group">
                  <div class="w-10 h-10 rounded-xl bg-surface-container-low border border-outline-variant/40 flex items-center justify-center overflow-hidden shadow-sm">
                      <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-8 h-8 object-contain group-hover:scale-105 transition-transform">
                  </div>
                  <span class="font-headline text-xl font-extrabold text-on-surface">Teks-Invoice</span>
              </router-link>
          </div>
          
          <!-- Center: Navigation -->
          <nav class="hidden md:flex items-center justify-center gap-2 flex-1" v-if="authStore.isAuthenticated">
              <router-link class="px-4 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-2"
                 :class="$route.name === 'dashboard' ? 'text-on-surface bg-surface-variant/30' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/dashboard">
                  <span class="material-symbols-outlined text-[18px]">dashboard</span>
                  <span>Dashboard</span>
              </router-link>
              <router-link class="px-4 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-2"
                 :class="$route.name?.includes('invoice') ? 'text-on-surface bg-surface-variant/30' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/invoices">
                  <span class="material-symbols-outlined text-[18px]">receipt_long</span>
                  <span>Invoices</span>
              </router-link>
              <router-link class="px-4 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-2"
                 :class="$route.name === 'clients' ? 'text-on-surface bg-surface-variant/30' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/clients">
                  <span class="material-symbols-outlined text-[18px]">groups</span>
                  <span>Clients</span>
              </router-link>
              <router-link class="px-4 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-2"
                 :class="$route.name === 'finance' ? 'text-on-surface bg-surface-variant/30' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/finance">
                  <span class="material-symbols-outlined text-[18px]">account_balance_wallet</span>
                  <span>Finance</span>
              </router-link>
          </nav>

          <!-- Right: User Controls -->
          <div class="flex items-center justify-end gap-3 w-1/4" v-if="authStore.isAuthenticated">
              <!-- Credits Pill -->
              <router-link to="/user/credits/history" class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 rounded-full font-label text-xs font-bold hover:bg-emerald-500/25 transition-colors" title="Available Credits">
                  <span class="material-symbols-outlined text-[16px]">hexagon</span>
                  <span>{{ authStore.credits || 0 }} Credits</span>
              </router-link>
              
              <!-- Profile Pill -->
              <router-link
                  to="/user/profile/setup"
                  class="flex items-center gap-2 py-1 pl-3 pr-1 rounded-full bg-surface-container-high border border-outline-variant/30 hover:bg-surface-variant/60 transition-all text-xs font-semibold text-on-surface shadow-sm"
                  :title="authStore.user?.name"
              >
                  <span class="hidden sm:inline font-body text-xs">{{ authStore.user?.name || 'My Account' }}</span>
                  <div class="w-7 h-7 rounded-full bg-emerald-700 text-white font-bold flex items-center justify-center uppercase text-xs shadow-sm">
                      {{ authStore.user?.name ? authStore.user.name[0] : 'U' }}
                  </div>
              </router-link>

              <!-- Logout Button -->
              <div class="pl-1 flex items-center border-l border-outline-variant/30 ml-1">
                  <button @click="handleLogout" type="button" class="text-on-surface-variant hover:text-error transition-colors p-1.5 rounded-lg hover:bg-error-container/30 cursor-pointer" title="Sign Out">
                      <span class="material-symbols-outlined text-[20px]">logout</span>
                  </button>
              </div>
              
              <!-- Mobile Hamburger Button -->
              <button type="button" @click="mobileMenuOpen = !mobileMenuOpen" class="md:hidden text-on-surface-variant hover:text-primary p-2 rounded-xl transition-colors cursor-pointer" title="Toggle Navigation Menu">
                  <span class="material-symbols-outlined text-[24px]">menu</span>
              </button>
          </div>

          <!-- Unauthenticated Controls -->
          <div class="flex items-center justify-end gap-3 w-1/4" v-else>
              <router-link to="/user/login" class="text-sm font-bold text-on-surface-variant hover:text-primary transition-colors hidden sm:block">Log in</router-link>
              <router-link to="/user/register" class="px-4 py-2 bg-primary text-on-primary rounded-xl text-sm font-bold shadow hover:shadow-md hover:-translate-y-0.5 transition-all">Get Started</router-link>
          </div>
      </div>
      <!-- Mobile Dropdown Navigation Drawer -->
      <div v-show="mobileMenuOpen" class="md:hidden border-t border-outline-variant/50 bg-surface/95 backdrop-blur-xl px-4 py-3 space-y-2 font-label text-sm font-medium animate-fade-in">
          <router-link to="/user/dashboard" @click="mobileMenuOpen = false" class="block px-3 py-2 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">Dashboard</router-link>
          <router-link to="/user/invoices" @click="mobileMenuOpen = false" class="block px-3 py-2 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">Invoices</router-link>
          <router-link to="/user/clients" @click="mobileMenuOpen = false" class="block px-3 py-2 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">Clients</router-link>
          <router-link to="/user/finance" @click="mobileMenuOpen = false" class="block px-3 py-2 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">Finance Tracker</router-link>
      </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const mobileMenuOpen = ref(false)

async function handleLogout() {
  await authStore.logout()
  router.push('/user/login')
}
</script>
