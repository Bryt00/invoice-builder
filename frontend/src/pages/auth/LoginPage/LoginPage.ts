import { defineComponent,  reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'

export default defineComponent({
  name: 'LoginPage',
  components: {
    FlashAlert
  },
  setup() {
    const authStore = useAuthStore()

    const router = useRouter()

    const route = useRoute()

    const { showFlash } = useFlash()

    const form = reactive({
      email: '',
      password: '',
    })

    const errors = ref<Record<string, string>>({})

    const showPassword = ref(false)

    async function handleSubmit() {
      errors.value = {}
      try {
        await authStore.login(form.email, form.password)
        showFlash('Welcome back!', 'success')
        
        if (authStore.user?.role?.name === 'Admin') {
          router.push('/user/admin/dashboard')
        } else {
          const redirectPath = String(route.query.redirect || '/user/dashboard')
          router.push(redirectPath)
        }
      } catch (err: any) {
        if (err.response?.status === 422 && typeof err.response.data?.error === 'object') {
          errors.value = err.response.data.error
        } else {
          showFlash(authStore.error || 'Failed to sign in', 'error')
        }
      }
    }
    return {
      FlashAlert,
      authStore,
      router,
      route,
      showFlash,
      form,
      errors,
      showPassword,
      handleSubmit
    }
  }
})
