<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <!-- Page Header & Filters -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">payments</span>
          Global Payments Ledger
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">Monitor all platform transactions, top-ups, and subscription payments.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
        <form @submit.prevent="fetchPayments(1)" class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
          <select v-model="statusFilter" @change="fetchPayments(1)" class="w-full sm:w-auto px-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface font-semibold focus:outline-none focus:border-amber-500 cursor-pointer">
            <option value="">All Statuses</option>
            <option value="successful">Successful</option>
            <option value="pending">Pending</option>
            <option value="failed">Failed</option>
            <option value="refunded">Refunded</option>
          </select>
          <button type="submit" class="bg-amber-500 hover:bg-amber-600 text-on-primary px-4 py-2 rounded-xl font-label text-sm font-bold transition-colors cursor-pointer shrink-0">Filter</button>
        </form>
      </div>
    </div>

    <!-- Data Table -->
    <div class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden shadow-xs relative min-h-[300px]">
      
      <!-- Loading Overlay -->
      <div v-if="loading" class="absolute inset-0 bg-surface/50 backdrop-blur-sm z-10 flex items-center justify-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
      </div>

      <div class="px-6 py-4 border-b border-outline-variant/40 flex justify-between items-center bg-surface-container-lowest/50">
        <h3 class="font-headline text-base sm:text-lg font-bold text-on-surface">
          Payment Transactions ({{ meta.total_count || 0 }})
        </h3>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
              <th class="px-6 py-3.5">Txn Ref</th>
              <th class="px-6 py-3.5">User / Customer</th>
              <th class="px-6 py-3.5">Amount</th>
              <th class="px-6 py-3.5">Gateway</th>
              <th class="px-6 py-3.5">Status</th>
              <th class="px-6 py-3.5">Date</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
            <tr v-for="pay in payments" :key="pay.id" class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-mono font-semibold text-on-surface text-xs truncate max-w-[150px]" :title="pay.transaction_ref || pay.reference || pay.id">
                {{ pay.transaction_ref || pay.reference || (pay.id ? pay.id.substring(0, 8) : 'N/A') }}
              </td>
              <td class="px-6 py-4 text-on-surface-variant">
                <router-link :to="`/user/admin/users/${pay.user_id}`" class="hover:text-amber-500 transition-colors font-mono text-xs">
                  {{ pay.user_id ? pay.user_id.substring(0, 8) + '...' : 'N/A' }}
                </router-link>
              </td>
              <td class="px-6 py-4 font-bold text-on-surface">{{ currencySymbol(pay.currency) }}{{ Number(pay.amount || 0).toFixed(2) }}</td>
              <td class="px-6 py-4">
                <span class="inline-flex items-center gap-1.5 px-2 py-1 bg-surface-container-high rounded-md text-xs font-semibold text-on-surface uppercase">
                  <span class="material-symbols-outlined text-[14px] text-amber-500">{{ pay.payment_method === 'stripe' ? 'credit_card' : 'account_balance' }}</span>
                  {{ pay.payment_method || 'Paystack' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span :class="getStatusBadge(pay.status)">
                  {{ pay.status ? (pay.status.charAt(0).toUpperCase() + pay.status.slice(1)) : 'Pending' }}
                </span>
              </td>
              <td class="px-6 py-4 text-on-surface-variant">{{ formatDate(pay.created_at) }}</td>
            </tr>
            <tr v-if="!loading && payments.length === 0">
              <td colspan="6" class="px-6 py-12 text-center text-on-surface-variant">No payments found matching your criteria.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="meta.total_count > meta.limit" class="px-6 py-4 border-t border-outline-variant/40 flex items-center justify-between bg-surface-container-lowest/30">
        <span class="text-sm text-on-surface-variant">
          Showing {{ ((meta.page - 1) * meta.limit) + 1 }} to {{ Math.min(meta.page * meta.limit, meta.total_count) }} of {{ meta.total_count }} records
        </span>
        <div class="flex items-center gap-2">
          <button @click="fetchPayments(meta.page - 1)" :disabled="meta.page === 1" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Previous</button>
          <button @click="fetchPayments(meta.page + 1)" :disabled="meta.page * meta.limit >= meta.total_count" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useToast } from '@/composables/useToast'
import dayjs from 'dayjs'
import type { Payment } from '@/types/payment'

const { showToast } = useToast()

const payments = ref<Payment[]>([])
const meta = ref({ page: 1, limit: 10, total_count: 0 })
const loading = ref(true)
const statusFilter = ref('')

onMounted(() => {
  fetchPayments(1)
})

const currencyMap: Record<string, string> = {
  USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', NGN: '₦'
}
function currencySymbol(code: string) {
  return currencyMap[code] || code
}

async function fetchPayments(page = 1) {
  loading.value = true
  try {
    const res = await api.get(`/admin/payments?page=${page}&limit=${meta.value.limit}&status=${encodeURIComponent(statusFilter.value)}`)
    payments.value = res.data.payments || []
    meta.value = res.data.meta || { page: 1, limit: 10, total_count: payments.value.length }
  } catch (err: any) {
    showToast('Failed to load payments', 'error')
  } finally {
    loading.value = false
  }
}

function formatDate(date?: string) {
  if (!date) return '-'
  return dayjs(date).format('MMM DD, YYYY HH:mm')
}

function getStatusBadge(status?: string) {
  const base = 'px-2.5 py-1 rounded-full text-xs font-semibold border '
  if (!status) return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
  switch(String(status).toLowerCase()) {
    case 'successful':
    case 'succeeded':
    case 'completed':
    case 'paid':
      return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
    case 'pending':
      return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
    case 'failed':
      return base + 'bg-rose-500/10 text-rose-600 border-rose-500/20'
    case 'refunded':
      return base + 'bg-purple-500/10 text-purple-600 border-purple-500/20'
    default:
      return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
  }
}
</script>
