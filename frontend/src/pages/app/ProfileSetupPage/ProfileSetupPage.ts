import { defineComponent,  reactive, ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'

export default defineComponent({
  name: 'ProfileSetupPage',
  setup() {
    const router = useRouter()

    const authStore = useAuthStore()

    const saving = ref(false)

    const uploadingLogo = ref(false)

    const currencies = ref([])

    const profile = ref(null)

    const isEditing = ref(false)

    const flashMessage = ref('')

    const flashType = ref('success')

    const showFlash = (msg, type = 'success') => {
      flashMessage.value = msg
      flashType.value = type
      setTimeout(() => { flashMessage.value = '' }, 5000)
    }

    const isProfileComplete = computed(() => {
      return !!authStore.user?.is_profile_complete
    })

    const form = reactive({
      name: '',
      company_name: '',
      role: '',
      address: '',
      tax_id: '',
      default_currency: 'USD',
      registration_number: '',
      registration_date: '',
      business_type: '',
      registered_address: '',
      logo_url: ''
    })

    onMounted(async () => {
      try {
        const res = await api.get('/profile')
        if (res.data?.profile) {
          profile.value = res.data.profile
          populateForm(profile.value)
        } else {
          // Default name to user's name if profile not setup
          form.name = authStore.user?.name || ''
        }
        if (res.data?.currencies) {
          currencies.value = res.data.currencies
        }
      } catch (err: any) {
        console.error(err)
      }
    })

    function populateForm(p) {
      form.name = authStore.user?.name || ''
      form.company_name = p.company_name || ''
      form.role = p.role || ''
      form.address = p.address || ''
      form.tax_id = p.tax_id || ''
      form.default_currency = p.default_currency || 'USD'
      form.registration_number = p.registration_number || ''
      form.registration_date = p.registration_date ? p.registration_date.substring(0, 10) : ''
      form.business_type = p.business_type || ''
      form.registered_address = p.registered_address || ''
      form.logo_url = p.logo_url || ''
    }

    function formatDate(dateStr) {
      if (!dateStr) return '-'
      const d = new Date(dateStr)
      if (isNaN(d.getTime())) return dateStr
      return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
    }

    async function handleLogoUpload(event) {
      const file = event.target.files[0]
      if (!file) return

      const formData = new FormData()
      formData.append('logo', file)

      uploadingLogo.value = true
      try {
        const res = await api.post('/profile/logo', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        if (res.data?.logo_url) {
          form.logo_url = res.data.logo_url
          showFlash('Logo uploaded successfully!', 'success')
        }
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to upload logo', 'error')
      } finally {
        uploadingLogo.value = false
        // Clear the input so the same file can be uploaded again if needed
        event.target.value = ''
      }
    }

    async function handleSubmit() {
      saving.value = true
      try {
        await api.put('/profile', form)
        await authStore.fetchCurrentUser() // Update user name/state
        
        showFlash('Business profile updated successfully!', 'success')
        
        const wasEditing = isEditing.value
        isEditing.value = false
        
        // Refetch profile to update display view
        const res = await api.get('/profile')
        if (res.data?.profile) {
          profile.value = res.data.profile
          populateForm(profile.value)
        }

        if (!wasEditing) {
          // If it was their first time setting it up, route them to dashboard
          router.push({ name: 'dashboard' })
        }
      } catch (err: any) {
        let errorMsg = 'Failed to update profile'
        if (err.response?.data?.error) {
          const apiErr = err.response.data.error
          if (typeof apiErr === 'object') {
            errorMsg = String(Object.values(apiErr)[0]) // Extract the first validation error string
          } else {
            errorMsg = apiErr
          }
        }
        showFlash(errorMsg, 'error')
      } finally {
        saving.value = false
      }
    }
    return {
      api,
      router,
      authStore,
      saving,
      uploadingLogo,
      currencies,
      profile,
      isEditing,
      flashMessage,
      flashType,
      showFlash,
      isProfileComplete,
      form,
      populateForm,
      formatDate,
      handleLogoUpload,
      handleSubmit
    }
  }
})
