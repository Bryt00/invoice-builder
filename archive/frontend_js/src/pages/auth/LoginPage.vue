<template>
  <main class="w-full max-w-[420px] mx-auto px-4 sm:px-6 z-10 my-auto">
    <!-- Header -->
    <div class="text-center mb-8">
      <div class="w-24 h-24 mx-auto mb-4 rounded-3xl bg-surface-container-lowest/90 border border-outline-variant/50 p-3 shadow-lg flex items-center justify-center">
        <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-full h-full object-contain" />
      </div>
      <h1 class="brand-title text-4xl font-bold mb-2 text-on-surface">Teks-Invoice</h1>
      <p class="font-body text-base text-on-surface-variant">Welcome back. Sign in to continue.</p>
    </div>

    <!-- Auth Card -->
    <div class="glass-card rounded-2xl p-8 sm:p-10">
      <FlashAlert />

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div>
          <label class="block font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider mb-2">
            Email Address
          </label>
          <div class="relative">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline text-lg">mail</span>
            <input
              v-model="form.email"
              type="email"
              required
              placeholder="name@company.com"
              class="w-full pl-10 pr-4 py-3 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-sm font-medium text-on-surface focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
            />
          </div>
        </div>

        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="block font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">
              Password
            </label>
            <router-link to="/user/forgot-password" class="text-xs font-semibold text-primary hover:underline">
              Forgot?
            </router-link>
          </div>
          <div class="relative">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline text-lg">lock</span>
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder="••••••••"
              class="w-full pl-10 pr-12 py-3 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-sm font-medium text-on-surface focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-outline text-lg hover:text-on-surface-variant transition-colors focus:outline-none flex items-center justify-center p-1 rounded-md"
              title="Toggle password visibility"
            >
              <span class="material-symbols-outlined">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
            </button>
          </div>
        </div>

        <button
          type="submit"
          :disabled="authStore.loading"
          class="w-full flex items-center justify-center gap-2 py-3 px-4 rounded-xl font-label text-sm font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 shadow-md"
        >
          <span v-if="!authStore.loading">Sign In</span>
          <span v-else>Signing in...</span>
          <span class="material-symbols-outlined text-[18px]">arrow_forward</span>
        </button>
      </form>

      <!-- Footer Links -->
      <div class="mt-8 text-center space-y-2.5">
        <p class="font-body text-xs text-on-surface-variant">
          Don't have an account?
          <router-link to="/user/register" class="font-semibold text-primary hover:underline ml-1">
            Sign Up
          </router-link>
        </p>
        <p class="font-body text-xs text-on-surface-variant">
          Didn't receive an activation email?
          <router-link to="/user/resend-verification" class="font-medium text-primary hover:underline ml-1">
            Resend Link
          </router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useFlash } from '../../composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const { showFlash } = useFlash()

const form = reactive({
  email: '',
  password: '',
})

const showPassword = ref(false)

async function handleSubmit() {
  try {
    await authStore.login(form.email, form.password)
    showFlash('Welcome back!', 'success')
    
    if (authStore.user?.role?.name === 'Admin') {
      router.push('/user/admin/dashboard')
    } else {
      const redirectPath = route.query.redirect || '/user/dashboard'
      router.push(redirectPath)
    }
  } catch (err) {
    showFlash(authStore.error || 'Failed to sign in', 'error')
  }
}
</script>
