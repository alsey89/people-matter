<template>
    <div class="relative w-full h-full flex flex-col gap-4">
        <div class="absolute bottom-4 right-4">
            <NBButtonSquare size="sm">
                <Icon name="material-symbols:add" class="h-8 w-8" />
            </NBButtonSquare>
        </div>
        <div class="w-full">
            <NBCard>
                <div v-if="companyStore.getCompanyData" class="grid grid-cols-12 gap-4">
                    <div class="col-span-3 my-auto p-4">
                        <AppLogo :src="companyStore.getCompanyLogoUrl" />
                    </div>
                    <div class="col-span-9 my-auto p-4 flex flex-col gap-2">
                        <h1 class="text-lg font-bold"> {{ companyStore.getCompanyName }} </h1>
                        <p>{{ companyStore.getFullAddress }}</p>
                        <p>phone: {{ companyStore.getCompanyPhone }}</p>
                        <p>email: {{ companyStore.getCompanyEmail }}</p>
                    </div>
                </div>
                <div v-else>
                    No Company Data
                </div>
            </NBCard>
        </div>
        <div class="w-full flex flex-col gap-2">
            <h1> Departments </h1>
            <div v-if="companyStore.companyDepartments" v-for="department in companyStore.companyDepartments"
                :key="department.id">
                <NBCard>
                    <h1> {{ department.name }} </h1>
                </NBCard>
            </div>
            <div v-else>
                No Data
            </div>
        </div>
        <div class="w-full flex flex-col gap-2">
            <h1> Titles </h1>
            <div v-if="companyStore.companyTitles" v-for="title in companyStore.companyTitles" :key="title.id">
                <NBCard>
                    <h1> {{ title.name }} </h1>
                </NBCard>
            </div>
            <div v-else>
                No Data
            </div>
        </div>
        <div class="w-full flex flex-col gap-2">
            <h1> Locations </h1>
            <div v-if="companyStore.companyLocations" v-for="location in companyStore.companyLocations" :key="location.id">
                <NBCard>
                    <h1> {{ location.name }} </h1>
                </NBCard>
            </div>
            <div v-else>
                No Data
            </div>
        </div>
    </div>
</template>

<script setup>

definePageMeta({
    title: 'Company',
    layout: 'default',
})

const companyStore = useCompanyStore()

onBeforeMount(() => {
    companyStore.fetchCompanyData()
})

</script>