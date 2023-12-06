<template>
    <div class="text-5xl text-red-500 font-bold underline">
        {{ userStore.getUserData }}
    </div>
    <button @click="userStore.signout" class="my-8 px-4 py-2.5 border border-gray-300 rounded-lg hover:bg-gray-200">
        SIGN OUT
    </button>
</template>

<script setup>
definePageMeta({
    title: 'Home',
    layout: 'default',
    middleware: ['auth-check', 'user-data-check'],
});
const userStore = useUserStore();
</script>