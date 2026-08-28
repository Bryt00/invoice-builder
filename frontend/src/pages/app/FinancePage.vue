<template>
  <div class="space-y-6">
    <!-- Full-Page Bouncing Coming Soon Banner -->
    <div class="w-full h-[75vh] bg-primary-container/20 border border-primary/30 rounded-3xl p-8 flex flex-col items-center justify-center shadow-sm animate-bounce">
        <span class="material-symbols-outlined text-6xl text-primary mb-4">rocket_launch</span>
        <span class="font-headline text-3xl sm:text-5xl font-black text-primary tracking-widest uppercase text-center">
            COMING SOON
        </span>
        <p class="font-body text-on-surface-variant text-center mt-4 max-w-md">
            We are working on a powerful new financial tracking experience tightly integrated with your invoices. Stay tuned!
        </p>
    </div>

    <!-- Hidden Finance Content for Future -->
    <div v-if="false">

    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
            <h2 class="font-headline text-3xl font-bold text-on-surface flex items-center gap-3">
                <span class="material-symbols-outlined text-primary text-3xl">account_balance_wallet</span>
                Income Tracker
            </h2>
            <p class="font-body text-sm text-on-surface-variant mt-1">Track your paid invoices and generated revenue.</p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
            <button @click="exportCsv"
               class="bg-surface-container-high hover:bg-surface-container-highest text-on-surface transition-colors px-4 py-2.5 rounded-xl font-label text-xs sm:text-sm font-semibold flex items-center gap-2 border border-outline-variant/50 shadow-sm cursor-pointer">
                <span class="material-symbols-outlined text-[18px]">download</span>
                <span>Export CSV</span>
            </button>
        </div>
    </header>

    <!-- Bento Grid Summary Cards -->
    <div class="grid grid-cols-1 gap-6 max-w-sm">
        <!-- Total Income Card -->
        <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 flex items-center justify-between">
            <div>
                <p class="font-label text-xs font-semibold uppercase tracking-wider text-on-surface-variant mb-1 flex items-center gap-1.5">
                    <span class="material-symbols-outlined text-emerald-500 text-[18px]">trending_up</span>
                    Total Paid Income
                </p>
                <h3 class="font-headline text-3xl font-extrabold text-emerald-600 dark:text-emerald-400">
                    {{ currencySymbol }}{{ summary.total_income?.toFixed(2) || '0.00' }}
                </h3>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 flex items-center justify-center shrink-0">
                <span class="material-symbols-outlined text-2xl">check_circle</span>
            </div>
        </div>

    </div>

    <!-- Filter & Search Toolbar -->
    <div class="glass-card rounded-2xl p-4 border border-outline-variant/60 flex flex-col md:flex-row gap-4 items-center justify-between">
        <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
            <!-- Category Filter -->
            <select v-model="filters.category_id" class="bg-surface-container-high text-on-surface border border-outline-variant/60 text-xs sm:text-sm font-label rounded-xl px-3 py-2 focus:ring-2 focus:ring-primary outline-none">
                <option value="all">All Categories</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }} ({{ cat.type }})</option>
            </select>

            <!-- Start Date -->
            <input type="date" v-model="filters.start_date" class="bg-surface-container-high text-on-surface border border-outline-variant/60 text-xs sm:text-sm font-label rounded-xl px-3 py-2 focus:ring-2 focus:ring-primary outline-none">

            <!-- End Date -->
            <input type="date" v-model="filters.end_date" class="bg-surface-container-high text-on-surface border border-outline-variant/60 text-xs sm:text-sm font-label rounded-xl px-3 py-2 focus:ring-2 focus:ring-primary outline-none">
        </div>

        <!-- Search Input -->
        <div class="relative w-full md:w-64">
            <span class="material-symbols-outlined absolute left-3 top-2.5 text-on-surface-variant text-[18px]">search</span>
            <input type="text" v-model="filters.search" placeholder="Search transactions..."
                   class="w-full bg-surface-container-high text-on-surface pl-9 pr-3 py-2 border border-outline-variant/60 text-xs sm:text-sm font-body rounded-xl focus:ring-2 focus:ring-primary outline-none">
        </div>
    </div>

    <!-- Financial Transactions Table -->
    <div class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden shadow-sm">
        <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse min-w-[800px]">
                <thead>
                    <tr class="bg-surface-container-low border-b border-outline-variant/60 font-label text-xs uppercase tracking-wider text-on-surface-variant">
                        <th class="py-3.5 px-4 sm:px-6">Date</th>
                        <th class="py-3.5 px-4 sm:px-6">Title & Description</th>
                        <th class="py-3.5 px-4 sm:px-6">Category</th>
                        <th class="py-3.5 px-4 sm:px-6">Payee / Payer</th>
                        <th class="py-3.5 px-4 sm:px-6 text-right">Amount</th>
                        <th class="py-3.5 px-4 sm:px-6 text-center">Receipt</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-outline-variant/40 font-body text-sm text-on-surface">
                    <tr v-if="loading">
                        <td colspan="7" class="py-12 text-center text-outline text-sm">Loading transactions...</td>
                    </tr>
                    <template v-else-if="filteredTransactions.length > 0">
                        <tr v-for="txn in filteredTransactions" :key="txn.id" class="hover:bg-surface-container-low/60 transition-colors">
                            <td class="py-4 px-4 sm:px-6 whitespace-nowrap text-on-surface-variant font-label text-xs">
                                {{ formatDate(txn.transaction_date) }}
                            </td>
                            <td class="py-4 px-4 sm:px-6">
                                <p class="font-semibold text-on-surface">{{ txn.title }}</p>
                                <p v-if="txn.description" class="text-xs text-on-surface-variant line-clamp-1">{{ txn.description }}</p>
                            </td>
                            <td class="py-4 px-4 sm:px-6 whitespace-nowrap">
                                <span :class="['inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold border', txn.type === 'income' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' : 'bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/20']">
                                    <span class="material-symbols-outlined text-[14px]">{{ txn.category?.icon || 'payments' }}</span>
                                    <span>{{ txn.category?.name || 'Uncategorized' }}</span>
                                </span>
                            </td>
                            <td class="py-4 px-4 sm:px-6 text-on-surface-variant">
                                {{ txn.payee_or_payer || '&mdash;' }}
                            </td>
                            <td :class="['py-4 px-4 sm:px-6 text-right whitespace-nowrap font-headline font-bold', txn.type === 'income' ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400']">
                                {{ txn.type === 'income' ? '+' : '-' }}{{ txn.currency || currencySymbol }}{{ txn.amount?.toFixed(2) }}
                            </td>
                            <td class="py-4 px-4 sm:px-6 text-center whitespace-nowrap">
                                <a v-if="txn.receipt_url" :href="txn.receipt_url" target="_blank" class="inline-flex items-center gap-1 text-xs font-semibold text-primary hover:underline" title="View Receipt">
                                    <span class="material-symbols-outlined text-[16px]">receipt</span>
                                    <span>View</span>
                                </a>
                                <span v-else class="text-xs text-on-surface-variant/50">&mdash;</span>
                            </td>
                        </tr>
                    </template>
                    <tr v-else>
                        <td colspan="6" class="py-12 text-center text-on-surface-variant">
                            <span class="material-symbols-outlined text-4xl text-outline mb-2">account_balance_wallet</span>
                            <p class="font-label text-base font-semibold">No income transactions found.</p>
                            <p class="font-body text-xs text-on-surface-variant/80 mt-1">When an invoice is marked as paid, it will automatically appear here.</p>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    
    </div> <!-- End hidden content -->
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const transactions = ref([])
const summary = ref({ total_income: 0, total_expenses: 0, net_profit: 0 })
const categories = ref([])
const loading = ref(true)

const { showFlash } = useFlash()

const currencySymbol = computed(() => authStore.currencySymbol)

const filters = reactive({
    category_id: 'all',
    start_date: '',
    end_date: '',
    search: ''
})

const filteredTransactions = computed(() => {
    let result = transactions.value
    
    if (filters.category_id !== 'all') {
        result = result.filter(t => t.category_id === filters.category_id)
    }
    if (filters.start_date) {
        result = result.filter(t => new Date(t.transaction_date) >= new Date(filters.start_date))
    }
    if (filters.end_date) {
        result = result.filter(t => new Date(t.transaction_date) <= new Date(filters.end_date))
    }
    if (filters.search) {
        const s = filters.search.toLowerCase()
        result = result.filter(t => 
            (t.title && t.title.toLowerCase().includes(s)) || 
            (t.description && t.description.toLowerCase().includes(s)) ||
            (t.payee_or_payer && t.payee_or_payer.toLowerCase().includes(s))
        )
    }
    return result
})

async function fetchData() {
    loading.value = true
    try {
        const [txnRes, sumRes, catRes] = await Promise.all([
            api.get('/finance/transactions'),
            api.get('/finance/summary'),
            api.get('/finance/categories'),
        ])

        if (txnRes.data?.transactions) transactions.value = txnRes.data.transactions
        if (sumRes.data?.summary) summary.value = sumRes.data.summary
        if (catRes.data?.categories) {
            categories.value = catRes.data.categories
        }
    } catch(err) {
        // Handle error silently
    } finally {
        loading.value = false
    }
}

onMounted(fetchData)

function formatDate(dateStr) {
    if (!dateStr) return '-'
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
}

async function exportCsv() {
    try {
        const res = await api.get('/finance/export', { responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([res.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'finance_export.csv')
        document.body.appendChild(link)
        link.click()
        link.parentNode.removeChild(link)
        showFlash('Export downloaded successfully!', 'success')
    } catch (err) {
        showFlash('Failed to export CSV', 'error')
    }
}
</script>
