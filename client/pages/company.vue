<template>
    <div v-auto-animate class="w-full h-full flex flex-col gap-4">
        <!-- !Tab Navigation -->
        <div v-auto-animate class="flex justify-start gap-0.5 items-center">
            <div v-for="tab in ['Company', 'Departments', 'Locations']" :key="tab" :class="getTabClasses(tab)"
                class="text-lg font-bold rounded-t-md cursor-pointer shadow-[2px_0px_0px_rgba(0,0,0,1)]"
                @click="setActiveTab(tab)">
                {{ tab }}
            </div>
            <div class="mt-auto flex-grow border-b-2 border-black"></div>
        </div>
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full max-h-32" />
        <!-- !Company -->
        <AppCompany v-if="activeTab == 'Company'" />
        <!-- !Department -->
        <div v-if="activeTab == 'Departments'" v-auto-animate class="w-full flex flex-col gap-2">
            <div class="flex justify-between border-b-2 border-black py-2">
                <div class=" items-center  text-lg font-bold"> Departments </div>
                <div v-auto-animate class="flex gap-4">
                    <!-- !toggle department action buttons -->
                    <NBButtonSquare
                        v-if="companyStore.getCompanyDepartments && companyStore.getCompanyDepartments?.length > 0"
                        @click="handleExpandDepartmentButtonClick" size="xs">
                        <Icon v-if="expandDepartment" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                        <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                    </NBButtonSquare>
                    <!-- !add department button -->
                    <NBButtonSquare @click="handleAddDepartmentButtonClick" size="xs">
                        <Icon v-if="showDepartmentForm" name="material-symbols:close" class="h-6 w-6" />
                        <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                    </NBButtonSquare>
                </div>
            </div>
            <!-- !department form -->
            <div v-auto-animate>
                <AppDepartmentForm @submit="handleDepartmentFormSubmit" v-if="showDepartmentForm"
                    :formData="departmentFormData" />
            </div>
            <!-- !department list -->
            <div v-if="companyStore.getCompanyDepartments?.length > 0" v-for="department in companyStore.companyDepartments"
                :key="department.ID">
                <NBCard v-auto-animate>
                    <NBCardHeader>
                        <div v-auto-animate class="flex justify-between items-center px-2">
                            <p> {{ department.name }} </p>
                            <div v-if="expandDepartment" class="flex gap-4">
                                <NBButtonSquare size="xs" @click.stop="handleEditDepartmentButtonClick(department)">
                                    <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                                </NBButtonSquare>
                                <NBButtonSquare size="xs" @click.stop="handleDeleteDepartmentButtonClick(department)">
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
                    <div class="m-auto">
                        No Data
                    </div>
                </NBCard>
            </div>
        </div>
        <!-- !Location -->
        <div v-if="activeTab == 'Locations'" v-auto-animate class="w-full flex flex-col gap-2">
            <div class="flex justify-between items-center border-b-2 border-black py-2">
                <h1 class="text-lg font-bold"> Locations </h1>
                <div v-auto-animate class="flex gap-4">
                    <!-- !toggle location actions button -->
                    <NBButtonSquare v-if="companyStore.getCompanyLocations && companyStore.getCompanyLocations?.length > 0"
                        @click="handleExpandLocationButtonClick" size="xs">
                        <Icon v-if="expandLocation" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                        <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                    </NBButtonSquare>
                    <!-- !add location button -->
                    <NBButtonSquare @click="handleAddLocationButtonClick" size="xs">
                        <Icon v-if="showLocationForm" name="material-symbols:close" class="h-6 w-6" />
                        <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                    </NBButtonSquare>
                </div>
            </div>
            <!-- !location form -->
            <div v-if="showLocationForm">
                <AppLocationForm @submit="handleLocationFormSubmit" :formData="locationFormData" />
            </div>
            <!-- !location list -->
            <div v-if="companyStore.getCompanyLocations?.length > 0" v-for="location in companyStore.companyLocations"
                :key="location.ID">
                <NBCard v-auto-animate>
                    <NBCardHeader>
                        <div v-auto-animate class="flex justify-between items-center px-2">
                            <div>{{ location.name }} <span v-if="location.isHeadOffice">[Head Office]</span></div>
                            <div v-if="expandLocation" class="flex gap-4">
                                <NBButtonSquare size="xs" @click.stop="handleEditLocationButtonClick(location)">
                                    <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                                </NBButtonSquare>
                                <NBButtonSquare size="xs" @click.stop="handleDeleteLocationButtonClick(location)">
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
                    <div class="m-auto">
                        No Data
                    </div>
                </NBCard>
            </div>
        </div>
    </div>
</template>

<script setup>
const companyStore = useCompanyStore()
const messageStore = useMessageStore()

// Tabs
const activeTab = ref('Company');
const setActiveTab = (tabName) => {
    activeTab.value = tabName;
};
const getTabClasses = (tabName) => {
    return {
        'mt-1 px-4 py-1.5 text-primary border-l border-t border-r border-black': activeTab.value === tabName,
        'mt-2 px-2 py-1 text-gray-500 border-gray-500 border-l border-t border-r border-b-2': activeTab.value !== tabName,
    };
};

//! Company -----------------------------


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
    showLocationForm.value = false
}

//! meta & loading -----------------------------
definePageMeta({
    title: 'Company',
    layout: 'default',
})
const selectedCompanyId = persistedState.sessionStorage.getItem('activeCompanyId')
onBeforeMount(() => {
    if (selectedCompanyId) {
        companyStore.fetchCompanyListAndExpandById(selectedCompanyId)
        return
    }
    companyStore.fetchCompanyListAndExpandDefault()
})

</script>