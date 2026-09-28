<template>
  <div class="fixed top-5 right-5 z-[9999] flex flex-col gap-3 max-w-sm w-full pointer-events-none px-4 sm:px-0">
    <transition-group
      enter-active-class="transform transition duration-300 ease-out"
      enter-from-class="translate-y-[-10px] opacity-0 scale-95"
      enter-to-class="translate-y-0 opacity-100 scale-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-90"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto flex items-center justify-between gap-3 px-4 py-3.5 rounded-2xl shadow-xl border backdrop-blur-xl transition-all font-body text-sm font-semibold"
        :class="[
          toast.type === 'error'
            ? 'bg-rose-950/90 border-rose-500/30 text-rose-100 shadow-rose-950/20'
            : toast.type === 'warning'
            ? 'bg-amber-950/90 border-amber-500/30 text-amber-100 shadow-amber-950/20'
            : toast.type === 'info'
            ? 'bg-slate-900/90 border-slate-700/50 text-slate-100 shadow-slate-950/20'
            : 'bg-emerald-950/90 border-emerald-500/30 text-emerald-100 shadow-emerald-950/20'
        ]"
      >
        <div class="flex items-center gap-3">
          <span class="material-symbols-outlined text-[20px] shrink-0" :class="{
            'text-rose-400': toast.type === 'error',
            'text-amber-400': toast.type === 'warning',
            'text-emerald-400': toast.type === 'success' || !toast.type,
            'text-sky-400': toast.type === 'info'
          }">
            {{ toast.type === 'error' ? 'error' : toast.type === 'warning' ? 'warning' : toast.type === 'info' ? 'info' : 'check_circle' }}
          </span>
          <span class="leading-tight">{{ toast.message }}</span>
        </div>
        <button
          @click="removeToast(toast.id)"
          type="button"
          class="p-1 rounded-lg hover:bg-white/10 text-white/70 hover:text-white transition-colors shrink-0 cursor-pointer"
        >
          <span class="material-symbols-outlined text-[16px]">close</span>
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'

const { toasts, removeToast } = useToast()
</script>
