<template>
  <main class="w-full max-w-[420px] mx-auto px-4 sm:px-6 z-10 my-auto text-center">
    <div class="glass-card rounded-2xl p-8 sm:p-10">
      <div v-if="loading" class="space-y-4">
        <span class="material-symbols-outlined text-4xl text-primary animate-spin">refresh</span>
        <h2 class="text-xl font-bold text-on-surface">Activating your account...</h2>
        <p class="text-sm text-on-surface-variant">Please wait while we verify your activation link.</p>
      </div>

      <div v-else-if="success" class="space-y-4">
        <span class="material-symbols-outlined text-5xl text-primary">check_circle</span>
        <h2 class="text-2xl font-bold text-on-surface">Account Activated!</h2>
        <p class="text-sm text-on-surface-variant">Your email has been verified. You can now sign in to your account.</p>
        <router-link
          to="/user/login"
          class="inline-block w-full py-3 px-4 rounded-xl font-label text-sm font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all shadow-md mt-4"
        >
          Sign In Now
        </router-link>
      </div>

      <div v-else class="space-y-4">
        <span class="material-symbols-outlined text-5xl text-error">error</span>
        <h2 class="text-2xl font-bold text-on-surface">Verification Failed</h2>
        <p class="text-sm text-error">{{ errorMessage }}</p>
        <router-link
          to="/user/resend-verification"
          class="inline-block w-full py-3 px-4 rounded-xl font-label text-sm font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all shadow-md mt-4"
        >
          Request New Activation Link
        </router-link>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/utils/api'

const route = useRoute()
const loading = ref(true)
const success = ref(false)
const errorMessage = ref('')

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    loading.value = false
    errorMessage.value = 'Invalid or missing activation token.'
    return
  }

  try {
    await api.get(`/auth/verify?token=${token}`)
    success.value = true
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'The activation link is invalid or has expired.'
  } finally {
    loading.value = false
  }
})
</script>
