<template>
  <div>
    <component :is="'style'">
      @media print {
        body * { display: none !important; }
        body::before {
            content: "Direct printing is disabled. Please use the 'Download PDF' feature on the dashboard.";
            display: block !important;
            text-align: center;
            font-size: 24px;
            font-weight: bold;
            font-family: sans-serif;
            margin-top: 50px;
        }
      }
    </component>

    <div v-if="loading" class="text-center py-12 text-outline">Loading invoice...</div>
    <div v-else-if="!invoice" class="text-center py-12 text-outline">Invoice not found.</div>
    <div v-else>
        <header class="no-print flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
            <div>
                <h2 class="font-headline text-3xl font-bold text-on-surface mb-1">Invoice {{ invoice.invoice_number }}</h2>
                <p class="font-body text-base text-on-surface-variant">Created on {{ formatDate(invoice.issue_date || invoice.created_at) }}</p>
            </div>
            <div class="flex items-center gap-3 flex-wrap">
                <router-link to="/user/invoices" class="px-4 py-2 rounded-xl border border-outline-variant/60 font-label text-sm font-medium text-on-surface-variant hover:bg-surface-container-low transition-colors">
                    Back to Invoices
                </router-link>
                <router-link v-if="invoice.status === 'draft'" :to="`/user/invoices/edit?id=${invoice.id}`" class="px-4 py-2 rounded-xl border border-primary/50 bg-primary-container/20 font-label text-sm font-semibold text-primary hover:bg-primary-container/40 flex items-center gap-2 transition-colors">
                    <span class="material-symbols-outlined text-[18px]">edit</span> Edit Draft
                </router-link>
                
                <button v-if="!invoice.is_paid && invoice.status !== 'paid'" @click="handleMarkPaid" class="px-4 py-2 rounded-xl bg-primary text-on-primary font-label text-sm font-semibold hover:bg-on-primary-fixed-variant flex items-center gap-2 cursor-pointer shadow-sm transition-colors">
                    <span class="material-symbols-outlined text-[18px]">check_circle</span> Mark as Paid
                </button>
                <template v-else>
                    <router-link v-if="receipt" :to="`/user/invoices/receipt/view?id=${receipt.id}`" class="px-4 py-2 rounded-xl bg-tertiary text-on-tertiary font-label text-sm font-semibold hover:bg-tertiary/90 flex items-center gap-2 cursor-pointer shadow-sm transition-colors">
                        <span class="material-symbols-outlined text-[18px]">receipt</span> View Receipt
                    </router-link>
                    <button v-else @click="generateReceipt" class="px-4 py-2 rounded-xl bg-primary text-on-primary font-label text-sm font-semibold hover:bg-on-primary-fixed-variant flex items-center gap-2 cursor-pointer shadow-sm transition-colors">
                        <span class="material-symbols-outlined text-[18px]">receipt_long</span> Generate Receipt (1 Credit)
                    </button>
                </template>

                <button @click="copyPublicLink" class="px-4 py-2 rounded-xl border border-outline-variant/60 font-label text-sm font-semibold text-on-surface hover:bg-surface-container-low flex items-center gap-2 transition-colors">
                    <span class="material-symbols-outlined text-[18px]">link</span> Copy Link
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
                        <span class="material-symbols-outlined text-[16px]">download</span> PDF (1 Credit)
                    </button>
                </div>
                
                <a v-if="invoice.public_token" :href="`/invoice/public/${invoice.public_token}`" target="_blank" class="bg-primary text-on-primary rounded-xl px-5 py-2 font-label text-sm font-semibold flex items-center gap-2 hover:bg-on-primary-fixed-variant transition-colors shadow-sm">
                    <span class="material-symbols-outlined text-[18px]">open_in_new</span> Public View
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
                    <span class="font-headline text-3xl sm:text-4xl font-black text-primary block">{{ invoice.invoice_number }}</span>
                    <div class="inline-block pt-1">
                        <span v-if="invoice.is_paid || invoice.status === 'paid'" class="px-3 py-1 rounded-full text-xs font-bold bg-primary-container/40 text-primary border border-primary/20 flex items-center gap-1">
                            <span class="material-symbols-outlined text-[14px]">check_circle</span> PAID
                        </span>
                        <span v-else-if="invoice.status === 'sent'" class="px-3 py-1 rounded-full text-xs font-bold bg-tertiary-fixed-dim/40 text-tertiary border border-tertiary/20">SENT</span>
                        <span v-else-if="invoice.status === 'overdue'" class="px-3 py-1 rounded-full text-xs font-bold bg-error-container/40 text-error border border-error/20">OVERDUE</span>
                        <span v-else class="px-3 py-1 rounded-full text-xs font-bold bg-surface-container-high text-on-surface-variant border border-outline-variant/40">DRAFT</span>
                    </div>
                    <p class="font-body text-xs sm:text-sm text-on-surface-variant pt-2">Issue Date: {{ formatDate(invoice.issue_date) }}</p>
                    <p class="font-body text-xs sm:text-sm text-on-surface-variant">Due Date: {{ formatDate(invoice.due_date) }}</p>
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
                        <span>Total Amount</span>
                        <span class="text-primary text-xl sm:text-2xl font-black">{{ invoice.currency }} {{ invoice.total?.toFixed(2) }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'

const route = useRoute()
const router = useRouter()
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
        const res = await api.get(`/invoices/view?id=${id}`)
        if (res.data?.invoice) {
            invoice.value = res.data.invoice
            profile.value = res.data.profile // We will need to make sure backend returns profile too
        }
    } catch(e) {
        showFlash('Failed to load invoice', 'error')
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

async function handleMarkPaid() {
    try {
        await api.post('/invoices/mark-paid', { id: invoice.value.id })
        invoice.value.status = 'paid'
        invoice.value.is_paid = true
        showFlash('Invoice marked as paid!', 'success')
    } catch (err) {
        showFlash('Failed to mark invoice as paid', 'error')
    }
}

async function generateReceipt() {
    try {
        const res = await api.post('/invoices/receipts', { invoice_id: invoice.value.id })
        showFlash('Receipt generated successfully!', 'success')
        router.push(`/user/invoices/receipt/view?id=${res.data.receipt.id}`)
    } catch (err) {
        showFlash(err.response?.data?.error || 'Failed to generate receipt', 'error')
    }
}

function printInvoice() {
    window.print()
}

async function downloadPDF() {
    try {
        const res = await api.get(`/invoices/download?id=${invoice.value.id}&size=${paperSize.value}`, { responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([res.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `${invoice.value.invoice_number}.pdf`)
        document.body.appendChild(link)
        link.click()
        link.parentNode.removeChild(link)
    } catch (err) {
        let msg = 'Failed to download PDF'
        if (err.response?.data instanceof Blob) {
            try {
                const text = await err.response.data.text()
                const json = JSON.parse(text)
                if (json.error) msg = json.error
            } catch(e) {}
        } else if (err.response?.data?.error) {
            msg = err.response.data.error
        }
        showFlash(msg, 'error')
    }
}

function copyPublicLink() {
    if(!invoice.value?.public_token) {
        showFlash('No public token found for this invoice', 'error')
        return
    }
    const url = window.location.origin + '/invoice/public/' + invoice.value.public_token
    navigator.clipboard.writeText(url).then(() => {
        showFlash('Public link copied to clipboard!', 'success')
    }).catch(() => {
        showFlash('Failed to copy link', 'error')
    })
}
</script>
