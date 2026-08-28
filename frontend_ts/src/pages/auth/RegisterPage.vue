<template>
  <main class="w-full max-w-[420px] mx-auto px-4 sm:px-6 z-10 my-auto">
    <!-- Header -->
    <div class="text-center mb-8">
      <div class="w-24 h-24 mx-auto mb-4 rounded-3xl bg-surface-container-lowest/90 border border-outline-variant/50 p-3 shadow-lg flex items-center justify-center">
        <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-full h-full object-contain" />
      </div>
      <h1 class="brand-title text-4xl font-bold mb-2 text-on-surface">Teks-Invoice</h1>
      <p class="font-body text-base text-on-surface-variant">Create your account to start invoicing.</p>
    </div>

    <!-- Auth Card -->
    <div class="glass-card rounded-2xl p-8 sm:p-10">
      <FlashAlert />

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider mb-1">
            Full Name
          </label>
          <input
            v-model="form.name"
            type="text"
            required
            placeholder="John Doe"
            class="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-sm font-medium text-on-surface focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>

        <div>
          <label class="block font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider mb-1">
            Email Address
          </label>
          <input
            v-model="form.email"
            type="email"
            required
            placeholder="name@company.com"
            class="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-sm font-medium text-on-surface focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>

        <div>
          <label class="block font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider mb-1">
            Password
          </label>
          <div class="relative">
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder="••••••••"
              class="w-full px-4 pr-12 py-2.5 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-sm font-medium text-on-surface focus:outline-none focus:ring-2 focus:ring-primary"
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

        <div>
          <label class="block font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider mb-1">
            Confirm Password
          </label>
          <div class="relative">
            <input
              v-model="form.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              required
              placeholder="••••••••"
              class="w-full px-4 pr-12 py-2.5 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-sm font-medium text-on-surface focus:outline-none focus:ring-2 focus:ring-primary"
            />
            <button
              type="button"
              @click="showConfirmPassword = !showConfirmPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-outline text-lg hover:text-on-surface-variant transition-colors focus:outline-none flex items-center justify-center p-1 rounded-md"
              title="Toggle password visibility"
            >
              <span class="material-symbols-outlined">{{ showConfirmPassword ? 'visibility_off' : 'visibility' }}</span>
            </button>
          </div>
        </div>

        <button
          type="submit"
          :disabled="authStore.loading"
          class="w-full py-3 px-4 mt-2 rounded-xl font-label text-sm font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all disabled:opacity-50 shadow-md"
        >
          <span v-if="!authStore.loading">Create Account</span>
          <span v-else>Creating Account...</span>
        </button>
      </form>

      <div class="mt-6 text-center">
        <p class="font-body text-sm text-on-surface-variant">
          Already have an account?
          <router-link to="/user/login" class="font-label text-sm font-semibold text-primary hover:underline ml-1">
            Sign In
          </router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'

const authStore = useAuthStore()
const router = useRouter()
const { showFlash } = useFlash()

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const showPassword = ref(false)
const showConfirmPassword = ref(false)

async function handleSubmit() {
  if (form.password !== form.confirmPassword) {
    showFlash('Passwords do not match', 'error')
    return
  }

  try {
    const res = await authStore.register(form.name, form.email, form.password)
    showFlash(res.message || 'Account created! Please check your email for activation link.', 'success', 8000)
    router.push('/user/login')
  } catch (err: any) {
    showFlash(authStore.error || 'Registration failed', 'error')
  }
}
</script>
