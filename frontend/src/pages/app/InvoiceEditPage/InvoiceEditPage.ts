import { defineComponent } from 'vue'
import { useInvoiceForm } from '@/composables/useInvoiceForm'
import InvoiceDocumentPreview from '@/components/invoice/InvoiceDocumentPreview.vue'

export default defineComponent({
  name: 'InvoiceEditPage',
  components: {
    InvoiceDocumentPreview
  },
  setup() {
    const {
        form,
        clients,
        currencies,
        profile,
        saving,
        currencySymbol,
        totals,
        onClientChange,
        addItem,
        removeItem,
        formatDate,
        saveDraft,
        handleSubmit
    } = useInvoiceForm({ mode: 'edit' })
    return {
      InvoiceDocumentPreview,
      form,
      clients,
      currencies,
      profile,
      saving,
      currencySymbol,
      totals,
      onClientChange,
      addItem,
      removeItem,
      formatDate,
      saveDraft,
      handleSubmit
    }
  }
})
