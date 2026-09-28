<template>
  <footer class="no-print bg-surface-container-lowest border-t border-outline-variant/40 pt-12 sm:pt-16 pb-8 relative z-10 mt-auto">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 grid grid-cols-2 md:grid-cols-4 gap-8 lg:gap-10 mb-10 sm:mb-12">
      
      <!-- Column 1: Brand -->
      <div class="col-span-2 md:col-span-1 space-y-3.5">
        <div class="flex items-center gap-2.5">
          <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-9 h-9 sm:w-10 sm:h-10 rounded-md object-contain shadow-xs">
          <span class="font-headline text-lg font-bold text-on-surface">Teks-Invoice</span>
        </div>
        <p class="font-body text-xs sm:text-sm text-on-surface-variant leading-relaxed max-w-xs">
          Precise billing for modern professionals. Create, verify, and dispatch beautiful PDF invoices in seconds.
        </p>
        <div class="flex items-center gap-2.5 pt-1">
          <!-- X / Twitter -->
          <a href="#" class="w-8 h-8 rounded-full bg-surface-container-low flex items-center justify-center text-on-surface hover:text-primary hover:bg-primary-container/40 transition-colors" aria-label="X (Twitter)">
            <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>
          </a>
          <!-- LinkedIn -->
          <a href="#" class="w-8 h-8 rounded-full bg-surface-container-low flex items-center justify-center text-on-surface hover:text-primary hover:bg-primary-container/40 transition-colors" aria-label="LinkedIn">
            <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
          </a>
          <!-- GitHub -->
          <a href="#" class="w-8 h-8 rounded-full bg-surface-container-low flex items-center justify-center text-on-surface hover:text-primary hover:bg-primary-container/40 transition-colors" aria-label="GitHub">
            <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24"><path fill-rule="evenodd" clip-rule="evenodd" d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.161 22 16.416 22 12c0-5.523-4.477-10-10-10z"/></svg>
          </a>
        </div>
      </div>

      <!-- Column 2: Product -->
      <div class="col-span-1">
        <h4 class="font-headline font-bold text-sm sm:text-base text-on-surface mb-3">Product</h4>
        <ul class="space-y-2 font-body text-xs sm:text-sm text-on-surface-variant">
          <li><router-link to="/#features" class="hover:text-primary transition-colors">Features</router-link></li>
          <li><router-link to="/#pricing" class="hover:text-primary transition-colors">Pricing</router-link></li>
          <li><router-link to="/user/register" class="hover:text-primary transition-colors">Sign Up Free</router-link></li>
          <li><router-link to="/user/login" class="hover:text-primary transition-colors">Sign In</router-link></li>
        </ul>
      </div>

      <!-- Column 3: Resources -->
      <div class="col-span-1">
        <h4 class="font-headline font-bold text-sm sm:text-base text-on-surface mb-3">Resources</h4>
        <ul class="space-y-2 font-body text-xs sm:text-sm text-on-surface-variant">
          <li><router-link to="/faq" class="hover:text-primary transition-colors">Help & FAQ</router-link></li>
          <li><router-link to="/contact" class="hover:text-primary transition-colors">Contact Support</router-link></li>
          <li><a href="#" class="hover:text-primary transition-colors">Templates</a></li>
          <li><a href="#" class="hover:text-primary transition-colors">System Status</a></li>
        </ul>
      </div>

      <!-- Column 4: Legal -->
      <div class="col-span-2 md:col-span-1">
        <h4 class="font-headline font-bold text-sm sm:text-base text-on-surface mb-3">Legal</h4>
        <ul class="space-y-2 font-body text-xs sm:text-sm text-on-surface-variant flex flex-wrap md:flex-col gap-x-4 gap-y-2">
          <li><router-link to="/terms" class="hover:text-primary transition-colors">Terms of Service</router-link></li>
          <li><router-link to="/privacy" class="hover:text-primary transition-colors">Privacy Policy</router-link></li>
          <li><router-link to="/refund-policy" class="hover:text-primary transition-colors">Refund Policy</router-link></li>
          <li><router-link to="/security" class="hover:text-primary transition-colors">Security</router-link></li>
        </ul>
      </div>
      
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 border-t border-outline-variant/40 pt-6 sm:pt-8 flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
      <span class="font-body text-xs text-on-surface-variant">
        &copy; {{ new Date().getFullYear() }} Teks-Invoice Inc. All rights reserved.
      </span>
      <div class="flex items-center gap-2 text-xs font-body text-on-surface-variant">
        <span class="material-symbols-outlined text-[16px]">language</span>
        <span>English (US)</span>
      </div>
    </div>
  </footer>
</template>

<script setup>
</script>
