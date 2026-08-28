<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <!-- Page Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">inventory_2</span>
          Pay-Per-Use Credit Bundles
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">Configure credit top-up packages available for user recharge on demand.</p>
      </div>
      <button @click="openCreateModal" class="inline-flex items-center gap-2 px-4 py-2.5 bg-amber-500 hover:bg-amber-400 text-on-primary font-label text-sm font-bold rounded-xl shadow-sm transition-all cursor-pointer">
        <span class="material-symbols-outlined text-[18px]">add</span>
        <span>New Credit Bundle</span>
      </button>
    </div>

    <!-- Package Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 relative min-h-[300px]">
      <!-- Loading Overlay -->
      <div v-if="loading" class="absolute inset-0 bg-surface/50 backdrop-blur-sm z-10 flex items-center justify-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
      </div>

      <template v-if="packages.length > 0">
        <div v-for="pkg in packages" :key="pkg.id" :class="[
          'glass-card rounded-2xl border p-6 flex flex-col justify-between shadow-xs transition-all',
          pkg.is_active ? 'border-amber-500/40 hover:border-amber-500/70' : 'border-outline-variant/40 opacity-75'
        ]">
          <div class="space-y-4">
            <div class="flex justify-between items-start">
              <div>
                <h3 class="font-headline text-xl font-bold text-on-surface">{{ pkg.name }}</h3>
                <span class="font-mono text-xs text-on-surface-variant">/{{ pkg.slug }}</span>
              </div>
              <span :class="[
                'px-2.5 py-1 text-xs font-semibold rounded-full border',
                pkg.is_active ? 'bg-emerald-500/10 text-emerald-600 border-emerald-500/30' : 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
              ]">
                {{ pkg.is_active ? 'Active' : 'Inactive' }}
              </span>
            </div>

            <div v-if="pkg.badge_tag">
              <span class="inline-block px-2.5 py-0.5 text-[10px] uppercase font-mono font-bold tracking-wider rounded-md bg-amber-500/20 text-amber-600 border border-amber-500/30">
                {{ pkg.badge_tag }}
              </span>
            </div>

            <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">{{ pkg.description }}</p>

            <div class="space-y-2 pt-2 border-t border-outline-variant/30 text-xs sm:text-sm font-body">
              <div class="flex justify-between items-center">
                <span class="text-on-surface-variant">Recharge Price:</span>
                <span class="font-headline font-extrabold text-base text-on-surface">
                  {{ currencySymbol(pkg.currency) }}{{ (pkg.price / 100).toFixed(2) }} 
                  <span class="text-xs font-mono text-on-surface-variant font-normal">({{ pkg.currency }})</span>
                </span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-on-surface-variant">Credits Granted:</span>
                <span class="font-bold text-amber-600 text-sm">{{ pkg.credits_granted }} Credits</span>
              </div>
            </div>
          </div>

          <div class="mt-6 pt-4 border-t border-outline-variant/30 flex items-center gap-2">
            <!-- Toggle Active -->
            <button @click="togglePackageStatus(pkg)" :class="[
              'flex-1 py-2 text-xs font-label font-bold rounded-xl transition-colors',
              pkg.is_active ? 'bg-rose-500/10 text-rose-600 border border-rose-500/30 hover:bg-rose-500/20' : 'bg-emerald-500/10 text-emerald-600 border border-emerald-500/30 hover:bg-emerald-500/20'
            ]">
              {{ pkg.is_active ? 'Deactivate' : 'Activate' }}
            </button>

            <!-- Edit Button -->
            <button @click="openEditModal(pkg)" class="p-2 text-on-surface-variant hover:text-amber-500 hover:bg-surface-container-high rounded-xl transition-colors" title="Edit Package">
              <span class="material-symbols-outlined text-[18px]">edit</span>
            </button>

            <!-- Delete Form -->
            <button @click="deletePackage(pkg.id)" class="p-2 text-on-surface-variant hover:text-rose-500 hover:bg-rose-500/10 rounded-xl transition-colors" title="Delete Package">
              <span class="material-symbols-outlined text-[18px]">delete</span>
            </button>
          </div>
        </div>
      </template>
      <div v-else-if="!loading" class="col-span-1 md:col-span-3 glass-card rounded-2xl p-12 text-center text-on-surface-variant space-y-3">
        <span class="material-symbols-outlined text-[48px] text-amber-500/60">inventory_2</span>
        <h3 class="font-headline text-lg font-bold text-on-surface">No Credit Bundles Defined</h3>
        <p class="text-sm">Click "New Credit Bundle" above to create pay-per-use recharge packages for users.</p>
      </div>
    </div>

    <!-- Create Package Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-surface border border-outline-variant/60 rounded-2xl max-w-lg w-full p-6 shadow-2xl space-y-5 animate-fade-in">
        <div class="flex justify-between items-center border-b border-outline-variant/40 pb-4">
          <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
            <span class="material-symbols-outlined text-amber-500">add_box</span>
            New Credit Top-Up Bundle
          </h3>
          <button @click="showCreateModal = false" class="text-on-surface-variant hover:text-on-surface">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>

        <form @submit.prevent="submitCreatePackage" class="space-y-4">
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Bundle Name</label>
            <input type="text" v-model="formCreate.name" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Slug ID</label>
            <input type="text" v-model="formCreate.slug" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-mono focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Description</label>
            <textarea v-model="formCreate.description" rows="2" class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none"></textarea>
          </div>
          <div class="grid grid-cols-3 gap-4">
            <div class="col-span-2">
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Price (Amount)</label>
              <input type="number" step="0.01" v-model.number="formCreate.price" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
            </div>
            <div>
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Currency</label>
              <select v-model="formCreate.currency" class="w-full px-3 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
                <option value="GHS">GHS (GH₵)</option>
                <option value="USD">USD ($)</option>
                <option value="EUR">EUR (€)</option>
                <option value="GBP">GBP (£)</option>
                <option value="NGN">NGN (₦)</option>
              </select>
            </div>
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Credits Granted</label>
            <input type="number" v-model.number="formCreate.credits_granted" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Badge Tag (Optional)</label>
            <input type="text" v-model="formCreate.badge_tag" class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div class="flex justify-end gap-3 pt-4 border-t border-outline-variant/40">
            <button type="button" @click="showCreateModal = false" class="px-4 py-2 bg-surface-container-high text-on-surface-variant hover:text-on-surface font-label text-xs font-semibold rounded-xl transition-colors">Cancel</button>
            <button type="submit" :disabled="saving" class="px-5 py-2 bg-amber-500 hover:bg-amber-400 text-on-primary font-label text-xs font-bold rounded-xl shadow-xs disabled:opacity-50">Save Bundle</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Package Modal -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-surface border border-outline-variant/60 rounded-2xl max-w-lg w-full p-6 shadow-2xl space-y-5 animate-fade-in">
        <div class="flex justify-between items-center border-b border-outline-variant/40 pb-4">
          <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
            <span class="material-symbols-outlined text-amber-500">edit_note</span>
            Edit Credit Top-Up Bundle
          </h3>
          <button @click="showEditModal = false" class="text-on-surface-variant hover:text-on-surface">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>

        <form @submit.prevent="submitEditPackage" class="space-y-4">
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Bundle Name</label>
            <input type="text" v-model="formEdit.name" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Slug ID</label>
            <input type="text" v-model="formEdit.slug" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-mono focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Description</label>
            <textarea v-model="formEdit.description" rows="2" class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none"></textarea>
          </div>
          <div class="grid grid-cols-3 gap-4">
            <div class="col-span-2">
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Price (Amount)</label>
              <input type="number" step="0.01" v-model.number="formEdit.price" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
            </div>
            <div>
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Currency</label>
              <select v-model="formEdit.currency" class="w-full px-3 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
                <option value="GHS">GHS (GH₵)</option>
                <option value="USD">USD ($)</option>
                <option value="EUR">EUR (€)</option>
                <option value="GBP">GBP (£)</option>
                <option value="NGN">NGN (₦)</option>
              </select>
            </div>
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Credits Granted</label>
            <input type="number" v-model.number="formEdit.credits_granted" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Badge Tag (Optional)</label>
            <input type="text" v-model="formEdit.badge_tag" class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div class="flex justify-end gap-3 pt-4 border-t border-outline-variant/40">
            <button type="button" @click="showEditModal = false" class="px-4 py-2 bg-surface-container-high text-on-surface-variant hover:text-on-surface font-label text-xs font-semibold rounded-xl transition-colors">Cancel</button>
            <button type="submit" :disabled="saving" class="px-5 py-2 bg-amber-500 hover:bg-amber-400 text-on-primary font-label text-xs font-bold rounded-xl shadow-xs disabled:opacity-50">Save Changes</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'

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
  } catch (err) {
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
  } catch (err) {
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
  } catch (err) {
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
  } catch (err) {
    showFlash('Failed to update package status', 'error')
  }
}

import { useConfirm } from '../../composables/useConfirm'

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
  } catch (err) {
    showFlash('Failed to delete package', 'error')
  }
}
</script>
