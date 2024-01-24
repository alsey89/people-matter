<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full max-h-32" />
        <!-- !Job Header -->
        <div class="min-h-max flex justify-between border-b-2 border-black py-2 pr-4">
            <div class="text-lg font-bold"> Jobs </div>
            <div v-auto-animate class="flex gap-4">
                <!-- !toggle job action buttons -->
                <NBButtonSquare v-if="jobStore.getAllJobs && jobStore.getAllJobs?.length > 0"
                    @click="handleExpandJobButtonClick" size="xs">
                    <Icon v-if="expandJob" name="solar:list-arrow-up-bold" class="h-6 w-6" />
                    <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                </NBButtonSquare>
                <!-- !add job button -->
                <NBButtonSquare @click="handleAddJobButtonClick" size="xs">
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
                    <div class="flex justify-between items-center px-2">
                        <p> {{ job.title }} </p>
                        <div v-if="expandJob" class="flex gap-4">
                            <NBButtonSquare size="xs" @click.stop="handleEditJobButtonClick(job)">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                            <NBButtonSquare size="xs" @click.stop="handleDeleteJobButtonClick(job)">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                        </div>
                    </div>
                </NBCardHeader>
                <div v-if="expandJob" class="px-2">
                    <div class="grid grid-cols-2 gap-4">
                        <div class="flex flex-col gap-2">
                            <p class="font-semibold border-b py-1">Description:</p>
                            <p>
                                <MDRender :content="job.description" />
                            </p>
                        </div>
                        <div class="flex flex-col gap-2">
                            <p class="font-semibold border-b py-1">Duties:</p>
                            <p>
                                <MDRender :content="job.duties" />
                            </p>
                        </div>
                        <div class="flex flex-col gap-2">
                            <p class="font-semibold border-b py-1">Qualifications:</p>
                            <p>
                                <MDRender :content="job.qualifications" />
                            </p>
                        </div>
                        <div class="flex flex-col gap-2">
                            <p class="font-semibold border-b py-1">Experience:</p>
                            <p>
                                <MDRender :content="job.experience" />
                            </p>
                        </div>
                    </div>
                    <div class="mt-4">
                        <p class="font-semibold">Salary Range:</p>
                        <p>{{ job.minSalary }} - {{ job.maxSalary }}</p>
                    </div>
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
//expand
const expandJob = ref(false);
const handleExpandJobButtonClick = () => {
    showJobForm.value = false;
    expandJob.value = !expandJob.value;
};
// add/edit
const showJobForm = ref(false);
const jobFormData = reactive({
    formType: "",
    jobId: null,

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
const handleAddJobButtonClick = () => {
    clearJobFormData();
    jobFormData.formType = "add";
    showJobForm.value = !showJobForm.value;
};
const handleEditJobButtonClick = (job) => {
    scrollTo({
        top: 0,
        left: 0,
        behavior: 'smooth'
    })
    populateJobFormData(job)
    showJobForm.value = !showJobForm.value;
};
const populateJobFormData = (job) => {
    jobFormData.formType = "edit";
    jobFormData.jobId = job.ID;

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
};
const clearJobFormData = () => {
    jobFormData.title = null;
    jobFormData.description = null;
    jobFormData.duties = null;
    jobFormData.qualifications = null;
    jobFormData.experience = null;
    jobFormData.departmentId = null;
    jobFormData.locationId = null;
    jobFormData.managerId = null;
    jobFormData.minSalary = null;
    jobFormData.maxSalary = null;
};
const handleFormSubmit = () => {
    if (jobFormData.formType == "add") {
        jobStore.createJob({ companyId: activeCompanyId, jobFormData: jobFormData });
    } else if (jobFormData.formType == "edit") {
        jobStore.updateJob({ companyId: activeCompanyId, jobId: jobFormData.jobId, jobFormData: jobFormData });
    }

    showJobForm.value = false;
};
// delete
const showConfirmationModal = ref(false);
const confirmationModalMessage = ref(null);
const handleModalConfirmEvent = ref(() => { });
const handleModalCancelEvent = () => {
    showConfirmationModal.value = false;
};
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
// life cycle functions
onBeforeMount(() => {
    jobStore.fetchJobList({ companyId: activeCompanyId });
});
</script>