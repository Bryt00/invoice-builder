<template>
  <div class="min-h-[85vh] flex items-center justify-center p-4 relative overflow-hidden">
    <!-- Ambient Glow Background -->
    <div class="fixed inset-0 pointer-events-none overflow-hidden z-0">
      <div class="absolute top-[-10%] left-[-10%] w-[45%] h-[45%] bg-amber-500/10 rounded-full blur-3xl"></div>
      <div class="absolute bottom-[-10%] right-[-10%] w-[45%] h-[45%] bg-primary/10 rounded-full blur-3xl"></div>
    </div>

    <!-- Login Container -->
    <main class="w-full max-w-md relative z-10 my-8">
      <!-- Header -->
      <div class="text-center mb-8 flex flex-col items-center">
        <div class="w-16 h-16 bg-amber-500/15 border border-amber-500/30 rounded-2xl flex items-center justify-center mb-4 text-amber-500 shadow-lg shadow-amber-500/10">
          <span class="material-symbols-outlined text-4xl">shield</span>
        </div>
        <h1 class="text-3xl font-extrabold font-headline text-on-surface tracking-tight">Admin Console</h1>
        <p class="text-on-surface-variant mt-2 text-sm font-body">Enterprise Root Management Portal</p>
      </div>

      <!-- Login Card -->
      <div class="glass-card rounded-2xl p-8 shadow-xl border border-outline-variant/60 relative overflow-hidden space-y-6">
        <!-- Top Accent Gradient -->
        <div class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-amber-500/30 via-amber-500 to-amber-500/30"></div>

        <!-- Error Alert -->
        <div v-if="error" class="p-4 rounded-xl bg-error-container/80 border border-error/40 text-on-error-container text-xs font-body font-medium flex items-center gap-2.5">
          <span class="material-symbols-outlined text-error text-[20px]">error</span>
          <span>{{ error }}</span>
        </div>

        <form @submit.prevent="handleAdminLogin" class="space-y-5 font-body">
          <!-- Email Input -->
          <div class="space-y-1.5">
            <label class="block text-xs font-label font-bold text-on-surface uppercase tracking-wider" for="admin-email">Administrator Email</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-on-surface-variant">
                <span class="material-symbols-outlined text-[20px]">mail</span>
              </div>
              <input
                v-model="email"
                type="email"
                id="admin-email"
                placeholder="admin@example.com"
                required
                autocomplete="email"
                class="block w-full pl-10 pr-3.5 py-3 bg-surface-container-low border border-outline-variant/60 rounded-xl text-on-surface placeholder-on-surface-variant/50 focus:outline-none focus:border-amber-500 text-sm font-medium transition-colors"
              />
            </div>
          </div>

          <!-- Password Input -->
          <div class="space-y-1.5">
            <label class="block text-xs font-label font-bold text-on-surface uppercase tracking-wider" for="admin-password">Secure Key / Password</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-on-surface-variant">
                <span class="material-symbols-outlined text-[20px]">lock</span>
              </div>
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                id="admin-password"
                placeholder="••••••••"
                required
                autocomplete="current-password"
                class="block w-full pl-10 pr-10 py-3 bg-surface-container-low border border-outline-variant/60 rounded-xl text-on-surface placeholder-on-surface-variant/50 focus:outline-none focus:border-amber-500 text-sm font-medium transition-colors"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-on-surface-variant hover:text-on-surface"
              >
                <span class="material-symbols-outlined text-[20px]">
                  {{ showPassword ? 'visibility_off' : 'visibility' }}
                </span>
              </button>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full flex justify-center items-center gap-2 py-3.5 px-4 rounded-xl text-sm font-bold text-on-primary bg-amber-500 hover:bg-amber-600 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 transition-all duration-200 shadow-md hover:shadow-lg disabled:opacity-50"
          >
            <span class="material-symbols-outlined text-[18px]">verified_user</span>
            <span>{{ loading ? 'Authenticating...' : 'Authenticate Admin' }}</span>
          </button>
        </form>
      </div>

      <!-- Back to Standard User Login -->
      <div class="mt-8 text-center">
        <router-link to="/user/login" class="inline-flex items-center justify-center text-xs font-label font-bold text-on-surface-variant hover:text-amber-500 transition-colors group">
          <span class="material-symbols-outlined text-[18px] mr-1.5 group-hover:-translate-x-1 transition-transform">arrow_back</span>
          <span>Return to standard user sign in</span>
        </router-link>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

const authStore = useAuthStore()
const router = useRouter()

async function handleAdminLogin() {
  loading.value = true
  error.value = ''

  try {
    await authStore.login(email.value, password.value)
    if (authStore.user?.role?.name === 'Admin') {
      router.push('/user/admin/dashboard')
    } else {
      error.value = 'Access Denied: Administrator credentials required.'
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Authentication failed'
  } finally {
    loading.value = false
  }
}
</script>
