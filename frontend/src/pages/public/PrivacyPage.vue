<template>
  <div class="max-w-3xl mx-auto px-4 py-12 space-y-6">
    <router-link to="/" class="text-sm font-semibold text-primary hover:underline">&larr; Back to Home</router-link>
    <h1 class="brand-title text-3xl font-bold">Privacy Policy</h1>
    <div class="glass-card rounded-2xl p-6 sm:p-8 space-y-4 text-sm text-on-surface-variant leading-relaxed">
      <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
        {{ customContent }}
      </div>
      <template v-else>
        <p>Your privacy is important to us. We collect minimal information required to provide invoicing services.</p>
        <p>1. Data Collection: We collect account email, business profile info, and invoice metadata.</p>
        <p>2. Data Use: Your financial data is confidential and never sold to third parties.</p>
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
    if (res.data?.settings?.legal_privacy) {
      customContent.value = res.data.settings.legal_privacy
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
