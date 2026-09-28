<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <!-- Page Header & Filters -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">webhook</span>
          Webhook Delivery Logs
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">Monitor webhook events from Stripe and other external integrations.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
        <form @submit.prevent="fetchWebhooks(1)" class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
          <select v-model="statusFilter" @change="fetchWebhooks(1)" class="w-full sm:w-auto px-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface font-semibold focus:outline-none focus:border-amber-500 cursor-pointer">
            <option value="">All Statuses</option>
            <option value="processed">Processed</option>
            <option value="failed">Failed</option>
            <option value="pending">Pending</option>
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
          Webhook Events ({{ meta.total_count || 0 }})
        </h3>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
              <th class="px-6 py-3.5">Event Type</th>
              <th class="px-6 py-3.5">Provider</th>
              <th class="px-6 py-3.5">Status</th>
              <th class="px-6 py-3.5">Received At</th>
              <th class="px-6 py-3.5 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
            <tr v-for="log in webhooks" :key="log.id" class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-mono text-xs font-semibold text-on-surface">{{ log.event_type }}</td>
              <td class="px-6 py-4">
                <span class="inline-flex items-center gap-1.5 px-2 py-1 bg-surface-container-high rounded-md text-xs font-semibold text-on-surface uppercase tracking-wider">
                  {{ log.provider }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span :class="getStatusBadge(log.status)">
                  {{ log.status.charAt(0).toUpperCase() + log.status.slice(1) }}
                </span>
                <p v-if="log.error_message" class="text-[10px] text-rose-500 mt-1 max-w-[200px] truncate" :title="log.error_message">{{ log.error_message }}</p>
              </td>
              <td class="px-6 py-4 font-mono text-xs text-on-surface-variant">{{ formatDate(log.created_at) }}</td>
              <td class="px-6 py-4 text-right">
                <button v-if="log.status === 'failed' || log.status === 'pending'" @click="replayWebhook(log.id)" class="p-1.5 text-amber-600 hover:text-amber-500 hover:bg-amber-500/10 rounded-lg transition-colors cursor-pointer" title="Replay Webhook">
                  <span class="material-symbols-outlined text-[18px]">replay</span>
                </button>
              </td>
            </tr>
            <tr v-if="!loading && webhooks.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-on-surface-variant">No webhook events found.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="meta.total_count > meta.limit" class="px-6 py-4 border-t border-outline-variant/40 flex items-center justify-between bg-surface-container-lowest/30">
        <span class="text-sm text-on-surface-variant">
          Showing {{ ((meta.page - 1) * meta.limit) + 1 }} to {{ Math.min(meta.page * meta.limit, meta.total_count) }} of {{ meta.total_count }} events
        </span>
        <div class="flex items-center gap-2">
          <button @click="fetchWebhooks(meta.page - 1)" :disabled="meta.page === 1" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Previous</button>
          <button @click="fetchWebhooks(meta.page + 1)" :disabled="meta.page * meta.limit >= meta.total_count" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
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

const webhooks = ref([])
const meta = ref({ page: 1, limit: 15, total_count: 0 })
const loading = ref(true)
const statusFilter = ref('')

onMounted(() => {
  fetchWebhooks(1)
})

async function fetchWebhooks(page = 1) {
  loading.value = true
  try {
    const res = await api.get(`/admin/webhooks?page=${page}&limit=${meta.value.limit}&status=${encodeURIComponent(statusFilter.value)}`)
    webhooks.value = res.data.webhooks || []
    meta.value = res.data.meta
  } catch (err) {
    showFlash('Failed to load webhook logs', 'error')
  } finally {
    loading.value = false
  }
}

async function replayWebhook(id) {
  try {
    await api.post(`/admin/webhooks/${id}/replay`)
    showFlash('Webhook queued for replay', 'success')
    fetchWebhooks(meta.value.page)
  } catch (err) {
    showFlash('Failed to replay webhook', 'error')
  }
}

function formatDate(date) {
  return dayjs(date).format('MMM DD, YYYY HH:mm:ss')
}

function getStatusBadge(status) {
  const base = 'px-2.5 py-1 rounded-full text-[10px] uppercase tracking-wider font-bold border '
  switch(status.toLowerCase()) {
    case 'processed':
      return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
    case 'pending':
      return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
    case 'failed':
      return base + 'bg-rose-500/10 text-rose-600 border-rose-500/20'
    default:
      return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
  }
}
</script>
