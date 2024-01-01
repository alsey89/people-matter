<template>
    <div v-if="messageStore.getMessage || messageStore.getError"
        :class="`fixed bottom-20 md:bottom-auto top-auto md:top-20 inset-x-0 mx-auto z-50 w-full max-w-xl px-8 py-2 rounded-lg shadow ${messageStore.getError ? 'border-red-500 bg-red-50' : 'border-green-500 bg-green-50'}`">
        <div class="flex justify-between py-3">
            <div class="flex">
                <div>
                    <svg xmlns="http://www.w3.org/2000/svg"
                        :class="`h-6 w-6 rounded-full ${messageStore.getError ? 'text-red-500' : 'text-green-500'}`"
                        viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd"
                            d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                            clipRule="evenodd" />
                    </svg>
                </div>
                <div class="self-center ml-3">
                    <span :class="`${messageStore.getError ? 'text-red-600' : 'text-green-600'} font-semibold`">
                        Message
                    </span>
                    <p :class="`${messageStore.getError ? 'text-red-600' : 'text-green-600'} mt-1`">
                        {{ messageStore.getMessage || messageStore.getError }}
                    </p>
                </div>
            </div>
            <button @click="messageStore.clearMessages"
                class="self-start `${messageStore.getError ? 'text-red-500' : 'text-green-500'}">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd"
                        d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                        clipRule="evenodd" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup>
const messageStore = useMessageStore();
</script>