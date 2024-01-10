<template>
    <NBMenuItemNavigation navigateTo="/" :class="{ 'border border-black rounded-sm': $route.path === '/' }"
        class="hover:bg-primary-light">
        Home
    </NBMenuItemNavigation>
    <NBMenuItemNavigation navigateTo="/user" :class="{ 'border border-black rounded-sm': $route.path === '/user' }"
        class="hover:bg-primary-light">
        User
    </NBMenuItemNavigation>
    <NBMenuItemNavigation navigateTo="/job" :class="{ 'border border-black rounded-sm': $route.path === '/job' }"
        class="hover:bg-primary-light">
        Job
    </NBMenuItemNavigation>
    <NBMenuItemNavigation navigateTo="/salary" :class="{ 'border border-black rounded-sm': $route.path === '/salary' }"
        class="hover:bg-primary-light">
        Salary
    </NBMenuItemNavigation>
    <NBMenuItem @click="userStore.signout" class="hover:bg-primary-light">
        Signout
    </NBMenuItem>
</template>
  
<script setup>
const userStore = useUserStore()

const { navigateTo } = defineProps({
    navigateTo: String,
    required: true,
})

</script>