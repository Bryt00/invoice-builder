<template>
  <div>
    <!-- Header Section -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
            <h2 class="font-headline text-3xl font-bold text-on-surface mb-1">Invoices Directory</h2>
            <p class="font-body text-base text-on-surface-variant">View, track, and manage all your generated invoices.</p>
        </div>
        <router-link to="/user/invoices/new" class="btn-auth-submit text-on-primary rounded-xl px-5 py-2.5 font-label text-sm font-semibold flex items-center gap-2 shrink-0 self-start md:self-auto bg-primary shadow-md hover:bg-on-primary-fixed-variant transition-all">
            <span class="material-symbols-outlined text-[20px]">add</span>
            Create Invoice
        </router-link>
    </header>

    <!-- Search & Filter Controls -->
    <div class="mb-4 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
        <div class="relative flex-1 max-w-md">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">search</span>
            <input type="text" v-model="searchQuery" @input="debouncedSearch"
                   placeholder="Search invoices by number or client..." 
                   class="w-full pl-9 pr-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface focus:outline-none focus:border-primary">
        </div>
    </div>

    <!-- Invoices Directory Table Container -->
    <div class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden">
        <div class="px-6 py-4 border-b border-outline-variant/40 flex flex-wrap justify-between items-center bg-surface-container-lowest/50 gap-4">
            <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-[20px]">receipt_long</span>
                All Dispatches
            </h3>
            <!-- Interactive Status Filter Tabs -->
            <div class="flex items-center gap-1 bg-surface-container-high/60 p-1 rounded-xl border border-outline-variant/40 text-xs font-label font-semibold">
                <button type="button" @click="setStatusFilter('all')" :class="['px-3 py-1.5 rounded-lg transition-all cursor-pointer', statusFilter === 'all' ? 'bg-surface text-primary shadow-xs' : 'text-on-surface-variant hover:text-on-surface']">All</button>
                <button type="button" @click="setStatusFilter('paid')" :class="['px-3 py-1.5 rounded-lg transition-all cursor-pointer', statusFilter === 'paid' ? 'bg-surface text-primary shadow-xs' : 'text-on-surface-variant hover:text-on-surface']">Paid</button>
                <button type="button" @click="setStatusFilter('sent')" :class="['px-3 py-1.5 rounded-lg transition-all cursor-pointer', statusFilter === 'sent' ? 'bg-surface text-primary shadow-xs' : 'text-on-surface-variant hover:text-on-surface']">Sent</button>
                <button type="button" @click="setStatusFilter('draft')" :class="['px-3 py-1.5 rounded-lg transition-all cursor-pointer', statusFilter === 'draft' ? 'bg-surface text-primary shadow-xs' : 'text-on-surface-variant hover:text-on-surface']">Drafts</button>
                <button type="button" @click="setStatusFilter('overdue')" :class="['px-3 py-1.5 rounded-lg transition-all cursor-pointer', statusFilter === 'overdue' ? 'bg-surface text-primary shadow-xs' : 'text-on-surface-variant hover:text-on-surface']">Overdue</button>
            </div>
        </div>
        
        <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse min-w-[800px]">
                <thead>
                    <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
                        <th class="px-6 py-3.5">Invoice #</th>
                        <th class="px-6 py-3.5">Client</th>
                        <th class="px-6 py-3.5">Issue Date</th>
                        <th class="px-6 py-3.5">Due Date</th>
                        <th class="px-6 py-3.5">Total Amount</th>
                        <th class="px-6 py-3.5">Status</th>
                        <th class="px-6 py-3.5 text-right">Actions</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
                    <tr v-if="loading">
                        <td colspan="7" class="px-6 py-12 text-center text-outline text-sm">
                            Loading invoices...
                        </td>
                    </tr>
                    <template v-else-if="invoices.length > 0">
                        <tr v-for="inv in invoices" :key="inv.id" class="hover:bg-surface-container-low/40 transition-colors">
                            <td class="px-6 py-4 font-semibold text-primary">
                                <router-link :to="`/user/invoices/view?id=${inv.id}`" class="hover:underline">{{ inv.invoice_number }}</router-link>
                            </td>
                            <td class="px-6 py-4 font-medium text-on-surface">
                                {{ inv.client?.name || 'Direct Client' }}
                            </td>
                            <td class="px-6 py-4 text-on-surface-variant">{{ formatDate(inv.issue_date || inv.created_at) }}</td>
                            <td class="px-6 py-4 text-on-surface-variant">{{ formatDate(inv.due_date) }}</td>
                            <td class="px-6 py-4 font-bold text-on-surface">
                                {{ inv.currency }} {{ inv.total?.toFixed(2) }}
                                <span class="ml-1 px-1.5 py-0.5 rounded text-[10px] font-semibold bg-surface-container-high text-on-surface-variant border border-outline-variant/40">{{ inv.currency }}</span>
                            </td>
                            <td class="px-6 py-4">
                                <span v-if="inv.status === 'paid'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 inline-flex items-center gap-1">
                                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Paid
                                </span>
                                <span v-else-if="inv.status === 'sent'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border border-indigo-500/20 inline-flex items-center gap-1">
                                    <span class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span> Sent
                                </span>
                                <span v-else-if="inv.status === 'overdue'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-rose-500/10 text-rose-600 dark:text-rose-400 border border-rose-500/20 inline-flex items-center gap-1">
                                    <span class="w-1.5 h-1.5 rounded-full bg-rose-500"></span> Overdue
                                </span>
                                <span v-else class="px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20 inline-flex items-center gap-1">
                                    <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span> Draft
                                </span>
                            </td>
                            <td class="px-6 py-4 text-right flex items-center justify-end gap-2">
                                <router-link v-if="inv.status === 'draft'" :to="`/user/invoices/edit?id=${inv.id}`" class="p-1.5 text-primary hover:bg-primary-container/30 transition-colors rounded-lg" title="Edit Draft">
                                    <span class="material-symbols-outlined text-[18px]">edit</span>
                                </router-link>
                                <router-link :to="`/user/invoices/view?id=${inv.id}`" class="p-1.5 text-outline hover:text-primary transition-colors rounded-lg hover:bg-surface-container-high" title="View Preview">
                                    <span class="material-symbols-outlined text-[18px]">visibility</span>
                                </router-link>
                                <a v-if="inv.public_token" :href="`/invoice/public/${inv.public_token}`" target="_blank" class="p-1.5 text-outline hover:text-primary transition-colors rounded-lg hover:bg-surface-container-high" title="Open Public Share Link">
                                    <span class="material-symbols-outlined text-[18px]">open_in_new</span>
                                </a>
                            </td>
                        </tr>
                    </template>
                    <tr v-else>
                        <td colspan="7" class="px-6 py-12 text-center text-on-surface-variant">
                            <span class="material-symbols-outlined text-4xl text-outline mb-2">description</span>
                            <p class="font-body text-base">No invoices found matching criteria.</p>
                            <p class="font-body text-xs text-outline mt-1">Try clearing search filters or click "Create Invoice" above.</p>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../../utils/api'

const router = useRouter()
const route = useRoute()

const invoices = ref([])
const loading = ref(true)
const searchQuery = ref(route.query.q || '')
const statusFilter = ref(route.query.status || 'all')

let searchTimeout = null

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
}

async function fetchInvoices() {
    loading.value = true
    try {
        const params = new URLSearchParams()
        if (searchQuery.value) params.append('q', searchQuery.value)
        if (statusFilter.value && statusFilter.value !== 'all') params.append('status', statusFilter.value)
        
        // Update URL to match filters without full reload
        router.replace({ query: { q: searchQuery.value || undefined, status: statusFilter.value !== 'all' ? statusFilter.value : undefined } })

        const res = await api.get(`/invoices?${params.toString()}`)
        if (res.data?.invoices) {
            invoices.value = res.data.invoices
        } else {
            invoices.value = []
        }
    } catch (err) {
        console.error('Failed to fetch invoices:', err)
        invoices.value = []
    } finally {
        loading.value = false
    }
}

function debouncedSearch() {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
        fetchInvoices()
    }, 300)
}

function setStatusFilter(status) {
    statusFilter.value = status
    fetchInvoices()
}

onMounted(() => {
    fetchInvoices()
})
</script>
