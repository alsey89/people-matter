<template>

    <section class="fixed top-0 left-0 right-0 w-full h-24 flex items-center bg-background border-b text-primary">
        <NuxtImg src="/logo.svg" alt="logo" class="w-48 h-full px-4 border-r" />

        <nav class="w-full h-full flex justify-center items-center gap-4">
            <Button variant="ghost" class="h-full text-secondary" :class="{
                'text-primary': currentPath === '/',
            }" @click="navigateTo('/')">
                Home
            </Button>
            <Button variant="ghost" class="h-full text-secondary" :class="{
                'text-primary': currentPath === '/about',
            }" @click="navigateTo('/about')">About</Button>
            <Button variant="ghost" class="h-full text-secondary" :class="{
                'text-primary': currentPath === '/pricing',
            }" @click="navigateTo('/pricing')">Pricing</Button>
        </nav>

        <Button class="w-24 h-full"> Try it now </Button>
    </section>

    <section class="fixed top-24 left-0 right-0 w-full h-[640px] border-b">
        <slot name="hero"></slot>
    </section>

    <section class="w-full h-screen flex justify-center items-center">
        <div class="w-full h-full max-w-[1280px] pt-[736px]">
            <slot name="default"></slot>
        </div>
    </section>

</template>

<script setup>
const route = useRoute();
const currentPath = ref(route.path);
</script>
