<template>
  <div class="max-w-3xl mx-auto px-4 py-12 space-y-6">
    <router-link to="/" class="text-sm font-semibold text-primary hover:underline">&larr; Back to Home</router-link>
    <h1 class="brand-title text-3xl font-bold text-on-surface">Refund Policy</h1>
    <div class="glass-card rounded-2xl p-6 sm:p-8 space-y-4 text-sm text-on-surface-variant leading-relaxed">
      <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
        {{ customContent }}
      </div>
      <template v-else>
        <p class="font-medium text-base text-on-surface">Thank you for using Teks-Invoice.</p>
        
        <h2 class="font-bold text-on-surface pt-2 text-base">1. Credit Packages & Digital Purchases</h2>
        <p>All credit purchases made on Teks-Invoice provide instant digital credits to your account balance for invoice finalization, PDF exports, and receipts.</p>

        <h2 class="font-bold text-on-surface pt-2 text-base">2. Refund Eligibility</h2>
        <p>Unused credit balances may be requested for a refund within 14 days of purchase if no credits from the package have been consumed. Once credits have been utilized for PDF generation or dispatch, that portion of the transaction is non-refundable.</p>

        <h2 class="font-bold text-on-surface pt-2 text-base">3. Service Disruption & System Errors</h2>
        <p>If a credit was deducted due to a failed generation or system error, please contact our support team. Verified system failures will result in instant credit reimbursement to your account.</p>

        <h2 class="font-bold text-on-surface pt-2 text-base">4. Requesting a Refund</h2>
        <p>To request a refund or credit adjustment, please reach out to our support team with your payment reference at <router-link to="/contact" class="text-primary hover:underline font-semibold">Contact Support</router-link>.</p>
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
    if (res.data?.settings?.legal_refund) {
      customContent.value = res.data.settings.legal_refund
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
