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

    const showPassword = ref(false)

    async function handleSubmit() {
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
        showFlash(authStore.error || 'Failed to sign in', 'error')
      }
    }
    return {
      FlashAlert,
      authStore,
      router,
      route,
      showFlash,
      form,
      showPassword,
      handleSubmit
    }
  }
})
