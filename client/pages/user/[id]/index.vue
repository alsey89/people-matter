<template>
    <!-- !Confirmation Modal -->
    <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
        @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full" />
    <!--! Content -->
    <div v-if="userStore.getSingleUserData" class="flex flex-col gap-4">
        <UICard>
            <h1 class="text-lg font-bold px-2">
                Contact Information
            </h1>
            <div class="flex flex-col gap-2 px-2">
                <div>
                    Email: {{ userStore.getSingleUserEmail || "No Data" }}
                </div>
                <div>
                    Mobile: {{ userStore.getSingleUserMobile || "No Data" }}
                </div>
                <div>
                    Address: {{ userStore.getSingleUserFullAddress || "No Data" }}
                </div>
            </div>
        </UICard>
        <UICard v-if="userStore.getSingleUserEmergencyContact">
            <h1 class="text-lg font-bold px-2">
                Emergency Contact Information
            </h1>
            <div v-if="userStore.getSingleUserEmergencyContact" class="flex flex-col gap-2 px-2">
                <div>
                    Name: {{ userStore.getSingleUserEmergencyContactFullName }}
                </div>
                <div>
                    Email: {{ userStore.getSingleUserEmergencyContactEmail || "No Data" }}
                </div>
                <div>
                    Mobile: {{ userStore.getSingleUserEmergencyContactMobile || "No Data" }}
                </div>
                <div>
                    Relation: {{ userStore.getSingleUserEmergencyContactRelation || "No Data" }}
                </div>
            </div>
            <div v-else class="m-auto">
                No Data
            </div>
        </UICard>
    </div>
    <UICard v-else>
        <div v-if="userStore.isLoading">
            Loading...
        </div>
        <div v-else class="m-auto">
            No Data
        </div>
    </UICard>
</template>

<script setup>
definePageMeta({
    title: "User Details",
    description: "User details",
    layout: "user-details",
});

const userStore = useUserStore();

const showConfirmationModal = ref(false);

const userFormData = reactive({
    userFormType: null,
    userId: null,

    // personal information
    firstName: '',
    middleName: '',
    lastName: '',
    email: '',
});

</script>

