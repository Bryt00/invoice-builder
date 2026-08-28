<template>
  <div class="w-full max-w-[1600px] mx-auto pb-24 space-y-6">
    <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
            <!-- Left Pane: Data Entry Form -->
            <section class="lg:col-span-7 flex flex-col gap-6 bg-surface p-6 sm:p-8 rounded-2xl border border-outline-variant/60 shadow-sm">
                <div class="flex justify-between items-end border-b border-outline-variant/40 pb-4">
                    <div>
                        <h1 class="font-headline text-2xl font-bold text-on-surface mb-1">Invoice Editor</h1>
                        <p class="font-body text-sm text-on-surface-variant">Drafting <span class="font-semibold text-primary">{{ form.invoice_number || 'Auto-generated' }}</span></p>
                    </div>
                    <div class="bg-surface-container-high px-3 py-1 rounded-full flex items-center gap-1.5 border border-outline-variant/40">
                        <span class="w-2 h-2 rounded-full bg-primary animate-pulse"></span>
                        <span class="font-label text-xs font-semibold text-on-surface-variant">Draft Mode</span>
                    </div>
                </div>

                <div class="bg-surface-container-lowest border border-outline-variant/60 rounded-xl p-6 shadow-sm flex flex-col gap-4">
                    <div class="flex items-center gap-2 text-secondary border-b border-outline-variant/40 pb-3 mb-1">
                        <span class="material-symbols-outlined text-primary text-[20px]">person</span>
                        <h2 class="font-label text-sm font-semibold text-on-surface">Client & Invoice Setup</h2>
                    </div>

                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Billed To (Saved Client)</label>
                            <select v-model="form.client_id" @change="onClientChange" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                <option value="">-- Select Saved Client --</option>
                                <option v-for="c in clients" :key="c.id" :value="c.id">{{ c.name }} ({{ c.email }})</option>
                            </select>
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Currency</label>
                            <select v-model="form.currency" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                <option v-for="curr in currencies" :key="curr.Code || curr" :value="curr.Code || curr">
                                    {{ curr.Flag ? curr.Flag + ' ' : '' }}{{ curr.Code || curr }} ({{ curr.Symbol || '$' }})
                                </option>
                            </select>
                        </div>
                    </div>

                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Client Email</label>
                            <input type="email" v-model="form.client_email" placeholder="client@company.com" required
                                   class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider flex items-center justify-between">
                                <span>Invoice Number</span>
                                <span class="text-[10px] font-normal text-outline lowercase">(auto-generated)</span>
                            </label>
                            <input type="text" v-model="form.invoice_number" readonly placeholder="Auto-generated"
                                   class="w-full px-3 py-2 bg-surface-container-low/60 border border-outline-variant/50 rounded-xl font-headline text-sm font-bold text-primary outline-none cursor-not-allowed select-none">
                        </div>
                    </div>

                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Issue Date</label>
                            <input type="date" v-model="form.issue_date" required
                                   class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Due Date</label>
                            <input type="date" v-model="form.due_date" required
                                   class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                        </div>
                    </div>

                    <div class="flex flex-col gap-1.5">
                        <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Client Billing Address (Optional)</label>
                        <textarea v-model="form.client_address" rows="2" placeholder="123 Client St, City, Country"
                                  class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary resize-none"></textarea>
                    </div>

                    <!-- Line Items Section -->
                    <div class="space-y-3 pt-2">
                        <div class="flex justify-between items-center">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Line Items</label>
                            <button type="button" @click="addItem" class="text-xs font-label font-semibold text-primary hover:underline flex items-center gap-1 cursor-pointer">
                                <span class="material-symbols-outlined text-[16px]">add_circle</span> Add Item
                            </button>
                        </div>

                        <div class="space-y-3">
                            <div v-for="(item, idx) in form.items" :key="idx" class="line-item-row grid grid-cols-12 gap-2 items-center p-2.5 bg-surface-container-low/50 rounded-xl border border-outline-variant/40">
                                <div class="col-span-6">
                                    <input type="text" v-model="item.description" placeholder="Description of service or product" required
                                           class="w-full px-2 py-1.5 bg-transparent border-0 border-b border-outline-variant/60 focus:border-primary outline-none font-body text-sm text-on-surface">
                                </div>
                                <div class="col-span-2">
                                    <input type="number" v-model.number="item.quantity" min="0.01" step="0.01" required
                                           class="w-full px-2 py-1.5 bg-transparent border-0 border-b border-outline-variant/60 focus:border-primary outline-none font-body text-sm text-on-surface text-center">
                                </div>
                                <div class="col-span-2">
                                    <input type="number" v-model.number="item.unit_price" min="0" step="0.01" required
                                           class="w-full px-2 py-1.5 bg-transparent border-0 border-b border-outline-variant/60 focus:border-primary outline-none font-body text-sm text-on-surface text-right">
                                </div>
                                <div class="col-span-2 flex items-center justify-between pl-1">
                                    <span class="font-headline text-xs font-semibold text-on-surface">{{ currencySymbol }}{{ ((item.quantity || 0) * (item.unit_price || 0)).toFixed(2) }}</span>
                                    <button type="button" @click="removeItem(idx)" class="p-1 text-outline hover:text-error transition-colors">
                                        <span class="material-symbols-outlined text-[16px]">close</span>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Taxes & Discounts -->
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-2">
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Tax Rate (%)</label>
                            <input type="number" v-model.number="form.tax_rate" min="0" max="100" step="0.1"
                                   class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Discount Amount</label>
                            <input type="number" v-model.number="form.discount_amount" min="0" step="0.01"
                                   class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                        </div>
                    </div>

                    <!-- Notes -->
                    <div class="flex flex-col gap-1.5 pt-2">
                        <label class="font-label text-xs font-semibold text-on-surface-variant uppercase tracking-wider">Invoice Notes / Instructions</label>
                        <textarea v-model="form.notes" rows="3" placeholder="Payment due within 14 days. Bank Details: ..."
                                  class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary resize-none"></textarea>
                    </div>
                </div>
            </section>

            <!-- Right Column: Live PDF Document Preview -->
            <section class="lg:col-span-5 flex flex-col justify-start">
                <div class="sticky top-24 space-y-3">
                    <div class="flex items-center justify-between px-1">
                        <h2 class="font-headline text-sm font-bold text-on-surface flex items-center gap-2">
                            <span class="material-symbols-outlined text-primary text-[18px]">visibility</span>
                            Live Document Preview
                        </h2>
                        <span class="font-body text-xs text-on-surface-variant">Updates in real-time</span>
                    </div>

                    <!-- Document Sheet Mockup -->
                    <div class="relative w-full max-w-[700px] bg-white rounded-xl shadow-2xl p-6 sm:p-10 text-on-surface flex flex-col justify-between overflow-hidden min-h-[640px] aspect-[1/1.35] border border-outline-variant/30 z-0">
                        
                        <!-- Watermark Backdrop -->
                        <div class="absolute inset-0 z-[-1] flex items-center justify-center pointer-events-none select-none opacity-[0.03]">
                            <img v-if="profile?.logo_url" :src="profile.logo_url" class="w-3/4 h-3/4 object-contain grayscale" alt="">
                            <span v-else class="material-symbols-outlined text-[400px]">receipt_long</span>
                        </div>

                        <div class="space-y-5">
                            <div class="flex justify-between items-start border-b border-outline-variant/40 pb-5">
                                <div>
                                    <img v-if="profile?.logo_url" :src="profile.logo_url" alt="Logo" class="max-h-14 w-auto mb-2 object-contain">
                                    <div v-else class="w-10 h-10 rounded-xl mb-2 bg-primary/10 text-primary flex items-center justify-center">
                                        <span class="material-symbols-outlined text-[20px]">receipt_long</span>
                                    </div>
                                    <h3 class="font-headline text-xl font-bold text-on-surface">
                                        {{ profile?.company_name || authStore.user?.name || 'Your Company' }}
                                    </h3>
                                    <p class="font-body text-xs sm:text-sm text-on-surface-variant max-w-xs mt-0.5 whitespace-pre-line">
                                        {{ profile?.address || '' }}
                                    </p>
                                </div>
                                <div class="text-right space-y-1">
                                    <span class="font-headline text-2xl sm:text-3xl font-black text-primary uppercase block">{{ form.invoice_number || 'INV-0001' }}</span>
                                    <p class="font-body text-xs sm:text-sm text-on-surface-variant">Issue: <span class="font-semibold text-on-surface">{{ formatDate(form.issue_date) }}</span></p>
                                    <p class="font-body text-xs sm:text-sm text-on-surface-variant">Due: <span class="font-semibold text-on-surface">{{ formatDate(form.due_date) }}</span></p>
                                </div>
                            </div>

                            <div class="p-4 sm:p-5 bg-surface-container-low/40 rounded-2xl border border-outline-variant/30 font-body text-xs sm:text-sm space-y-1">
                                <span class="font-label text-xs font-bold text-primary uppercase tracking-wider block">Billed To</span>
                                <p class="font-bold text-base text-on-surface">{{ form.client_email || 'client@company.com' }}</p>
                                <p class="text-xs sm:text-sm text-on-surface-variant whitespace-pre-line">{{ form.client_address || 'Client Address...' }}</p>
                            </div>

                            <div>
                                <table class="w-full text-left font-body text-xs sm:text-sm">
                                    <thead>
                                        <tr class="border-b-2 border-outline-variant/60 font-label text-xs uppercase text-on-surface-variant">
                                            <th class="py-2.5 w-1/2">Description</th>
                                            <th class="py-2.5 text-center">Qty</th>
                                            <th class="py-2.5 text-right">Price</th>
                                            <th class="py-2.5 text-right">Amount</th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-outline-variant/30">
                                        <tr v-for="(item, idx) in form.items" :key="idx">
                                            <td class="py-2.5 font-medium text-sm">{{ item.description || 'Item description' }}</td>
                                            <td class="py-2.5 text-center text-sm">{{ item.quantity || 0 }}</td>
                                            <td class="py-2.5 text-right text-sm">{{ currencySymbol }}{{ (item.unit_price || 0).toFixed(2) }}</td>
                                            <td class="py-2.5 text-right font-semibold text-sm">{{ currencySymbol }}{{ ((item.quantity || 0) * (item.unit_price || 0)).toFixed(2) }}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>

                        <div class="pt-4 border-t border-outline-variant/40 flex justify-between items-end gap-4 mt-6">
                            <div class="w-1/2 space-y-1">
                                <span class="font-label text-xs font-bold text-primary uppercase tracking-wider block">Notes</span>
                                <p class="font-body text-xs sm:text-sm text-on-surface-variant italic whitespace-pre-line">{{ form.notes || 'Thank you for your business!' }}</p>
                            </div>
                            <div class="w-52 bg-surface-container-low/60 p-4 rounded-xl border border-outline-variant/40 space-y-2 font-body text-xs sm:text-sm shrink-0">
                                <div class="flex justify-between text-on-surface-variant">
                                    <span>Subtotal</span>
                                    <span class="font-semibold text-on-surface">{{ currencySymbol }}{{ totals.subtotal.toFixed(2) }}</span>
                                </div>
                                <div class="flex justify-between text-on-surface-variant">
                                    <span>Tax</span>
                                    <span class="font-semibold text-on-surface">{{ currencySymbol }}{{ totals.tax.toFixed(2) }}</span>
                                </div>
                                <div class="flex justify-between text-on-surface-variant">
                                    <span>Discount</span>
                                    <span class="font-semibold text-on-surface">{{ currencySymbol }}{{ totals.discount.toFixed(2) }}</span>
                                </div>
                                <div class="flex justify-between items-center pt-2.5 border-t border-outline-variant/40 font-headline text-base font-bold text-on-surface">
                                    <span>Total Due</span>
                                    <span class="text-primary font-black text-lg sm:text-xl">{{ currencySymbol }}{{ totals.total.toFixed(2) }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </div>

        <!-- Floating Bottom Contextual Action Bar -->
        <div class="fixed bottom-6 left-1/2 -translate-x-1/2 z-40 flex items-center justify-center">
            <div class="bg-inverse-surface/95 text-inverse-on-surface rounded-full px-5 py-2.5 flex items-center gap-3 shadow-2xl border border-outline/30 backdrop-blur-md">
                <button type="button" @click="saveDraft" :disabled="saving" class="flex items-center gap-1.5 px-4 py-2 rounded-full hover:bg-white/10 transition-colors font-label text-xs font-semibold cursor-pointer disabled:opacity-50">
                    <span class="material-symbols-outlined text-[18px]">draft</span> Save Draft
                </button>
                <div class="w-px h-5 bg-outline/40"></div>
                <router-link to="/user/invoices" class="flex items-center gap-1.5 px-4 py-2 rounded-full hover:bg-white/10 transition-colors font-label text-xs font-semibold">
                    <span class="material-symbols-outlined text-[18px]">arrow_back</span> Cancel
                </router-link>
                <button type="submit" :disabled="saving" class="bg-primary-container text-on-primary-container hover:bg-primary hover:text-on-primary transition-colors px-6 py-2 rounded-full font-label text-xs font-bold flex items-center gap-2 shadow-md cursor-pointer disabled:opacity-50">
                    <span class="material-symbols-outlined text-[18px]">send</span> Save &amp; Dispatch
                    <span class="bg-black/20 px-2 py-0.5 rounded text-[10px] font-mono">-1 Credit</span>
                </button>
            </div>
        </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'

const router = useRouter()
const authStore = useAuthStore()
const { showFlash } = useFlash()

const clients = ref([])
const currencies = ref([])
const profile = ref(null)
const saving = ref(false)

const today = new Date().toISOString().split('T')[0]
const nextWeek = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]

const form = reactive({
    client_id: '',
    client_email: '',
    client_address: '',
    currency: 'USD',
    invoice_number: '',
    issue_date: today,
    due_date: nextWeek,
    tax_rate: 0,
    discount_amount: 0,
    notes: '',
    items: [{ description: 'UI/UX Design Retainer', quantity: 1, unit_price: 1200.00 }]
})

const currencySymbol = computed(() => {
    const cur = currencies.value.find(c => c.Code === form.currency)
    return cur?.Symbol || '$'
})

const totals = computed(() => {
    let subtotal = 0
    form.items.forEach(item => {
        subtotal += (item.quantity || 0) * (item.unit_price || 0)
    })
    
    let tax = 0
    if (form.tax_rate > 0) {
        tax = subtotal * (form.tax_rate / 100)
    }
    
    let discount = Number(form.discount_amount || 0)
    let total = subtotal + tax - discount
    if (total < 0) total = 0
    
    return { subtotal, tax, discount, total }
})

onMounted(async () => {
    const [cRes, pRes] = await Promise.all([
        api.get('/clients').catch(() => ({ data: { clients: [] } })),
        api.get('/profile').catch(() => ({ data: {} })),
    ])
    if (cRes.data?.clients) clients.value = cRes.data.clients
    if (pRes.data?.profile) {
        profile.value = pRes.data.profile
        if (profile.value.default_currency) form.currency = profile.value.default_currency
    }
    if (pRes.data?.currencies) currencies.value = pRes.data.currencies

    // Generate mockup ID
    form.invoice_number = 'INV-' + Date.now().toString().slice(-6)
})

function onClientChange() {
    const client = clients.value.find(c => c.id === form.client_id)
    if (client) {
        form.client_email = client.email || ''
        form.client_address = client.address || ''
    } else {
        form.client_email = ''
        form.client_address = ''
    }
}

function addItem() {
    form.items.push({ description: '', quantity: 1, unit_price: 0 })
}

function removeItem(idx) {
    if (form.items.length > 1) {
        form.items.splice(idx, 1)
    }
}

function formatDate(dateStr) {
    if (!dateStr) return '-'
    const d = new Date(dateStr)
    // ensure valid date before formatting
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
}

async function saveDraft() {
    await submitForm(true)
}

async function handleSubmit() {
    await submitForm(false)
}

async function submitForm(isDraft) {
    saving.value = true
    try {
        const payload = {
            ...form,
            save_as_draft: isDraft,
            action: isDraft ? 'draft' : 'dispatch'
        }
        if (!payload.id) delete payload.id
        await api.post('/invoices', payload)
        showFlash(isDraft ? 'Invoice saved as draft!' : 'Invoice generated and dispatched!', 'success')
        router.push('/user/invoices')
    } catch (err) {
        showFlash(err.response?.data?.error || 'Failed to process invoice', 'error')
    } finally {
        saving.value = false
    }
}
</script>
