<template>
  <div>
    <!-- Header Section -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
            <h2 class="font-headline text-3xl font-bold text-on-surface mb-1">Clients Directory</h2>
            <p class="font-body text-base text-on-surface-variant">Manage your saved client billing details and addresses.</p>
        </div>
        <button @click="toggleAddClient" class="btn-auth-submit text-on-primary bg-primary rounded-xl px-5 py-2.5 font-label text-sm font-semibold flex items-center gap-2 shrink-0 self-start md:self-auto cursor-pointer shadow-md hover:bg-on-primary-fixed-variant transition-colors">
            <span class="material-symbols-outlined text-[20px]">person_add</span>
            Add New Client
        </button>
    </header>

    <!-- Add Client Form Accordion -->
    <div v-show="isAddingClient" class="glass-card rounded-2xl p-6 mb-8 border border-outline-variant/60">
        <h3 class="font-headline text-lg font-bold text-on-surface mb-4 flex items-center gap-2">
            <span class="material-symbols-outlined text-primary text-[20px]">person_add</span>
            New Client Details
        </h3>
        
        <form @submit.prevent="handleCreateClient" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="block font-label text-sm font-medium text-on-surface mb-1.5">Full Name</label>
                    <div class="relative">
                        <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">person</span>
                        <input type="text" v-model="form.name" placeholder="Client Name" required
                               class="w-full pl-10 pr-3 py-2 bg-surface-container-lowest/80 border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                    </div>
                </div>
                <div>
                    <label class="block font-label text-sm font-medium text-on-surface mb-1.5">Email Address</label>
                    <div class="relative">
                        <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">mail</span>
                        <input type="email" v-model="form.email" placeholder="client@company.com" required
                               class="w-full pl-10 pr-3 py-2 bg-surface-container-lowest/80 border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                    </div>
                </div>
                <div>
                    <label class="block font-label text-sm font-medium text-on-surface mb-1.5">Phone Number</label>
                    <div class="relative">
                        <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">call</span>
                        <input type="tel" v-model="form.phone" placeholder="+1 (555) 000-0000"
                               class="w-full pl-10 pr-3 py-2 bg-surface-container-lowest/80 border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                    </div>
                </div>
                <div>
                    <label class="block font-label text-sm font-medium text-on-surface mb-1.5">Company Name</label>
                    <div class="relative">
                        <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">domain</span>
                        <input type="text" v-model="form.company" placeholder="Acme Inc."
                               class="w-full pl-10 pr-3 py-2 bg-surface-container-lowest/80 border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                    </div>
                </div>
                <div>
                    <label class="block font-label text-sm font-medium text-on-surface mb-1.5">Tax ID / VAT No.</label>
                    <div class="relative">
                        <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">receipt</span>
                        <input type="text" v-model="form.tax_id" placeholder="VAT-987654"
                               class="w-full pl-10 pr-3 py-2 bg-surface-container-lowest/80 border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                    </div>
                </div>
                <div>
                    <label class="block font-label text-sm font-medium text-on-surface mb-1.5">Billing Address</label>
                    <textarea v-model="form.address" rows="1" placeholder="123 Business Rd, Suite 100"
                              class="w-full px-3 py-2 bg-surface-container-lowest/80 border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary"></textarea>
                </div>
            </div>
            <div class="flex justify-end gap-3 pt-2">
                <button type="button" @click="toggleAddClient" class="px-4 py-2 rounded-xl border border-outline-variant/60 font-label text-sm font-medium text-on-surface-variant hover:bg-surface-container-low cursor-pointer transition-colors">
                    Cancel
                </button>
                <button type="submit" :disabled="saving" class="btn-auth-submit bg-primary px-5 py-2 rounded-xl font-label text-sm font-semibold text-on-primary shadow-sm hover:bg-on-primary-fixed-variant transition-colors disabled:opacity-50">
                    Save Client
                </button>
            </div>
        </form>
    </div>

    <!-- Clients Directory Table -->
    <div class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden">
        <div class="px-6 py-4 border-b border-outline-variant/40 flex justify-between items-center bg-surface-container-lowest/50">
            <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-[20px]">contacts</span>
                Saved Clients
            </h3>
        </div>
        <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse min-w-[700px]">
                <thead>
                    <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
                        <th class="px-6 py-3.5">Client Name</th>
                        <th class="px-6 py-3.5">Company</th>
                        <th class="px-6 py-3.5">Email</th>
                        <th class="px-6 py-3.5">Phone</th>
                        <th class="px-6 py-3.5 text-right">Actions</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
                    <tr v-if="loading">
                        <td colspan="5" class="px-6 py-12 text-center text-outline text-sm">Loading clients...</td>
                    </tr>
                    <template v-else-if="clients.length > 0">
                        <template v-for="client in clients" :key="client.id">
                            <!-- Display Row -->
                            <tr class="hover:bg-surface-container-low/40 transition-colors">
                                <td class="px-6 py-4 font-semibold text-on-surface flex items-center gap-3">
                                    <div class="w-8 h-8 rounded-full bg-primary-container/30 text-primary flex items-center justify-center text-xs font-bold uppercase">
                                        {{ client.name.charAt(0) }}
                                    </div>
                                    {{ client.name }}
                                </td>
                                <td class="px-6 py-4 text-on-surface-variant">{{ client.company || '-' }}</td>
                                <td class="px-6 py-4 text-on-surface-variant">{{ client.email }}</td>
                                <td class="px-6 py-4 text-on-surface-variant">{{ client.phone || '-' }}</td>
                                <td class="px-6 py-4 text-right">
                                    <div class="flex items-center justify-end gap-1">
                                        <button @click="toggleEditClient(client.id)" class="p-1.5 text-outline hover:text-primary transition-colors rounded-lg hover:bg-primary-container/30 cursor-pointer" title="Edit Client">
                                            <span class="material-symbols-outlined text-[18px]">edit</span>
                                        </button>
                                        <button @click="deleteClient(client.id)" class="p-1.5 text-outline hover:text-error transition-colors rounded-lg hover:bg-error-container/30 cursor-pointer" title="Delete Client">
                                            <span class="material-symbols-outlined text-[18px]">delete</span>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                            
                            <!-- Inline Edit Row -->
                            <tr v-if="editingClientId === client.id" class="bg-surface-container-lowest/60">
                                <td colspan="5" class="px-6 py-4">
                                    <form @submit.prevent="handleEditClient(client.id)" class="space-y-4">
                                        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                                            <div>
                                                <label class="block font-label text-xs font-medium text-on-surface-variant mb-1">Name</label>
                                                <input type="text" v-model="editForm.name" required
                                                       class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                            </div>
                                            <div>
                                                <label class="block font-label text-xs font-medium text-on-surface-variant mb-1">Email</label>
                                                <input type="email" v-model="editForm.email" required
                                                       class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                            </div>
                                            <div>
                                                <label class="block font-label text-xs font-medium text-on-surface-variant mb-1">Phone</label>
                                                <input type="tel" v-model="editForm.phone"
                                                       class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                            </div>
                                            <div>
                                                <label class="block font-label text-xs font-medium text-on-surface-variant mb-1">Company</label>
                                                <input type="text" v-model="editForm.company"
                                                       class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                            </div>
                                            <div>
                                                <label class="block font-label text-xs font-medium text-on-surface-variant mb-1">Tax ID</label>
                                                <input type="text" v-model="editForm.tax_id"
                                                       class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                            </div>
                                            <div>
                                                <label class="block font-label text-xs font-medium text-on-surface-variant mb-1">Address</label>
                                                <input type="text" v-model="editForm.address"
                                                       class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/70 rounded-xl font-body text-sm text-on-surface outline-none focus:border-primary">
                                            </div>
                                        </div>
                                        <div class="flex justify-end gap-3 pt-1">
                                            <button type="button" @click="editingClientId = null" class="px-4 py-2 rounded-xl border border-outline-variant/60 font-label text-xs font-medium text-on-surface-variant hover:bg-surface-container-low cursor-pointer transition-colors">
                                                Cancel
                                            </button>
                                            <button type="submit" :disabled="saving" class="btn-auth-submit bg-primary px-5 py-2 rounded-xl font-label text-xs font-semibold text-on-primary cursor-pointer shadow-sm hover:bg-on-primary-fixed-variant transition-colors disabled:opacity-50">
                                                Save Changes
                                            </button>
                                        </div>
                                    </form>
                                </td>
                            </tr>
                        </template>
                    </template>
                    <tr v-else>
                        <td colspan="5" class="px-6 py-12 text-center text-on-surface-variant">
                            <span class="material-symbols-outlined text-4xl text-outline mb-2">person_off</span>
                            <p class="font-body text-base">No saved clients yet.</p>
                            <p class="font-body text-xs text-outline">Click "Add New Client" above to create your first client entry.</p>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'

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
  } catch (err) {
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
  } catch (err) {
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
    } catch (err) {
        showFlash(err.response?.data?.error || 'Failed to update client', 'error')
    } finally {
        saving.value = false
    }
}

import { useConfirm } from '../../composables/useConfirm'

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
    } catch (err) {
        showFlash(err.response?.data?.error || 'Failed to delete client', 'error')
    }
}
</script>
