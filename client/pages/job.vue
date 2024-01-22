<template>
    <!-- !Confirmation Modal -->
    <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
        @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full max-h-32" />
    <!-- !Job Header -->
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <div class="min-h-max flex justify-between border-b-2 border-black py-2">
            <div class="text-lg font-bold"> Jobs </div>
            <div v-auto-animate class="flex gap-4">
                <!-- !toggle job action buttons -->
                <NBButtonSquare v-if="jobStore.getAllJobs && jobStore.getAllJobs?.length > 0"
                    @click="handleExpandJobButtonClick" size="sm">
                    <Icon v-if="expandJob" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                    <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                </NBButtonSquare>
                <!-- !add job button -->
                <NBButtonSquare @click="handleAddJobButtonClick" size="sm">
                    <Icon v-if="showJobForm" name="material-symbols:close" class="h-6 w-6" />
                    <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                </NBButtonSquare>
            </div>
        </div>
        <!-- !Job Form -->
        <div v-auto-animate>
            <AppJobForm v-if="showJobForm" :jobFormData="jobFormData" @submit="handleFormSubmit" />
        </div>
        <!-- !Job List -->
        <div v-auto-animate class="flex flex-col gap-4">
            <div v-if="jobStore.getAllJobs && jobStore.getAllJobs?.length > 0">
                <div v-if="expandJob" v-auto-animate class="flex flex-col gap-4">
                    <div v-for="job in jobStore.getAllJobs" :key="job.id" :job="job">
                        {{ job.name }}
                    </div>
                </div>
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
</template>

<script setup>
const jobStore = useJobStore();

const jobFormData = reactive({
    companyId: parseInt(persistedState.sessionStorage.getItem('activeCompanyId')),

    title: null,
    description: null,
    duties: null,
    qualifications: null,
    experience: null,
    departmentId: null,
    locationId: null,
    managerId: null,
    minSalary: null,
    maxSalary: null,
});

const showConfirmationModal = ref(false);
const confirmationModalMessage = ref(null);

const showJobForm = ref(false);
const expandJob = ref(false);
const handleAddJobButtonClick = () => {
    showJobForm.value = !showJobForm.value;
};
const handleExpandJobButtonClick = () => {
    expandJob.value = !expandJob.value;
};

const handleFormSubmit = () => {
    jobStore.createJob({ jobFormData: jobFormData });
};

</script>