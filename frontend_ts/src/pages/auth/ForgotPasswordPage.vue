<template>
  <main class="w-full max-w-[420px] mx-auto px-4 sm:px-6 z-10 my-auto">
    <div class="text-center mb-8">
      <div class="w-24 h-24 mx-auto mb-4 rounded-3xl bg-surface-container-lowest/90 border border-outline-variant/50 p-3 shadow-lg flex items-center justify-center">
        <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-full h-full object-contain" />
      </div>
      <h1 class="brand-title text-4xl font-bold mb-2 text-on-surface">Teks-Invoice</h1>
      <p class="font-body text-base text-on-surface-variant">Enter your email to receive a password reset link.</p>
    </div>

    <div class="glass-card rounded-2xl p-8 sm:p-10">
      <FlashAlert />

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider mb-2">
            Email Address
          </label>
          <input
            v-model="email"
            type="email"
            required
            placeholder="name@company.com"
            class="w-full px-4 py-3 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-sm font-medium text-on-surface focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 px-4 rounded-xl font-label text-sm font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all disabled:opacity-50 shadow-md"
        >
          <span v-if="!loading">Send Reset Link</span>
          <span v-else>Sending...</span>
        </button>
      </form>

      <div class="mt-6 text-center">
        <router-link to="/user/login" class="font-label text-sm font-semibold text-primary hover:underline">
          Back to Sign In
        </router-link>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'

const email = ref('')
const loading = ref(false)
const { showFlash } = useFlash()

async function handleSubmit() {
  loading.value = true
  try {
    const res = await api.post('/auth/forgot-password', { email: email.value })
    showFlash(res.data?.message || 'Password reset link sent to your email.', 'success', 8000)
  } catch (err: any) {
    showFlash('Failed to send reset link', 'error')
  } finally {
    loading.value = false
  }
}
</script>
