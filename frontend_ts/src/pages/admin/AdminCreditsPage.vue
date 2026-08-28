<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <!-- Page Header & Filters -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">toll</span>
          Credits History Ledger
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">Audit log of all credit purchases, manual grants, and usage deductions.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
        <form @submit.prevent="fetchCredits(1)" class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
          <select v-model="typeFilter" @change="fetchCredits(1)" class="w-full sm:w-auto px-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface font-semibold focus:outline-none focus:border-amber-500 cursor-pointer">
            <option value="">All Transactions</option>
            <option value="purchase">Purchases (Top-ups)</option>
            <option value="deduction">Deductions (Usage)</option>
            <option value="grant">Manual Grants</option>
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
          Credit Transactions ({{ meta.total_count || 0 }})
        </h3>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
              <th class="px-6 py-3.5">User</th>
              <th class="px-6 py-3.5">Type</th>
              <th class="px-6 py-3.5">Amount</th>
              <th class="px-6 py-3.5">Reason / Package</th>
              <th class="px-6 py-3.5">Date</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
            <tr v-for="txn in credits" :key="txn.id" class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-semibold text-on-surface">
                <router-link :to="`/user/admin/users/${txn.user_id}`" class="hover:text-amber-500 transition-colors">
                  {{ txn.user_id }}
                </router-link>
              </td>
              <td class="px-6 py-4">
                <span :class="getTypeBadge(txn.type)">
                  {{ txn.type ? txn.type.charAt(0).toUpperCase() + txn.type.slice(1) : 'Unknown' }}
                </span>
              </td>
              <td class="px-6 py-4 font-bold" :class="txn.amount > 0 ? 'text-emerald-600' : 'text-rose-600'">
                {{ txn.amount > 0 ? '+' : '' }}{{ txn.amount }}
              </td>
              <td class="px-6 py-4 text-on-surface-variant">{{ txn.description || 'N/A' }}</td>
              <td class="px-6 py-4 text-on-surface-variant">{{ formatDate(txn.created_at) }}</td>
            </tr>
            <tr v-if="!loading && credits.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-on-surface-variant">No credit transactions found.</td>
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
          <button @click="fetchCredits(meta.page - 1)" :disabled="meta.page === 1" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Previous</button>
          <button @click="fetchCredits(meta.page + 1)" :disabled="meta.page * meta.limit >= meta.total_count" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'
import dayjs from 'dayjs'

const { showFlash } = useFlash()

const credits = ref([])
const meta = ref({ page: 1, limit: 10, total_count: 0 })
const loading = ref(true)
const typeFilter = ref('')

onMounted(() => {
  fetchCredits(1)
})

async function fetchCredits(page = 1) {
  loading.value = true
  try {
    const res = await api.get(`/admin/credits?page=${page}&limit=${meta.value.limit}&type=${encodeURIComponent(typeFilter.value)}`)
    credits.value = res.data.credits || []
    meta.value = res.data.meta
  } catch (err: any) {
    showFlash('Failed to load credit history', 'error')
  } finally {
    loading.value = false
  }
}

function formatDate(date) {
  return dayjs(date).format('MMM DD, YYYY HH:mm')
}

function getTypeBadge(type) {
  const base = 'px-2.5 py-1 rounded-full text-xs font-semibold border '
  if (!type) return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
  switch(type.toLowerCase()) {
    case 'purchase':
    case 'grant':
      return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
    case 'deduction':
      return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
    default:
      return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
  }
}
</script>
