<template>
  <div class="space-y-6">
    <!-- Header -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
            <h2 class="font-headline text-3xl font-bold text-on-surface">Credit Transaction History</h2>
            <p class="font-body text-base text-on-surface-variant">Full ledger of credit purchases, admin grants, and invoice deductions.</p>
        </div>
        <router-link to="/user/dashboard" class="btn-auth-submit bg-primary text-on-primary rounded-xl px-5 py-2.5 font-label text-sm font-semibold flex items-center gap-2 shrink-0 self-start md:self-auto hover:bg-on-primary-fixed-variant transition-colors shadow-sm">
            <span class="material-symbols-outlined text-[20px]">add_card</span>
            Top Up Credits
        </router-link>
    </header>

    <!-- Stats Row -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div class="glass-card p-5 rounded-2xl border border-outline-variant/60">
            <p class="font-label text-xs text-on-surface-variant font-medium uppercase tracking-wider mb-1">Available Credits</p>
            <h3 class="font-headline text-3xl font-extrabold text-primary">{{ stats.balance || 0 }}</h3>
        </div>
        <div class="glass-card p-5 rounded-2xl border border-outline-variant/60">
            <p class="font-label text-xs text-on-surface-variant font-medium uppercase tracking-wider mb-1">Total Purchased</p>
            <h3 class="font-headline text-3xl font-bold text-on-surface">{{ stats.total_purchased || 0 }}</h3>
        </div>
        <div class="glass-card p-5 rounded-2xl border border-outline-variant/60">
            <p class="font-label text-xs text-on-surface-variant font-medium uppercase tracking-wider mb-1">Total Consumed</p>
            <h3 class="font-headline text-3xl font-bold text-on-surface-variant">{{ stats.total_used || 0 }}</h3>
        </div>
    </div>

    <!-- Ledger Table Card -->
    <div class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden">
        <div class="px-6 py-4 border-b border-outline-variant/40 bg-surface-container-lowest/50 flex items-center justify-between">
            <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-[20px]">receipt_long</span>
                Transaction Log
            </h3>
            <span class="font-body text-xs text-on-surface-variant">Showing latest activity</span>
        </div>

        <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse min-w-[600px]">
                <thead>
                    <tr class="border-b border-outline-variant/40 bg-surface-container-low/40 text-xs font-label font-bold text-on-surface-variant uppercase tracking-wider">
                        <th class="py-3.5 px-6">Date & Time</th>
                        <th class="py-3.5 px-6">Type</th>
                        <th class="py-3.5 px-6">Description</th>
                        <th class="py-3.5 px-6 text-right">Credits</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
                    <tr v-if="loading">
                        <td colspan="4" class="py-12 text-center text-outline text-sm">Loading transaction history...</td>
                    </tr>
                    <template v-else-if="history.length > 0">
                        <tr v-for="item in history" :key="item.id" class="hover:bg-surface-container-low/40 transition-colors">
                            <td class="py-4 px-6 text-xs text-on-surface-variant whitespace-nowrap">{{ formatDate(item.created_at) }}</td>
                            <td class="py-4 px-6 whitespace-nowrap">
                                <span v-if="item.type === 'purchase'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-primary-container/40 text-primary border border-primary/20">Purchase</span>
                                <span v-else-if="item.type === 'grant'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-tertiary-container/40 text-tertiary border border-tertiary/20">Bonus / Grant</span>
                                <span v-else class="px-2.5 py-1 rounded-full text-xs font-semibold bg-surface-container-high text-on-surface-variant border border-outline-variant/40">Usage</span>
                            </td>
                            <td class="py-4 px-6 font-medium">{{ item.description || item.reason || 'Transaction' }}</td>
                            <td :class="['py-4 px-6 text-right font-headline font-bold', item.amount > 0 ? 'text-primary' : 'text-on-surface-variant']">
                                {{ item.amount > 0 ? '+' : '' }}{{ item.amount }}
                            </td>
                        </tr>
                    </template>
                    <tr v-else>
                        <td colspan="4" class="py-12 text-center text-on-surface-variant text-xs">
                            <span class="material-symbols-outlined text-4xl text-outline mb-2">history</span>
                            <p class="font-body text-base font-semibold">No transaction history found.</p>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'

const { showFlash } = useFlash()
const stats = ref({ balance: 0, total_purchased: 0, total_used: 0 })
const history = ref([])
const loading = ref(true)

onMounted(async () => {
    loading.value = true
    try {
        const [balRes, histRes] = await Promise.all([
            api.get('/credits/balance'),
            api.get('/credits/history'),
        ])

        if (balRes.data?.stats) {
            stats.value = balRes.data.stats
        } else if (balRes.data?.balance !== undefined) {
            // Fallback if API returns just balance
            stats.value.balance = balRes.data.balance
        }
        
        if (histRes.data?.history) {
            history.value = histRes.data.history
        } else if (histRes.data?.transactions) {
            history.value = histRes.data.transactions
        }
    } catch (err: any) {
        showFlash('Failed to load credit history', 'error')
    } finally {
        loading.value = false
    }
})

function formatDate(dateStr: any) {
    if (!dateStr) return '-'
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    
    // Format like "Jan 02, 2006 • 15:04"
    const optsDate: Intl.DateTimeFormatOptions = { month: 'short', day: '2-digit', year: 'numeric' }
    const optsTime: Intl.DateTimeFormatOptions = { hour: '2-digit', minute: '2-digit', hour12: false }
    
    return `${d.toLocaleDateString('en-US', optsDate)} • ${d.toLocaleTimeString('en-US', optsTime)}`
}
</script>
