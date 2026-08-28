<template>
  <div class="min-h-screen bg-background p-4 sm:p-8 flex items-center justify-center">
    <div v-if="loading" class="text-outline text-sm">Loading invoice...</div>

    <div v-else-if="invoice" class="relative glass-card max-w-3xl w-full rounded-2xl p-6 sm:p-10 space-y-6 z-0 overflow-hidden">
      
      <!-- Watermark Backdrop -->
      <div class="absolute inset-0 z-[-1] flex items-center justify-center pointer-events-none select-none opacity-[0.03]">
          <img v-if="profile?.logo_url" :src="profile.logo_url" class="w-3/4 h-3/4 object-contain grayscale" alt="">
          <span v-else class="material-symbols-outlined text-[400px]">receipt_long</span>
      </div>

      <div class="flex justify-between items-start border-b border-outline-variant/40 pb-6">
        <div>
          <h1 class="brand-title text-2xl font-bold text-on-surface">Invoice #{{ invoice.invoice_number }}</h1>
          <p class="text-xs text-outline mt-1">Issued: {{ new Date(invoice.created_at).toLocaleDateString() }}</p>
        </div>
        <a :href="`/api/v1/invoices/public/download?token=${token}`" target="_blank" class="px-4 py-2 bg-primary text-on-primary font-semibold text-xs rounded-xl shadow-md">
          Download PDF
        </a>
      </div>

      <div class="grid grid-cols-2 gap-4 text-sm">
        <div>
          <h3 class="text-xs uppercase text-outline font-semibold">From</h3>
          <p class="font-bold text-on-surface mt-1">{{ profile?.company_name }}</p>
          <p class="text-on-surface-variant whitespace-pre-line">{{ profile?.address }}</p>
        </div>
        <div>
          <h3 class="text-xs uppercase text-outline font-semibold">Bill To</h3>
          <p class="font-bold text-on-surface mt-1">{{ invoice.client?.name }}</p>
          <p class="text-on-surface-variant">{{ invoice.client?.company }}</p>
        </div>
      </div>

      <table class="w-full text-left text-sm mt-6">
        <thead>
          <tr class="border-b border-outline-variant/40 text-xs text-outline uppercase">
            <th class="pb-2">Description</th>
            <th class="pb-2">Qty</th>
            <th class="pb-2">Price</th>
            <th class="pb-2 text-right">Amount</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-outline-variant/20">
          <tr v-for="item in (invoice.items || invoice.line_items || [])" :key="item.id">
            <td class="py-3 font-medium text-on-surface">{{ item.description }}</td>
            <td class="py-3 text-outline">{{ item.quantity }}</td>
            <td class="py-3 text-outline">{{ invoice.currency }} {{ item.unit_price?.toFixed(2) }}</td>
            <td class="py-3 font-semibold text-right text-on-surface">{{ invoice.currency }} {{ item.amount?.toFixed(2) }}</td>
          </tr>
        </tbody>
      </table>

      <div class="border-t border-outline-variant/40 pt-4 flex justify-between items-center text-lg font-bold">
        <span>Total Amount Due:</span>
        <span class="text-primary">{{ invoice.currency }} {{ invoice.total?.toFixed(2) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/utils/api'
import type { BusinessProfile } from '@/types/user'

const route = useRoute()
const invoice = ref<any>(null)
const profile = ref<BusinessProfile | null>(null)
const loading = ref(true)
const token = ref<string>('')

onMounted(async () => {
  token.value = String(route.params.token || route.query.token || '')
  if (!token.value) return

  try {
    const res = await api.get(`/invoices/public?token=${token.value}`)
    if (res.data?.invoice) invoice.value = res.data.invoice
    if (res.data?.profile) profile.value = res.data.profile
  } finally {
    loading.value = false
  }
})
</script>
