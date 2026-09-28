<template>
  <div>
    <!-- Print-only styling rules -->
    <component :is="'style'">
      @media print {
        body * { visibility: hidden; }
        #printable-invoice, #printable-invoice * { visibility: visible; }
        #printable-invoice { position: absolute; left: 0; top: 0; width: 100%; }
        .no-print { display: none !important; }
      }
    </component>

    <div v-if="loading" class="text-center py-12 text-outline">Loading receipt...</div>
    <div v-else-if="!receipt" class="text-center py-12 text-outline">Receipt not found.</div>
    <div v-else>
        <header class="no-print flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
            <div>
                <h2 class="font-headline text-3xl font-bold text-on-surface mb-1">Receipt {{ receipt.receipt_number }}</h2>
                <p class="font-body text-base text-on-surface-variant">Issued on {{ formatDate(receipt.issued_at || receipt.created_at) }}</p>
            </div>
            <div class="flex items-center gap-3 flex-wrap">
                <router-link to="/user/invoices" class="px-4 py-2 rounded-xl border border-outline-variant/60 font-label text-sm font-medium text-on-surface-variant hover:bg-surface-container-low transition-colors">
                    Back to Invoices
                </router-link>
                <router-link :to="`/user/invoices/view?id=${invoice.id}`" class="px-4 py-2 rounded-xl border border-outline-variant/60 font-label text-sm font-semibold text-on-surface hover:bg-surface-container-low flex items-center gap-2 transition-colors">
                    <span class="material-symbols-outlined text-[18px]">receipt</span> View Original Invoice
                </router-link>
                
                <button @click="dispatchReceiptEmail" class="px-4 py-2 rounded-xl border border-outline-variant/60 font-label text-sm font-semibold text-on-surface hover:bg-surface-container-low flex items-center gap-2 transition-colors cursor-pointer">
                    <span class="material-symbols-outlined text-[18px]">send</span> Email Dispatch
                </button>
                
                <div class="flex items-center gap-1 bg-surface-container-lowest/80 p-1 border border-outline-variant/60 rounded-xl">
                    <select v-model="paperSize" class="pl-2 pr-6 py-1 bg-transparent border-0 font-label text-xs font-semibold text-on-surface appearance-none cursor-pointer focus:ring-0 outline-none">
                        <option value="a4">📄 A4 Standard</option>
                        <option value="pos_80">🧾 80mm POS Thermal Bill</option>
                        <option value="pos_58">🧾 58mm POS Mini Slip</option>
                        <option value="a5">📄 A5 Compact</option>
                        <option value="letter">📄 US Letter</option>
                        <option value="legal">📄 US Legal</option>
                    </select>
                    <button @click="downloadPDF" class="px-3 py-1.5 rounded-lg bg-primary text-on-primary font-label text-xs font-semibold hover:bg-on-primary-fixed-variant flex items-center gap-1.5 transition-colors cursor-pointer border-0">
                        <span class="material-symbols-outlined text-[16px]">download</span> PDF
                    </button>
                </div>
                
                <a v-if="invoice.public_token" :href="`/invoice/public/${invoice.public_token}`" target="_blank" class="bg-primary text-on-primary rounded-xl px-5 py-2 font-label text-sm font-semibold flex items-center gap-2 hover:bg-on-primary-fixed-variant transition-colors shadow-sm">
                    <span class="material-symbols-outlined text-[18px]">open_in_new</span> Original Public Invoice
                </a>
            </div>
        </header>

        <!-- Printable Invoice Document -->
        <div id="printable-invoice" class="relative max-w-[840px] mx-auto glass-card rounded-2xl p-8 sm:p-12 border border-outline-variant/60 bg-white space-y-8 z-0 overflow-hidden">
            
            <!-- Watermark Backdrop -->
            <div class="absolute inset-0 z-[-1] flex items-center justify-center pointer-events-none select-none opacity-[0.03]">
                <img v-if="profile?.logo_url" :src="profile.logo_url" class="w-3/4 h-3/4 object-contain grayscale" alt="">
                <span v-else class="material-symbols-outlined text-[400px]">receipt_long</span>
            </div>

            <!-- Invoice Header Row -->
            <div class="flex justify-between items-start gap-6 pb-6 border-b border-outline-variant/40">
                <div>
                    <img v-if="profile?.logo_url" :src="profile.logo_url" alt="Company Logo" class="max-h-16 w-auto mb-3 object-contain">
                    <div v-else class="w-12 h-12 rounded-xl mb-3 bg-primary/10 text-primary flex items-center justify-center">
                        <span class="material-symbols-outlined text-[24px]">receipt_long</span>
                    </div>
                    
                    <h3 class="font-headline text-2xl sm:text-3xl font-bold text-on-surface">
                        {{ profile?.company_name || 'My Business' }}
                    </h3>
                    <p class="font-body text-xs sm:text-sm text-on-surface-variant max-w-xs mt-1 whitespace-pre-line">
                        {{ profile?.address }}
                    </p>
                </div>

                <div class="text-right space-y-1">
                    <span class="font-headline text-3xl sm:text-4xl font-black text-primary block">RECEIPT</span>
                    <p class="font-body text-sm sm:text-base text-on-surface-variant block">{{ receipt.receipt_number }}</p>
                    <div class="inline-block pt-1">
                        <span class="px-3 py-1 rounded-full text-xs font-bold bg-primary-container/40 text-primary border border-primary/20 flex items-center gap-1">
                            <span class="material-symbols-outlined text-[14px]">check_circle</span> PAID
                        </span>
                    </div>
                    <p class="font-body text-xs sm:text-sm text-on-surface-variant pt-2">Date Paid: {{ formatDate(receipt.issued_at) }}</p>
                    <p class="font-body text-xs sm:text-sm text-on-surface-variant">For Invoice: {{ invoice.invoice_number }}</p>
                </div>
            </div>

            <!-- Client Information -->
            <div class="p-4 sm:p-5 bg-surface-container-low/40 rounded-2xl border border-outline-variant/30 font-body text-xs sm:text-sm text-on-surface space-y-1">
                <h4 class="font-label text-xs font-bold text-primary uppercase tracking-wider mb-1">Billed To</h4>
                <p class="font-bold text-base text-on-surface">{{ invoice.client?.name || 'Direct Client' }}</p>
                <p v-if="invoice.client?.email" class="text-xs sm:text-sm text-on-surface-variant">{{ invoice.client.email }}</p>
                <p v-if="invoice.client?.address" class="text-xs sm:text-sm text-on-surface-variant whitespace-pre-line">{{ invoice.client.address }}</p>
            </div>

            <!-- Line Items Table -->
            <div>
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="border-b border-outline-variant/40 font-label text-xs sm:text-sm uppercase text-on-surface-variant font-bold bg-surface-container-low/40">
                            <th class="px-4 py-3.5">Description</th>
                            <th class="px-4 py-3.5 text-center">Qty</th>
                            <th class="px-4 py-3.5 text-right">Unit Price</th>
                            <th class="px-4 py-3.5 text-right">Amount</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-outline-variant/30 font-body text-sm sm:text-base text-on-surface">
                        <tr v-for="item in invoice.line_items" :key="item.id">
                            <td class="px-4 py-4 sm:py-4.5 font-medium text-sm sm:text-base">{{ item.description }}</td>
                            <td class="px-4 py-4 sm:py-4.5 text-center text-sm sm:text-base">{{ item.quantity }}</td>
                            <td class="px-4 py-4 sm:py-4.5 text-right text-sm sm:text-base">{{ invoice.currency }} {{ item.unit_price?.toFixed(2) }}</td>
                            <td class="px-4 py-4 sm:py-4.5 text-right font-semibold text-sm sm:text-base">{{ invoice.currency }} {{ item.amount?.toFixed(2) }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Invoice Calculation Summary -->
            <div class="flex flex-col sm:flex-row justify-between items-start gap-6 pt-4 border-t border-outline-variant/40">
                <div class="w-full sm:w-1/2 space-y-1.5">
                    <h5 class="font-label text-xs font-bold text-primary uppercase">Notes & Terms</h5>
                    <p class="font-body text-xs sm:text-sm text-on-surface-variant whitespace-pre-line">{{ invoice.notes || 'Thank you for your business!' }}</p>
                </div>

                <div class="w-full sm:w-1/2 bg-surface-container-low/40 p-4 rounded-xl border border-outline-variant/30 space-y-2 font-body text-sm sm:text-base">
                    <div class="flex justify-between text-on-surface-variant">
                        <span>Subtotal</span>
                        <span class="font-semibold text-on-surface">{{ invoice.currency }} {{ invoice.subtotal?.toFixed(2) }}</span>
                    </div>
                    <div v-if="invoice.tax > 0" class="flex justify-between text-on-surface-variant">
                        <span>Tax</span>
                        <span class="font-semibold text-on-surface">{{ invoice.currency }} {{ invoice.tax?.toFixed(2) }}</span>
                    </div>
                    <div v-if="invoice.discount > 0" class="flex justify-between text-on-surface-variant">
                        <span>Discount</span>
                        <span class="font-semibold text-on-surface">-{{ invoice.currency }} {{ invoice.discount?.toFixed(2) }}</span>
                    </div>
                    <div class="flex justify-between items-center pt-3 border-t border-outline-variant/40 font-headline text-lg sm:text-xl font-bold text-on-surface">
                        <span>Total Paid</span>
                        <span class="text-primary text-xl sm:text-2xl font-black">{{ receipt.currency }} {{ receipt.amount?.toFixed(2) }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'

const route = useRoute()
const { showFlash } = useFlash()
const invoice = ref(null)
const profile = ref(null)
const receipt = ref(null)
const loading = ref(true)
const paperSize = ref('a4')

onMounted(async () => {
    const id = route.query.id
    if (!id) return

    try {
        const res = await api.get(`/invoices/receipts/view?id=${id}`)
        if (res.data?.receipt) {
            receipt.value = res.data.receipt
            invoice.value = receipt.value.invoice
            
            // Fetch profile separately
            const pRes = await api.get('/profile').catch(() => ({ data: {} }))
            if (pRes.data?.profile) {
                profile.value = pRes.data.profile
            }
        }
    } catch(e) {
        showFlash('Failed to load receipt', 'error')
    } finally {
        loading.value = false
    }
})

function formatDate(dateStr) {
    if (!dateStr) return '-'
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
}

function printInvoice() {
    window.print()
}

async function downloadPDF() {
    try {
        const res = await api.get(`/invoices/receipts/download?id=${receipt.value.id}&size=${paperSize.value}`, { responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([res.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `${receipt.value.receipt_number}.pdf`)
        document.body.appendChild(link)
        link.click()
        if (link.parentNode) link.parentNode.removeChild(link)
    } catch (err) {
        showFlash('Failed to download PDF', 'error')
    }
}

async function dispatchReceiptEmail() {
    try {
        const targetEmail = receipt.value?.invoice?.client?.email || ''
        await api.post('/invoices/receipts/dispatch', { receipt_id: receipt.value.id, email: targetEmail })
        showFlash(`Payment receipt dispatched successfully to ${targetEmail || 'client email'}!`, 'success')
    } catch (err) {
        showFlash(err.response?.data?.error || 'Failed to dispatch receipt email', 'error')
    }
}
</script>
