<template>
    <DefaultLayout>
        <div v-auto-animate class="w-full h-full flex flex-col gap-4">
            <!-- !Company Tabs -->
            <app-company-tabs />
            <!-- !Content -->
            <slot />
        </div>
    </DefaultLayout>
</template>

<script setup>
import DefaultLayout from '~/layouts/default.vue'

const companyStore = useCompanyStore()
onBeforeMount(() => {
    companyStore.fetchCompany()
});
</script>