<template>
    <div class="w-full h-full flex flex-col gap-4">
        <!-- !Company -->
        <div class="w-full flex flex-col gap-2">
            <div class="flex justify-between items-center border-b-2 border-black py-2">
                <h1 class="text-lg font-bold"> Company </h1>
                <NBButtonSquare @click="showCompanyForm = !showCompanyForm" size="sm">
                    <Icon v-if="showCompanyForm" name="material-symbols:close" class="h-6 w-6" />
                    <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                </NBButtonSquare>
            </div>
            <div v-auto-animate class="w-full">
                <!-- !New Company Form -->
                <NBCard v-if="showCompanyForm" class="w-full">
                    <form @submit.prevent="handleCompanyFormSubmit" class="w-full flex flex-col gap-4 p-2">
                        <div class="flex flex-wrap md:flex-nowrap justify-between gap-4">
                            <div class="w-full md:w-3/4 flex-flex-col gap-4">
                                <div class="flex flex-col">
                                    <label for="name">Name*</label>
                                    <input v-model="companyName" type="text" name="name" id="name"
                                        class="focus:outline-primary p-1" required />
                                </div>
                                <div class="flex flex-col">
                                    <label for="phone">Phone</label>
                                    <input v-model="companyPhone" type="text" name="phone" id="phone"
                                        class="focus:outline-primary p-1" />
                                </div>
                                <div class="flex flex-col">
                                    <label for="website"> Website </label>
                                    <input v-model="companyWebsite" type="text" name="name" id="website"
                                        class="focus:outline-primary p-1" />
                                </div>
                            </div>
                            <div class="w-full md:w-1/4 flex justify-center">
                                <!-- todo: logo upload -->
                                <NBButtonSquare type="button" size="sm" textSize="sm" class="w-full block md:hidden">
                                    <div class="flex justify-center items-center gap-2">
                                        <Icon name="material-symbols:upload" class="h-6 w-6" />
                                        Upload Photo
                                    </div>
                                </NBButtonSquare>
                                <AppLogo src="https://static.thenounproject.com/png/4974686-200.png" shape="square"
                                    class="hidden md:block w-32 border-2 border-black hover:cursor-pointer" />
                            </div>
                        </div>
                        <div class="flex flex-col">
                            <label for="address">Address</label>
                            <input v-model="companyAddress" type="text" name="address" id="address"
                                class="focus:outline-primary p-1" />
                        </div>
                        <div class="flex flex-wrap md:flex-nowrap gap-2 md:gap-4">
                            <div class="w-full md:w-1/2 flex flex-col gap-4">
                                <div class="flex flex-col">
                                    <label for="city">City</label>
                                    <input v-model="companyCity" type="text" name="city" id="city"
                                        class="focus:outline-primary p-1" />
                                </div>
                                <div class="flex flex-col">
                                    <label for="state">State</label>
                                    <input v-model="companyState" type="text" name="state" id="state"
                                        class="focus:outline-primary p-1" />
                                </div>
                            </div>
                            <div class="w-full md:w-1/2 flex flex-col gap-4">
                                <div class="flex flex-col">
                                    <label for="country">Country</label>
                                    <input v-model="companyCountry" type="text" name="country" id="country"
                                        class="focus:outline-primary p-1" />
                                </div>
                                <div class="flex flex-col">
                                    <label for="postal_code">Postal Code</label>
                                    <input v-model="companyPostalCode" type="text" name="postal_code" id="postal_code"
                                        class="focus:outline-primary p-1" />
                                </div>
                            </div>
                        </div>
                        <div class="w-full flex justify-end">
                            <NBButtonSquare type="submit" size="sm" textSize="md"
                                class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">
                                Save
                            </NBButtonSquare>
                        </div>
                    </form>
                </NBCard>
                <!-- !Company Data -->
                <NBCard v-else>
                    <div v-if="companyStore.getCompanyData" class="relative flex">
                        <div class="w-1/3 flex justify-center items-center p-4">
                            <AppLogo :src="companyStore.getCompanyLogoUrl" shape="square" class="max-h-full max-w-full" />
                        </div>
                        <div class="w-9/12 my-auto flex flex-col gap-2">
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
        </div>
        <div class="pb-2">
            <!-- !Department -->
            <div v-auto-animate class="w-full flex flex-col gap-2">
                <div class="flex justify-between items-center border-b-2 border-black py-2">
                    <h1 class="text-lg font-bold"> Departments </h1>
                    <NBButtonSquare @click="showDepartmentForm = !showDepartmentForm" size="sm">
                        <Icon v-if="showDepartmentForm" name="material-symbols:close" class="h-6 w-6" />
                        <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                    </NBButtonSquare>
                </div>
                <div v-if="showDepartmentForm">
                    <NBCard>
                        <template v-slot:header>
                            <div class="flex justify-between">
                                Add New Department
                            </div>
                        </template>
                        <form action="submit.prevent">
                            FORM
                        </form>
                    </NBCard>
                </div>
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
            <!-- !Title -->
            <div v-auto-animate class="w-full flex flex-col gap-2">
                <div class="flex justify-between items-center border-b-2 border-black py-2">
                    <h1 class="text-lg font-bold"> Titles </h1>
                    <NBButtonSquare @click="showTitleForm = !showTitleForm" size="sm">
                        <Icon v-if="showTitleForm" name="material-symbols:close" class="h-6 w-6" />
                        <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                    </NBButtonSquare>
                </div>
                <div v-if="showTitleForm">
                    <NBCard>
                        <template v-slot:header>
                            <div class="flex justify-between">
                                Add New Title
                            </div>
                        </template>
                        <form action="submit.prevent">
                            FORM
                        </form>
                    </NBCard>
                </div>
                <div v-if="companyStore.companyTitles" v-for="title in companyStore.companyTitles" :key="title.id">
                    <NBCard>
                        <h1> {{ title.name }} </h1>
                    </NBCard>
                </div>
                <div v-else>
                    No Data
                </div>
            </div>
            <!-- !Location -->
            <div v-auto-animate class="w-full flex flex-col gap-2">
                <div class="flex justify-between items-center border-b-2 border-black py-2">
                    <h1 class="text-lg font-bold"> Locations </h1>
                    <NBButtonSquare @click="showLocationForm = !showLocationForm" size="sm">
                        <Icon v-if="showLocationForm" name="material-symbols:close" class="h-6 w-6" />
                        <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                    </NBButtonSquare>
                </div>

                <div v-if="showLocationForm">
                    <NBCard>
                        <template v-slot:header>
                            <div class="flex justify-between">
                                Add New Location
                            </div>
                        </template>
                        <form action="submit.prevent">
                            FORM
                        </form>
                    </NBCard>
                </div>
                <div v-if="companyStore.companyLocations" v-for="location in companyStore.companyLocations"
                    :key="location.id">
                    <NBCard>
                        <h1> {{ location.name }} </h1>
                    </NBCard>
                </div>
                <div v-else>
                    No Data
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
const companyStore = useCompanyStore()

//! Company -----------------------------
const showCompanyForm = ref(false)
//v-models
const companyName = ref('')
const companyPhone = ref('')
const companyWebsite = ref('')
const companyLogoUrl = ref('')
const companyAddress = ref('')
const companyCity = ref('')
const companyState = ref('')
const companyCountry = ref('')
const companyPostalCode = ref('')
//methods
const handleCompanyFormSubmit = async () => {
    await companyStore.createCompany({
        companyName: companyName.value,
        companyPhone: companyPhone.value,
        companyWebsite: companyWebsite.value,
        companyLogoUrl: companyLogoUrl.value,
        companyAddress: companyAddress.value,
        companyCity: companyCity.value,
        companyState: companyState.value,
        companyCountry: companyCountry.value,
        companyPostalCode: companyPostalCode.value,
    })
    //close form
    showCompanyForm.value = false
    //clear refs
    companyName.value = ''
    companyPhone.value = ''
    companyWebsite.value = ''
    companyLogoUrl.value = ''
    companyAddress.value = ''
    companyCity.value = ''
    companyState.value = ''
    companyCountry.value = ''
    companyPostalCode.value = ''
}

//! Department -----------------------------
const showDepartmentForm = ref(false)


const showTitleForm = ref(false)
const showLocationForm = ref(false)


//! meta & loading -----------------------------
definePageMeta({
    title: 'Company',
    layout: 'default',
})
onBeforeMount(() => {
    companyStore.fetchCompanyData()
})

</script>