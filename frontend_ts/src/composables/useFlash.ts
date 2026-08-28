import { ref } from 'vue'

export type FlashType = 'info' | 'success' | 'error'

const message = ref('')
const type = ref<FlashType>('info')

export function useFlash() {
  function showFlash(msg: string, msgType: FlashType = 'info', duration = 5000) {
    message.value = msg
    type.value = msgType
    if (duration > 0) {
      setTimeout(() => {
        clearFlash()
      }, duration)
    }
  }

  function clearFlash() {
    message.value = ''
    type.value = 'info'
  }

  return {
    message,
    type,
    showFlash,
    clearFlash,
  }
}
