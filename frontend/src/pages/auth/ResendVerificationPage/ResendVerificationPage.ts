import { defineComponent,  ref } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'

export default defineComponent({
  name: 'ResendVerificationPage',
  components: {
    FlashAlert
  },
  setup() {
    const email = ref('')

    const loading = ref(false)

    const { showFlash } = useFlash()

    async function handleSubmit() {
      loading.value = true
      try {
        const res = await api.post('/auth/resend-verification', { email: email.value })
        showFlash(res.data?.message || 'If an account exists, a new activation email has been sent.', 'success', 8000)
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to resend activation link', 'error', 6000)
      } finally {
        loading.value = false
      }
    }
    return {
      FlashAlert,
      email,
      loading,
      showFlash,
      handleSubmit
    }
  }
})
