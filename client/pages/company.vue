<template>
    <div v-auto-animate class="w-full h-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full h-44" />
        <!-- !Company -->
        <div class="w-full flex flex-col gap-2">
            <div class="flex justify-between items-center border-b-2 border-black py-2">
                <h1 class="text-lg font-bold"> Company </h1>
                <div v-auto-animate class="flex gap-4">
                    <!-- !switch company button -->
                    <NBButtonSquare v-if="companyStore.getCompanyList.length > 1" @click="handleShowCompanyListButtonClick"
                        size="sm">
                        <Icon v-if="showCompanyList" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                        <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                    </NBButtonSquare>
                    <!-- !add company button -->
                    <NBButtonSquare @click="handleAddCompanyButtonClick" size="sm">
                        <Icon v-if="showCompanyForm" name="material-symbols:close" class="h-6 w-6" />
                        <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                    </NBButtonSquare>
                </div>
            </div>
            <div v-auto-animate class="w-full">
                <!-- !Company List -->
                <NBCard v-if="showCompanyList && !showCompanyForm" class="w-full">
                    <div v-if="companyStore.getCompanyList" v-for="company in companyStore.getCompanyList"
                        :key="company.ID">
                        <div @click="handleSelectCompany(company)"
                            class="flex justify-between items-center p-2 hover:cursor-pointer hover:bg-secondary-bg">
                            <div class="w-full flex gap-4 items-center">
                                <AppLogo :src="company.logoUrl" shape="square" class="w-16 h-16" />
                                <div class="flex flex-col">
                                    <h1 class="text-lg font-bold"> {{ company.name }} </h1>
                                    <p> {{ company.phone }} </p>
                                    <p> {{ company.address }} </p>
                                </div>
                                <div class="border-b-2 border-black">
                                </div>
                            </div>
                            <div class="flex gap-4">
                                <NBButtonSquare size="sm" @click.stop="handleEditCompanyButtonClick(company)">
                                    <Icon name="material-symbols:edit" class="h-6 w-6" />
                                </NBButtonSquare>
                                <NBButtonSquare size="sm" @click.stop="handleDeleteCompanyButtonClick(company)">
                                    <Icon name="material-symbols:delete" class="h-6 w-6" />
                                </NBButtonSquare>
                            </div>
                        </div>
                    </div>
                </NBCard>
                <!-- !New Company Form -->
                <NBCard v-else-if="showCompanyForm && !showCompanyList" class="w-full">
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
                    <div v-if="companyStore.getCompanyData" class="relative flex flex-wrap">
                        <div class="w-full md:w-3/12 flex justify-center items-center p-4">
                            <AppLogo :src="companyStore.getCompanyLogoUrl" shape="square" class="w-40 h-40" />
                        </div>
                        <div class="w-full md:w-9/12 my-auto flex flex-col gap-2 text-center md:text-left ">
                            <h1 class="text-lg font-bold"> {{ companyStore.getCompanyName }} </h1>
                            <p>phone: {{ companyStore.getCompanyPhone }}</p>
                            <p>email: {{ companyStore.getCompanyEmail }}</p>
                            <p>{{ companyStore.getFullAddress }}</p>
                        </div>
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
                    :key="department.ID">
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
                <div v-if="companyStore.companyTitles" v-for="title in companyStore.companyTitles" :key="title.ID">
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
                    :key="location.ID">
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
const messageStore = useMessageStore()

//! Company -----------------------------
const showCompanyForm = ref(false)
const showCompanyList = ref(false)
const showConfirmationModal = ref(false)
const confirmationModalMessage = ref("")
//form refs/ v-models
const formType = ref(null)
const companyId = ref(null)
const companyName = ref(null)
const companyPhone = ref(null)
const companyWebsite = ref(null)
const companyLogoUrl = ref(null)
const companyAddress = ref(null)
const companyCity = ref(null)
const companyState = ref(null)
const companyCountry = ref(null)
const companyPostalCode = ref(null)
//methods
const handleCompanyFormSubmit = async () => {
    const companyData = {};

    // List of all company fields
    const fields = [
        { ref: companyId, key: 'companyId' },
        { ref: companyName, key: 'companyName' },
        { ref: companyPhone, key: 'companyPhone' },
        { ref: companyWebsite, key: 'companyWebsite' },
        { ref: companyLogoUrl, key: 'companyLogoUrl' },
        { ref: companyAddress, key: 'companyAddress' },
        { ref: companyCity, key: 'companyCity' },
        { ref: companyState, key: 'companyState' },
        { ref: companyCountry, key: 'companyCountry' },
        { ref: companyPostalCode, key: 'companyPostalCode' },
    ];

    // Add field to the companyData object only if it has a value
    fields.forEach(({ ref, key }) => {
        if (ref.value) {
            companyData[key] = ref.value;
        }
    });

    if (formType.value === "edit") {
        await companyStore.updateCompany(companyData);
    } else if (formType.value === "add") {
        await companyStore.createCompany(companyData);
    }

    showCompanyForm.value = false;
};
const handleSelectCompany = async (company) => {
    const companyId = company.ID
    if (!companyId) {
        console.error("No company ID")
        messageStore.setError("Error selecting company")
        return
    }
    await companyStore.fetchOneCompanyData(companyId)
    showCompanyList.value = false
    showCompanyForm.value = false
}
const handleModalConfirmEvent = ref(null) //! stored function to be called when confirmation modal is confirmed
const handleModalCancelEvent = () => {
    showConfirmationModal.value = false
}
const handleDeleteCompanyButtonClick = (company) => {
    const companyId = company.ID
    if (!companyId) {
        console.error("No company ID")
        messageStore.setError("Error deleting company")
        return
    }
    confirmationModalMessage.value = `Are you sure you want to delete ${company.name}? This action cannot be undone.`
    showConfirmationModal.value = true

    // store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        // showCompanyList.value = false;
        showCompanyForm.value = false;
        await companyStore.deleteCompany(companyId);
        handleModalConfirmEvent.value = null;
    };
}
const handleShowCompanyListButtonClick = () => {
    showCompanyForm.value = false
    showCompanyList.value = !showCompanyList.value
}
const handleAddCompanyButtonClick = () => {
    clearForm()
    formType.value = "add"
    showCompanyList.value = false
    showCompanyForm.value = !showCompanyForm.value
}
const handleEditCompanyButtonClick = (company) => {
    populateForm(company)
    formType.value = "edit"
    showCompanyList.value = false
    showCompanyForm.value = true
}
const clearForm = () => {
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
const populateForm = (company) => {
    companyId.value = company.ID
    companyName.value = company.name
    companyPhone.value = company.phone
    companyWebsite.value = company.website
    companyLogoUrl.value = company.logoUrl
    companyAddress.value = company.address
    companyCity.value = company.city
    companyState.value = company.state
    companyCountry.value = company.country
    companyPostalCode.value = company.postalCode
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
    companyStore.fetchDefaultCompanyData()
})

</script>