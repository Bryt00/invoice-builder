<template>
  <header class="bg-background/95 backdrop-blur-md w-full sticky top-0 z-50 border-b border-outline-variant/30 shadow-xs transition-all">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex justify-between items-center relative">
          <!-- Left: Logo & Brand -->
          <div class="flex items-center shrink-0">
              <router-link :to="authStore.isAuthenticated ? '/user/dashboard' : '/'" class="flex items-center gap-2.5 sm:gap-3 group">
                  <div class="w-10 h-10 sm:w-11 sm:h-11 rounded-xl bg-surface-container-low/80 border border-outline-variant/40 flex items-center justify-center overflow-hidden shadow-xs group-hover:border-primary/50 transition-all shrink-0 p-1.5">
                      <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-full h-full object-contain group-hover:scale-110 transition-transform">
                  </div>
                  <span class="font-headline text-lg sm:text-xl font-extrabold text-on-surface whitespace-nowrap tracking-tight">Teks-Invoice</span>
              </router-link>
          </div>
          
          <!-- Center: Navigation (Authenticated) -->
          <nav class="hidden md:flex items-center justify-center gap-1.5 flex-1" v-if="authStore.isAuthenticated">
              <router-link class="px-3.5 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-1.5"
                 :class="$route.name === 'dashboard' ? 'text-on-surface bg-surface-variant/40' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/dashboard">
                  <span class="material-symbols-outlined text-[18px]">dashboard</span>
                  <span>Dashboard</span>
              </router-link>
              <router-link class="px-3.5 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-1.5"
                 :class="String($route.name || '').includes('invoice') ? 'text-on-surface bg-surface-variant/40' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/invoices">
                  <span class="material-symbols-outlined text-[18px]">receipt_long</span>
                  <span>Invoices</span>
              </router-link>
              <router-link class="px-3.5 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-1.5"
                 :class="$route.name === 'clients' ? 'text-on-surface bg-surface-variant/40' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/clients">
                  <span class="material-symbols-outlined text-[18px]">groups</span>
                  <span>Clients</span>
              </router-link>
              <router-link class="px-3.5 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all flex items-center gap-1.5"
                 :class="$route.name === 'finance' ? 'text-on-surface bg-surface-variant/40' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-variant/20'"
                 to="/user/finance">
                  <span class="material-symbols-outlined text-[18px]">account_balance_wallet</span>
                  <span>Finance</span>
              </router-link>
          </nav>

          <!-- Center: Navigation (Unauthenticated Public) -->
          <nav class="hidden md:flex items-center justify-center gap-6 flex-1 text-sm font-semibold text-on-surface-variant" v-else>
              <a href="/#features" class="hover:text-primary transition-colors">Features</a>
              <a href="/#how-it-works" class="hover:text-primary transition-colors">How It Works</a>
              <a href="/#pricing" class="hover:text-primary transition-colors">Pricing</a>
              <router-link to="/faq" class="hover:text-primary transition-colors">FAQ</router-link>
              <router-link to="/contact" class="hover:text-primary transition-colors">Contact</router-link>
          </nav>

          <!-- Right: User Controls (Authenticated) -->
          <div class="flex items-center justify-end gap-2 sm:gap-3 shrink-0" v-if="authStore.isAuthenticated">
              <!-- Theme Toggle Button -->
              <button
                  type="button"
                  @click="toggleTheme"
                  class="w-9 h-9 rounded-xl flex items-center justify-center text-on-surface-variant hover:text-primary hover:bg-surface-variant/40 transition-colors cursor-pointer"
                  :title="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
                  aria-label="Toggle theme"
              >
                  <span class="material-symbols-outlined text-[20px] transition-transform duration-300" :class="isDark ? 'rotate-180 text-amber-400' : 'text-slate-600'">
                      {{ isDark ? 'light_mode' : 'dark_mode' }}
                  </span>
              </button>

              <!-- Credits Pill -->
              <router-link to="/user/credits/history" class="inline-flex items-center gap-1 px-2.5 sm:px-3 py-1.5 bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 rounded-full font-label text-xs font-bold hover:bg-emerald-500/25 transition-colors" title="Available Credits">
                  <span class="material-symbols-outlined text-[16px]">hexagon</span>
                  <span>{{ authStore.credits || 0 }} <span class="hidden sm:inline">Credits</span></span>
              </router-link>
              
              <!-- Profile Pill -->
              <router-link
                  to="/user/profile/setup"
                  class="flex items-center gap-2 py-1 pl-2 sm:pl-3 pr-1 rounded-full bg-surface-container-high border border-outline-variant/30 hover:bg-surface-variant/60 transition-all text-xs font-semibold text-on-surface shadow-xs"
                  :title="authStore.user?.name"
              >
                  <span class="hidden sm:inline font-body text-xs">{{ authStore.user?.name || 'My Account' }}</span>
                  <div class="w-7 h-7 rounded-full bg-emerald-700 text-white font-bold flex items-center justify-center uppercase text-xs shadow-xs">
                      {{ authStore.user?.name ? authStore.user.name[0] : 'U' }}
                  </div>
              </router-link>

              <!-- Logout Button (Desktop) -->
              <div class="hidden sm:flex pl-1 items-center border-l border-outline-variant/30 ml-1">
                  <button @click="handleLogout" type="button" class="text-on-surface-variant hover:text-error transition-colors p-1.5 rounded-lg hover:bg-error-container/30 cursor-pointer" title="Sign Out">
                      <span class="material-symbols-outlined text-[20px]">logout</span>
                  </button>
              </div>
              
              <!-- Mobile Hamburger Button -->
              <button type="button" @click="mobileMenuOpen = !mobileMenuOpen" class="md:hidden text-on-surface-variant hover:text-primary p-1.5 rounded-xl transition-colors cursor-pointer" :title="mobileMenuOpen ? 'Close Menu' : 'Open Navigation Menu'">
                  <span class="material-symbols-outlined text-[24px]">{{ mobileMenuOpen ? 'close' : 'menu' }}</span>
              </button>
          </div>

          <!-- Right: Unauthenticated Controls -->
          <div class="flex items-center justify-end gap-2 sm:gap-3 shrink-0" v-else>
              <!-- Theme Toggle Button -->
              <button
                  type="button"
                  @click="toggleTheme"
                  class="w-9 h-9 rounded-xl flex items-center justify-center text-on-surface-variant hover:text-primary hover:bg-surface-variant/40 transition-colors cursor-pointer"
                  :title="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
                  aria-label="Toggle theme"
              >
                  <span class="material-symbols-outlined text-[20px] transition-transform duration-300" :class="isDark ? 'rotate-180 text-amber-400' : 'text-slate-600'">
                      {{ isDark ? 'light_mode' : 'dark_mode' }}
                  </span>
              </button>

              <router-link to="/user/login" class="text-xs sm:text-sm font-bold text-on-surface-variant hover:text-primary transition-colors px-2 py-1.5">Log in</router-link>
              <router-link to="/user/register" class="px-3 sm:px-4 py-2 bg-primary text-on-primary rounded-xl text-xs sm:text-sm font-bold shadow hover:shadow-md hover:-translate-y-0.5 active:translate-y-0 transition-all shrink-0">Get Started</router-link>
              
              <!-- Mobile Hamburger Button for Visitors -->
              <button type="button" @click="mobileMenuOpen = !mobileMenuOpen" class="md:hidden text-on-surface-variant hover:text-primary p-1.5 rounded-xl transition-colors cursor-pointer ml-0.5" :title="mobileMenuOpen ? 'Close Menu' : 'Open Navigation Menu'">
                  <span class="material-symbols-outlined text-[24px]">{{ mobileMenuOpen ? 'close' : 'menu' }}</span>
              </button>
          </div>
      </div>

      <!-- Mobile Dropdown Navigation Drawer (Authenticated) -->
      <div v-show="mobileMenuOpen && authStore.isAuthenticated" class="md:hidden border-t border-outline-variant/30 bg-background/98 backdrop-blur-xl px-4 py-3 space-y-1.5 font-label text-sm font-medium shadow-xl">
          <router-link to="/user/dashboard" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px]">dashboard</span>
              <span>Dashboard</span>
          </router-link>
          <router-link to="/user/invoices" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px]">receipt_long</span>
              <span>Invoices</span>
          </router-link>
          <router-link to="/user/clients" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px]">groups</span>
              <span>Clients</span>
          </router-link>
          <router-link to="/user/finance" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px]">account_balance_wallet</span>
              <span>Finance Tracker</span>
          </router-link>

          <!-- Mobile Theme Switcher Row -->
          <div class="pt-2 border-t border-outline-variant/30 flex items-center justify-between px-3.5 py-2">
              <span class="text-xs font-semibold text-on-surface-variant flex items-center gap-2">
                  <span class="material-symbols-outlined text-[18px]">{{ isDark ? 'dark_mode' : 'light_mode' }}</span>
                  <span>Theme ({{ isDark ? 'Dark' : 'Light' }})</span>
              </span>
              <button
                  type="button"
                  @click="toggleTheme"
                  class="px-2.5 py-1 text-xs rounded-lg font-bold border border-outline-variant/50 bg-surface-container-low text-on-surface flex items-center gap-1.5 shadow-2xs hover:bg-surface-variant/40"
              >
                  <span class="material-symbols-outlined text-[15px]">{{ isDark ? 'light_mode' : 'dark_mode' }}</span>
                  <span>Switch</span>
              </button>
          </div>

          <div class="pt-2 border-t border-outline-variant/30 flex items-center justify-between px-3">
              <router-link to="/user/profile/setup" @click="mobileMenuOpen = false" class="text-xs font-semibold text-on-surface-variant hover:text-primary flex items-center gap-1.5">
                  <span class="material-symbols-outlined text-[18px]">settings</span>
                  <span>Profile Settings</span>
              </router-link>
              <button @click="handleLogout" type="button" class="text-xs font-semibold text-error hover:underline flex items-center gap-1">
                  <span class="material-symbols-outlined text-[18px]">logout</span>
                  <span>Sign Out</span>
              </button>
          </div>
      </div>

      <!-- Mobile Dropdown Navigation Drawer (Unauthenticated) -->
      <div v-show="mobileMenuOpen && !authStore.isAuthenticated" class="md:hidden border-t border-outline-variant/30 bg-background/98 backdrop-blur-xl px-4 py-4 space-y-1.5 font-label text-sm font-medium shadow-xl">
          <a href="/#features" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px] text-primary">auto_awesome</span>
              <span>Features</span>
          </a>
          <a href="/#how-it-works" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px] text-primary">filter_3</span>
              <span>How It Works</span>
          </a>
          <a href="/#pricing" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px] text-primary">payments</span>
              <span>Pricing & Credits</span>
          </a>
          <router-link to="/faq" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px] text-primary">help_outline</span>
              <span>FAQ & Help</span>
          </router-link>
          <router-link to="/contact" @click="mobileMenuOpen = false" class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-on-surface-variant hover:text-primary hover:bg-surface-container-low transition-colors">
              <span class="material-symbols-outlined text-[20px] text-primary">mail</span>
              <span>Contact Support</span>
          </router-link>

          <!-- Mobile Theme Switcher Row (Visitor) -->
          <div class="pt-2 border-t border-outline-variant/30 flex items-center justify-between px-3.5 py-2">
              <span class="text-xs font-semibold text-on-surface-variant flex items-center gap-2">
                  <span class="material-symbols-outlined text-[18px]">{{ isDark ? 'dark_mode' : 'light_mode' }}</span>
                  <span>Theme ({{ isDark ? 'Dark' : 'Light' }})</span>
              </span>
              <button
                  type="button"
                  @click="toggleTheme"
                  class="px-2.5 py-1 text-xs rounded-lg font-bold border border-outline-variant/50 bg-surface-container-low text-on-surface flex items-center gap-1.5 shadow-2xs hover:bg-surface-variant/40"
              >
                  <span class="material-symbols-outlined text-[15px]">{{ isDark ? 'light_mode' : 'dark_mode' }}</span>
                  <span>Switch</span>
              </button>
          </div>

          <div class="pt-3 border-t border-outline-variant/30 grid grid-cols-2 gap-2 mt-2">
              <router-link to="/user/login" @click="mobileMenuOpen = false" class="w-full py-2.5 text-center rounded-xl border border-outline-variant/60 font-semibold text-on-surface hover:bg-surface-container-low transition-colors">
                  Log In
              </router-link>
              <router-link to="/user/register" @click="mobileMenuOpen = false" class="w-full py-2.5 text-center rounded-xl bg-primary text-on-primary font-semibold shadow hover:shadow-md transition-all">
                  Get Started
              </router-link>
          </div>
      </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useTheme } from '@/composables/useTheme'

const authStore = useAuthStore()
const router = useRouter()
const mobileMenuOpen = ref(false)
const { isDark, toggleTheme } = useTheme()

async function handleLogout() {
  await authStore.logout()
  router.push('/user/login')
}
</script>
