import { defineComponent,  ref, onMounted } from 'vue'
import AppNavbar from '@/components/common/AppNavbar.vue'
import AppFooter from '@/components/common/AppFooter.vue'
import api from '@/utils/api'

export default defineComponent({
  name: 'LandingPage',
  components: {
    AppNavbar,
    AppFooter
  },
  setup() {
    const packages = ref([
      {
        id: 1,
        name: 'Starter Pack',
        price: 1500,
        credits_granted: 15,
        features: ['15 Invoice PDF Downloads', 'Client Management', 'Finance Ledger Sync'],
      },
      {
        id: 2,
        name: 'Professional Pack',
        price: 3500,
        credits_granted: 50,
        badge_tag: 'Most Popular',
        features: ['50 Invoice PDF Downloads', 'Client Management', 'Finance Ledger Sync', 'Custom Business Branding'],
      },
      {
        id: 3,
        name: 'Agency Bundle',
        price: 9500,
        credits_granted: 150,
        features: ['150 Invoice PDF Downloads', 'Unlimited Clients', 'Full Financial Stats & CSV Export', 'Priority Support'],
      },
    ])

    onMounted(async () => {
      try {
        const res = await api.get('/credits/packages')
        if (res.data?.packages && res.data.packages.length > 0) {
          packages.value = res.data.packages.map((pkg: any) => ({
            ...pkg,
            features: pkg.description ? pkg.description.split('\n').filter(Boolean) : ['Invoice PDF Generation', 'Client Management'],
          }))
        }
      } catch (err: any) {
        // Fallback to static package defaults
      }
    })
    return {
      AppNavbar,
      AppFooter,
      packages
    }
  }
})
