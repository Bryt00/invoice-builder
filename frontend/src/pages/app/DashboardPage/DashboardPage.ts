import { defineComponent,  ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import api from '@/utils/api'

export default defineComponent({
  name: 'DashboardPage',
  setup() {
    const route = useRoute()

    const router = useRouter()

    const authStore = useAuthStore()

    const { showToast } = useToast()

    const loading = ref(true)

    const stats = ref({ balance: 0, total_purchased: 0, total_used: 0 })

    const recentInvoices = ref([])

    const loadingPackages = ref(true)

    const packages = ref([])

    const selectedPackage = ref('')

    const purchasing = ref(false)

    const usagePercent = computed(() => {
      const pur = stats.value.total_purchased || 0
      const usd = stats.value.total_used || 0
      if (pur === 0) return 0
      return Math.round((usd / pur) * 100)
    })

    function formatDate(dateStr) {
      if (!dateStr) return '-'
      return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
    }

    onMounted(async () => {
      const reference = route.query.reference || route.query.trxref
      if (reference) {
        try {
          await api.get(`/credits/topup/verify?reference=${reference}`)
          showToast('Payment verified! Your credits have been added.', 'success')
          
          // Clean up the URL
          const newQuery = { ...route.query }
          delete newQuery.reference
          delete newQuery.trxref
          router.replace({ query: newQuery })
        } catch (err: any) {
          console.error('Verification failed', err)
          showToast('Payment verification failed.', 'error')
        }
      }

      try {
        const [creditRes, invoicesRes, packagesRes] = await Promise.all([
          api.get('/credits/balance').catch(() => ({ data: { stats: {} } })),
          api.get('/invoices?limit=3').catch(() => ({ data: { invoices: [] } })),
          api.get('/credits/packages').catch(() => ({ data: { packages: [] } })),
        ])

        if (creditRes.data?.stats) stats.value = creditRes.data.stats
        if (invoicesRes.data?.invoices) recentInvoices.value = invoicesRes.data.invoices
        if (packagesRes.data?.packages) {
          packages.value = packagesRes.data.packages
          if (packages.value.length > 0) {
            selectedPackage.value = packages.value[0].id
          }
        }
      } finally {
        loading.value = false
        loadingPackages.value = false
      }
    })

    async function handleTopup() {
        if (!selectedPackage.value) return
        purchasing.value = true
        try {
            const res = await api.post('/credits/topup/initialize', { package_id: selectedPackage.value })
            if (res.data?.authorization_url) {
                window.location.href = res.data.authorization_url
            }
        } catch (err: any) {
            showToast(err.response?.data?.error || 'Failed to initialize payment', 'error')
        } finally {
            purchasing.value = false
        }
    }
    return {
      route,
      router,
      authStore,
      showToast,
      loading,
      stats,
      recentInvoices,
      loadingPackages,
      packages,
      selectedPackage,
      purchasing,
      usagePercent,
      formatDate,
      handleTopup
    }
  }
})
