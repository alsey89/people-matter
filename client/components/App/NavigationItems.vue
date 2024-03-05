<template>
    <NuxtLink to="/" :class="{ 'border-l-4 border-l-primary': activeRoute === '/' }" class="p-2 hover:bg-secondary-bg">
        Home
    </NuxtLink>
    <NuxtLink to="/company" :class="{ 'border-l-4 border-l-primary': activeRoute === '/company' }"
        class="p-2 hover:bg-secondary-bg">
        Company
    </NuxtLink>
    <NuxtLink to="/user/list" :class="{ 'border-l-4 border-l-primary': activeRoute === '/user/list' }"
        class="p-2 hover:bg-secondary-bg">
        Users
    </NuxtLink>
    <NuxtLink to="/user" :class="{ 'border-l-4 border-l-primary': activeRoute === '/user' }"
        class="p-2 hover:bg-secondary-bg">
        Profile
    </NuxtLink>
    <Button @click="userStore.signout" class="mt-auto hover:bg-secondary-bg">
        Signout
    </Button>
</template>

<script setup>
const userStore = useUserStore()

const { to } = defineProps({
    to: String,
    required: true,
})

const router = useRouter();
const activeRoute = computed(() => {
    return router.currentRoute.value.path
})

</script>