import { defineComponent,  ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/utils/api'

export default defineComponent({
  name: 'PublicInvoicePage',
  setup() {
    const route = useRoute()

    const invoice = ref(null)

    const profile = ref(null)

    const loading = ref(true)

    const token = ref('')

    onMounted(async () => {
      token.value = String(route.params.token || route.query.token || '')
      if (!token.value) return

      try {
        const res = await api.get(`/invoices/public?token=${token.value}`)
        if (res.data?.invoice) invoice.value = res.data.invoice
        if (res.data?.profile) profile.value = res.data.profile
      } finally {
        loading.value = false
      }
    })
    return {
      api,
      route,
      invoice,
      profile,
      loading,
      token
    }
  }
})
