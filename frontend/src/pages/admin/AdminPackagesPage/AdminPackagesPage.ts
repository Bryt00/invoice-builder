import { defineComponent,  ref, reactive, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import { useConfirm } from '@/composables/useConfirm'

export default defineComponent({
  name: 'AdminPackagesPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()

    const packages = ref([])

    const loading = ref(true)

    const saving = ref(false)

    const showCreateModal = ref(false)

    const showEditModal = ref(false)

    const formCreate = reactive({ name: '', slug: '', description: '', price: 0, currency: 'GHS', credits_granted: 0, badge_tag: '', is_active: true })

    const formEdit = reactive({ id: '', name: '', slug: '', description: '', price: 0, currency: 'GHS', credits_granted: 0, badge_tag: '', is_active: true })

    onMounted(() => {
      fetchPackages()
    })

    const currencyMap = {
      USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', NGN: '₦'
    }

    function currencySymbol(code) {
      return currencyMap[code] || code
    }

    async function fetchPackages() {
      loading.value = true
      try {
        const res = await api.get('/admin/packages')
        packages.value = res.data.packages || []
      } catch (err: any) {
        showFlash('Failed to load packages', 'error')
      } finally {
        loading.value = false
      }
    }

    function openCreateModal() {
      formCreate.name = ''
      formCreate.slug = ''
      formCreate.description = ''
      formCreate.price = 0
      formCreate.currency = 'GHS'
      formCreate.credits_granted = 0
      formCreate.badge_tag = ''
      formCreate.is_active = true
      showCreateModal.value = true
    }

    async function submitCreatePackage() {
      saving.value = true
      try {
        await api.post('/admin/packages', formCreate)
        showFlash('Package created successfully', 'success')
        showCreateModal.value = false
        fetchPackages()
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to create package', 'error')
      } finally {
        saving.value = false
      }
    }

    function openEditModal(pkg) {
      formEdit.id = pkg.id
      formEdit.name = pkg.name
      formEdit.slug = pkg.slug
      formEdit.description = pkg.description
      formEdit.price = pkg.price / 100 // Convert back to float for input
      formEdit.currency = pkg.currency
      formEdit.credits_granted = pkg.credits_granted
      formEdit.badge_tag = pkg.badge_tag
      formEdit.is_active = pkg.is_active
      showEditModal.value = true
    }

    async function submitEditPackage() {
      saving.value = true
      try {
        await api.put('/admin/packages', formEdit)
        showFlash('Package updated successfully', 'success')
        showEditModal.value = false
        fetchPackages()
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to update package', 'error')
      } finally {
        saving.value = false
      }
    }

    async function togglePackageStatus(pkg) {
      const newStatus = !pkg.is_active
      try {
        await api.put('/admin/packages', {
          ...pkg,
          price: pkg.price / 100, // API expects float
          is_active: newStatus
        })
        pkg.is_active = newStatus
        showFlash(`Package ${newStatus ? 'activated' : 'deactivated'} successfully`, 'success')
      } catch (err: any) {
        showFlash('Failed to update package status', 'error')
      }
    }

    const { askConfirm } = useConfirm()

    async function deletePackage(id) {
      const ok = await askConfirm({
        title: 'Delete Credit Package',
        message: 'Are you sure you want to permanently delete this credit bundle?',
        confirmText: 'Delete Package',
        type: 'danger'
      })
      if (!ok) return

      try {
        await api.delete(`/admin/packages/${id}`)
        showFlash('Package deleted successfully', 'success')
        fetchPackages()
      } catch (err: any) {
        showFlash('Failed to delete package', 'error')
      }
    }
    return {
      FlashAlert,
      showFlash,
      packages,
      loading,
      saving,
      showCreateModal,
      showEditModal,
      formCreate,
      formEdit,
      currencyMap,
      currencySymbol,
      fetchPackages,
      openCreateModal,
      submitCreatePackage,
      openEditModal,
      submitEditPackage,
      togglePackageStatus,
      askConfirm,
      deletePackage
    }
  }
})
