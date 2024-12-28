<template>
</template>

<script setup>
import { useAuthStore } from '@/stores/Auth';
import { onBeforeMount } from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter();

onBeforeMount(() => {
    const authStore = useAuthStore();
    authStore.signout(router);
});
</script>