<template>
    <AppJobForm :formData="formData" @submit="handleFormSubmit" />
</template>

<script setup>

const formData = reactive({
    title: null,
    department: null,
    location: null,
    manager: null,
});

const handleFormSubmit = () => {
    console.log(formData)
};

</script>