<template>
    <ul class="w-full flex flex-col gap-4">
        {{ activeRoute }}
        <li :class="`px-2 py-1 rounded-md hover:bg-gray-200 ${activeRoute === '/profile' ? 'bg-gray-200' : ''}`">
            <NuxtLink to="/profile">Profile</NuxtLink>
        </li>

        <li :class="`px-2 py-1 rounded-md hover:bg-gray-200 ${activeRoute === '/profile/contact' ? 'bg-gray-200' : ''}`">
            <NuxtLink to="/profile/contact">Contact</NuxtLink>
        </li>

        <li :class="`px-2 py-1 rounded-md hover:bg-gray-200 ${activeRoute === '/profile/emergency' ? 'bg-gray-200' : ''}`">
            <NuxtLink to="/profile/emergency">Emergency</NuxtLink>
        </li>

        <li :class="`px-2 py-1 rounded-md hover:bg-gray-200 ${activeRoute === '/profile/job' ? 'bg-gray-200' : ''}`">
            <NuxtLink to="/profile/job">Job Info</NuxtLink>
        </li>

        <li :class="`px-2 py-1 rounded-md hover:bg-gray-200 ${activeRoute === '/profile/salary' ? 'bg-gray-200' : ''}`">
            <NuxtLink to="/profile/salary">Salary Info</NuxtLink>
        </li>
    </ul>
</template>

<script setup>
const route = useRoute();

const activeRoute = computed(() => route.fullPath);
</script>

