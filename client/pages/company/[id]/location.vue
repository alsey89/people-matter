<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full" />
        <!-- !Location Header -->
        <div class="min-h-max flex justify-between border-b-2 border-black py-2 pr-4">
            <h1 class="text-lg font-bold"> Locations </h1>
            <div v-auto-animate class="flex gap-4">
                <!-- !add location button -->
                <UIButtonSquare @click="handleAddLocationButtonClick" size="xs">
                    <Icon v-if="showLocationForm" name="material-symbols:close" class="text-primary-dark h-6 w-6" />
                    <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                </UIButtonSquare>
            </div>
        </div>
        <!-- !location form -->
        <div v-if="showLocationForm">
            <AppCompanyLocationForm @submit="handleLocationFormSubmit" :formData="locationFormData" />
        </div>
        <!-- !location list -->
        <div v-if="companyStore.getCompanyLocations?.length > 0" v-for="location in companyStore.companyLocations"
            :key="location.ID">
            <UICard v-auto-animate>
                <UICardHeader>
                    <div v-auto-animate class="flex justify-between items-center px-2">
                        <div>{{ location.name }} <span v-if="location.isHeadOffice">[Head Office]</span></div>
                        <div class="flex gap-4">
                            <UIButtonSquare size="xs" @click.stop="handleEditLocationButtonClick(location)">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </UIButtonSquare>
                            <UIButtonSquare size="xs" @click.stop="handleDeleteLocationButtonClick(location)">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </UIButtonSquare>
                        </div>
                    </div>
                </UICardHeader>
                <div class="px-2">
                    <p>Phone: {{ location.phone || "No number listed" }}</p>
                    <p>Address: {{ location.address }}, {{ location.city }}, {{ location.state }}, {{
                        location.country }} {{
        location.postalCode }}</p>
                </div>
            </UICard>
        </div>
        <div v-else>
            <UICard>
                <div class="m-auto">
                    No Data
                </div>
            </UICard>
        </div>
    </div>
</template>

<script setup>
definePageMeta({
    title: "Locations",
    description: "Manage locations",
    layout: "company",
})

const companyStore = useCompanyStore()
const messageStore = useMessageStore()

//add/edit location
const showLocationForm = ref(false)
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
const handleLocationFormSubmit = async () => {
    showLocationForm.value = false;
    if (locationFormData.locationFormType === "edit") {
        await companyStore.updateLocation({ companyId: companyStore.getCompanyId, locationFormData: locationFormData });
    } else if (locationFormData.locationFormType === "add") {
        await companyStore.createLocation({ companyId: companyStore.getCompanyId, locationFormData: locationFormData });
    } else {
        console.error("No location form type")
        messageStore.setError("Error submitting location form")
        return
    }
};
const handleAddLocationButtonClick = () => {
    clearLocationForm()
    locationFormData.locationFormType = "add"
    showLocationForm.value = !showLocationForm.value
};
const handleEditLocationButtonClick = (location) => {
    locationFormData.locationFormType = "edit"
    populateLocationForm(location)
    showLocationForm.value = !showLocationForm.value
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

// delete location
const showConfirmationModal = ref(false)
const confirmationModalMessage = ref("")
const handleModalConfirmEvent = ref(null) //! stored function to be called when confirmation modal is confirmed

const handleDeleteLocationButtonClick = (location) => {
    const locationId = location.ID
    if (!locationId) {
        console.error("No location ID")
        messageStore.setError("Error deleting location")
        return
    }

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
const handleModalCancelEvent = () => {
    showConfirmationModal.value = false
};
</script>
