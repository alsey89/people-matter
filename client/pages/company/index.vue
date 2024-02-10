<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full" />
        <!-- !Company Header -->
        <div class="min-h-max flex justify-between border-b-2 border-black py-2 pr-4">
            <h1 class="text-lg font-bold"> Company </h1>
            <!-- !add company button -->
            <NBButtonSquare @click="handleAddCompanyButtonClick" size="xs">
                <Icon v-if="showCompanyForm" name="material-symbols:close" class="h-6 w-6" />
                <Icon v-else name="material-symbols:add" class="h-6 w-6" />
            </NBButtonSquare>
        </div>
        <div v-auto-animate class="w-full flex flex-col gap-4">
            <!-- !New Company Form -->
            <AppCompanyForm v-if="showCompanyForm" :formData="companyFormData" @submit="handleCompanyFormSubmit" />
            <!-- !Company List -->
            <div v-if="companyStore.getCompanyList && companyStore.getCompanyList.length > 0"
                v-for="company in companyStore.getCompanyList" :key="company.ID">
                <NBCard class="w-full">
                    <div class="flex justify-between items-center p-2">
                        <div class="flex gap-4 items-center">
                            <AppLogo :src="companyStore.getCompanyLogoUrl"
                                :class="company.ID === companyStore.getCompanyId ? 'w-24 h-24' : 'w-12 h-12'"
                                shape="square" />
                            <div class="flex flex-col overflow-x-hidden">
                                <h1 class="text-sm md:text-lg text-nowrap font-bold"> {{ company.name }} </h1>
                                <p class="text-sm md:text-lg text-nowrap"> {{ company.website }} </p>
                                <div v-if="company.ID === companyStore.getCompanyId">
                                    <p>phone: {{ companyStore.getCompanyPhone }}</p>
                                    <p>email: {{ companyStore.getCompanyEmail }}</p>
                                    <p>
                                        {{ companyStore.getFullAddress }}
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="flex flex-col md:flex-row flex-wrap md:flex-nowrap gap-2 md:gap-4">
                            <NBButtonSquare size="xs" @click.stop="handleSelectCompany(company)">
                                <Icon name="material-symbols:check" class="h-6 w-6 "
                                    :class="company.ID === companyStore.getCompanyId ? 'bg-primary text-white' : 'text-black hover:text-primary'" />
                            </NBButtonSquare>
                            <NBButtonSquare size="xs" @click.stop="handleEditCompanyButtonClick(company)">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                            <NBButtonSquare size="xs" @click.stop="handleDeleteCompanyButtonClick(company)">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                        </div>
                    </div>
                </NBCard>
            </div>
            <div v-else>
                <NBCard>
                    <div v-if="companyStore.isLoading">
                        Loading...
                    </div>
                    <div v-else class="m-auto">
                        No Data
                    </div>
                </NBCard>
            </div>
        </div>
    </div>
</template>

<script setup>
//meta
definePageMeta({
    title: "Company",
    description: "Company",
    layout: "company",
});
//stores
const companyStore = useCompanyStore()
const messageStore = useMessageStore()

// select company
const handleSelectCompany = async (company) => {
    const companyId = company.ID
    if (!companyId) {
        console.error("No company ID")
        messageStore.setError("Error selecting company")
        return
    }
    if (companyStore.getCompanyId === companyId) {
        showCompanyForm.value = false
        return
    }
    await companyStore.fetchCompanyListAndExpandById({ companyId: companyId })
    showCompanyForm.value = false
};
// add/edit company
const showCompanyForm = ref(false)
const companyFormData = reactive({
    companyFormType: null,
    companyId: null,
    companyName: null,
    companyPhone: null,
    companyEmail: null,
    companyWebsite: null,
    companyLogoUrl: ref("defaultLogo.png"),
    companyAddress: null,
    companyCity: null,
    companyState: null,
    companyCountry: null,
    companyPostalCode: null,
})

const handleAddCompanyButtonClick = () => {
    clearCompanyForm();
    companyFormData.companyFormType = "add"
    showCompanyForm.value = !showCompanyForm.value
};
const handleEditCompanyButtonClick = (company) => {
    populateCompanyForm(company)
    companyFormData.companyFormType = "edit"
    showCompanyForm.value = true
};
const clearCompanyForm = () => {
    companyFormData.companyName = ''
    companyFormData.companyPhone = ''
    companyFormData.companyWebsite = ''
    companyFormData.companyEmail = ''
    companyFormData.companyLogoUrl = 'defaultLogo.png'
    companyFormData.companyAddress = ''
    companyFormData.companyCity = ''
    companyFormData.companyState = ''
    companyFormData.companyCountry = ''
    companyFormData.companyPostalCode = ''
};
const populateCompanyForm = (company) => {
    companyFormData.companyId = company.ID
    companyFormData.companyName = company.name
    companyFormData.companyPhone = company.phone
    companyFormData.companyEmail = company.email
    companyFormData.companyWebsite = company.website
    companyFormData.companyLogoUrl = company.logoUrl
    companyFormData.companyAddress = company.address
    companyFormData.companyCity = company.city
    companyFormData.companyState = company.state
    companyFormData.companyCountry = company.country
    companyFormData.companyPostalCode = company.postalCode
};
const handleCompanyFormSubmit = async () => {
    showCompanyForm.value = false;
    if (companyFormData.companyFormType === "edit") {
        await companyStore.updateCompany({ companyFormData: companyFormData });
    } else if (companyFormData.companyFormType === "add") {
        await companyStore.createCompany({ companyFormData: companyFormData });
    } else {
        console.error("No company form type")
        messageStore.setError("Error submitting company form")
        return
    }
};
// delete
const showConfirmationModal = ref(false)
const confirmationModalMessage = ref("")
const handleModalConfirmEvent = ref(null)

const handleModalCancelEvent = () => {
    showConfirmationModal.value = false;
    handleModalConfirmEvent.value = null; //! clear the stored function
};
const handleDeleteCompanyButtonClick = (company) => {
    const companyId = company.ID
    if (!companyId) {
        console.error("No company ID")
        messageStore.setError("Error deleting company")
        return
    }
    confirmationModalMessage.value = `Are you sure you want to delete ${company.name}? All data associated with this company (Departments, Locations, Jobs, etc.) will be deleted. This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showCompanyForm.value = false;
        await companyStore.deleteCompany({ companyId: companyId });
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
//watchers
//show company form when company list is empty
watch(() => companyStore.getCompanyList, (newList) => {
    if (!newList || newList.length === 0) {
        showCompanyForm.value = true;
        companyFormData.companyFormType = "add"
    } else {
        showCompanyForm.value = false;
        clearCompanyForm();
    }
});

</script>