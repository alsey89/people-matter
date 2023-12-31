export const useMessageStore = defineStore("message-store", {
  state: () => ({
    error: null,
    message: null,
  }),
  getters: {
    getError: (state) => state.error,
    getMessage: (state) => state.message,
  },
  actions: {
    setMessage(incomingMessage) {
      this.message = incomingMessage;
      setTimeout(() => {
        this.message = null;
      }, 3000);
    },
    setError(incomingError) {
      this.error = incomingError;
      setTimeout(() => {
        this.error = null;
      }, 3000);
    },
    clearMessages() {
      this.message = null;
      this.error = null;
    },
  },
});
