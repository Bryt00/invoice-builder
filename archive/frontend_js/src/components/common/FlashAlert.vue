<template>
  <div
    v-if="message"
    :class="[
      'p-4 mb-6 text-sm rounded-xl flex items-start gap-3 border transition-all duration-200',
      type === 'error'
        ? 'bg-error-container/60 border-error/20 text-on-error-container'
        : type === 'success'
        ? 'bg-primary-container/50 border-primary/20 text-on-surface'
        : 'bg-surface-variant/50 border-outline/20 text-on-surface'
    ]"
    role="alert"
  >
    <span class="material-symbols-outlined text-xl flex-shrink-0 mt-0.5">
      {{ type === 'error' ? 'warning' : type === 'success' ? 'check_circle' : 'info' }}
    </span>
    <div class="flex-1 font-body text-sm font-medium">
      {{ message }}
    </div>
    <button @click="clearFlash" class="text-outline hover:text-on-surface">
      <span class="material-symbols-outlined text-sm">close</span>
    </button>
  </div>
</template>

<script setup>
import { useFlash } from '../../composables/useFlash'

const { message, type, clearFlash } = useFlash()
</script>
