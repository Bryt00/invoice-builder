import { defineComponent, ref } from 'vue'
import { useInvoiceForm } from '@/composables/useInvoiceForm'
import InvoiceDocumentPreview from '@/components/invoice/InvoiceDocumentPreview.vue'

export default defineComponent({
  name: 'InvoiceCreatePage',
  components: {
    InvoiceDocumentPreview
  },
  setup() {
    const {
        activeMobileTab,
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
        onClientChange,
        addItem,
        removeItem,
        formatDate,
        saveDraft,
        handleSubmit,
        discardRestoredDraft
    } = useInvoiceForm({ mode: 'create' })
    return {
      activeMobileTab,
      InvoiceDocumentPreview,
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
      onClientChange,
      addItem,
      removeItem,
      formatDate,
      saveDraft,
      handleSubmit,
      discardRestoredDraft
    }
  }
})
