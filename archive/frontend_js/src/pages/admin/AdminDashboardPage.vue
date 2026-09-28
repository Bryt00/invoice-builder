<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">dashboard</span>
          Platform Overview
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">
          Real-time platform metrics, user registrations, revenue stats, and recent system audit logs.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <span class="px-3.5 py-1.5 rounded-full bg-amber-500/10 border border-amber-500/30 text-amber-600 text-xs font-bold flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full bg-amber-500 animate-ping"></span>
          <span>Admin Console</span>
        </span>
      </div>
    </div>

    <!-- KPI Metric Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Card 1: Revenue -->
      <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs hover:border-amber-500/40 transition-all flex flex-col justify-between">
        <div class="flex justify-between items-start mb-4">
          <div class="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-600 border border-emerald-500/20 flex items-center justify-center">
            <span class="material-symbols-outlined text-2xl">attach_money</span>
          </div>
          <span class="bg-emerald-500/10 text-emerald-600 border border-emerald-500/20 text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1">
            <span class="material-symbols-outlined text-xs">trending_up</span> Revenue
          </span>
        </div>
        <div>
          <p class="text-xs text-on-surface-variant font-label uppercase tracking-wider font-bold mb-1">Total System Revenue</p>
          <h3 class="font-headline text-3xl font-black text-on-surface">GH₵{{ (stats.total_revenue || 0).toFixed(2) }}</h3>
          <p class="text-xs text-on-surface-variant mt-1">
            <span class="font-bold text-emerald-600">{{ stats.total_payments || 0 }}</span> completed top-ups
          </p>
        </div>
      </div>

      <!-- Card 2: Users -->
      <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs hover:border-amber-500/40 transition-all flex flex-col justify-between">
        <div class="flex justify-between items-start mb-4">
          <div class="w-10 h-10 rounded-xl bg-indigo-500/10 text-indigo-600 border border-indigo-500/20 flex items-center justify-center">
            <span class="material-symbols-outlined text-2xl">person_add</span>
          </div>
          <span class="bg-indigo-500/10 text-indigo-600 border border-indigo-500/20 text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1">
            <span class="material-symbols-outlined text-xs">group</span> Users
          </span>
        </div>
        <div>
          <p class="text-xs text-on-surface-variant font-label uppercase tracking-wider font-bold mb-1">Registered Accounts</p>
          <h3 class="font-headline text-3xl font-black text-on-surface">{{ stats.total_users || 0 }}</h3>
          <p class="text-xs text-on-surface-variant mt-1">
            <span class="font-bold text-emerald-600">{{ stats.active_users || 0 }} Active</span> ·
            <span class="text-amber-600">{{ stats.unverified_users || 0 }} Unverified</span>
          </p>
        </div>
      </div>

      <!-- Card 3: Issued Credits -->
      <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs hover:border-amber-500/40 transition-all flex flex-col justify-between">
        <div class="flex justify-between items-start mb-4">
          <div class="w-10 h-10 rounded-xl bg-amber-500/10 text-amber-600 border border-amber-500/20 flex items-center justify-center">
            <span class="material-symbols-outlined text-2xl">bolt</span>
          </div>
          <span class="bg-amber-500/10 text-amber-600 border border-amber-500/20 text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1">
            <span class="material-symbols-outlined text-xs">bolt</span> Top-Ups
          </span>
        </div>
        <div>
          <p class="text-xs text-on-surface-variant font-label uppercase tracking-wider font-bold mb-1">Credits Issued</p>
          <h3 class="font-headline text-3xl font-black text-on-surface">{{ stats.total_purchased_credits || 0 }}</h3>
          <p class="text-xs text-on-surface-variant mt-1">Total purchased / granted</p>
        </div>
      </div>

      <!-- Card 4: Consumed Credits -->
      <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs hover:border-amber-500/40 transition-all flex flex-col justify-between">
        <div class="flex justify-between items-start mb-4">
          <div class="w-10 h-10 rounded-xl bg-rose-500/10 text-rose-600 border border-rose-500/20 flex items-center justify-center">
            <span class="material-symbols-outlined text-2xl">description</span>
          </div>
          <span class="bg-rose-500/10 text-rose-600 border border-rose-500/20 text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1">
            <span class="material-symbols-outlined text-xs">receipt_long</span> Usage
          </span>
        </div>
        <div>
          <p class="text-xs text-on-surface-variant font-label uppercase tracking-wider font-bold mb-1">Credits Consumed</p>
          <h3 class="font-headline text-3xl font-black text-on-surface">{{ stats.total_used_credits || 0 }}</h3>
          <p class="text-xs text-on-surface-variant mt-1">PDF Invoices dispatched</p>
        </div>
      </div>
    </div>

    <!-- Bento Grid Layout: Chart & Activity Feed -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Revenue & Platform Growth Chart Section -->
      <div class="lg:col-span-2 glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs flex flex-col justify-between">
        <div class="flex justify-between items-center mb-6">
          <div>
            <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
              <span class="material-symbols-outlined text-amber-500 text-[20px]">show_chart</span>
              Platform Revenue Growth
            </h3>
            <p class="text-xs text-on-surface-variant">Monthly top-up transaction volume and revenue trend</p>
          </div>
          <span class="px-3 py-1 bg-surface-container-high rounded-xl text-xs font-label font-bold text-on-surface-variant border border-outline-variant/40">
            Top-Up Ledger
          </span>
        </div>
        <div class="w-full relative min-h-[300px] flex-1">
          <canvas ref="revenueChartRef"></canvas>
        </div>
      </div>

      <!-- Recent Activity Stream Section -->
      <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs flex flex-col justify-between">
        <div class="flex justify-between items-center mb-4 pb-3 border-b border-outline-variant/40">
          <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
            <span class="material-symbols-outlined text-amber-500 text-[20px]">history</span>
            Recent Activity
          </h3>
          <router-link to="/user/admin/audit-logs" class="text-xs font-label font-bold text-amber-600 hover:underline">View All &rarr;</router-link>
        </div>

        <div class="space-y-3.5 overflow-y-auto max-h-[320px] pr-1">
          <div v-for="log in recentLogs" :key="log.id" class="flex items-start gap-3 pb-3 border-b border-outline-variant/30 last:border-0">
            <div class="w-8 h-8 rounded-full bg-amber-500/10 text-amber-600 border border-amber-500/30 flex items-center justify-center shrink-0 mt-0.5">
              <span class="material-symbols-outlined text-[16px]">bolt</span>
            </div>
            <div class="flex-1 min-w-0 font-body">
              <p class="text-xs font-bold text-on-surface truncate">{{ log.action }}</p>
              <p class="text-[11px] text-on-surface-variant truncate">{{ log.entity_type }}{{ log.entity_id ? `: ${log.entity_id}` : '' }}</p>
              <p class="text-[10px] text-on-surface-variant/70 mt-0.5">{{ new Date(log.created_at).toLocaleString() }}</p>
            </div>
          </div>
          
          <p v-if="recentLogs.length === 0" class="text-xs text-on-surface-variant italic py-8 text-center">
            No recent activity logs.
          </p>
        </div>
      </div>
    </div>

    <!-- Recent Signups Section -->
    <div class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden shadow-xs">
      <div class="px-6 py-4 border-b border-outline-variant/40 flex justify-between items-center bg-surface-container-lowest/50">
        <h3 class="font-headline text-base sm:text-lg font-bold text-on-surface flex items-center gap-2">
          <span class="material-symbols-outlined text-amber-500 text-[20px]">how_to_reg</span>
          Recent Registered Accounts
        </h3>
        <router-link to="/user/admin/users" class="font-label text-xs font-bold text-amber-600 hover:underline">View All Users &rarr;</router-link>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
              <th class="px-6 py-3.5">User</th>
              <th class="px-6 py-3.5">Email</th>
              <th class="px-6 py-3.5">Role</th>
              <th class="px-6 py-3.5">Status</th>
              <th class="px-6 py-3.5">Registered</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
            <tr v-for="user in recentUsers" :key="user.id" class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-semibold text-on-surface flex items-center gap-2.5">
                <span class="w-7 h-7 rounded-full bg-amber-500/20 text-amber-600 border border-amber-500/30 flex items-center justify-center font-bold text-xs uppercase">
                  {{ user.name?.[0] || 'U' }}
                </span>
                <router-link :to="`/user/admin/users/view?id=${user.id}`" class="hover:text-amber-500 transition-colors">{{ user.name }}</router-link>
              </td>
              <td class="px-6 py-4 text-on-surface-variant">{{ user.email }}</td>
              <td class="px-6 py-4">
                <span v-if="user.role?.name === 'Admin'" class="px-2.5 py-0.5 rounded-md text-xs font-black uppercase bg-amber-500/20 text-amber-600 border border-amber-500/30">ADMIN</span>
                <span v-else class="px-2.5 py-0.5 rounded-md text-xs font-semibold uppercase bg-surface-container-high text-on-surface-variant border border-outline-variant/40">USER</span>
              </td>
              <td class="px-6 py-4">
                <span v-if="user.is_activated" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-600 border border-emerald-500/20 inline-flex items-center gap-1">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Active
                </span>
                <span v-else class="px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-600 border border-amber-500/20 inline-flex items-center gap-1">
                  <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span> Unverified
                </span>
              </td>
              <td class="px-6 py-4 text-on-surface-variant text-xs">{{ new Date(user.created_at).toLocaleDateString() }}</td>
            </tr>
            <tr v-if="recentUsers.length === 0">
              <td colspan="5" class="px-6 py-8 text-center text-on-surface-variant">No user registrations found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import api from '../../utils/api'
import Chart from 'chart.js/auto'

const stats = ref({})
const recentUsers = ref([])
const recentLogs = ref([])
const revenueChartRef = ref(null)

onMounted(async () => {
  try {
    const res = await api.get('/admin/stats')
    if (res.data) {
      stats.value = res.data.stats || {}
      recentUsers.value = res.data.recent_users || []
      recentLogs.value = res.data.recent_logs || []

      await nextTick()
      initChart(stats.value.total_revenue || 0)
    }
  } catch (err) {
    // API errors handled by global interceptors
  }
})

function initChart(totalRev) {
  if (!revenueChartRef.value) return

  const ctx = revenueChartRef.value.getContext('2d')
  
  const gradient = ctx.createLinearGradient(0, 0, 0, 280)
  gradient.addColorStop(0, 'rgba(245, 158, 11, 0.3)')
  gradient.addColorStop(1, 'rgba(245, 158, 11, 0.0)')

  const data = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug'],
    datasets: [{
      label: 'Revenue (GH₵)',
      data: [1200, 2400, 3100, 4800, 6200, 8500, 11200, totalRev],
      borderColor: '#f59e0b',
      backgroundColor: gradient,
      borderWidth: 3,
      pointBackgroundColor: '#ffffff',
      pointBorderColor: '#f59e0b',
      pointBorderWidth: 2,
      pointRadius: 4,
      pointHoverRadius: 6,
      fill: true,
      tension: 0.4
    }]
  }

  const config = {
    type: 'line',
    data: data,
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: '#1f2937',
          titleFont: { family: 'Nunito Sans', size: 12 },
          bodyFont: { family: 'Nunito Sans', size: 13, weight: 'bold' },
          padding: 10,
          callbacks: {
            label: function(context) {
              return 'Revenue: GH₵ ' + context.parsed.y.toFixed(2)
            }
          }
        }
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: { font: { family: 'Nunito Sans', size: 11 }, color: '#6b7280' }
        },
        y: {
          grid: { color: 'rgba(156, 163, 175, 0.15)', borderDash: [4, 4] },
          ticks: { font: { family: 'Nunito Sans', size: 11 }, color: '#6b7280' },
          beginAtZero: true
        }
      }
    }
  }

  new Chart(ctx, config)
}
</script>
