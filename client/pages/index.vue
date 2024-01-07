<template>
    <div class="flex justify-center items-center">
        <NBCard>
            <div class="flex flex-col gap-4">
                <div v-for="(value, key) in userStore.getUserData" :key="key">
                    <NBText class="text-center">
                        {{ key }}: {{ value }}
                    </NBText>
                </div>
            </div>
            <NBButtonSquare @click="userStore.signout">
                <NBText class="px-4"> Sign Out </NBText>
            </NBButtonSquare>
        </NBCard>
    </div>
</template>

<script setup>

definePageMeta({
    title: 'Home',
    layout: 'default',
    middleware: ['auth-check', 'user-data-check'],
});
const userStore = useUserStore();

</script>