<template>
    <div v-if="userStore.getAllUsersData" v-for="(user, index) in userStore.getAllUsersData" :key="user.id" class="w-full">
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
</template>

<script setup>
definePageMeta({
    title: 'Users',
    layout: 'users',
})

const userStore = useUserStore();
onBeforeMount(() => {
    userStore.fetchAllUsersData();
});

</script>

