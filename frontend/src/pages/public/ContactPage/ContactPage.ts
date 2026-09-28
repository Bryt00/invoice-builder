import { defineComponent,  ref, onMounted } from 'vue'
import PublicPageLayout from '@/layouts/PublicPageLayout.vue'
import api from '@/utils/api'

export default defineComponent({
  name: 'ContactPage',
  components: {
    PublicPageLayout
  },
  setup() {
    const supportEmail = ref('support@teks-invoice.com')

    onMounted(async () => {
      try {
        const res = await api.get('/public/settings')
        if (res.data?.settings?.support_email) {
          supportEmail.value = res.data.settings.support_email
        }
      } catch (err: any) {
        // Keep default
      }
    })

    const submitForm = () => {
      alert('Thanks for your message! This is a demo form, but in a real app it would submit to our backend.')
    }
    return {
      PublicPageLayout,
      supportEmail,
      submitForm
    }
  }
})
