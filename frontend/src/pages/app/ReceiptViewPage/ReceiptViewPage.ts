import { defineComponent,  ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'

export default defineComponent({
  name: 'ReceiptViewPage',
  setup() {
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
        } catch (err: any) {
            showFlash('Failed to download PDF', 'error')
        }
    }

    async function dispatchReceiptEmail() {
        try {
            const targetEmail = receipt.value?.invoice?.client?.email || ''
            await api.post('/invoices/receipts/dispatch', { receipt_id: receipt.value.id, email: targetEmail })
            showFlash(`Payment receipt dispatched successfully to ${targetEmail || 'client email'}!`, 'success')
        } catch (err: any) {
            showFlash(err.response?.data?.error || 'Failed to dispatch receipt email', 'error')
        }
    }
    return {
      route,
      showFlash,
      invoice,
      profile,
      receipt,
      loading,
      paperSize,
      formatDate,
      printInvoice,
      downloadPDF,
      dispatchReceiptEmail
    }
  }
})
