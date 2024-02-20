<template>
    <div class="w-full flex justify-center items-center">
        <UICard class="min-w-full">
            This is a placeholder for the home page.
            You are currently logged in as:
            <div>
                {{ userStore.getCurrentUserFullName }}
                {{ userStore.getCurrentUserTitle }}
                {{ userStore.getCurrentUserDepartment }}
            </div>
            <!-- <div class="flex flex-col gap-4">
                <div v-for="(value, key) in userStore.getCurrentUserData" :key="key">
                    <UIText class="text-center">
                        {{ key }}: {{ value }}
                    </UIText>
                </div>
            </div> -->
            <UIButtonSquare @click="userStore.signout">
                <UIText class="px-4"> Sign Out </UIText>
            </UIButtonSquare>
        </UICard>
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