import { defineComponent,  ref, reactive, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import { useConfirm } from '@/composables/useConfirm'

export default defineComponent({
  name: 'ClientsPage',
  setup() {
    const clients = ref([])

    const loading = ref(true)

    const saving = ref(false)

    const isAddingClient = ref(false)

    const editingClientId = ref(null)

    const { showFlash } = useFlash()

    const form = reactive({
      name: '',
      email: '',
      company: '',
      phone: '',
      tax_id: '',
      address: ''
    })

    const editForm = reactive({
      name: '',
      email: '',
      company: '',
      phone: '',
      tax_id: '',
      address: ''
    })

    async function fetchClients() {
      loading.value = true
      try {
        const res = await api.get('/clients')
        if (res.data?.clients) {
          clients.value = res.data.clients
        } else {
          clients.value = []
        }
      } catch (err: any) {
        showFlash('Failed to fetch clients', 'error')
      } finally {
        loading.value = false
      }
    }

    onMounted(fetchClients)

    function toggleAddClient() {
      isAddingClient.value = !isAddingClient.value
      if (!isAddingClient.value) {
        resetForm()
      }
    }

    function resetForm() {
      form.name = ''
      form.email = ''
      form.company = ''
      form.phone = ''
      form.tax_id = ''
      form.address = ''
    }

    async function handleCreateClient() {
      saving.value = true
      try {
        await api.post('/clients', form)
        showFlash('Client added successfully!', 'success')
        isAddingClient.value = false
        resetForm()
        fetchClients()
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to add client', 'error')
      } finally {
        saving.value = false
      }
    }

    function toggleEditClient(id) {
        if (editingClientId.value === id) {
            editingClientId.value = null
            return
        }
        const client = clients.value.find(c => c.id === id)
        if (client) {
            editForm.name = client.name || ''
            editForm.email = client.email || ''
            editForm.company = client.company || ''
            editForm.phone = client.phone || ''
            editForm.tax_id = client.tax_id || ''
            editForm.address = client.address || ''
            editingClientId.value = id
        }
    }

    async function handleEditClient(id) {
        saving.value = true
        try {
            await api.put('/clients', { id, ...editForm })
            showFlash('Client updated successfully!', 'success')
            editingClientId.value = null
            fetchClients()
        } catch (err: any) {
            showFlash(err.response?.data?.error || 'Failed to update client', 'error')
        } finally {
            saving.value = false
        }
    }

    const { askConfirm } = useConfirm()

    async function deleteClient(id) {
        const ok = await askConfirm({
            title: 'Delete Client',
            message: 'Are you sure you want to delete this client? This action cannot be undone.',
            confirmText: 'Delete Client',
            type: 'danger'
        })
        if (!ok) return

        try {
            await api.post('/clients/delete', { id })
            showFlash('Client deleted.', 'success')
            fetchClients()
        } catch (err: any) {
            showFlash(err.response?.data?.error || 'Failed to delete client', 'error')
        }
    }
    return {
      clients,
      loading,
      saving,
      isAddingClient,
      editingClientId,
      showFlash,
      form,
      editForm,
      fetchClients,
      toggleAddClient,
      resetForm,
      handleCreateClient,
      toggleEditClient,
      handleEditClient,
      askConfirm,
      deleteClient
    }
  }
})
