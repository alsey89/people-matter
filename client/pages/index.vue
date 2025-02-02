<template>
    <NuxtLayout name="landing">

        <template #hero>
            <!-- Hero -->
            <section class="w-full h-full flex justify-center items-center">
                <div class="w-1/2 h-full border-b">
                    test
                </div>
                <div class="w-0.5 h-full bg-border">

                </div>
                <div class="w-1/2 h-full border-b">
                    test
                </div>
            </section>
        </template>

        <template #default>
            <div class="w-full h-full bg-blue-500">
                test
            </div>
        </template>


    </NuxtLayout>
</template>

<script setup>
definePageMeta({
    title: 'Home',
    description: 'This is a starter template example.',
    layout: 'landing'
})
const isConfirmed = ref(false)
</script>