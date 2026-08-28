import { ref } from 'vue'

const message = ref('')
const type = ref('info') // 'info', 'success', 'error'

export function useFlash() {
  function showFlash(msg, msgType = 'info', duration = 5000) {
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
