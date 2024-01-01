<template>
    <div class="w-120 flex justify-center">
        <NBCardHorizontal class="w-full">
            <button @click="showUserData = !showUserData"> show
            </button>
            <div v-if="showUserData" v-for="(value, key) in userStore.userData" :key="key">
                <div class="flex flex-row justify-between">
                    <NBText class="text-gray-500">{{ key }}</NBText>
                    <NBText class="text-gray-700">{{ value }}</NBText>
                </div>
                <hr class="border-gray-300" />
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

const showUserData = ref(false);
</script>