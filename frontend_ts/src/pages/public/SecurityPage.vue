<template>
  <div class="max-w-3xl mx-auto px-4 py-12 space-y-6">
    <router-link to="/" class="text-sm font-semibold text-primary hover:underline">&larr; Back to Home</router-link>
    <h1 class="brand-title text-3xl font-bold text-on-surface">Security & Compliance</h1>
    <div class="glass-card rounded-2xl p-6 sm:p-8 space-y-4 text-sm text-on-surface-variant leading-relaxed">
      <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
        {{ customContent }}
      </div>
      <template v-else>
        <p class="font-medium text-base text-on-surface">Security is at the heart of Teks-Invoice.</p>
        
        <h2 class="font-bold text-on-surface pt-2 text-base">1. Data Encryption & Transport</h2>
        <p>All data sent between your browser and our servers is encrypted using industry-standard TLS 1.3 encryption. Data at rest is encrypted in isolated, secure PostgreSQL databases.</p>

        <h2 class="font-bold text-on-surface pt-2 text-base">2. Payment & Credit Processing</h2>
        <p>Teks-Invoice does not store your credit card or mobile money credentials. All financial transactions are processed securely through certified payment gateways (Paystack).</p>

        <h2 class="font-bold text-on-surface pt-2 text-base">3. Authentication & Access Control</h2>
        <p>Accounts use secure JWT token authentication with bcrypt password hashing. Role-based access control guarantees that your invoices, client directory, and financial data are strictly accessible to your authenticated session.</p>

        <h2 class="font-bold text-on-surface pt-2 text-base">4. Security Disclosures</h2>
        <p>If you discover a vulnerability or security issue, please contact our engineering team immediately at <router-link to="/contact" class="text-primary hover:underline font-semibold">Contact Support</router-link>.</p>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const customContent = ref('')

onMounted(async () => {
  try {
    const res = await api.get('/public/settings')
    if (res.data?.settings?.legal_security) {
      customContent.value = res.data.settings.legal_security
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
