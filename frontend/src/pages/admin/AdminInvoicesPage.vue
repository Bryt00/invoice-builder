<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <!-- Page Header & Filters -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">receipt_long</span>
          Global Invoice Ledger
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">View and monitor all invoices generated across the platform.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
        <form @submit.prevent="fetchInvoices(1)" class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
          <div class="relative w-full sm:w-48">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">search</span>
            <input type="text" v-model="searchQuery" placeholder="Search invoices..." class="w-full pl-9 pr-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface focus:outline-none focus:border-amber-500">
          </div>
          <select v-model="statusFilter" @change="fetchInvoices(1)" class="w-full sm:w-auto px-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface font-semibold focus:outline-none focus:border-amber-500 cursor-pointer">
            <option value="">All Statuses</option>
            <option value="draft">Draft</option>
            <option value="pending">Pending</option>
            <option value="paid">Paid</option>
            <option value="overdue">Overdue</option>
          </select>
          <button type="submit" class="bg-amber-500 hover:bg-amber-600 text-on-primary px-4 py-2 rounded-xl font-label text-sm font-bold transition-colors cursor-pointer shrink-0">Search</button>
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
          Invoice Records ({{ meta.total_count || 0 }})
        </h3>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
              <th class="px-6 py-3.5">Invoice ID</th>
              <th class="px-6 py-3.5">Owner / User</th>
              <th class="px-6 py-3.5">Client</th>
              <th class="px-6 py-3.5">Amount</th>
              <th class="px-6 py-3.5">Status</th>
              <th class="px-6 py-3.5">Issued Date</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
            <tr v-for="inv in invoices" :key="inv.id" class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-mono font-semibold text-on-surface">#{{ inv.invoice_number }}</td>
              <td class="px-6 py-4 text-on-surface-variant">
                <router-link :to="`/user/admin/users/${inv.user_id}`" class="hover:text-amber-500 transition-colors" title="View User Profile">
                  {{ inv.user_id }} <!-- We ideally want user name, but we only have user_id on invoice in MVP unless we joined. We display ID for now -->
                </router-link>
              </td>
              <td class="px-6 py-4 font-semibold text-on-surface">{{ inv.client_name }}</td>
              <td class="px-6 py-4 font-bold text-on-surface">{{ currencySymbol(inv.currency) }}{{ (inv.total_amount / 100).toFixed(2) }}</td>
              <td class="px-6 py-4">
                <span :class="getStatusBadge(inv.status)">
                  {{ inv.status.charAt(0).toUpperCase() + inv.status.slice(1) }}
                </span>
              </td>
              <td class="px-6 py-4 text-on-surface-variant">{{ formatDate(inv.issue_date) }}</td>
            </tr>
            <tr v-if="!loading && invoices.length === 0">
              <td colspan="6" class="px-6 py-12 text-center text-on-surface-variant">No invoices found matching your criteria.</td>
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
          <button @click="fetchInvoices(meta.page - 1)" :disabled="meta.page === 1" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Previous</button>
          <button @click="fetchInvoices(meta.page + 1)" :disabled="meta.page * meta.limit >= meta.total_count" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'
import dayjs from 'dayjs'

const { showFlash } = useFlash()

const invoices = ref([])
const meta = ref({ page: 1, limit: 10, total_count: 0 })
const loading = ref(true)
const searchQuery = ref('')
const statusFilter = ref('')

onMounted(() => {
  fetchInvoices(1)
})

const currencyMap = {
  USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', NGN: '₦'
}
function currencySymbol(code) {
  return currencyMap[code] || code
}

async function fetchInvoices(page = 1) {
  loading.value = true
  try {
    const res = await api.get(`/admin/invoices?page=${page}&limit=${meta.value.limit}&search=${encodeURIComponent(searchQuery.value)}&status=${encodeURIComponent(statusFilter.value)}`)
    invoices.value = res.data.invoices || []
    meta.value = res.data.meta
  } catch (err) {
    showFlash('Failed to load invoices', 'error')
  } finally {
    loading.value = false
  }
}

function formatDate(date) {
  if (!date || date === '0001-01-01T00:00:00Z') return 'N/A'
  return dayjs(date).format('MMM DD, YYYY')
}

function getStatusBadge(status) {
  const base = 'px-2.5 py-1 rounded-full text-xs font-semibold border '
  switch(status.toLowerCase()) {
    case 'paid':
      return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
    case 'pending':
      return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
    case 'overdue':
      return base + 'bg-rose-500/10 text-rose-600 border-rose-500/20'
    case 'draft':
    default:
      return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
  }
}
</script>
