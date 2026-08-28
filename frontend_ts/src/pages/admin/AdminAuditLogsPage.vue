<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <!-- Page Header & Filters -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">history</span>
          System Audit Logs
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">Review security events, logins, and administrative actions.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
        <form @submit.prevent="fetchLogs(1)" class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
          <div class="relative w-full sm:w-64">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">search</span>
            <input type="text" v-model="searchQuery" placeholder="Search by action or IP..." class="w-full pl-9 pr-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface focus:outline-none focus:border-amber-500">
          </div>
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
          Event Logs ({{ meta.total_count || 0 }})
        </h3>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
              <th class="px-6 py-3.5">Timestamp</th>
              <th class="px-6 py-3.5">User</th>
              <th class="px-6 py-3.5">Action</th>
              <th class="px-6 py-3.5">IP Address</th>
              <th class="px-6 py-3.5">User Agent</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
            <tr v-for="log in logs" :key="log.id" class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-mono text-xs text-on-surface-variant whitespace-nowrap">{{ formatDate(log.created_at) }}</td>
              <td class="px-6 py-4 font-semibold text-on-surface">
                <router-link v-if="log.user_id" :to="`/user/admin/users/${log.user_id}`" class="hover:text-amber-500 transition-colors">
                  {{ log.user_id }}
                </router-link>
                <span v-else class="text-on-surface-variant italic">System / Unknown</span>
              </td>
              <td class="px-6 py-4">
                <span class="inline-flex items-center gap-1.5 px-2 py-1 bg-surface-container-high rounded-md text-xs font-semibold text-on-surface">
                  {{ log.action }}
                </span>
              </td>
              <td class="px-6 py-4 font-mono text-xs text-on-surface-variant">{{ log.ip_address }}</td>
              <td class="px-6 py-4 text-xs text-on-surface-variant truncate max-w-[200px]" :title="log.user_agent">{{ log.user_agent }}</td>
            </tr>
            <tr v-if="!loading && logs.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-on-surface-variant">No audit logs found.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="meta.total_count > meta.limit" class="px-6 py-4 border-t border-outline-variant/40 flex items-center justify-between bg-surface-container-lowest/30">
        <span class="text-sm text-on-surface-variant">
          Showing {{ ((meta.page - 1) * meta.limit) + 1 }} to {{ Math.min(meta.page * meta.limit, meta.total_count) }} of {{ meta.total_count }} logs
        </span>
        <div class="flex items-center gap-2">
          <button @click="fetchLogs(meta.page - 1)" :disabled="meta.page === 1" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Previous</button>
          <button @click="fetchLogs(meta.page + 1)" :disabled="meta.page * meta.limit >= meta.total_count" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
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

const logs = ref([])
const meta = ref({ page: 1, limit: 15, total_count: 0 })
const loading = ref(true)
const searchQuery = ref('')

onMounted(() => {
  fetchLogs(1)
})

async function fetchLogs(page = 1) {
  loading.value = true
  try {
    const res = await api.get(`/admin/audit-logs?page=${page}&limit=${meta.value.limit}&search=${encodeURIComponent(searchQuery.value)}`)
    logs.value = res.data.logs || []
    meta.value = res.data.meta
  } catch (err: any) {
    showFlash('Failed to load audit logs', 'error')
  } finally {
    loading.value = false
  }
}

function formatDate(date) {
  return dayjs(date).format('MMM DD, YYYY HH:mm:ss')
}
</script>
