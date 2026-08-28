<template>
  <footer class="no-print bg-surface-container-lowest/90 backdrop-blur-md w-full py-6 px-6 border-t border-outline-variant/40 relative z-10 mt-auto">
      <div class="flex flex-col md:flex-row justify-between items-center max-w-7xl mx-auto gap-4">
          <div class="flex items-center gap-2.5">
              <img src="/src/assets/brand_logo.png" alt="Teks-Invoice Logo" class="w-10 h-10 rounded-md object-contain shadow-sm">
              <span class="font-headline text-xs sm:text-sm font-bold text-on-surface">&copy; {{ new Date().getFullYear() }} Teks-Invoice Inc. All rights reserved.</span>
          </div>
          <div class="flex flex-wrap justify-center gap-x-6 gap-y-2 font-body text-xs text-on-surface-variant">
              <router-link class="hover:text-primary transition-colors" to="/terms">Terms of Service</router-link>
              <router-link class="hover:text-primary transition-colors" to="/privacy">Privacy Policy</router-link>
              <router-link class="hover:text-primary transition-colors" to="/refund-policy">Refund Policy</router-link>
              <router-link class="hover:text-primary transition-colors" to="/security">Security</router-link>
              <router-link class="hover:text-primary transition-colors" to="/faq">FAQ</router-link>
              <router-link class="hover:text-primary transition-colors" to="/contact">Support</router-link>
          </div>
      </div>
  </footer>
</template>

<script setup lang="ts">
</script>
