import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Layouts
import AuthLayout from '../layouts/AuthLayout.vue'
import AppLayout from '../layouts/AppLayout.vue'
import AdminLayout from '../layouts/AdminLayout.vue'

// Pages
import LandingPage from '../pages/public/LandingPage.vue'
import LoginPage from '../pages/auth/LoginPage.vue'
import RegisterPage from '../pages/auth/RegisterPage.vue'
import ForgotPasswordPage from '../pages/auth/ForgotPasswordPage.vue'
import ResetPasswordPage from '../pages/auth/ResetPasswordPage.vue'
import ResendVerificationPage from '../pages/auth/ResendVerificationPage.vue'
import VerifyPage from '../pages/auth/VerifyPage.vue'

import DashboardPage from '../pages/app/DashboardPage.vue'
import ProfileSetupPage from '../pages/app/ProfileSetupPage.vue'
import InvoicesPage from '../pages/app/InvoicesPage.vue'
import InvoiceCreatePage from '../pages/app/InvoiceCreatePage.vue'
import InvoiceEditPage from '../pages/app/InvoiceEditPage.vue'
import InvoiceViewPage from '../pages/app/InvoiceViewPage.vue'
import ReceiptViewPage from '../pages/app/ReceiptViewPage.vue'
import ClientsPage from '../pages/app/ClientsPage.vue'
import FinancePage from '../pages/app/FinancePage.vue'
import CreditsPage from '../pages/app/CreditsPage.vue'
import AdminDashboardPage from '../pages/admin/AdminDashboardPage.vue'
import AdminUsersPage from '../pages/admin/AdminUsersPage.vue'
import AdminUserViewPage from '../pages/admin/AdminUserViewPage.vue'
import AdminPackagesPage from '../pages/admin/AdminPackagesPage.vue'
import AdminInvoicesPage from '../pages/admin/AdminInvoicesPage.vue'
import AdminPaymentsPage from '../pages/admin/AdminPaymentsPage.vue'
import AdminCreditsPage from '../pages/admin/AdminCreditsPage.vue'
import AdminAuditLogsPage from '../pages/admin/AdminAuditLogsPage.vue'
import AdminWebhooksPage from '../pages/admin/AdminWebhooksPage.vue'
import AdminSettingsPage from '../pages/admin/AdminSettingsPage.vue'
import AdminLoginPage from '../pages/admin/AdminLoginPage.vue'

import TermsPage from '../pages/public/TermsPage.vue'
import PrivacyPage from '../pages/public/PrivacyPage.vue'
import RefundPolicyPage from '../pages/public/RefundPolicyPage.vue'
import SecurityPage from '../pages/public/SecurityPage.vue'
import ContactPage from '../pages/public/ContactPage.vue'
import FaqPage from '../pages/public/FaqPage.vue'
import PublicInvoicePage from '../pages/public/PublicInvoicePage.vue'
import NotFoundPage from '../pages/public/NotFoundPage.vue'
import ServerErrorPage from '../pages/public/ServerErrorPage.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'landing',
    component: LandingPage,
  },
  {
    path: '/terms',
    name: 'terms',
    component: TermsPage,
  },
  {
    path: '/privacy',
    name: 'privacy',
    component: PrivacyPage,
  },
  {
    path: '/refund-policy',
    name: 'refund-policy',
    component: RefundPolicyPage,
  },
  {
    path: '/refund',
    name: 'refund',
    component: RefundPolicyPage,
  },
  {
    path: '/security',
    name: 'security',
    component: SecurityPage,
  },
  {
    path: '/contact',
    name: 'contact',
    component: ContactPage,
  },
  {
    path: '/faq',
    name: 'faq',
    component: FaqPage,
  },
  {
    path: '/500',
    name: 'server-error',
    component: ServerErrorPage,
  },
  {
    path: '/invoice/public/:token?',
    name: 'public-invoice',
    component: PublicInvoicePage,
  },

  // Auth Group
  {
    path: '/user',
    component: AuthLayout,
    children: [
      { path: 'login', name: 'login', component: LoginPage, meta: { guestOnly: true } },
      { path: '/admin/login', name: 'admin-login', component: AdminLoginPage, meta: { guestOnly: true } },
      { path: 'register', name: 'register', component: RegisterPage, meta: { guestOnly: true } },
      { path: 'forgot-password', name: 'forgot-password', component: ForgotPasswordPage, meta: { guestOnly: true } },
      { path: 'reset-password', name: 'reset-password', component: ResetPasswordPage, meta: { guestOnly: true } },
      { path: 'resend-verification', name: 'resend-verification', component: ResendVerificationPage, meta: { guestOnly: true } },
      { path: 'verify', name: 'verify', component: VerifyPage },
    ],
  },

  // App Group (Protected)
  {
    path: '/user',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      { path: 'dashboard', name: 'dashboard', component: DashboardPage },
      { path: 'profile/setup', name: 'profile-setup', component: ProfileSetupPage },
      { path: 'invoices', name: 'invoices', component: InvoicesPage },
      { path: 'invoices/new', name: 'invoice-create', component: InvoiceCreatePage, meta: { requiresProfile: true } },
      { path: 'invoices/edit', name: 'invoice-edit', component: InvoiceEditPage, meta: { requiresProfile: true } },
      { path: 'invoices/view', name: 'invoice-view', component: InvoiceViewPage },
      { path: 'invoices/receipt/view', name: 'receipt-view', component: ReceiptViewPage },
      { path: 'clients', name: 'clients', component: ClientsPage, meta: { requiresProfile: true } },
      { path: 'finance', name: 'finance', component: FinancePage, meta: { requiresProfile: true } },
      { path: 'credits/history', name: 'credits-history', component: CreditsPage },
    ],
  },

  // Admin Group (Protected)
  {
    path: '/user/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      { path: 'dashboard', name: 'admin-dashboard', component: AdminDashboardPage },
      { path: 'users', name: 'admin-users', component: AdminUsersPage },
      { path: 'users/:id', name: 'admin-user-view', component: AdminUserViewPage },
      { path: 'packages', name: 'admin-packages', component: AdminPackagesPage },
      { path: 'invoices', name: 'admin-invoices', component: AdminInvoicesPage },
      { path: 'payments', name: 'admin-payments', component: AdminPaymentsPage },
      { path: 'credits', name: 'admin-credits', component: AdminCreditsPage },
      { path: 'audit-logs', name: 'admin-audit-logs', component: AdminAuditLogsPage },
      { path: 'webhooks', name: 'admin-webhooks', component: AdminWebhooksPage },
      { path: 'settings', name: 'admin-settings', component: AdminSettingsPage },
    ],
  },

  // Catch-all 404
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFoundPage,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  if (authStore.accessToken && !authStore.user) {
    await authStore.fetchCurrentUser()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return next({ name: 'login', query: { redirect: to.fullPath } })
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    if (authStore.user?.role?.name?.toLowerCase() === 'admin') {
      return next({ name: 'admin-dashboard' })
    }
    return next({ name: 'dashboard' })
  }

  if (to.meta.requiresProfile && authStore.isAuthenticated && !authStore.isProfileComplete) {
    return next({ name: 'profile-setup' })
  }

  if (to.meta.requiresAdmin && (!authStore.isAuthenticated || authStore.user?.role?.name?.toLowerCase() !== 'admin')) {
    return next({ name: 'dashboard' })
  }

  next()
})

export default router
