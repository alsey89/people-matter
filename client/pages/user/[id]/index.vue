<template>
    <div>
        THIS PAGE SHOWS THE USER WITH ID: {{ route?.params?.id }}
    </div>
</template>

<script setup>
definePageMeta({
    title: "User Details",
    description: "User details",
    layout: "user-details",
});

const route = useRouter().currentRoute.value
</script>