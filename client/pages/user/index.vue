<template>
    <div class="relative w-full flex flex-col gap-2">
        <div class="absolute bottom-4 right-4">
            <NBButtonSquare size="sm">
                <Icon name="material-symbols:add" class="h-8 w-8" />
            </NBButtonSquare>
        </div>
        <div v-for="user in userStore.getAllUsersData" class="w-full">
            <div class="w-full flex bg-secondary-bg">
                <NBCard>
                    <div class="flex flex-wrap justify-between">
                        <div>
                            <h1>{{ user.firstName }} {{ user.lastName }}</h1>
                            <p>{{ user.email }}</p>
                        </div>
                        <div>
                            <p>{{ user.title }}</p>
                            <p>{{ user.department }}</p>
                        </div>
                    </div>

                </NBCard>
            </div>
        </div>
    </div>
</template>

<script setup>
definePageMeta({
    title: 'Users',
    layout: 'default',
})

const userStore = useUserStore();
onBeforeMount(() => {
    userStore.fetchAllUsersData();
});
</script>