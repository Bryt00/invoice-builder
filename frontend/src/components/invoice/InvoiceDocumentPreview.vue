<template>
  <component
    :is="activeTemplateComponent"
    :invoice="invoice"
    :profile="profile"
    :totals="totals"
    :currency-symbol="currencySymbol"
    :format-date="formatDateFn"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { InvoiceTemplate } from '@/types/invoice'
import ModernTemplate from './templates/ModernTemplate.vue'
import MinimalistTemplate from './templates/MinimalistTemplate.vue'
import ClassicTemplate from './templates/ClassicTemplate.vue'

interface Props {
  template?: InvoiceTemplate | string
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
  currencySymbol?: string
  formatDate?: (dateStr: any) => string
}

const props = withDefaults(defineProps<Props>(), {
  template: 'modern',
  currencySymbol: '$',
  profile: null
})

const defaultFormatDate = (dateStr: any) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return String(dateStr)
  return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
}

const formatDateFn = computed(() => props.formatDate || defaultFormatDate)

const activeTemplateComponent = computed(() => {
  switch (props.template) {
    case 'minimalist':
      return MinimalistTemplate
    case 'classic':
      return ClassicTemplate
    case 'modern':
    default:
      return ModernTemplate
  }
})
</script>
