import { defineComponent,  ref, onMounted } from 'vue'
import PublicPageLayout from '@/layouts/PublicPageLayout.vue'
import api from '@/utils/api'

export default defineComponent({
  name: 'RefundPolicyPage',
  components: {
    PublicPageLayout
  },
  setup() {
    const customContent = ref('')

    onMounted(async () => {
      try {
        const res = await api.get('/public/settings')
        if (res.data?.settings?.legal_refund) {
          customContent.value = res.data.settings.legal_refund
        }
      } catch (err: any) {
        // Fallback to default template
      }
    })
    return {
      PublicPageLayout,
      customContent
    }
  }
})
