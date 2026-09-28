import { defineComponent,  ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { InvoiceTemplate } from '@/types/invoice'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import InvoiceDocumentPreview from '@/components/invoice/InvoiceDocumentPreview.vue'

export default defineComponent({
  name: 'InvoiceViewPage',
  components: {
    InvoiceDocumentPreview
  },
  setup() {
    const route = useRoute()

    const router = useRouter()

    const { showFlash } = useFlash()

    const invoice = ref<any>(null)

    const profile = ref<any>(null)

    const receipt = ref<any>(null)

    const loading = ref(true)

    const paperSize = ref('a4')

    const currentTemplate = ref<InvoiceTemplate>('modern')

    const pdfPreviewModalOpen = ref(false)

    const loadingPdfPreview = ref(false)

    const pdfPreviewUrl = ref<string | null>(null)

    const displayNotes = computed(() => {
        return (invoice.value?.notes || 'Thank you for your business!')
            .replace(/[Template:s*(Modern|Minimalist|Classic)]/gi, '')
            .trim()
    })

    const invoiceDoc = computed(() => {
        if (!invoice.value) return null
        return {
            ...invoice.value,
            client_name: invoice.value.client?.name,
            client_email: invoice.value.client?.email,
            client_address: invoice.value.client?.address,
            notes: displayNotes.value,
            items: invoice.value.line_items || invoice.value.items || []
        }
    })

    const viewTotals = computed(() => ({
        subtotal: invoice.value?.subtotal || 0,
        tax: invoice.value?.tax || 0,
        discount: invoice.value?.discount || 0,
        total: invoice.value?.total || 0
    }))

    onMounted(async () => {
        const id = route.query.id
        if (!id) return

        try {
            const res = await api.get('/invoices/view?id=' + id)
            if (res.data?.invoice) {
                invoice.value = res.data.invoice
                profile.value = res.data.profile

                // Detect template
                if (invoice.value.notes?.includes('[Template: Minimalist]')) {
                    currentTemplate.value = 'minimalist'
                } else if (invoice.value.notes?.includes('[Template: Classic]')) {
                    currentTemplate.value = 'classic'
                } else if (invoice.value.notes?.includes('[Template: Modern]')) {
                    currentTemplate.value = 'modern'
                } else {
                    currentTemplate.value = (localStorage.getItem('inv_template_' + invoice.value.id) as any) || 'modern'
                }
            }
        } catch(e) {
            showFlash('Failed to load invoice', 'error')
        } finally {
            loading.value = false
        }
    })

    onUnmounted(() => {
        if (pdfPreviewUrl.value) {
            URL.revokeObjectURL(pdfPreviewUrl.value)
        }
    })

    function setTemplate(tpl: InvoiceTemplate) {
        currentTemplate.value = tpl
        if (invoice.value?.id) {
            localStorage.setItem('inv_template_' + invoice.value.id, tpl)
        }
    }

    function formatDate(dateStr: any) {
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
        } catch (err: any) {
            showFlash('Failed to mark invoice as paid', 'error')
        }
    }

    async function dispatchEmail() {
        try {
            const targetEmail = invoice.value.client?.email || ''
            await api.post('/invoices/dispatch', { invoice_id: invoice.value.id, email: targetEmail })
            showFlash('Invoice dispatched successfully to ' + (targetEmail || 'client email') + '!', 'success')
        } catch (err: any) {
            showFlash(err.response?.data?.error || 'Failed to dispatch invoice email', 'error')
        }
    }

    async function generateReceipt() {
        try {
            const res = await api.post('/invoices/receipts', { invoice_id: invoice.value.id })
            showFlash('Receipt generated successfully!', 'success')
            router.push('/user/invoices/receipt/view?id=' + res.data.receipt.id)
        } catch (err: any) {
            showFlash(err.response?.data?.error || 'Failed to generate receipt', 'error')
        }
    }

    async function openPdfPreview() {
        pdfPreviewModalOpen.value = true
        loadingPdfPreview.value = true
        try {
            const res = await api.get('/invoices/download?id=' + invoice.value.id + '&size=' + paperSize.value, { responseType: 'blob' })
            if (pdfPreviewUrl.value) {
                URL.revokeObjectURL(pdfPreviewUrl.value)
            }
            const blob = new Blob([res.data], { type: 'application/pdf' })
            pdfPreviewUrl.value = URL.createObjectURL(blob)
        } catch (err: any) {
            let msg = 'Failed to load PDF preview'
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
        } finally {
            loadingPdfPreview.value = false
        }
    }

    function closePdfPreview() {
        pdfPreviewModalOpen.value = false
        if (pdfPreviewUrl.value) {
            URL.revokeObjectURL(pdfPreviewUrl.value)
            pdfPreviewUrl.value = null
        }
    }

    async function downloadPDF() {
        try {
            const res = await api.get('/invoices/download?id=' + invoice.value.id + '&size=' + paperSize.value, { responseType: 'blob' })
            const blob = new Blob([res.data], { type: 'application/pdf' })
            const url = window.URL.createObjectURL(blob)
            const a = document.createElement('a')
            a.href = url
            a.download = 'Invoice-' + invoice.value.invoice_number + '.pdf'
            document.body.appendChild(a)
            a.click()
            window.URL.revokeObjectURL(url)
            document.body.removeChild(a)
            showFlash('Invoice downloaded successfully!', 'success')
        } catch (err: any) {
            showFlash('Failed to download invoice PDF', 'error')
        }
    }

    function copyPublicLink() {
        const token = invoice.value?.public_token
        if (!token) {
            showFlash('Public link unavailable for this invoice', 'error')
            return
        }
        const publicUrl = window.location.origin + '/invoice/public/' + token
        navigator.clipboard.writeText(publicUrl).then(() => {
            showFlash('Public invoice link copied to clipboard!', 'success')
        }).catch(() => {
            showFlash('Failed to copy public link', 'error')
        })
    }
    return {
      InvoiceDocumentPreview,
      route,
      router,
      showFlash,
      invoice,
      profile,
      receipt,
      loading,
      paperSize,
      currentTemplate,
      pdfPreviewModalOpen,
      loadingPdfPreview,
      pdfPreviewUrl,
      displayNotes,
      invoiceDoc,
      viewTotals,
      setTemplate,
      formatDate,
      handleMarkPaid,
      dispatchEmail,
      generateReceipt,
      openPdfPreview,
      closePdfPreview,
      downloadPDF,
      copyPublicLink
    }
  }
})
