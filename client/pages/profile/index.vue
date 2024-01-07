<template>
    PROFILE HOME
</template>

<script setup>
definePageMeta({
    title: 'Profile',
    layout: 'profile',
    middleware: ['auth-check', 'user-data-check']
});

const userStore = useUserStore();
</script>