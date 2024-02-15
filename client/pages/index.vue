<template>
    <div class="w-full flex justify-center items-center">
        <NBCard class="min-w-full">
            This is a placeholder for the home page.
            You are currently logged in as:
            <div>
                {{ userStore.getCurrentUserFullName }}
                {{ userStore.getCurrentUserTitle }}
                {{ userStore.getCurrentUserDepartment }}
            </div>
            <!-- <div class="flex flex-col gap-4">
                <div v-for="(value, key) in userStore.getCurrentUserData" :key="key">
                    <NBText class="text-center">
                        {{ key }}: {{ value }}
                    </NBText>
                </div>
            </div> -->
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
    middleware: [],
});

const userStore = useUserStore();
onBeforeMount(() => {
    userStore.fetchCurrentUserData();
});
</script>