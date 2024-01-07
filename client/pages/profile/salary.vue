<template>
    SALARY INFO
</template>

<script setup>
definePageMeta({
    title: 'Salary Information',
    layout: 'profile',
    middleware: ['auth-check', 'user-data-check']
});


</script>