<template>
  <div class="min-h-screen bg-background text-on-surface flex flex-col selection:bg-primary-container/40 selection:text-primary">
    <!-- Navbar -->
    <AppNavbar />

    <main class="flex-1">
      <!-- Hero Section -->
      <section class="relative pt-8 sm:pt-12 md:pt-16 pb-12 sm:pb-16 lg:pb-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto overflow-hidden">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 sm:gap-12 items-center relative z-10">
          <!-- Hero Text -->
          <div class="flex flex-col gap-5 sm:gap-6">
            <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-bold w-fit shadow-xs">
              <span class="material-symbols-outlined text-[16px]">auto_awesome</span>
              <span>Next-Gen Invoicing Platform</span>
            </div>

            <h1 class="font-headline text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-extrabold text-on-surface leading-[1.2] sm:leading-tight tracking-tight">
              Precise Billing for <span class="brand-title text-primary block sm:inline">Modern Professionals.</span>
            </h1>

            <p class="font-body text-sm sm:text-base md:text-lg text-on-surface-variant max-w-lg leading-relaxed">
              Frictionless invoice creation builder. Create, verify, track, and dispatch beautiful PDF invoices to clients in seconds. No monthly commitments — pay only for what you send.
            </p>

            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 sm:gap-4 pt-1">
              <router-link
                to="/user/register"
                class="w-full sm:w-auto px-7 py-3.5 rounded-xl font-label text-sm sm:text-base font-semibold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-all shadow-md hover:shadow-lg flex items-center justify-center gap-2.5 active:scale-[0.98]"
              >
                <span>Start Invoicing Free</span>
                <span class="material-symbols-outlined text-[18px]">arrow_forward</span>
              </router-link>

              <div class="flex items-center justify-center sm:justify-start gap-1.5 font-body text-xs font-medium text-on-surface-variant">
                <span class="material-symbols-outlined text-primary text-[18px]">check_circle</span>
                <span>No credit card required</span>
              </div>
            </div>
          </div>

          <!-- Hero Abstract Invoice Mockup Card -->
          <div class="relative w-full rounded-2xl glass-card p-4 sm:p-6 border border-outline-variant/60 shadow-xl sm:shadow-2xl flex flex-col justify-between overflow-hidden">
            <!-- Mockup Top Bar -->
            <div class="flex items-center justify-between pb-3 sm:pb-4 border-b border-outline-variant/40 gap-2">
              <div class="flex items-center gap-1.5 sm:gap-2 min-w-0">
                <span class="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-error/70 shrink-0"></span>
                <span class="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-tertiary-fixed-dim shrink-0"></span>
                <span class="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-primary/70 shrink-0"></span>
                <span class="font-mono-data text-[11px] sm:text-xs text-on-surface-variant ml-1 sm:ml-2 font-semibold truncate">INV-2026-0042.pdf</span>
              </div>
              <span class="px-2.5 py-0.5 sm:px-3 sm:py-1 rounded-full bg-primary-container/60 text-on-primary-container text-[10px] sm:text-xs font-bold flex items-center gap-1 sm:gap-1.5 shrink-0">
                <span class="w-1.5 h-1.5 sm:w-2 sm:h-2 rounded-full bg-primary animate-pulse"></span>
                <span>Ready to Send</span>
              </span>
            </div>

            <!-- Mockup Body -->
            <div class="py-4 space-y-3.5">
              <div class="flex justify-between items-start gap-2">
                <div>
                  <h3 class="font-headline text-base sm:text-lg font-bold text-on-surface leading-snug">Acme Studio Design</h3>
                  <p class="font-body text-[11px] sm:text-xs text-on-surface-variant">Client: TechCorp Global Ltd.</p>
                </div>
                <div class="text-right shrink-0">
                  <span class="font-headline text-xl sm:text-2xl font-bold text-primary block">$3,450.00</span>
                  <p class="font-body text-[11px] sm:text-xs text-on-surface-variant">Due in 14 days</p>
                </div>
              </div>

              <!-- Line Items Mockup -->
              <div class="bg-surface-container-low/90 rounded-xl p-3 sm:p-3.5 space-y-2 border border-outline-variant/30 text-xs font-body">
                <div class="flex justify-between items-center font-semibold text-on-surface gap-2">
                  <span class="truncate">Brand Strategy & Design System</span>
                  <span class="shrink-0 font-medium">$2,500.00</span>
                </div>
                <div class="flex justify-between items-center text-on-surface-variant gap-2">
                  <span class="truncate">Interactive Web Component Suite</span>
                  <span class="shrink-0 font-medium">$950.00</span>
                </div>
              </div>
            </div>

            <!-- Mockup Bottom Bar -->
            <div class="pt-3 border-t border-outline-variant/40 flex flex-wrap sm:flex-nowrap items-center justify-between gap-3">
              <span class="font-body text-xs text-on-surface-variant flex items-center gap-1.5">
                <span class="material-symbols-outlined text-[16px] text-primary">verified</span>
                <span>Verified Public Link Included</span>
              </span>
              <router-link
                to="/user/register"
                class="text-xs font-semibold px-3.5 py-2 bg-primary text-on-primary rounded-lg flex items-center gap-1.5 hover:bg-on-primary-fixed-variant transition-colors ml-auto sm:ml-0 shrink-0"
              >
                <span>Send Invoice</span>
                <span class="material-symbols-outlined text-[14px]">send</span>
              </router-link>
            </div>
          </div>
        </div>
      </section>

      <!-- Trust Metrics Bar -->
      <section class="border-y border-outline-variant/30 bg-surface-container-lowest/60 py-6 sm:py-8 px-4 sm:px-6">
        <div class="max-w-7xl mx-auto grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6 text-center">
          <div class="p-2 sm:p-3">
            <div class="font-headline text-2xl sm:text-3xl font-extrabold text-primary">10k+</div>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant mt-1">Invoices Dispatched</p>
          </div>
          <div class="p-2 sm:p-3">
            <div class="font-headline text-2xl sm:text-3xl font-extrabold text-primary">0</div>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant mt-1">Monthly Subscriptions</p>
          </div>
          <div class="p-2 sm:p-3">
            <div class="font-headline text-2xl sm:text-3xl font-extrabold text-primary">&lt; 30s</div>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant mt-1">Quick Generation</p>
          </div>
          <div class="p-2 sm:p-3">
            <div class="font-headline text-2xl sm:text-3xl font-extrabold text-primary">100%</div>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant mt-1">Pay-Per-Use Model</p>
          </div>
        </div>
      </section>

      <!-- How It Works Section -->
      <section class="py-14 sm:py-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto space-y-10 sm:space-y-12" id="how-it-works">
        <div class="text-center max-w-2xl mx-auto space-y-3">
          <h2 class="font-headline text-2xl sm:text-3xl md:text-4xl font-bold text-on-surface">Invoice in 3 Simple Steps</h2>
          <p class="font-body text-sm sm:text-base text-on-surface-variant">We've stripped away the complexity. Focus on your work, and let us handle the billing pipeline.</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8 relative">
          <!-- Connecting Line (Desktop) -->
          <div class="hidden md:block absolute top-1/2 left-[10%] right-[10%] h-0.5 bg-gradient-to-r from-primary-container/40 via-primary/30 to-primary-container/40 -translate-y-1/2 z-0"></div>
          
          <!-- Step 1 -->
          <div class="relative z-10 glass-card rounded-2xl p-6 sm:p-8 flex flex-col items-center text-center space-y-3 sm:space-y-4 hover:shadow-lg transition-all border border-outline-variant/60 bg-surface-container-lowest">
            <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-full bg-primary text-on-primary flex items-center justify-center font-headline text-xl sm:text-2xl font-bold shadow-md ring-4 ring-primary-container/40">1</div>
            <h3 class="font-headline text-lg sm:text-xl font-bold text-on-surface">Draft & Build</h3>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">Use our clean, live-preview editor to add line items, taxes, and notes instantly.</p>
          </div>
          <!-- Step 2 -->
          <div class="relative z-10 glass-card rounded-2xl p-6 sm:p-8 flex flex-col items-center text-center space-y-3 sm:space-y-4 hover:shadow-lg transition-all border border-outline-variant/60 bg-surface-container-lowest">
            <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-full bg-primary text-on-primary flex items-center justify-center font-headline text-xl sm:text-2xl font-bold shadow-md ring-4 ring-primary-container/40">2</div>
            <h3 class="font-headline text-lg sm:text-xl font-bold text-on-surface">Generate & Verify</h3>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">Create a pixel-perfect PDF and a secure public verification link with one click.</p>
          </div>
          <!-- Step 3 -->
          <div class="relative z-10 glass-card rounded-2xl p-6 sm:p-8 flex flex-col items-center text-center space-y-3 sm:space-y-4 hover:shadow-lg transition-all border border-outline-variant/60 bg-surface-container-lowest">
            <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-full bg-primary text-on-primary flex items-center justify-center font-headline text-xl sm:text-2xl font-bold shadow-md ring-4 ring-primary-container/40">3</div>
            <h3 class="font-headline text-lg sm:text-xl font-bold text-on-surface">Send & Get Paid</h3>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">Dispatch the invoice directly to your client's email inbox and get paid faster.</p>
          </div>
        </div>
      </section>

      <!-- Features Bento Grid -->
      <section class="py-14 sm:py-20 px-4 sm:px-6 lg:px-8 bg-surface-container-lowest/70 border-y border-outline-variant/40" id="features">
        <div class="max-w-7xl mx-auto space-y-10 sm:space-y-12">
          <div class="text-center max-w-2xl mx-auto space-y-3">
            <h2 class="font-headline text-2xl sm:text-3xl md:text-4xl font-bold text-on-surface">Everything you need, nothing you don't.</h2>
            <p class="font-body text-sm sm:text-base text-on-surface-variant">A focused toolset designed specifically for freelancers, consultants, and teams who value speed and clarity.</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8">
            <!-- Feature 1 -->
            <div class="col-span-1 md:col-span-2 glass-card rounded-2xl p-6 sm:p-8 flex flex-col justify-between hover:shadow-xl transition-all">
              <div>
                <div class="w-11 h-11 sm:w-12 sm:h-12 rounded-xl bg-primary/10 text-primary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[26px]">edit_document</span>
                </div>
                <h3 class="font-headline text-xl sm:text-2xl font-bold text-on-surface mb-2">Live Invoice Builder</h3>
                <p class="font-body text-sm sm:text-base text-on-surface-variant leading-relaxed">Create invoices instantly. What you type is precisely what your client receives in their inbox and PDF export.</p>
              </div>
            </div>

            <!-- Feature 2 -->
            <div class="col-span-1 glass-card rounded-2xl p-6 sm:p-8 flex flex-col justify-between hover:shadow-xl transition-all">
              <div>
                <div class="w-11 h-11 sm:w-12 sm:h-12 rounded-xl bg-secondary/10 text-secondary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[26px]">picture_as_pdf</span>
                </div>
                <h3 class="font-headline text-lg sm:text-xl font-bold text-on-surface mb-2">Pixel-Perfect PDF</h3>
                <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">Generate standardized, high-density PDFs instantly for print, archive, or email attachment.</p>
              </div>
            </div>

            <!-- Feature 3 -->
            <div class="col-span-1 glass-card rounded-2xl p-6 sm:p-8 flex flex-col justify-between hover:shadow-xl transition-all">
              <div>
                <div class="w-11 h-11 sm:w-12 sm:h-12 rounded-xl bg-tertiary/10 text-tertiary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[26px]">mark_email_read</span>
                </div>
                <h3 class="font-headline text-lg sm:text-xl font-bold text-on-surface mb-2">Secure Link Verification</h3>
                <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">Automated account verification links eliminate spoofing and ensure clean database record hygiene.</p>
              </div>
            </div>

            <!-- Feature 4 -->
            <div class="col-span-1 md:col-span-2 glass-card rounded-2xl p-6 sm:p-8 flex flex-col justify-between hover:shadow-xl transition-all bg-gradient-to-br from-primary-container/20 to-transparent border-primary/20">
              <div>
                <div class="w-11 h-11 sm:w-12 sm:h-12 rounded-xl bg-primary/10 text-primary flex items-center justify-center mb-4">
                  <span class="material-symbols-outlined text-[26px]">token</span>
                </div>
                <h3 class="font-headline text-xl sm:text-2xl font-bold text-on-surface mb-2">Pay-Per-Use Credit Pricing</h3>
                <p class="font-body text-sm sm:text-base text-on-surface-variant leading-relaxed">No recurring subscriptions draining your profits during quiet months. Purchase credits when needed and use them whenever you choose.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Pricing Section -->
      <section class="py-14 sm:py-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto space-y-10 sm:space-y-12" id="pricing">
        <div class="text-center max-w-lg mx-auto space-y-3">
          <h2 class="font-headline text-2xl sm:text-3xl md:text-4xl font-bold text-on-surface">Simple, Transparent Pricing</h2>
          <p class="font-body text-sm sm:text-base text-on-surface-variant">No subscriptions. No hidden fees. Pay only for the invoices you dispatch.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8 max-w-5xl mx-auto">
          <div
            v-for="pkg in packages"
            :key="pkg.id"
            :class="[
              'glass-card rounded-2xl p-6 sm:p-8 flex flex-col justify-between border relative transition-all',
              pkg.badge_tag ? 'border-2 border-primary shadow-xl md:scale-[1.02]' : 'border-outline-variant/60'
            ]"
          >
            <span
              v-if="pkg.badge_tag"
              class="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary text-on-primary text-[10px] font-bold px-3 py-0.5 rounded-full uppercase tracking-widest shadow-sm whitespace-nowrap"
            >
              {{ pkg.badge_tag }}
            </span>

            <div>
              <h3 class="font-headline text-lg sm:text-xl font-bold text-on-surface mb-1">{{ pkg.name }}</h3>
              <div class="mb-5 mt-3">
                <span class="font-headline text-3xl sm:text-4xl font-extrabold text-on-surface">
                  GH₵{{ (pkg.price / 100).toFixed(0) }}
                </span>
                <span class="font-body text-xs text-on-surface-variant"> / {{ pkg.credits_granted }} Credits</span>
              </div>

              <ul class="space-y-2.5 font-body text-xs sm:text-sm text-on-surface-variant mb-6">
                <li v-for="(feat, fIdx) in pkg.features" :key="fIdx" class="flex items-center gap-2">
                  <span class="material-symbols-outlined text-primary text-[18px] shrink-0">check</span>
                  <span>{{ feat }}</span>
                </li>
              </ul>
            </div>

            <router-link
              to="/user/register"
              :class="[
                'w-full py-3 sm:py-3.5 text-center font-label text-sm font-semibold rounded-xl transition-all shadow-md active:scale-[0.98]',
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

      <!-- Testimonials Section -->
      <section class="py-14 sm:py-20 px-4 sm:px-6 lg:px-8 bg-surface-container-lowest/70 border-y border-outline-variant/40" id="testimonials">
        <div class="max-w-7xl mx-auto space-y-10 sm:space-y-12">
          <div class="text-center max-w-2xl mx-auto space-y-3">
            <h2 class="font-headline text-2xl sm:text-3xl md:text-4xl font-bold text-on-surface">Trusted by Professionals</h2>
            <p class="font-body text-sm sm:text-base text-on-surface-variant">See what independent creators and businesses are saying about Teks-Invoice.</p>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-5 sm:gap-6">
            <div class="glass-card rounded-2xl p-5 sm:p-6 border border-outline-variant/60 flex flex-col justify-between space-y-4">
              <div>
                <div class="flex items-center gap-0.5 text-amber-500 mb-2">
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                </div>
                <p class="font-body text-xs sm:text-sm text-on-surface-variant italic leading-relaxed">"This platform changed how I bill. It's incredibly fast, the PDFs look highly professional, and the pay-as-you-go credit system is perfect for my freelance workload."</p>
              </div>
              <div class="flex items-center gap-3 pt-2">
                <div class="w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-secondary-container text-on-secondary-container flex items-center justify-center font-bold text-xs sm:text-sm shrink-0">SD</div>
                <div>
                  <h4 class="font-headline text-xs sm:text-sm font-bold text-on-surface">Sarah D.</h4>
                  <p class="font-body text-[11px] sm:text-xs text-on-surface-variant">Freelance Designer</p>
                </div>
              </div>
            </div>

            <div class="glass-card rounded-2xl p-5 sm:p-6 border border-outline-variant/60 flex flex-col justify-between space-y-4">
              <div>
                <div class="flex items-center gap-0.5 text-amber-500 mb-2">
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                </div>
                <p class="font-body text-xs sm:text-sm text-on-surface-variant italic leading-relaxed">"The secure public link feature is brilliant. My clients love being able to just click a link and see their invoice without logging into anything."</p>
              </div>
              <div class="flex items-center gap-3 pt-2">
                <div class="w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-tertiary-container text-on-tertiary-container flex items-center justify-center font-bold text-xs sm:text-sm shrink-0">MR</div>
                <div>
                  <h4 class="font-headline text-xs sm:text-sm font-bold text-on-surface">Mark R.</h4>
                  <p class="font-body text-[11px] sm:text-xs text-on-surface-variant">Software Consultant</p>
                </div>
              </div>
            </div>

            <div class="glass-card rounded-2xl p-5 sm:p-6 border border-outline-variant/60 flex flex-col justify-between space-y-4">
              <div>
                <div class="flex items-center gap-0.5 text-amber-500 mb-2">
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                  <span class="material-symbols-outlined text-[18px]">star</span>
                </div>
                <p class="font-body text-xs sm:text-sm text-on-surface-variant italic leading-relaxed">"I was tired of paying $15/month for invoicing software I only used 3 times a month. Teks-Invoice's credit model is exactly what small agencies need."</p>
              </div>
              <div class="flex items-center gap-3 pt-2">
                <div class="w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-primary-container text-on-primary-container flex items-center justify-center font-bold text-xs sm:text-sm shrink-0">ET</div>
                <div>
                  <h4 class="font-headline text-xs sm:text-sm font-bold text-on-surface">Elena T.</h4>
                  <p class="font-body text-[11px] sm:text-xs text-on-surface-variant">Creative Agency Owner</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- FAQ Summary Section -->
      <section class="py-14 sm:py-20 px-4 sm:px-6 lg:px-8 max-w-3xl mx-auto space-y-8 sm:space-y-10" id="faq">
        <div class="text-center space-y-3">
          <h2 class="font-headline text-2xl sm:text-3xl font-bold text-on-surface">Have Questions?</h2>
          <p class="font-body text-sm sm:text-base text-on-surface-variant">Here are the most common questions we get asked.</p>
        </div>
        <div class="space-y-3 sm:space-y-4">
          <div class="glass-card rounded-2xl p-5 sm:p-6 border border-outline-variant/60">
            <h3 class="font-headline text-base sm:text-lg font-bold text-on-surface mb-1.5 sm:mb-2">Is there a monthly subscription fee?</h3>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">No. We operate strictly on a pay-as-you-go credit system. You purchase credits, and they never expire.</p>
          </div>
          <div class="glass-card rounded-2xl p-5 sm:p-6 border border-outline-variant/60">
            <h3 class="font-headline text-base sm:text-lg font-bold text-on-surface mb-1.5 sm:mb-2">How many credits does an invoice cost?</h3>
            <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed">Generating and dispatching a finalized PDF invoice costs 1 credit. Drafts are free.</p>
          </div>
        </div>
        <div class="text-center pt-2 sm:pt-4">
          <router-link to="/faq" class="font-label text-sm font-semibold text-primary hover:underline inline-flex items-center gap-1">
            <span>Read all FAQs</span>
            <span>&rarr;</span>
          </router-link>
        </div>
      </section>

      <!-- Final CTA Section (Floating Card) -->
      <section class="px-4 sm:px-6 lg:px-8 py-10 sm:py-16 max-w-7xl mx-auto">
        <div class="relative rounded-3xl bg-primary text-on-primary text-center px-6 py-12 sm:py-16 sm:px-12 overflow-hidden shadow-2xl">
          <!-- Subtle decorative blur circles -->
          <div class="absolute -top-24 -right-24 w-72 h-72 rounded-full bg-white/10 blur-2xl pointer-events-none"></div>
          <div class="absolute -bottom-24 -left-24 w-72 h-72 rounded-full bg-black/10 blur-2xl pointer-events-none"></div>
          
          <div class="max-w-2xl mx-auto space-y-6 sm:space-y-8 relative z-10">
            <h2 class="font-headline text-2xl sm:text-4xl md:text-5xl font-extrabold tracking-tight leading-tight">
              Ready to streamline your invoicing?
            </h2>
            <p class="font-body text-sm sm:text-base md:text-lg text-primary-container/90 max-w-xl mx-auto leading-relaxed">
              Join thousands of modern professionals who get paid faster and spend less time on billing.
            </p>
            <div class="pt-2">
              <router-link
                to="/user/register"
                class="w-full sm:w-auto inline-flex items-center justify-center gap-2.5 px-8 py-4 rounded-xl font-label text-sm sm:text-base font-bold bg-surface text-primary hover:bg-surface-bright active:scale-[0.98] transition-all shadow-xl"
              >
                <span>Start For Free Today</span>
                <span class="material-symbols-outlined text-[18px]">arrow_forward</span>
              </router-link>
            </div>
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
