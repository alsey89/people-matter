<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !CONFIRMATION MODAL -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full" />
        <!-- !COMPANY HEADER -->
        <div class="min-h-max flex justify-between border-b-2 border-black py-2 pr-4">
            <h1 class="text-lg font-bold"> Company </h1>
            <!-- !Add Company Button -->
            <Button v-if="!companyStore.getCompanyData" @click="handleAddCompanyButtonClick" size="xs">
                <Icon v-if="showCompanyForm" name="material-symbols:close" class="text-primary-dark h-6 w-6" />
                <Icon v-else name="material-symbols:add" class="h-6 w-6" />
            </Button>
            <Button v-if="showCompanyForm" @click="showCompanyForm = false" size="xs">
                <Icon name="material-symbols:close" class="text-primary-dark h-6 w-6" />
            </Button>
        </div>
        <div v-auto-animate class="w-full flex flex-col gap-4">
            <!-- !New Company Form -->
            <AppCompanyForm v-if="showCompanyForm" :formData="companyFormData" @submit="handleCompanyFormSubmit" />
            <!-- !Company List -->
            <div v-if="companyStore.getCompanyData">
                <Card class="w-full">
                    <div class="flex justify-between items-center p-4">
                        <div class="flex gap-4 items-center">
                            <AppLogo :src="companyStore.getCompanyLogoUrl || '/defaultLogo.png'" class="w-24 h-24"
                                shape="square" />
                            <div class="flex flex-col overflow-x-hidden">
                                <h1 class="text-sm md:text-lg text-nowrap font-bold"> {{ companyStore.getCompanyName }}
                                </h1>
                                <p class="text-sm md:text-lg text-nowrap"> {{ companyStore.getCompanyWebsite }} </p>
                                <div>
                                    <p>phone: {{ companyStore.getCompanyPhone }}</p>
                                    <p>email: {{ companyStore.getCompanyEmail }}</p>
                                    <p>
                                        {{ companyStore.getFullAddress }}
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="flex flex-col md:flex-row flex-wrap md:flex-nowrap gap-2 md:gap-4">
                            <Button size="xs" @click.stop="handleEditCompanyButtonClick">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </Button>
                            <Button size="xs" @click.stop="handleDeleteCompanyButtonClick">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </Button>
                        </div>
                    </div>
                </Card>
            </div>
            <div v-else>
                <Card>
                    <div v-if="companyStore.isLoading">
                        Loading...
                    </div>
                    <div v-else class="m-auto">
                        No Data
                    </div>
                </Card>
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

// add/edit company
const showCompanyForm = ref(false)
const companyFormData = reactive({
    companyFormType: null,
    companyId: null,
    companyName: null,
    companyPhone: null,
    companyEmail: null,
    companyWebsite: null,
    companyLogoUrl: null,
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
    showCompanyForm.value = !showCompanyForm.value
};
const clearCompanyForm = () => {
    companyFormData.companyName = ''
    companyFormData.companyPhone = ''
    companyFormData.companyWebsite = ''
    companyFormData.companyEmail = ''
    companyFormData.companyLogoUrl = ''
    companyFormData.companyAddress = ''
    companyFormData.companyCity = ''
    companyFormData.companyState = ''
    companyFormData.companyCountry = ''
    companyFormData.companyPostalCode = ''
};
const populateCompanyForm = () => {
    companyFormData.companyId = companyStore.getCompanyId
    companyFormData.companyName = companyStore.getCompanyName
    companyFormData.companyPhone = companyStore.getCompanyPhone
    companyFormData.companyEmail = companyStore.getCompanyEmail
    companyFormData.companyWebsite = companyStore.getCompanyWebsite
    companyFormData.companyLogoUrl = companyStore.getCompanyLogoUrl
    companyFormData.companyAddress = companyStore.getCompanyAddress
    companyFormData.companyCity = companyStore.getCompanyCity
    companyFormData.companyState = companyStore.getCompanyState
    companyFormData.companyCountry = companyStore.getCompanyCountry
    companyFormData.companyPostalCode = companyStore.getCompanyPostalCode
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
const handleDeleteCompanyButtonClick = () => {
    const companyId = companyStore.getCompanyId
    if (!companyId) {
        console.error("No company ID")
        messageStore.setError("Error deleting company")
        return
    }
    confirmationModalMessage.value = `Are you sure you want to delete ${companyStore.getCompanyName}? All data associated with this company (Departments, Locations, Jobs, etc.) will be deleted. This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showCompanyForm.value = false;
        await companyStore.deleteCompany({ companyId: companyId });
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
</script>