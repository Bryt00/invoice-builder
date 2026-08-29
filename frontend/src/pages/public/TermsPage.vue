<template>
  <div class="max-w-3xl mx-auto px-4 py-12 space-y-6">
    <router-link to="/" class="text-sm font-semibold text-primary hover:underline">&larr; Back to Home</router-link>
    <h1 class="brand-title text-3xl font-bold">Terms of Service</h1>
    <div class="glass-card rounded-2xl p-6 sm:p-8 space-y-4 text-sm text-on-surface-variant leading-relaxed">
      <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
        {{ customContent }}
      </div>
      <template v-else>
        <p>Welcome to Teks-Invoice. By using our service, you agree to these terms.</p>
        <p>1. Account Terms: You must provide a valid email address and maintain account security.</p>
        <p>2. Payment & Credits: Credits purchased are non-refundable unless specified by law.</p>
        <p>3. Service Availability: We strive for 99.9% uptime but do not guarantee uninterrupted access.</p>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../utils/api'

const customContent = ref('')

onMounted(async () => {
  try {
    const res = await api.get('/public/settings')
    if (res.data?.settings?.legal_terms) {
      customContent.value = res.data.settings.legal_terms
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
