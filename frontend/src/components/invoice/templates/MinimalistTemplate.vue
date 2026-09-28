<template>
  <div class="relative w-full bg-white text-slate-900 flex flex-col justify-between overflow-hidden rounded-none shadow-sm p-6 sm:p-10 border border-slate-300 font-sans aspect-[1/1.35] min-h-[640px] z-0 transition-all">
    <!-- Watermark Backdrop -->
    <div class="absolute inset-0 z-[-1] flex items-center justify-center pointer-events-none select-none opacity-[0.03]">
      <img v-if="profile?.logo_url" :src="profile.logo_url" class="w-3/4 h-3/4 object-contain grayscale" alt="" />
      <span v-else class="material-symbols-outlined text-[400px]">receipt_long</span>
    </div>

    <div class="space-y-5">
      <!-- Invoice Header -->
      <div class="flex justify-between items-start gap-4 pb-5 border-b border-slate-200">
        <div>
          <img v-if="profile?.logo_url" :src="profile.logo_url" alt="Logo" class="max-h-14 w-auto mb-2 object-contain" />
          <div v-else class="w-10 h-10 rounded mb-2 flex items-center justify-center bg-slate-100 text-slate-800">
            <span class="material-symbols-outlined text-[20px]">receipt_long</span>
          </div>
          <h3 class="text-xl font-bold font-headline text-slate-900">
            {{ profile?.company_name || 'Your Company' }}
          </h3>
          <p class="font-body text-xs sm:text-sm text-slate-500 max-w-xs mt-0.5 whitespace-pre-line">
            {{ profile?.address || '' }}
          </p>
        </div>
        <div class="text-right space-y-1">
          <span class="text-2xl sm:text-3xl font-black uppercase block text-slate-900 font-sans tracking-tight">
            {{ invoice.invoice_number || 'INV-0001' }}
          </span>
          <div v-if="invoice.status" class="inline-block pt-1">
            <span v-if="invoice.is_paid || invoice.status === 'paid'" class="px-2.5 py-0.5 rounded-none text-xs font-bold bg-slate-200 text-slate-900 border border-slate-400">PAID</span>
            <span v-else-if="invoice.status === 'sent'" class="px-2.5 py-0.5 rounded-none text-xs font-bold bg-amber-500/20 text-amber-700 border border-amber-500/30">SENT</span>
            <span v-else-if="invoice.status === 'overdue'" class="px-2.5 py-0.5 rounded-none text-xs font-bold bg-red-500/20 text-red-700 border border-red-500/30">OVERDUE</span>
            <span v-else class="px-2.5 py-0.5 rounded-none text-xs font-bold bg-slate-100 text-slate-700 border border-slate-300">DRAFT</span>
          </div>
          <p class="font-body text-xs sm:text-sm text-slate-500">Issue: <span class="font-semibold text-slate-800">{{ formatDate(invoice.issue_date) }}</span></p>
          <p class="font-body text-xs sm:text-sm text-slate-500">Due: <span class="font-semibold text-slate-800">{{ formatDate(invoice.due_date) }}</span></p>
        </div>
      </div>

      <!-- Billed To Section -->
      <div class="p-4 sm:p-5 font-body text-xs sm:text-sm space-y-1 bg-white border-l-2 border-slate-900 pl-3 rounded-none">
        <span class="text-xs font-bold uppercase tracking-wider block text-slate-700">Billed To</span>
        <p class="font-bold text-base text-slate-900">{{ invoice.client_name || invoice.client_email || 'client@company.com' }}</p>
        <p v-if="invoice.client_name && invoice.client_email" class="text-xs text-slate-500">{{ invoice.client_email }}</p>
        <p class="text-xs sm:text-sm text-slate-600 whitespace-pre-line">{{ invoice.client_address || 'Client Address...' }}</p>
      </div>

      <!-- Line Items Table -->
      <div>
        <table class="w-full text-left font-body text-xs sm:text-sm">
          <thead>
            <tr class="font-label text-xs uppercase border-b border-t border-slate-900 text-slate-800 tracking-wider">
              <th class="py-2.5 w-1/2">Description</th>
              <th class="py-2.5 text-center">Qty</th>
              <th class="py-2.5 text-right">Price</th>
              <th class="py-2.5 text-right">Amount</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="(item, idx) in invoice.items" :key="idx">
              <td class="py-2.5 font-medium text-sm text-slate-900">{{ item.description || 'Item description' }}</td>
              <td class="py-2.5 text-center text-sm text-slate-600">{{ item.quantity || 0 }}</td>
              <td class="py-2.5 text-right text-sm text-slate-600">{{ currencySymbol }}{{ (item.unit_price || 0).toFixed(2) }}</td>
              <td class="py-2.5 text-right font-semibold text-sm text-slate-900">{{ currencySymbol }}{{ ((item.quantity || 0) * (item.unit_price || 0)).toFixed(2) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Invoice Calculation Summary -->
    <div class="pt-4 flex justify-between items-end gap-4 mt-6 border-t border-slate-200">
      <div class="w-1/2 space-y-1">
        <span class="text-xs font-bold uppercase tracking-wider block text-slate-700">Notes</span>
        <p class="font-body text-xs sm:text-sm text-slate-500 italic whitespace-pre-line">{{ invoice.notes || 'Thank you for your business!' }}</p>
      </div>
      <div class="w-52 p-4 space-y-2 font-body text-xs sm:text-sm shrink-0 bg-transparent border-t border-slate-900 rounded-none">
        <div class="flex justify-between text-slate-600">
          <span>Subtotal</span>
          <span class="font-semibold text-slate-900">{{ currencySymbol }}{{ totals.subtotal.toFixed(2) }}</span>
        </div>
        <div class="flex justify-between text-slate-600">
          <span>Tax</span>
          <span class="font-semibold text-slate-900">{{ currencySymbol }}{{ totals.tax.toFixed(2) }}</span>
        </div>
        <div class="flex justify-between text-slate-600">
          <span>Discount</span>
          <span class="font-semibold text-slate-900">{{ currencySymbol }}{{ totals.discount.toFixed(2) }}</span>
        </div>
        <div class="flex justify-between items-center pt-2.5 text-slate-900 border-t border-slate-900 font-bold">
          <span>Total Due</span>
          <span class="font-black text-lg sm:text-xl text-slate-900">{{ currencySymbol }}{{ totals.total.toFixed(2) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface InvoiceProps {
  invoice: {
    invoice_number?: string
    issue_date?: string
    due_date?: string
    client_name?: string
    client_email?: string
    client_address?: string
    notes?: string
    items?: Array<{ description?: string; quantity?: number; unit_price?: number }>
    status?: string
    is_paid?: boolean
  }
  profile?: {
    company_name?: string
    logo_url?: string
    address?: string
  } | null
  totals: {
    subtotal: number
    tax: number
    discount: number
    total: number
  }
  currencySymbol: string
  formatDate: (dateStr: any) => string
}

defineProps<InvoiceProps>()
</script>
