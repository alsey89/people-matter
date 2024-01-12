<template>
    <NBMenuItemNavigation navigateTo="/" :class="{ 'border-l-2 border-l-primary': activeRoute === '/' }"
        class="hover:bg-secondary-bg">
        Home
    </NBMenuItemNavigation>
    <NBMenuItemNavigation navigateTo="/user" :class="{ 'border-l-2 border-l-primary': activeRoute === '/user' }"
        class="hover:bg-secondary-bg">
        User
    </NBMenuItemNavigation>
    <NBMenuItemNavigation navigateTo="/job" :class="{ 'border-l-2 border-l-primary': activeRoute === '/job' }"
        class="hover:bg-secondary-bg">
        Job
    </NBMenuItemNavigation>
    <NBMenuItemNavigation navigateTo="/salary" :class="{ 'border-l-2 border-l-primary': activeRoute === '/salary' }"
        class="hover:bg-secondary-bg">
        Salary
    </NBMenuItemNavigation>
    <NBMenuItem @click="userStore.signout" class="hover:bg-secondary-bg">
        Signout
    </NBMenuItem>
</template>
  
<script setup>
const userStore = useUserStore()

const { navigateTo } = defineProps({
    navigateTo: String,
    required: true,
})

const router = useRouter();
const activeRoute = computed(() => {
    return router.currentRoute.value.path
})

</script>