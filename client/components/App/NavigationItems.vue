<template>
    <UIMenuItemNavigation navigateTo="/" :class="{ 'border-l-2 border-l-primary': activeRoute === '/' }"
        class="hover:bg-secondary-bg">
        Home
    </UIMenuItemNavigation>
    <UIMenuItemNavigation navigateTo="/company" :class="{ 'border-l-2 border-l-primary': activeRoute === '/company' }"
        class="hover:bg-secondary-bg">
        Company
    </UIMenuItemNavigation>
    <UIMenuItemNavigation navigateTo="/user/list" :class="{ 'border-l-2 border-l-primary': activeRoute === '/user' }"
        class="hover:bg-secondary-bg">
        Users
    </UIMenuItemNavigation>
    <UIMenuItemNavigation navigateTo="/user" :class="{ 'border-l-2 border-l-primary': activeRoute === '/salary' }"
        class="hover:bg-secondary-bg">
        Profile
    </UIMenuItemNavigation>
    <UIMenuItem @click="userStore.signout" class="hover:bg-secondary-bg">
        Signout
    </UIMenuItem>
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