import { defineComponent,  ref, reactive, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'

export default defineComponent({
  name: 'AdminSettingsPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()

    const loading = ref(true)

    const saving = ref(false)

    const form = reactive({
      maintenance_mode: false,
      support_email: '',
      default_signup_credits: '0',
      legal_terms: '',
      legal_privacy: '',
      legal_refund: '',
      legal_security: ''
    })

    onMounted(() => {
      fetchSettings()
    })

    async function fetchSettings() {
      loading.value = true
      try {
        const res = await api.get('/admin/settings')
        const s = res.data.settings
        if (s) {
          form.maintenance_mode = s.MaintenanceMode === true || s.MaintenanceMode === 'true'
          form.support_email = s.SupportContactEmail || ''
          form.default_signup_credits = s.DefaultSignupBonus || '0'
          form.legal_terms = s.LegalTerms || ''
          form.legal_privacy = s.LegalPrivacy || ''
          form.legal_refund = s.LegalRefund || ''
          form.legal_security = s.LegalSecurity || ''
        }
      } catch (err: any) {
        showFlash('Failed to load system settings', 'error')
      } finally {
        loading.value = false
      }
    }

    async function saveSettings() {
      saving.value = true
      try {
        await api.put('/admin/settings', {
          maintenance_mode: form.maintenance_mode ? 'true' : 'false',
          support_email: form.support_email,
          default_signup_credits: String(form.default_signup_credits),
          legal_terms: form.legal_terms,
          legal_privacy: form.legal_privacy,
          legal_refund: form.legal_refund,
          legal_security: form.legal_security
        })
        showFlash('System settings updated successfully', 'success')
      } catch (err: any) {
        showFlash('Failed to update system settings', 'error')
      } finally {
        saving.value = false
      }
    }
    return {
      FlashAlert,
      showFlash,
      loading,
      saving,
      form,
      fetchSettings,
      saveSettings
    }
  }
})
