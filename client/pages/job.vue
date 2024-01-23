<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full max-h-32" />
        <!-- !Job Header -->
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
        <div v-if="jobStore.getAllJobs?.length > 0" v-for="job in jobStore.getAllJobs" :key="job.ID">
            <NBCard v-auto-animate>
                <NBCardHeader>
                    <div v-auto-animate class="flex justify-between items-center px-2">
                        <p> {{ job.title }} </p>
                        <div v-if="expandJob" class="flex gap-4">
                            <NBButtonSquare size="sm" @click.stop="handleEditJobButtonClick(job)">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                            <NBButtonSquare size="sm" @click.stop="handleDeleteJobButtonClick(job)">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                        </div>
                    </div>
                </NBCardHeader>
                <div v-if="expandJob" class="px-2">
                    {{ job.description }}
                    {{ job.duties }}
                    {{ job.qualifications }}
                    {{ job.experience }}
                    {{ job.departmentId }}
                    {{ job.locationId }}
                    {{ job.managerId }}
                    {{ job.minSalary }}
                    {{ job.maxSalary }}
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
</template>

<script setup>
const jobStore = useJobStore();
const messageStore = useMessageStore();
const activeCompanyId = persistedState.sessionStorage.getItem('activeCompanyId');

const jobFormData = reactive({
    companyId: parseInt(activeCompanyId),

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
const handleEditJobButtonClick = (job) => {
    jobFormData.title = job.title;
    jobFormData.description = job.description;
    jobFormData.duties = job.duties;
    jobFormData.qualifications = job.qualifications;
    jobFormData.experience = job.experience;
    jobFormData.departmentId = job.departmentId;
    jobFormData.locationId = job.locationId;
    jobFormData.managerId = job.managerId;
    jobFormData.minSalary = job.minSalary;
    jobFormData.maxSalary = job.maxSalary;
    showJobForm.value = true;
};

const handleModalConfirmEvent = ref(() => { });
const handleDeleteJobButtonClick = (job) => {
    const jobId = job.ID;
    if (!jobId) {
        console.error('No job ID');
        messageStore.setError('Error deleting job');
        return;
    }
    confirmationModalMessage.value = `Are you sure you want to delete ${job.title}? This action cannot be undone.`;
    showConfirmationModal.value = true;

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showJobForm.value = false;
        await jobStore.deleteJob({ companyId: activeCompanyId, jobId: jobId });
        // if there is only one job left, hide the job list
        if (jobStore.getAllJobs.length == 1) {
            showJobList.value = false;
        }
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
const handleModalCancelEvent = () => {
    showConfirmationModal.value = false;
};

const handleFormSubmit = () => {
    jobStore.createJob({ companyId: activeCompanyId, jobFormData: jobFormData });
    showJobForm.value = false;
};

onBeforeMount(() => {
    jobStore.fetchJobList({ companyId: activeCompanyId });
});
</script>