<template>
  <div class="min-h-screen bg-background text-on-surface flex flex-col selection:bg-primary-container/40 selection:text-primary">
    <!-- Navbar -->
    <AppNavbar />

    <main class="flex-1">
      <!-- Hero Section -->
      <section class="relative pt-16 pb-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto overflow-hidden">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center relative z-10">
          <!-- Hero Text -->
          <div class="flex flex-col gap-6">
            <div class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-primary-container/40 border border-primary/20 text-primary text-xs font-bold w-fit">
              <span class="material-symbols-outlined text-[18px]">auto_awesome</span>
              <span>Next-Gen Invoicing Platform</span>
            </div>

            <h1 class="font-headline text-4xl sm:text-5xl lg:text-6xl font-extrabold text-on-surface leading-tight tracking-tight">
              Precise Billing for <span class="brand-title text-primary">Modern Professionals.</span>
            </h1>

            <p class="font-body text-base sm:text-lg text-on-surface-variant max-w-lg leading-relaxed">
              Frictionless invoice creation builder. Create, verify, track, and dispatch beautiful PDF invoices to clients in seconds. No monthly commitments — pay only for what you send.
            </p>

            <div class="flex flex-wrap items-center gap-4 pt-2">
              <router-link
                to="/user/register"
                class="px-7 py-3.5 rounded-xl font-label text-base font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all shadow-lg hover:shadow-xl flex items-center gap-3"
              >
                <span>Start Invoicing Free</span>
                <span class="material-symbols-outlined text-[18px]">arrow_forward</span>
              </router-link>

              <span class="font-body text-xs font-medium text-on-surface-variant flex items-center gap-1.5">
                <span class="material-symbols-outlined text-primary text-[18px]">check_circle</span>
                <span>No credit card required</span>
              </span>
            </div>
          </div>

          <!-- Hero Abstract Invoice Mockup Card -->
          <div class="relative w-full aspect-[4/3] rounded-2xl overflow-hidden glass-card p-6 border border-outline-variant/60 shadow-2xl flex flex-col justify-between">
            <!-- Mockup Top Bar -->
            <div class="flex items-center justify-between pb-4 border-b border-outline-variant/40">
              <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded-full bg-error/70"></span>
                <span class="w-3 h-3 rounded-full bg-tertiary-fixed-dim"></span>
                <span class="w-3 h-3 rounded-full bg-primary/70"></span>
                <span class="font-mono-data text-xs text-on-surface-variant ml-2 font-semibold">INV-2026-0042.pdf</span>
              </div>
              <span class="px-3 py-1 rounded-full bg-primary-container/60 text-on-primary-container text-xs font-bold flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-primary animate-pulse"></span>
                <span>Ready to Send</span>
              </span>
            </div>

            <!-- Mockup Body -->
            <div class="py-4 space-y-4">
              <div class="flex justify-between items-start">
                <div>
                  <h3 class="font-headline text-lg font-bold text-on-surface">Acme Studio Design</h3>
                  <p class="font-body text-xs text-on-surface-variant">Client: TechCorp Global Ltd.</p>
                </div>
                <div class="text-right">
                  <span class="font-headline text-2xl font-bold text-primary">$3,450.00</span>
                  <p class="font-body text-xs text-on-surface-variant">Due in 14 days</p>
                </div>
              </div>

              <!-- Line Items Mockup -->
              <div class="bg-surface-container-low/80 rounded-xl p-3.5 space-y-2 border border-outline-variant/30 text-xs font-body">
                <div class="flex justify-between font-semibold text-on-surface">
                  <span>Brand Strategy & Design System</span>
                  <span>$2,500.00</span>
                </div>
                <div class="flex justify-between text-on-surface-variant">
                  <span>Interactive Web Component Suite</span>
                  <span>$950.00</span>
                </div>
              </div>
            </div>

            <!-- Mockup Bottom Bar -->
            <div class="pt-3 border-t border-outline-variant/40 flex items-center justify-between">
              <span class="font-body text-xs text-on-surface-variant flex items-center gap-1.5">
                <span class="material-symbols-outlined text-[16px] text-primary">verified</span>
                <span>Verified Public Link Included</span>
              </span>
              <router-link
                to="/user/register"
                class="text-xs font-semibold px-4 py-2 bg-primary text-on-primary rounded-lg flex items-center gap-1 hover:bg-on-primary-fixed-variant transition-colors"
              >
                <span>Send Invoice</span>
                <span class="material-symbols-outlined text-[14px]">send</span>
              </router-link>
            </div>
          </div>
        </div>
      </section>

      <!-- Features Bento Grid -->
      <section class="py-20 px-4 sm:px-6 lg:px-8 bg-surface-container-lowest/70 border-y border-outline-variant/40" id="features">
        <div class="max-w-7xl mx-auto space-y-12">
          <div class="text-center max-w-2xl mx-auto space-y-3">
            <h2 class="font-headline text-3xl sm:text-4xl font-bold text-on-surface">Everything you need, nothing you don't.</h2>
            <p class="font-body text-base text-on-surface-variant">A focused toolset designed specifically for freelancers, consultants, and teams who value speed and clarity.</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
            <!-- Feature 1 -->
            <div class="col-span-1 md:col-span-2 glass-card rounded-2xl p-8 flex flex-col justify-between hover:shadow-xl transition-all">
              <div>
                <div class="w-12 h-12 rounded-xl bg-primary/10 text-primary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[28px]">edit_document</span>
                </div>
                <h3 class="font-headline text-2xl font-bold text-on-surface mb-2">Live Invoice Builder</h3>
                <p class="font-body text-base text-on-surface-variant">Create invoices instantly. What you type is precisely what your client receives in their inbox and PDF export.</p>
              </div>
            </div>

            <!-- Feature 2 -->
            <div class="col-span-1 glass-card rounded-2xl p-8 flex flex-col justify-between hover:shadow-xl transition-all">
              <div>
                <div class="w-12 h-12 rounded-xl bg-secondary/10 text-secondary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[28px]">picture_as_pdf</span>
                </div>
                <h3 class="font-headline text-xl font-bold text-on-surface mb-2">Pixel-Perfect PDF</h3>
                <p class="font-body text-sm text-on-surface-variant">Generate standardized, high-density PDFs instantly for print, archive, or email attachment.</p>
              </div>
            </div>

            <!-- Feature 3 -->
            <div class="col-span-1 glass-card rounded-2xl p-8 flex flex-col justify-between hover:shadow-xl transition-all">
              <div>
                <div class="w-12 h-12 rounded-xl bg-tertiary/10 text-tertiary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[28px]">mark_email_read</span>
                </div>
                <h3 class="font-headline text-xl font-bold text-on-surface mb-2">Secure Link Verification</h3>
                <p class="font-body text-sm text-on-surface-variant">Automated account verification links eliminate spoofing and ensure clean database record hygiene.</p>
              </div>
            </div>

            <!-- Feature 4 -->
            <div class="col-span-1 md:col-span-2 glass-card rounded-2xl p-8 flex flex-col justify-between hover:shadow-xl transition-all bg-gradient-to-br from-primary-container/20 to-transparent border-primary/20">
              <div>
                <div class="w-12 h-12 rounded-xl bg-primary/10 text-primary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[28px]">token</span>
                </div>
                <h3 class="font-headline text-2xl font-bold text-on-surface mb-2">Pay-Per-Use Credit Pricing</h3>
                <p class="font-body text-base text-on-surface-variant">No recurring subscriptions draining your profits during quiet months. Purchase credits when needed and use them whenever you choose.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Pricing Section -->
      <section class="py-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto space-y-12" id="pricing">
        <div class="text-center max-w-lg mx-auto space-y-3">
          <h2 class="font-headline text-3xl sm:text-4xl font-bold text-on-surface">Simple, Transparent Pricing</h2>
          <p class="font-body text-base text-on-surface-variant">No subscriptions. No hidden fees. Pay only for the invoices you dispatch.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          <div
            v-for="pkg in packages"
            :key="pkg.id"
            :class="[
              'glass-card rounded-2xl p-8 flex flex-col justify-between border relative transition-all',
              pkg.badge_tag ? 'border-2 border-primary shadow-xl scale-[1.02]' : 'border-outline-variant/60'
            ]"
          >
            <span
              v-if="pkg.badge_tag"
              class="absolute -top-3.5 left-1/2 -translate-x-1/2 bg-primary text-on-primary text-[10px] font-bold px-3 py-0.5 rounded-full uppercase tracking-widest shadow-sm"
            >
              {{ pkg.badge_tag }}
            </span>

            <div>
              <h3 class="font-headline text-xl font-bold text-on-surface mb-1">{{ pkg.name }}</h3>
              <div class="mb-6 mt-4">
                <span class="font-headline text-4xl font-extrabold text-on-surface">
                  GH₵{{ (pkg.price / 100).toFixed(0) }}
                </span>
                <span class="font-body text-xs text-on-surface-variant"> / {{ pkg.credits_granted }} Credits</span>
              </div>

              <ul class="space-y-3 font-body text-sm text-on-surface-variant mb-6">
                <li v-for="(feat, fIdx) in pkg.features" :key="fIdx" class="flex items-center gap-2">
                  <span class="material-symbols-outlined text-primary text-[18px]">check</span>
                  <span>{{ feat }}</span>
                </li>
              </ul>
            </div>

            <router-link
              to="/user/register"
              :class="[
                'w-full py-3 text-center font-label text-sm font-semibold rounded-xl transition-all shadow-md',
                pkg.badge_tag
                  ? 'bg-primary text-on-primary hover:bg-on-primary-fixed-variant'
                  : 'border border-primary text-primary hover:bg-primary hover:text-on-primary'
              ]"
            >
              Get Started with {{ pkg.name }}
            </router-link>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppNavbar from '../../components/common/AppNavbar.vue'
import AppFooter from '../../components/common/AppFooter.vue'
import api from '../../utils/api'

const packages = ref([
  {
    id: 1,
    name: 'Starter Pack',
    price: 1500,
    credits_granted: 15,
    features: ['15 Invoice PDF Downloads', 'Client Management', 'Finance Ledger Sync'],
  },
  {
    id: 2,
    name: 'Professional Pack',
    price: 3500,
    credits_granted: 50,
    badge_tag: 'Most Popular',
    features: ['50 Invoice PDF Downloads', 'Client Management', 'Finance Ledger Sync', 'Custom Business Branding'],
  },
  {
    id: 3,
    name: 'Agency Bundle',
    price: 9500,
    credits_granted: 150,
    features: ['150 Invoice PDF Downloads', 'Unlimited Clients', 'Full Financial Stats & CSV Export', 'Priority Support'],
  },
])

onMounted(async () => {
  try {
    const res = await api.get('/credits/packages')
    if (res.data?.packages && res.data.packages.length > 0) {
      packages.value = res.data.packages.map((pkg) => ({
        ...pkg,
        features: pkg.description ? pkg.description.split('\n').filter(Boolean) : ['Invoice PDF Generation', 'Client Management'],
      }))
    }
  } catch (err) {
    // Fallback to static package defaults
  }
})
</script>
