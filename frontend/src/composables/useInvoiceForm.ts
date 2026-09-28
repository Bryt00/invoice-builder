import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { InvoiceTemplate } from '@/types/invoice'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import { useToast } from '@/composables/useToast'

export interface UseInvoiceFormOptions {
  mode: 'create' | 'edit'
}

export function useInvoiceForm(options: UseInvoiceFormOptions = { mode: 'create' }) {
  const router = useRouter()
  const route = useRoute()
  const authStore = useAuthStore()
  const { showFlash } = useFlash()
  const { showToast } = useToast()

  const clients = ref<any[]>([])
  const currencies = ref<any[]>([])
  const profile = ref<any>(null)
  const saving = ref(false)
  const isOnline = ref(typeof navigator !== 'undefined' ? navigator.onLine : true)
  const lastSavedAt = ref('')
  const restoredDraftNotice = ref(false)
  const restoredDraftTime = ref('')

  const DRAFT_STORAGE_KEY = 'teks_offline_invoice_draft'

  const today = new Date().toISOString().split('T')[0]
  const nextWeek = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]

  const form = reactive({
    id: '',
    client_id: '',
    client_email: '',
    client_address: '',
    currency: 'USD',
    invoice_number: '',
    template: 'modern' as InvoiceTemplate,
    issue_date: today,
    due_date: nextWeek,
    tax_rate: 0,
    discount_amount: 0,
    notes: '',
    items: (options.mode === 'create'
      ? [{ description: 'UI/UX Design Retainer', quantity: 1, unit_price: 1200.00 }]
      : []) as any[]
  })

  const currencySymbol = computed(() => {
    const cur = currencies.value.find(c => (c.Code || c) === form.currency)
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

    const discount = Number(form.discount_amount || 0)
    let total = subtotal + tax - discount
    if (total < 0) total = 0

    return { subtotal, tax, discount, total }
  })

  // Auto-save draft watcher with debounce (create mode only)
  let debounceTimer: any = null
  if (options.mode === 'create') {
    watch(
      () => ({ ...form }),
      () => {
        clearTimeout(debounceTimer)
        debounceTimer = setTimeout(() => {
          saveDraftToLocalStorage()
        }, 800)
      },
      { deep: true }
    )
  }

  function saveDraftToLocalStorage() {
    try {
      const payload = {
        form,
        savedAt: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      }
      localStorage.setItem(DRAFT_STORAGE_KEY, JSON.stringify(payload))
      lastSavedAt.value = payload.savedAt
    } catch (e) {
      console.warn('Failed to save draft to localStorage', e)
    }
  }

  function restoreDraftFromLocalStorage() {
    try {
      const raw = localStorage.getItem(DRAFT_STORAGE_KEY)
      if (!raw) return
      const parsed = JSON.parse(raw)
      if (parsed?.form) {
        Object.assign(form, parsed.form)
        restoredDraftTime.value = parsed.savedAt || 'earlier session'
        restoredDraftNotice.value = true
      }
    } catch (e) {
      console.warn('Failed to restore draft', e)
    }
  }

  function discardRestoredDraft() {
    localStorage.removeItem(DRAFT_STORAGE_KEY)
    restoredDraftNotice.value = false
    showToast('Offline draft discarded. Starting fresh.', 'info')
  }

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

  function removeItem(idx: number) {
    if (form.items.length > 1) {
      form.items.splice(idx, 1)
    }
  }

  function formatDate(dateStr: any) {
    if (!dateStr) return '-'
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
  }

  async function saveDraft() {
    await submitForm(true)
  }

  async function handleSubmit() {
    await submitForm(false)
  }

  async function submitForm(isDraft: boolean) {
    if (options.mode === 'create' && !isOnline.value) {
      saveDraftToLocalStorage()
      showToast('Offline: Invoice draft saved securely on your device. Finalize once reconnected.', 'info')
      return
    }

    saving.value = true
    try {
      // Embed template tag into notes to preserve template choice cross-session
      let finalNotes = form.notes || ''
      finalNotes = finalNotes.replace(/\[Template:\s*(Modern|Minimalist|Classic)\]/gi, '').trim()
      const tplTag = `[Template: ${form.template.charAt(0).toUpperCase() + form.template.slice(1)}]`
      finalNotes = finalNotes ? `${finalNotes}\n\n${tplTag}` : tplTag

      const payload: any = {
        ...form,
        notes: finalNotes,
        template: form.template,
        save_as_draft: isDraft,
        action: isDraft ? 'draft' : 'finalize'
      }

      let res: any
      if (options.mode === 'edit' && payload.id) {
        res = await api.put('/invoices', payload)
      } else {
        if (!payload.id) delete payload.id
        res = await api.post('/invoices', payload)
      }

      const inv = res.data?.invoice
      const targetId = inv?.id || form.id

      if (targetId) {
        localStorage.setItem(`inv_template_${targetId}`, form.template)
      }

      if (options.mode === 'create') {
        localStorage.removeItem(DRAFT_STORAGE_KEY)
      }

      const successMsg = options.mode === 'create'
        ? (isDraft ? 'Invoice saved as draft!' : 'Invoice finalized and saved!')
        : (isDraft ? 'Invoice updated as draft!' : 'Invoice finalized and saved!')

      showFlash(successMsg, 'success')

      if (!isDraft && targetId) {
        router.push(`/user/invoices/view?id=${targetId}`)
      } else {
        router.push('/user/invoices')
      }
    } catch (err: any) {
      showFlash(err.response?.data?.error || 'Failed to process invoice', 'error')
    } finally {
      saving.value = false
    }
  }

  onMounted(async () => {
    if (options.mode === 'create') {
      window.addEventListener('online', () => {
        isOnline.value = true
        showToast('Back online! Ready to sync.', 'success')
      })
      window.addEventListener('offline', () => {
        isOnline.value = false
        showToast('You are now offline. Drafts are saved to your browser.', 'warning')
      })
      restoreDraftFromLocalStorage()
    }

    const [cRes, pRes] = await Promise.all([
      api.get('/clients').catch(() => ({ data: { clients: [] } })),
      api.get('/profile').catch(() => ({ data: {} })),
    ])

    if (cRes.data?.clients) clients.value = cRes.data.clients
    if (pRes.data?.profile) {
      profile.value = pRes.data.profile
      if (profile.value.default_currency && !form.currency) {
        form.currency = profile.value.default_currency
      }
    }
    if (pRes.data?.currencies) currencies.value = pRes.data.currencies

    if (options.mode === 'create') {
      if (!form.invoice_number) {
        form.invoice_number = 'INV-' + Date.now().toString().slice(-6)
      }
    } else if (options.mode === 'edit') {
      const id = route.query.id
      if (id) {
        try {
          const res = await api.get(`/invoices/view?id=${id}`)
          if (res.data?.invoice) {
            const inv = res.data.invoice
            if (inv.status !== 'draft') {
              showFlash('Finalized invoices cannot be edited. Please create a new invoice.', 'warning')
              router.replace(`/user/invoices/view?id=${inv.id}`)
              return
            }
            form.id = inv.id
            form.client_id = inv.client_id || ''
            form.client_email = inv.client?.email || ''
            form.client_address = inv.client?.address || ''
            form.currency = inv.currency
            form.invoice_number = inv.invoice_number
            form.issue_date = inv.issue_date ? inv.issue_date.split('T')[0] : today
            form.due_date = inv.due_date ? inv.due_date.split('T')[0] : nextWeek
            form.tax_rate = inv.tax_rate || (inv.tax > 0 ? (inv.tax / inv.subtotal * 100) : 0)
            form.discount_amount = inv.discount || 0
            form.items = inv.line_items || []

            let tpl: InvoiceTemplate = 'modern'
            if (inv.notes?.includes('[Template: Minimalist]')) tpl = 'minimalist'
            else if (inv.notes?.includes('[Template: Classic]')) tpl = 'classic'
            else if (inv.notes?.includes('[Template: Modern]')) tpl = 'modern'
            else {
              tpl = (localStorage.getItem(`inv_template_${inv.id}`) as any) || 'modern'
            }
            form.template = tpl
            form.notes = (inv.notes || '').replace(/\[Template:\s*(Modern|Minimalist|Classic)\]/gi, '').trim()
          }
        } catch (err: any) {
          showFlash('Failed to load invoice', 'error')
        }
      } else {
        form.invoice_number = 'INV-' + Date.now().toString().slice(-6)
      }
    }
  })

  return {
    form,
    clients,
    currencies,
    profile,
    saving,
    isOnline,
    lastSavedAt,
    restoredDraftNotice,
    restoredDraftTime,
    currencySymbol,
    totals,
    authStore,
    onClientChange,
    addItem,
    removeItem,
    formatDate,
    saveDraft,
    handleSubmit,
    discardRestoredDraft
  }
}
