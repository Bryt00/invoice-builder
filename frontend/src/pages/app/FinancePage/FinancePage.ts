import { defineComponent, ref, reactive, onMounted, onBeforeUnmount, computed, nextTick } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import { useAuthStore } from '@/stores/auth'
import Chart from 'chart.js/auto'

export default defineComponent({
  name: 'FinancePage',
  setup() {
    const authStore = useAuthStore()
    const { showFlash } = useFlash()

    // Data State
    const transactions = ref<any[]>([])
    const categories = ref<any[]>([])
    const analytics = ref<any>({
      summary: {
        total_income: 0,
        total_expenses: 0,
        net_profit: 0,
        profit_margin_pct: 0
      },
      monthly_trend: [],
      category_distribution: [],
      receivables: {
        unpaid_total: 0,
        unpaid_count: 0,
        overdue_total: 0,
        overdue_count: 0
      }
    })

    const loading = ref(true)
    const currencySymbol = computed(() => authStore.currencySymbol || '$')

    // Canvas references & Chart instances
    const cashFlowCanvasRef = ref<HTMLCanvasElement | null>(null)
    const catCanvasRef = ref<HTMLCanvasElement | null>(null)
    let cashFlowChartInstance: any = null
    let catChartInstance: any = null

    // Filter toolbar state
    const filters = reactive({
      type: 'all',
      category_id: 'all',
      date_preset: 'all',
      start_date: '',
      end_date: '',
      search: ''
    })

    // Pagination state
    const currentPage = ref(1)
    const pageSize = ref(10)

    // Transaction modal state
    const isTxnModalOpen = ref(false)
    const editingTxnId = ref<string | null>(null)
    const savingTxn = ref(false)
    const txnForm = reactive({
      id: '',
      title: '',
      amount: null as number | null,
      type: 'income',
      category_id: '',
      payee_or_payer: '',
      transaction_date: '',
      description: '',
      currency: ''
    })

    // Category modal state
    const isCatModalOpen = ref(false)
    const savingCat = ref(false)
    const catForm = reactive({
      name: '',
      type: 'expense',
      icon: 'payments'
    })

    const availableIcons = [
      'payments', 'receipt_long', 'work', 'laptop_mac', 'group',
      'campaign', 'inventory_2', 'wifi', 'flight_takeoff', 'policy',
      'savings', 'account_balance', 'local_shipping', 'restaurant',
      'shopping_bag', 'build', 'credit_card', 'domain'
    ]

    // Delete confirmation modal state
    const isDeleteModalOpen = ref(false)
    const txnToDelete = ref<any>(null)
    const deletingTxn = ref(false)

    // Computed filtered categories based on transaction form type
    const formCategories = computed(() => {
      return categories.value.filter(c => c.type === txnForm.type)
    })

    // Filtered transactions for the ledger
    const filteredTransactions = computed(() => {
      let result = transactions.value

      if (filters.type !== 'all') {
        result = result.filter(t => t.type === filters.type)
      }
      if (filters.category_id !== 'all') {
        result = result.filter(t => t.category_id === filters.category_id)
      }
      if (filters.start_date) {
        const start = new Date(filters.start_date)
        result = result.filter(t => new Date(t.transaction_date) >= start)
      }
      if (filters.end_date) {
        const end = new Date(filters.end_date)
        end.setHours(23, 59, 59, 999)
        result = result.filter(t => new Date(t.transaction_date) <= end)
      }
      if (filters.search) {
        const s = filters.search.toLowerCase()
        result = result.filter(t =>
          (t.title && t.title.toLowerCase().includes(s)) ||
          (t.description && t.description.toLowerCase().includes(s)) ||
          (t.payee_or_payer && t.payee_or_payer.toLowerCase().includes(s)) ||
          (t.category?.name && t.category.name.toLowerCase().includes(s))
        )
      }
      return result
    })

    // Paginated transactions
    const paginatedTransactions = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      return filteredTransactions.value.slice(start, start + pageSize.value)
    })

    const totalPages = computed(() => {
      return Math.ceil(filteredTransactions.value.length / pageSize.value) || 1
    })

    // Pagination actions
    function prevPage() {
      if (currentPage.value > 1) {
        currentPage.value--
      }
    }

    function nextPage() {
      if (currentPage.value < totalPages.value) {
        currentPage.value++
      }
    }

    function goToPage(page: number) {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
      }
    }

    // Filter controls
    function setTypeFilter(type: string) {
      filters.type = type
      currentPage.value = 1
    }

    function setDatePreset(preset: string) {
      filters.date_preset = preset
      currentPage.value = 1
      const now = new Date()

      if (preset === 'this_month') {
        const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
        filters.start_date = firstDay.toISOString().split('T')[0]
        filters.end_date = now.toISOString().split('T')[0]
      } else if (preset === 'last_30_days') {
        const past30 = new Date()
        past30.setDate(now.getDate() - 30)
        filters.start_date = past30.toISOString().split('T')[0]
        filters.end_date = now.toISOString().split('T')[0]
      } else if (preset === 'this_year') {
        const firstDayOfYear = new Date(now.getFullYear(), 0, 1)
        filters.start_date = firstDayOfYear.toISOString().split('T')[0]
        filters.end_date = now.toISOString().split('T')[0]
      } else if (preset === 'all') {
        filters.start_date = ''
        filters.end_date = ''
      }
    }

    // Chart rendering
    function renderCharts() {
      const isDark = document.documentElement.classList.contains('dark')
      const textColor = isDark ? '#94a3b8' : '#64748b'
      const gridColor = isDark ? 'rgba(255, 255, 255, 0.06)' : 'rgba(0, 0, 0, 0.06)'

      // 1. Cash Flow Trend Chart
      if (cashFlowCanvasRef.value) {
        if (cashFlowChartInstance) {
          cashFlowChartInstance.destroy()
        }

        const trendData = analytics.value.monthly_trend || []
        const labels = trendData.map((d: any) => d.month)
        const incomeData = trendData.map((d: any) => d.income)
        const expenseData = trendData.map((d: any) => d.expenses)

        const ctx = cashFlowCanvasRef.value.getContext('2d')
        if (ctx) {
          const incomeGradient = ctx.createLinearGradient(0, 0, 0, 260)
          incomeGradient.addColorStop(0, 'rgba(16, 185, 129, 0.85)')
          incomeGradient.addColorStop(1, 'rgba(16, 185, 129, 0.35)')

          const expenseGradient = ctx.createLinearGradient(0, 0, 0, 260)
          expenseGradient.addColorStop(0, 'rgba(244, 63, 94, 0.85)')
          expenseGradient.addColorStop(1, 'rgba(244, 63, 94, 0.35)')

          cashFlowChartInstance = new Chart(cashFlowCanvasRef.value, {
            type: 'bar',
            data: {
              labels: labels.length ? labels : ['No Data'],
              datasets: [
                {
                  label: 'Income',
                  data: incomeData.length ? incomeData : [0],
                  backgroundColor: incomeGradient,
                  borderRadius: 8,
                  barPercentage: 0.65,
                  categoryPercentage: 0.6
                },
                {
                  label: 'Expenses',
                  data: expenseData.length ? expenseData : [0],
                  backgroundColor: expenseGradient,
                  borderRadius: 8,
                  barPercentage: 0.65,
                  categoryPercentage: 0.6
                }
              ]
            },
            options: {
              responsive: true,
              maintainAspectRatio: false,
              plugins: {
                legend: {
                  position: 'top',
                  align: 'end',
                  labels: {
                    color: textColor,
                    usePointStyle: true,
                    boxWidth: 8,
                    font: { family: 'inherit', size: 12, weight: 600 }
                  }
                },
                tooltip: {
                  backgroundColor: isDark ? '#1e293b' : '#ffffff',
                  titleColor: isDark ? '#f8fafc' : '#0f172a',
                  bodyColor: isDark ? '#cbd5e1' : '#334155',
                  borderColor: isDark ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)',
                  borderWidth: 1,
                  padding: 12,
                  boxPadding: 6,
                  usePointStyle: true,
                  callbacks: {
                    label: (context: any) => {
                      const val = context.raw || 0
                      return ` ${context.dataset.label}: ${currencySymbol.value}${val.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
                    }
                  }
                }
              },
              scales: {
                x: {
                  grid: { display: false },
                  ticks: { color: textColor, font: { family: 'inherit', size: 11 } }
                },
                y: {
                  grid: { color: gridColor },
                  ticks: {
                    color: textColor,
                    font: { family: 'inherit', size: 11 },
                    callback: (value: any) => `${currencySymbol.value}${value >= 1000 ? (value / 1000).toFixed(0) + 'k' : value}`
                  }
                }
              }
            }
          })
        }
      }

      // 2. Category Distribution Doughnut Chart
      if (catCanvasRef.value) {
        if (catChartInstance) {
          catChartInstance.destroy()
        }

        const catData = analytics.value.category_distribution || []
        const palette = [
          '#6366f1', '#ec4899', '#f59e0b', '#10b981', '#3b82f6',
          '#8b5cf6', '#14b8a6', '#f43f5e', '#06b6d4', '#64748b'
        ]

        const labels = catData.map((c: any) => c.name)
        const amounts = catData.map((c: any) => c.amount)
        const colors = catData.map((_: any, idx: number) => palette[idx % palette.length])

        if (catData.length === 0) {
          labels.push('No Expenses')
          amounts.push(1)
          colors.push(isDark ? '#334155' : '#e2e8f0')
        }

        catChartInstance = new Chart(catCanvasRef.value, {
          type: 'doughnut',
          data: {
            labels,
            datasets: [
              {
                data: amounts,
                backgroundColor: colors,
                borderWidth: 0,
                hoverOffset: 6
              }
            ]
          },
          options: {
            responsive: true,
            maintainAspectRatio: false,
            cutout: '72%',
            plugins: {
              legend: {
                display: false
              },
              tooltip: {
                backgroundColor: isDark ? '#1e293b' : '#ffffff',
                titleColor: isDark ? '#f8fafc' : '#0f172a',
                bodyColor: isDark ? '#cbd5e1' : '#334155',
                borderColor: isDark ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)',
                borderWidth: 1,
                padding: 10,
                boxPadding: 6,
                callbacks: {
                  label: (context: any) => {
                    if (catData.length === 0) return ' No expense data'
                    const val = context.raw || 0
                    const item = catData[context.dataIndex]
                    const pct = item?.percentage ? item.percentage.toFixed(1) : '0'
                    return ` ${context.label}: ${currencySymbol.value}${val.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })} (${pct}%)`
                  }
                }
              }
            }
          }
        })
      }
    }

    // Fetch all finance dashboard data
    async function fetchData() {
      loading.value = true
      try {
        const [txnRes, catRes, anaRes] = await Promise.all([
          api.get('/finance/transactions'),
          api.get('/finance/categories'),
          api.get('/finance/analytics')
        ])

        if (txnRes.data?.transactions) {
          transactions.value = txnRes.data.transactions
        }
        if (catRes.data?.categories) {
          categories.value = catRes.data.categories
        }
        if (anaRes.data?.analytics) {
          analytics.value = anaRes.data.analytics
        }

        await nextTick()
        renderCharts()
      } catch (err: any) {
        showFlash('Failed to load finance data', 'error')
      } finally {
        loading.value = false
      }
    }

    onMounted(fetchData)

    onBeforeUnmount(() => {
      if (cashFlowChartInstance) cashFlowChartInstance.destroy()
      if (catChartInstance) catChartInstance.destroy()
    })

    // Formatters
    function formatDate(dateStr: string) {
      if (!dateStr) return '-'
      const d = new Date(dateStr)
      if (isNaN(d.getTime())) return dateStr
      return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
    }

    function formatNumber(val: number | undefined) {
      if (val === undefined || val === null || isNaN(val)) return '0.00'
      return val.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    }

    // Modal: Open Create Transaction
    function openCreateTxnModal(type: 'income' | 'expense' = 'income') {
      editingTxnId.value = null
      txnForm.id = ''
      txnForm.title = ''
      txnForm.amount = null
      txnForm.type = type
      txnForm.payee_or_payer = ''
      txnForm.transaction_date = new Date().toISOString().split('T')[0]
      txnForm.description = ''
      txnForm.currency = authStore.currencySymbol ? '' : 'USD'

      // Pick first matching category if available
      const matchingCats = categories.value.filter(c => c.type === type)
      txnForm.category_id = matchingCats.length > 0 ? matchingCats[0].id : ''

      isTxnModalOpen.value = true
    }

    // Modal: Open Edit Transaction
    function openEditTxnModal(txn: any) {
      editingTxnId.value = txn.id
      txnForm.id = txn.id
      txnForm.title = txn.title || ''
      txnForm.amount = txn.amount
      txnForm.type = txn.type || 'income'
      txnForm.category_id = txn.category_id || ''
      txnForm.payee_or_payer = txn.payee_or_payer || ''
      txnForm.transaction_date = txn.transaction_date ? txn.transaction_date.split('T')[0] : new Date().toISOString().split('T')[0]
      txnForm.description = txn.description || ''
      txnForm.currency = txn.currency || ''

      isTxnModalOpen.value = true
    }

    function closeTxnModal() {
      isTxnModalOpen.value = false
      editingTxnId.value = null
    }

    function switchTxnType(type: 'income' | 'expense') {
      txnForm.type = type
      const matching = categories.value.filter(c => c.type === type)
      txnForm.category_id = matching.length > 0 ? matching[0].id : ''
    }

    // Save Transaction (Create / Update)
    async function saveTransaction() {
      if (!txnForm.title || !txnForm.title.trim()) {
        showFlash('Please enter a transaction title', 'error')
        return
      }
      if (!txnForm.amount || txnForm.amount <= 0) {
        showFlash('Amount must be greater than 0', 'error')
        return
      }
      if (!txnForm.category_id) {
        showFlash('Please select a category', 'error')
        return
      }

      savingTxn.value = true
      try {
        if (editingTxnId.value) {
          await api.put('/finance/transactions', {
            id: editingTxnId.value,
            title: txnForm.title,
            amount: Number(txnForm.amount),
            type: txnForm.type,
            category_id: txnForm.category_id,
            payee_or_payer: txnForm.payee_or_payer,
            transaction_date: txnForm.transaction_date,
            description: txnForm.description,
            currency: txnForm.currency
          })
          showFlash('Transaction updated successfully', 'success')
        } else {
          await api.post('/finance/transactions', {
            title: txnForm.title,
            amount: Number(txnForm.amount),
            type: txnForm.type,
            category_id: txnForm.category_id,
            payee_or_payer: txnForm.payee_or_payer,
            transaction_date: txnForm.transaction_date,
            description: txnForm.description,
            currency: txnForm.currency
          })
          showFlash('Transaction recorded successfully', 'success')
        }

        closeTxnModal()
        await fetchData()
      } catch (err: any) {
        const msg = err.response?.data?.message || err.response?.data?.error || 'Failed to save transaction'
        showFlash(msg, 'error')
      } finally {
        savingTxn.value = false
      }
    }

    // Modal: Category Quick Add
    function openCreateCatModal(type: 'income' | 'expense' = 'expense') {
      catForm.name = ''
      catForm.type = type
      catForm.icon = 'payments'
      isCatModalOpen.value = true
    }

    function closeCatModal() {
      isCatModalOpen.value = false
    }

    async function saveCategory() {
      if (!catForm.name || !catForm.name.trim()) {
        showFlash('Please enter a category name', 'error')
        return
      }

      savingCat.value = true
      try {
        const res = await api.post('/finance/categories', {
          name: catForm.name.trim(),
          type: catForm.type,
          icon: catForm.icon
        })

        if (res.data?.category) {
          categories.value.push(res.data.category)
          // If transaction modal is open, auto-select this new category
          if (isTxnModalOpen.value && txnForm.type === catForm.type) {
            txnForm.category_id = res.data.category.id
          }
        }

        showFlash('Category created successfully', 'success')
        closeCatModal()
      } catch (err: any) {
        const msg = err.response?.data?.message || 'Failed to create category'
        showFlash(msg, 'error')
      } finally {
        savingCat.value = false
      }
    }

    // Modal: Delete Transaction Confirmation
    function confirmDeleteTxn(txn: any) {
      txnToDelete.value = txn
      isDeleteModalOpen.value = true
    }

    function closeDeleteModal() {
      isDeleteModalOpen.value = false
      txnToDelete.value = null
    }

    async function deleteTransaction() {
      if (!txnToDelete.value) return

      deletingTxn.value = true
      try {
        await api.post('/finance/transactions/delete', { id: txnToDelete.value.id })
        showFlash('Transaction deleted', 'success')
        closeDeleteModal()
        await fetchData()
      } catch (err: any) {
        showFlash('Failed to delete transaction', 'error')
      } finally {
        deletingTxn.value = false
      }
    }

    // CSV Export
    async function exportCsv() {
      try {
        const params = new URLSearchParams()
        if (filters.type !== 'all') params.append('type', filters.type)
        if (filters.category_id !== 'all') params.append('category', filters.category_id)
        if (filters.start_date) params.append('start_date', filters.start_date)
        if (filters.end_date) params.append('end_date', filters.end_date)
        if (filters.search) params.append('search', filters.search)

        const res = await api.get(`/finance/export?${params.toString()}`, { responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([res.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `finance_ledger_${new Date().toISOString().split('T')[0]}.csv`)
        document.body.appendChild(link)
        link.click()
        link.parentNode?.removeChild(link)
        showFlash('Ledger CSV exported successfully', 'success')
      } catch (err: any) {
        showFlash('Failed to export CSV', 'error')
      }
    }

    return {
      authStore,
      currencySymbol,
      transactions,
      categories,
      formCategories,
      analytics,
      loading,
      filters,
      filteredTransactions,
      paginatedTransactions,
      currentPage,
      pageSize,
      totalPages,
      prevPage,
      nextPage,
      goToPage,
      setTypeFilter,
      setDatePreset,
      cashFlowCanvasRef,
      catCanvasRef,
      formatDate,
      formatNumber,
      exportCsv,
      // Modals
      isTxnModalOpen,
      editingTxnId,
      savingTxn,
      txnForm,
      openCreateTxnModal,
      openEditTxnModal,
      closeTxnModal,
      switchTxnType,
      saveTransaction,
      // Category modal
      isCatModalOpen,
      savingCat,
      catForm,
      availableIcons,
      openCreateCatModal,
      closeCatModal,
      saveCategory,
      // Delete modal
      isDeleteModalOpen,
      txnToDelete,
      deletingTxn,
      confirmDeleteTxn,
      closeDeleteModal,
      deleteTransaction
    }
  }
})
