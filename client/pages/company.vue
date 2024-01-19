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
                                <NBButtonSquare v-if="companyStore.getCompanyList.length > 1" size="sm"
                                    @click.stop="handleDeleteCompanyButtonClick(company)">
                                    <Icon name="material-symbols:delete" class="h-6 w-6" />
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
                                    <label for="email">Email</label>
                                    <input v-model="companyEmail" type="email" name="email" id="email"
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
                                <AppLogo :src="`${companyLogoUrl}`" shape="square"
                                    class="hidden md:block w-full border-2 border-black hover:cursor-pointer" />
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
                                Submit
                            </NBButtonSquare>
                        </div>
                    </form>
                </NBCard>
                <!-- !Company Data -->
                <div v-else>
                    <NBCard v-if="companyStore.getCompanyData">
                        <div class="relative flex flex-wrap">
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
                            v-if="companyStore.getCompanyDepartments && companyStore.getCompanyDepartments.length > 0"
                            @click="handleShowDepartmentActionButtonClick" size="sm">
                            <Icon v-if="showCompanyList" name="solar:list-arrow-up-bold" class="h-6 w-6" />
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
                <div v-if="companyStore.getCompanyDepartments.length > 0"
                    v-for="department in companyStore.companyDepartments" :key="department.ID">
                    <NBCard v-auto-animate>
                        <NBCardHeader>
                            <div v-auto-animate class="flex justify-between px-2">
                                {{ department.name }}
                                <div v-if="showDepartmentActions" class="flex gap-4">
                                    <NBButtonSquare size="sm" @click.stop="handleEditDepartmentButtonClick(department)">
                                        <Icon name="material-symbols:edit" class="h-6 w-6" />
                                    </NBButtonSquare>
                                    <NBButtonSquare size="sm" @click.stop="handleDeleteDepartmentButtonClick(department)">
                                        <Icon name="material-symbols:delete" class="h-6 w-6" />
                                    </NBButtonSquare>
                                </div>
                            </div>
                        </NBCardHeader>
                        <div v-if="showDepartmentActions" class="px-2">
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
                        <NBButtonSquare v-if="companyStore.getCompanyTitles && companyStore.getCompanyTitles.length > 0"
                            @click="handleShowTitleActionButtonClick" size="sm">
                            <Icon v-if="showTitleForm" name="solar:list-arrow-up-bold" class="h-6 w-6" />
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
                <div v-if="companyStore.getCompanyTitles.length > 0" v-for="title in companyStore.companyTitles"
                    :key="title.ID">
                    <NBCard v-auto-animate>
                        <NBCardHeader>
                            <div v-auto-animate class="flex justify-between px-2">
                                {{ title.name }}
                                <div v-if="showTitleActions" class="flex gap-4">
                                    <NBButtonSquare size="sm" @click.stop="handleEditTitleButtonClick(title)">
                                        <Icon name="material-symbols:edit" class="h-6 w-6" />
                                    </NBButtonSquare>
                                    <NBButtonSquare size="sm" @click.stop="handleDeleteTitleButtonClick(title)">
                                        <Icon name="material-symbols:delete" class="h-6 w-6" />
                                    </NBButtonSquare>
                                </div>
                            </div>
                        </NBCardHeader>
                        <div v-if="showTitleActions" class="px-2">
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
                            v-if="companyStore.getCompanyLocations && companyStore.getCompanyLocations.length > 0"
                            @click="handleShowLocationActionButtonClick" size="sm">
                            <Icon v-if="showLocationForm" name="solar:list-arrow-up-bold" class="h-6 w-6" />
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
                <div v-if="companyStore.getCompanyLocations.length > 0" v-for="location in companyStore.companyLocations"
                    :key="location.ID">
                    <NBCard v-auto-animate>
                        <NBCardHeader>
                            <div v-auto-animate class="flex justify-between px-2">
                                <div>{{ location.name }} <span v-if="location.isHeadOffice">[Head Office]</span></div>
                                <div v-if="showLocationActions" class="flex gap-4">
                                    <NBButtonSquare size="sm" @click.stop="handleEditLocationButtonClick(location)">
                                        <Icon name="material-symbols:edit" class="h-6 w-6" />
                                    </NBButtonSquare>
                                    <NBButtonSquare size="sm" @click.stop="handleDeleteLocationButtonClick(location)">
                                        <Icon name="material-symbols:delete" class="h-6 w-6" />
                                    </NBButtonSquare>
                                </div>
                            </div>
                        </NBCardHeader>
                        <div v-if="showLocationActions" class="px-2">
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
const companyFormType = ref(null)
const companyId = ref(null)
const companyName = ref(null)
const companyPhone = ref(null)
const companyEmail = ref(null)
const companyWebsite = ref(null)
const companyLogoUrl = ref("defaultLogo.png")
const companyAddress = ref(null)
const companyCity = ref(null)
const companyState = ref(null)
const companyCountry = ref(null)
const companyPostalCode = ref(null)
//methods
const handleCompanyFormSubmit = async () => {
    const companyData = {
        companyId: companyId.value,
        companyName: companyName.value,
        companyPhone: companyPhone.value,
        companyEmail: companyEmail.value,
        companyWebsite: companyWebsite.value,
        companyLogoUrl: companyLogoUrl.value,
        companyAddress: companyAddress.value,
        companyCity: companyCity.value,
        companyState: companyState.value,
        companyCountry: companyCountry.value,
        companyPostalCode: companyPostalCode.value
    };

    if (companyFormType.value === "edit") {
        await companyStore.updateCompany(companyData);
    } else if (companyFormType.value === "add") {
        await companyStore.createCompany(companyData);
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
    companyFormType.value = "add"
    showCompanyForm.value = !showCompanyForm.value
};
const handleEditCompanyButtonClick = (company) => {
    populateCompanyForm(company)
    companyFormType.value = "edit"
    showCompanyList.value = false
    showCompanyForm.value = true
};

const clearCompanyForm = () => {
    companyName.value = ''
    companyPhone.value = ''
    companyWebsite.value = ''
    companyLogoUrl.value = 'defaultLogo.png'
    companyAddress.value = ''
    companyCity.value = ''
    companyState.value = ''
    companyCountry.value = ''
    companyPostalCode.value = ''
};
const populateCompanyForm = (company) => {
    companyId.value = company.ID
    companyName.value = company.name
    companyPhone.value = company.phone
    companyEmail.value = company.email
    companyWebsite.value = company.website
    companyLogoUrl.value = company.logoUrl
    companyAddress.value = company.address
    companyCity.value = company.city
    companyState.value = company.state
    companyCountry.value = company.country
    companyPostalCode.value = company.postalCode
};

//! Department -----------------------------
const showDepartmentForm = ref(false)
const showDepartmentActions = ref(false)
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
const handleShowDepartmentActionButtonClick = () => {
    showDepartmentForm.value = false
    showDepartmentActions.value = !showDepartmentActions.value
};
const handleAddDepartmentButtonClick = () => {
    showDepartmentActions.value = false
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
const showTitleActions = ref(false)

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
const handleShowTitleActionButtonClick = () => {
    showTitleForm.value = false
    showTitleActions.value = !showTitleActions.value
};
const handleAddTitleButtonClick = () => {
    showTitleActions.value = false
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
const showLocationActions = ref(false)

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
const handleShowLocationActionButtonClick = () => {
    showLocationForm.value = false
    showLocationActions.value = !showLocationActions.value
};
const handleAddLocationButtonClick = () => {
    showLocationActions.value = false
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
    showDepartmentActions.value = false
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