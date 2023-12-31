<template>
    <div class="flex justify-center">
        <NBCardHorizontal class="w-1/2">
            <div class="text-wrap">
                {{ userStore.getUserData }}
            </div>
        </NBCardHorizontal>
    </div>
    <NBButtonSquare @click="userStore.signout">
        <NBText class="px-4"> Sign Out </NBText>
    </NBButtonSquare>
</template>

<script setup>

definePageMeta({
    title: 'Home',
    layout: 'default',
    middleware: ['auth-check', 'user-data-check'],
});
const userStore = useUserStore();
</script>