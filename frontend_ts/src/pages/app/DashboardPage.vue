<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
            <div class="flex items-center gap-3 mb-1">
                <div class="group flex items-center gap-3">
                    <h2 class="font-headline text-3xl font-bold text-on-surface group-hover:text-primary transition-colors">Welcome back, {{ authStore.user?.name || 'User' }}</h2>
                </div>
            </div>
            <p class="font-body text-base text-on-surface-variant">Here is a quick overview of your account, billing credits, and invoice dispatches.</p>
        </div>
        <router-link to="/user/invoices/new" class="btn-auth-submit text-on-primary rounded-xl px-5 py-2.5 font-label text-sm font-semibold flex items-center gap-2 shrink-0 self-start md:self-auto shadow-md hover:bg-on-primary-fixed-variant transition-all bg-primary">
            <span class="material-symbols-outlined text-[20px]">add</span>
            Create New Invoice
        </router-link>
    </header>

    <!-- Bento Grid Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <!-- Credits Card -->
        <div class="md:col-span-2 glass-card rounded-2xl p-6 border border-outline-variant/60 relative overflow-hidden flex flex-col justify-between min-h-[190px]">
            <div class="relative z-10 flex justify-between items-start">
                <div>
                    <p class="font-label text-sm font-medium text-on-surface-variant mb-1 flex items-center gap-2">
                        <span class="material-symbols-outlined text-primary text-[20px]">token</span>
                        Available Credits
                    </p>
                    <h3 class="font-headline text-5xl font-extrabold text-on-surface">{{ stats.balance || 0 }}</h3>
                </div>
            </div>
            <div class="relative z-10 mt-6 pt-4 border-t border-outline-variant/40 flex items-center justify-between">
                <p class="font-body text-sm text-on-surface-variant">1 credit is used per invoice email dispatch or PDF export.</p>
                <a class="auth-link font-label text-sm font-semibold text-primary flex items-center gap-1 hover:underline" href="#topup">
                    Top Up Credits <span class="material-symbols-outlined text-[16px]">arrow_forward</span>
                </a>
            </div>
        </div>

        <!-- Credit Usage Statistics Card -->
        <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 flex flex-col justify-between min-h-[190px]">
            <div>
                <p class="font-label text-sm font-medium text-on-surface-variant mb-1 flex items-center gap-2">
                    <span class="material-symbols-outlined text-primary text-[20px]">pie_chart</span>
                    Credit Usage Rate
                </p>
                <h3 class="font-headline text-3xl font-bold text-on-surface mb-1">{{ usagePercent }}% <span class="text-sm font-normal text-on-surface-variant">used</span></h3>
                <p class="font-body text-xs text-on-surface-variant">{{ stats.total_used || 0 }} of {{ stats.total_purchased || 0 }} total credits consumed</p>
            </div>
            <div class="pt-4 border-t border-outline-variant/40">
                <div class="w-full bg-surface-container-high rounded-full h-2 overflow-hidden mb-2">
                    <div class="bg-primary h-2 rounded-full transition-all duration-500" :style="{ width: usagePercent + '%' }"></div>
                </div>
                <p class="font-body text-xs text-on-surface-variant flex items-center justify-between">
                    <span>Available: {{ stats.balance || 0 }} credits</span>
                    <a href="#topup" class="text-primary hover:underline font-medium">Buy More &rarr;</a>
                </p>
            </div>
        </div>
    </div>

    <!-- Main Grid: Recent Invoices & Top Up Bundles -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Recent Invoices List -->
        <div class="lg:col-span-2 glass-card rounded-2xl border border-outline-variant/60 overflow-hidden flex flex-col">
            <div class="px-6 py-4 border-b border-outline-variant/40 flex justify-between items-center bg-surface-container-lowest/50">
                <router-link to="/user/invoices" class="group flex items-center gap-2">
                    <h3 class="font-headline text-lg font-bold text-on-surface group-hover:text-primary transition-colors flex items-center gap-2">
                        <span class="material-symbols-outlined text-primary text-[20px]">history</span>
                        Recent Invoices
                    </h3>
                </router-link>
                <router-link to="/user/invoices" class="font-body text-xs text-on-surface-variant hover:text-primary transition-colors">Showing latest {{ recentInvoices.length }} &rarr;</router-link>
            </div>
            
            <div v-if="loading" class="py-8 text-center text-outline text-sm">
                Loading invoices...
            </div>
            <div v-else-if="recentInvoices.length === 0" class="px-6 py-8 text-center text-on-surface-variant">
                <span class="material-symbols-outlined text-3xl text-outline mb-1">description</span>
                <p class="font-body text-xs">No recent invoices found.</p>
            </div>
            <div v-else class="divide-y divide-outline-variant/30 flex-grow">
                <div v-for="inv in recentInvoices" :key="inv.id" class="px-6 py-4 hover:bg-surface-container-low/50 transition-colors flex items-center justify-between group">
                    <router-link :to="inv.status === 'draft' ? `/user/invoices/edit?id=${inv.id}` : `/user/invoices/view?id=${inv.id}`" class="flex items-center gap-4 flex-1">
                        <div :class="['w-10 h-10 rounded-xl flex items-center justify-center shrink-0', inv.status === 'paid' ? 'bg-primary-container/30 text-primary' : inv.status === 'sent' ? 'bg-tertiary-container/30 text-tertiary' : 'bg-surface-container-high text-on-surface-variant']">
                            <span class="material-symbols-outlined text-[20px]">{{ inv.status === 'paid' ? 'check_circle' : inv.status === 'sent' ? 'send' : 'edit_note' }}</span>
                        </div>
                        <div>
                            <p class="font-label text-sm font-semibold text-on-surface group-hover:text-primary transition-colors">{{ inv.client?.name || 'Direct Client' }}</p>
                            <p class="font-body text-xs text-on-surface-variant">#{{ inv.invoice_number }} • {{ formatDate(inv.created_at) }}</p>
                        </div>
                    </router-link>
                    <div class="flex items-center gap-4">
                        <span class="font-headline text-sm font-bold text-on-surface">{{ inv.currency }} {{ inv.total?.toFixed(2) }}</span>
                        <span v-if="inv.status === 'paid'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-primary-container/40 text-primary border border-primary/20">Paid</span>
                        <span v-else-if="inv.status === 'sent'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-tertiary-fixed-dim/40 text-on-tertiary-container border border-tertiary/20">Sent</span>
                        <span v-else-if="inv.status === 'overdue'" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-error-container/40 text-error border border-error/20">Overdue</span>
                        <router-link v-else :to="`/user/invoices/edit?id=${inv.id}`" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-surface-container-high text-primary hover:bg-primary-container/30 border border-outline-variant/40 flex items-center gap-1">
                            <span class="material-symbols-outlined text-[14px]">edit</span> Draft
                        </router-link>
                    </div>
                </div>
            </div>
            
            <div class="px-6 py-3 border-t border-outline-variant/40 bg-surface-container-lowest/50 text-center">
                <router-link to="/user/invoices" class="font-label text-sm font-semibold text-primary hover:underline">View All Invoices &rarr;</router-link>
            </div>
        </div>

        <!-- Top Up Credit Section -->
        <form @submit.prevent="handleTopup" class="lg:col-span-1 glass-card rounded-2xl p-6 border border-outline-variant/60 flex flex-col justify-between" id="topup">
            <div>
                <h3 class="font-headline text-lg font-bold text-on-surface mb-1 flex items-center gap-2">
                    <span class="material-symbols-outlined text-primary text-[20px]">add_card</span>
                    Top Up Credits
                </h3>
                <p class="font-body text-xs text-on-surface-variant mb-5">Select a package to replenish your billing balance via Paystack.</p>

                <div class="space-y-3">
                    <div v-if="loadingPackages" class="text-center py-6 text-outline text-sm">Loading packages...</div>
                    <template v-else-if="packages.length > 0">
                        <label v-for="(pkg, idx) in packages" :key="pkg.id" :class="['flex items-center justify-between p-3.5 rounded-xl cursor-pointer transition-all', pkg.badge_tag ? 'border-2 border-primary bg-primary-container/10' : 'border border-outline-variant/60 bg-surface-container-lowest/60 hover:border-primary']">
                            <div class="flex items-center gap-3">
                                <input type="radio" v-model="selectedPackage" :value="pkg.id" class="text-primary focus:ring-primary accent-primary w-4 h-4">
                                <div>
                                    <p class="font-label text-sm font-semibold text-on-surface flex items-center gap-2">
                                        {{ pkg.credits_granted }} Credits
                                        <span v-if="pkg.badge_tag" class="px-2 py-0.5 text-[10px] uppercase font-bold bg-primary text-on-primary rounded-full">{{ pkg.badge_tag }}</span>
                                    </p>
                                    <p class="font-body text-xs text-on-surface-variant">
                                        {{ pkg.currency }} {{ (pkg.price / 100).toFixed(2) }}
                                        <span v-if="pkg.description" class="text-primary font-semibold">({{ pkg.description }})</span>
                                    </p>
                                </div>
                            </div>
                        </label>
                    </template>
                    <div v-else class="text-center py-6">
                        <span class="material-symbols-outlined text-3xl text-outline mb-1">credit_card_off</span>
                        <p class="font-body text-xs text-on-surface-variant">No credit packages available at the moment.</p>
                    </div>
                </div>
            </div>

            <button v-if="packages.length > 0" type="submit" :disabled="purchasing" class="bg-primary hover:bg-on-primary-fixed-variant w-full mt-6 py-3 rounded-xl block text-center font-label text-sm font-semibold text-on-primary cursor-pointer disabled:opacity-50">
                <span v-if="!purchasing">Purchase Selected Package</span>
                <span v-else>Processing...</span>
            </button>
        </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

const loading = ref(true)
const stats = ref({ balance: 0, total_purchased: 0, total_used: 0 })
const recentInvoices = ref([])

const loadingPackages = ref(true)
const packages = ref([])
const selectedPackage = ref('')
const purchasing = ref(false)

const usagePercent = computed(() => {
  const pur = stats.value.total_purchased || 0
  const usd = stats.value.total_used || 0
  if (pur === 0) return 0
  return Math.round((usd / pur) * 100)
})

function formatDate(dateStr: any) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
}

onMounted(async () => {
  const reference = route.query.reference || route.query.trxref
  if (reference) {
    try {
      await api.get(`/credits/topup/verify?reference=${reference}`)
      showToast('Payment verified! Your credits have been added.', 'success')
      
      // Clean up the URL
      const newQuery = { ...route.query }
      delete newQuery.reference
      delete newQuery.trxref
      router.replace({ query: newQuery })
    } catch (err: any) {
      console.error('Verification failed', err)
      showToast('Payment verification failed.', 'error')
    }
  }

  try {
    const [creditRes, invoicesRes, packagesRes] = await Promise.all([
      api.get('/credits/balance').catch(() => ({ data: { stats: {} } })),
      api.get('/invoices?limit=3').catch(() => ({ data: { invoices: [] } })),
      api.get('/credits/packages').catch(() => ({ data: { packages: [] } })),
    ])

    if (creditRes.data?.stats) stats.value = creditRes.data.stats
    if (invoicesRes.data?.invoices) recentInvoices.value = invoicesRes.data.invoices
    if (packagesRes.data?.packages) {
      packages.value = packagesRes.data.packages
      if (packages.value.length > 0) {
        selectedPackage.value = packages.value[0].id
      }
    }
  } finally {
    loading.value = false
    loadingPackages.value = false
  }
})

async function handleTopup() {
    if (!selectedPackage.value) return
    purchasing.value = true
    try {
        const res = await api.post('/credits/topup/initialize', { package_id: selectedPackage.value })
        if (res.data?.authorization_url) {
            window.location.href = res.data.authorization_url
        }
    } catch (err: any) {
        showToast(err.response?.data?.error || 'Failed to initialize payment', 'error')
    } finally {
        purchasing.value = false
    }
}
</script>
