<template>
    <div v-auto-animate class="w-full h-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full max-h-32" />
        <!-- !Company -->
        <div class="w-full flex flex-col gap-2">
            <div class="flex justify-between items-center border-b-2 border-black py-2">
                <h1 class="text-lg font-bold"> Company </h1>
                <div v-auto-animate class="flex gap-4">
                    <!-- !switch company button -->
                    <NBButtonSquare @click="handleShowCompanyListButtonClick" size="sm">
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
                        <div class="flex justify-between items-center p-2">
                            <div class="flex gap-4 items-center">
                                <AppLogo :src="company.logoUrl" shape="square" class="w-12 h-12" />
                                <div class="flex flex-col overflow-x-hidden">
                                    <h1 class="text-sm md:text-lg text-nowrap font-bold"> {{ company.name }} </h1>
                                    <p class="text-sm md:text-lg text-nowrap"> {{ company.website }} </p>
                                </div>
                                <div class="border-b-2 border-black">
                                </div>
                            </div>
                            <div class="flex gap-2 md:gap-4">
                                <NBButtonSquare size="sm" @click.stop="handleSelectCompany(company)">
                                    <Icon name="material-symbols:check" class="h-6 w-6 hover:text-primary" />
                                </NBButtonSquare>
                                <NBButtonSquare size="sm" @click.stop="handleEditCompanyButtonClick(company)">
                                    <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                                </NBButtonSquare>
                                <NBButtonSquare v-if="companyStore.getCompanyList?.length > 1" size="sm"
                                    @click.stop="handleDeleteCompanyButtonClick(company)">
                                    <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                                </NBButtonSquare>
                                <NBButtonSquare v-else size="sm"
                                    @click.stop="messageStore.setError('Cannot delete last company!')">
                                    <Icon name="material-symbols:delete" class="h-6 w-6 text-gray-400" />
                                </NBButtonSquare>
                            </div>
                        </div>
                    </div>
                </NBCard>
                <!-- !New Company Form -->
                <AppCompanyForm v-else-if="showCompanyForm && !showCompanyList" :formData="companyFormData"
                    @submit="handleCompanyFormSubmit" />
                <!-- !Company Data -->
                <div v-else>
                    <NBCard v-if="companyStore.getCompanyData">
                        <div class="relative flex flex-wrap">
                            <div class="w-full md:w-2/12 flex justify-center items-center p-4">
                                <AppLogo :src="companyStore.getCompanyLogoUrl" shape="square" class="w-24 h-24" />
                            </div>
                            <div class="w-full md:w-10/12 my-auto flex flex-col gap-2 text-center md:text-left ">
                                <h1 class="text-lg font-bold"> {{ companyStore.getCompanyName }} </h1>
                                <p>phone: {{ companyStore.getCompanyPhone }}</p>
                                <p>email: {{ companyStore.getCompanyEmail }}</p>
                                <p>{{ companyStore.getFullAddress }}</p>
                            </div>
                        </div>
                    </NBCard>
                    <NBCard v-else>
                        <div class="flex justify-center items-center">
                            No Data
                        </div>
                    </NBCard>
                </div>
            </div>
        </div>
        <div class="flex flex-col gap-4 pb-2">
            <!-- !Department -->
            <div v-auto-animate class="w-full flex flex-col gap-2">
                <div class="flex justify-between border-b-2 border-black py-2">
                    <div class=" items-center  text-lg font-bold"> Departments </div>
                    <div v-auto-animate class="flex gap-4">
                        <!-- !toggle department action buttons -->
                        <NBButtonSquare
                            v-if="companyStore.getCompanyDepartments && companyStore.getCompanyDepartments?.length > 0"
                            @click="handleExpandDepartmentButtonClick" size="sm">
                            <Icon v-if="expandDepartment" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                            <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                        </NBButtonSquare>
                        <!-- !add department button -->
                        <NBButtonSquare @click="handleAddDepartmentButtonClick" size="sm">
                            <Icon v-if="showDepartmentForm" name="material-symbols:close" class="h-6 w-6" />
                            <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                        </NBButtonSquare>
                    </div>
                </div>
                <!-- !department form -->
                <div v-auto-animate>
                    <AppCompanyDepartmentForm @submit="handleDepartmentFormSubmit" v-if="showDepartmentForm"
                        :formData="departmentFormData" />
                </div>
                <!-- !department list -->
                <div v-if="companyStore.getCompanyDepartments?.length > 0"
                    v-for="department in companyStore.companyDepartments" :key="department.ID">
                    <NBCard v-auto-animate>
                        <NBCardHeader>
                            <div v-auto-animate class="flex justify-between items-center px-2">
                                <p> {{ department.name }} </p>
                                <div v-if="expandDepartment" class="flex gap-4">
                                    <NBButtonSquare size="sm" @click.stop="handleEditDepartmentButtonClick(department)">
                                        <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                                    </NBButtonSquare>
                                    <NBButtonSquare size="sm" @click.stop="handleDeleteDepartmentButtonClick(department)">
                                        <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                                    </NBButtonSquare>
                                </div>
                            </div>
                        </NBCardHeader>
                        <div v-if="expandDepartment" class="px-2">
                            {{ department.description }}
                        </div>
                    </NBCard>
                </div>
                <div v-else>
                    <NBCard>
                        <div class="flex justify-center items-center">
                            No Data
                        </div>
                    </NBCard>
                </div>
            </div>
            <!-- !Title -->
            <div v-auto-animate class="w-full flex flex-col gap-2">
                <div class="flex justify-between items-center border-b-2 border-black py-2">
                    <h1 class="text-lg font-bold"> Titles </h1>
                    <div v-auto-animate class="flex gap-4">
                        <!-- !toggle title actions button -->
                        <NBButtonSquare v-if="companyStore.getCompanyTitles && companyStore.getCompanyTitles?.length > 0"
                            @click="handleExpandTitleButtonClick" size="sm">
                            <Icon v-if="expandTitle" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                            <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                        </NBButtonSquare>
                        <!-- !add title button -->
                        <NBButtonSquare @click="handleAddTitleButtonClick" size="sm">
                            <Icon v-if="showTitleForm" name="material-symbols:close" class="h-6 w-6" />
                            <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                        </NBButtonSquare>
                    </div>
                </div>
                <!-- !title form -->
                <div v-auto-animate>
                    <AppCompanyTitleForm @submit="handleTitleFormSubmit" v-if="showTitleForm" :formData="titleFormData" />
                </div>
                <!-- !title list -->
                <div v-if="companyStore.getCompanyTitles?.length > 0" v-for="title in companyStore.companyTitles"
                    :key="title.ID">
                    <NBCard v-auto-animate>
                        <NBCardHeader>
                            <div v-auto-animate class="flex justify-between items-center px-2">
                                {{ title.name }}
                                <div v-if="expandTitle" class="flex gap-4">
                                    <NBButtonSquare size="sm" @click.stop="handleEditTitleButtonClick(title)">
                                        <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                                    </NBButtonSquare>
                                    <NBButtonSquare size="sm" @click.stop="handleDeleteTitleButtonClick(title)">
                                        <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                                    </NBButtonSquare>
                                </div>
                            </div>
                        </NBCardHeader>
                        <div v-if="expandTitle" class="px-2">
                            {{ title.description }}
                        </div>
                    </NBCard>
                </div>
                <div v-else>
                    <NBCard>
                        <div class="flex justify-center items-center">
                            No Data
                        </div>
                    </NBCard>
                </div>
            </div>
            <!-- !Location -->
            <div v-auto-animate class="w-full flex flex-col gap-2">
                <div class="flex justify-between items-center border-b-2 border-black py-2">
                    <h1 class="text-lg font-bold"> Locations </h1>
                    <div v-auto-animate class="flex gap-4">
                        <!-- !toggle location actions button -->
                        <NBButtonSquare
                            v-if="companyStore.getCompanyLocations && companyStore.getCompanyLocations?.length > 0"
                            @click="handleExpandLocationButtonClick" size="sm">
                            <Icon v-if="expandLocation" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                            <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                        </NBButtonSquare>
                        <!-- !add location button -->
                        <NBButtonSquare @click="handleAddLocationButtonClick" size="sm">
                            <Icon v-if="showLocationForm" name="material-symbols:close" class="h-6 w-6" />
                            <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                        </NBButtonSquare>
                    </div>
                </div>
                <!-- !location form -->
                <div v-if="showLocationForm">
                    <AppCompanyLocationForm @submit="handleLocationFormSubmit" :formData="locationFormData" />
                </div>
                <!-- !location list -->
                <div v-if="companyStore.getCompanyLocations?.length > 0" v-for="location in companyStore.companyLocations"
                    :key="location.ID">
                    <NBCard v-auto-animate>
                        <NBCardHeader>
                            <div v-auto-animate class="flex justify-between items-center px-2">
                                <div>{{ location.name }} <span v-if="location.isHeadOffice">[Head Office]</span></div>
                                <div v-if="expandLocation" class="flex gap-4">
                                    <NBButtonSquare size="sm" @click.stop="handleEditLocationButtonClick(location)">
                                        <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                                    </NBButtonSquare>
                                    <NBButtonSquare size="sm" @click.stop="handleDeleteLocationButtonClick(location)">
                                        <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                                    </NBButtonSquare>
                                </div>
                            </div>
                        </NBCardHeader>
                        <div v-if="expandLocation" class="px-2">
                            <p>Phone: {{ location.phone || "No number listed" }}</p>
                            <p>Address: {{ location.address }}, {{ location.city }}, {{ location.state }}, {{
                                location.country }} {{
        location.postalCode }}</p>
                        </div>
                    </NBCard>
                </div>
                <div v-else>
                    <NBCard>
                        <div class="flex justify-center items-center">
                            No Data
                        </div>
                    </NBCard>
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
//refs/ v-models
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
//methods
const handleCompanyFormSubmit = async () => {
    if (companyFormData.companyFormType === "edit") {
        await companyStore.updateCompany({ companyFormData: companyFormData });
    } else if (companyFormData.companyFormType === "add") {
        await companyStore.createCompany({ companyFormData: companyFormData });
    }
    showCompanyForm.value = null;
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
};
const handleDeleteCompanyButtonClick = (company) => {
    const companyId = company.ID
    if (!companyId) {
        console.error("No company ID")
        messageStore.setError("Error deleting company")
        return
    }
    confirmationModalMessage.value = `Are you sure you want to delete ${company.name}? This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showCompanyForm.value = false;
        await companyStore.deleteCompany(companyId);
        // if there is only one company left, hide the company list
        if (companyStore.getCompanyList.length == 1) {
            showCompanyList.value = false
        }
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
const handleShowCompanyListButtonClick = () => {
    showCompanyForm.value = false
    showCompanyList.value = !showCompanyList.value
};
const handleAddCompanyButtonClick = () => {
    clearCompanyForm();
    companyFormData.companyFormType = "add"
    showCompanyForm.value = !showCompanyForm.value
};
const handleEditCompanyButtonClick = (company) => {
    populateCompanyForm(company)
    companyFormData.companyFormType = "edit"
    showCompanyList.value = false
    showCompanyForm.value = true
};

const clearCompanyForm = () => {
    companyFormData.companyName = ''
    companyFormData.companyPhone = ''
    companyFormData.companyWebsite = ''
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

//! Department -----------------------------
const showDepartmentForm = ref(false)
const expandDepartment = ref(false)
// form refs/ v-models
const departmentFormData = reactive({
    departmentFormType: null,
    departmentId: null,
    departmentName: null,
    departmentDescription: null,
})
// methods
const handleDepartmentFormSubmit = async () => {
    if (departmentFormData.departmentFormType == "edit") {
        await companyStore.updateDepartment({ companyId: companyStore.getCompanyId, departmentFormData: departmentFormData });
    } else if (departmentFormData.departmentFormType == "add") {
        await companyStore.createDepartment({ companyId: companyStore.getCompanyId, departmentFormData: departmentFormData });
    } else {
        console.error("No department form type")
        messageStore.setError("Error submitting department form")
        return
    }
    showDepartmentForm.value = false;
};
const handleExpandDepartmentButtonClick = () => {
    showDepartmentForm.value = false
    expandDepartment.value = !expandDepartment.value
};
const handleAddDepartmentButtonClick = () => {
    expandDepartment.value = false
    clearDepartmentForm()
    departmentFormData.departmentFormType = "add"
    showDepartmentForm.value = !showDepartmentForm.value
};
const handleEditDepartmentButtonClick = (department) => {
    departmentFormData.departmentFormType = "edit"
    populateDepartmentForm(department)
    showDepartmentForm.value = !showDepartmentForm.value
};
const handleDeleteDepartmentButtonClick = (department) => {
    const departmentId = department.ID
    if (!departmentId) {
        console.error("No department ID")
        messageStore.setError("Error deleting department")
        return
    }

    closeAllForms()

    confirmationModalMessage.value = `Are you sure you want to delete ${department.name}? This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showDepartmentForm.value = false;
        await companyStore.deleteDepartment({ companyId: companyStore.getCompanyId, departmentId: departmentId });
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
}

const clearDepartmentForm = () => {
    departmentFormData.departmentId = null
    departmentFormData.departmentName = ''
    departmentFormData.departmentDescription = ''
};
const populateDepartmentForm = (department) => {
    departmentFormData.departmentId = department.ID
    departmentFormData.departmentName = department.name
    departmentFormData.departmentDescription = department.description
};

//! Titles -----------------------------
const showTitleForm = ref(false)
const expandTitle = ref(false)

// refs / v-models
const titleFormData = reactive({
    titleFormType: null,
    titleId: null,
    titleName: null,
    titleDescription: null,
})

// methods
const handleTitleFormSubmit = async () => {
    if (titleFormData.titleFormType == "edit") {
        await companyStore.updateTitle({ companyId: companyStore.getCompanyId, titleFormData: titleFormData });
    } else if (titleFormData.titleFormType == "add") {
        await companyStore.createTitle({ companyId: companyStore.getCompanyId, titleFormData: titleFormData });
    } else {
        console.error("No title form type")
        messageStore.setError("Error submitting title form")
        return
    }
    showTitleForm.value = false;
};
const handleExpandTitleButtonClick = () => {
    showTitleForm.value = false
    expandTitle.value = !expandTitle.value
};
const handleAddTitleButtonClick = () => {
    expandTitle.value = false
    clearTitleForm()
    titleFormData.titleFormType = "add"
    showTitleForm.value = !showTitleForm.value
};
const handleEditTitleButtonClick = (title) => {
    titleFormData.titleFormType = "edit"
    populateTitleForm(title)
    showTitleForm.value = !showTitleForm.value
};
const handleDeleteTitleButtonClick = (title) => {
    const titleId = title.ID
    if (!titleId) {
        console.error("No title ID")
        messageStore.setError("Error deleting title")
        return
    }

    closeAllForms()

    confirmationModalMessage.value = `Are you sure you want to delete ${title.name}? This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showTitleForm.value = false;
        await companyStore.deleteTitle({ companyId: companyStore.getCompanyId, titleId: titleId });
        handleModalConfirmEvent.value = null; //! clear the stored function
    };

}
const populateTitleForm = (title) => {
    titleFormData.titleId = title.ID
    titleFormData.titleName = title.name
    titleFormData.titleDescription = title.description
};
const clearTitleForm = () => {
    titleFormData.titleId = null
    titleFormData.titleName = ''
    titleFormData.titleDescription = ''
};

//! Locations -----------------------------
const showLocationForm = ref(false)
const expandLocation = ref(false)

// form refs/ v-models
const locationFormData = reactive({
    locationFormType: null,
    locationId: null,
    locationName: null,
    locationIsHeadOffice: false,
    locationPhone: null,
    locationAddress: null,
    locationCity: null,
    locationState: null,
    locationCountry: null,
    locationPostalCode: null,
})

// methods
const handleLocationFormSubmit = async () => {
    if (locationFormData.locationFormType === "edit") {
        await companyStore.updateLocation({ companyId: companyStore.getCompanyId, locationFormData: locationFormData });
    } else if (locationFormData.locationFormType === "add") {
        await companyStore.createLocation({ companyId: companyStore.getCompanyId, locationFormData: locationFormData });
    } else {
        console.error("No location form type")
        messageStore.setError("Error submitting location form")
        return
    }
    showLocationForm.value = false;
};
const handleExpandLocationButtonClick = () => {
    showLocationForm.value = false
    expandLocation.value = !expandLocation.value
};
const handleAddLocationButtonClick = () => {
    expandLocation.value = false
    clearLocationForm()
    locationFormData.locationFormType = "add"
    showLocationForm.value = !showLocationForm.value
};
const handleEditLocationButtonClick = (location) => {
    locationFormData.locationFormType = "edit"
    populateLocationForm(location)
    showLocationForm.value = !showLocationForm.value
};
const handleDeleteLocationButtonClick = (location) => {
    const locationId = location.ID
    if (!locationId) {
        console.error("No location ID")
        messageStore.setError("Error deleting location")
        return
    }

    closeAllForms()

    confirmationModalMessage.value = `Are you sure you want to delete ${location.name}? This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showLocationForm.value = false;
        await companyStore.deleteLocation({ companyId: companyStore.getCompanyId, locationId: locationId });
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
const populateLocationForm = (location) => {
    locationFormData.locationId = location.ID;
    locationFormData.locationName = location.name;
    locationFormData.locationIsHeadOffice = location.isHeadOffice;
    locationFormData.locationPhone = location.phone;
    locationFormData.locationAddress = location.address;
    locationFormData.locationCity = location.city
    locationFormData.locationState = location.state
    locationFormData.locationCountry = location.country
    locationFormData.locationPostalCode = location.postalCode
};
const clearLocationForm = () => {
    locationFormData.locationId = null;
    locationFormData.locationName = '';
    locationFormData.locationIsHeadOffice = false;
    locationFormData.locationPhone = null;
    locationFormData.locationAddress = null;
    locationFormData.locationCity = null;
    locationFormData.locationState = null;
    locationFormData.locationCountry = null;
    locationFormData.locationPostalCode = null;
};



//! common
const showConfirmationModal = ref(false)
const confirmationModalMessage = ref("")
const handleModalConfirmEvent = ref(null) //! stored function to be called when confirmation modal is confirmed
//! methods
const handleModalCancelEvent = () => {
    showConfirmationModal.value = false
};
const closeAllForms = () => {
    showCompanyForm.value = false
    showCompanyList.value = false
    showDepartmentForm.value = false
    expandDepartment.value = false
    showTitleForm.value = false
    showLocationForm.value = false
}

//! meta & loading -----------------------------
definePageMeta({
    title: 'Company',
    layout: 'default',
})
onBeforeMount(() => {
    companyStore.fetchDefaultCompanyData()
})

</script>