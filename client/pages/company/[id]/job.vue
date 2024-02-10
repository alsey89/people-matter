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
                <NBButtonSquare @click="handleToggleAllJobsButtonClick" size="xs">
                    <Icon v-if="expandedJobIndex == 'all'" name="solar:list-arrow-up-bold" class="h-6 w-6" />
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
            <AppCompanyJobForm v-if="showJobForm" :jobFormData="jobFormData" @submit="handleFormSubmit" />
        </div>
        <!-- !Job List -->
        <div v-if="companyStore.getCompanyJobs?.length > 0" v-for="(job, index) in companyStore.getCompanyJobs"
            :key="job.ID">
            <NBCard v-auto-animate>
                <NBCardHeader>
                    <div v-auto-animate @click="handleExpandJobButtonClick(index)"
                        class="flex justify-between items-center px-2 hover:cursor-pointer">
                        <p> {{ job.title }} </p>
                        <div class="flex gap-4">
                            <NBButtonSquare size="xs" @click.stop="handleEditJobButtonClick(job)">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                            <NBButtonSquare size="xs" @click.stop="handleDeleteJobButtonClick(job)">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </NBButtonSquare>
                        </div>
                    </div>
                </NBCardHeader>
                <div v-if="expandedJobIndex == index || expandedJobIndex == 'all'" class="px-2">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div class="flex flex-col flex-wrap gap-2">
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
definePageMeta({
    title: "Job",
    description: "Job",
    layout: "company",
});

const companyStore = useCompanyStore();
const messageStore = useMessageStore();


const activeCompanyId = persistedState.sessionStorage.getItem('activeCompanyId');
// show/hide details
const expandedJobIndex = ref("all")
const handleExpandJobButtonClick = (index) => {
    if (expandedJobIndex.value == index) {
        expandedJobIndex.value = null
        return
    }
    expandedJobIndex.value = index
}
const handleToggleAllJobsButtonClick = () => {
    if (expandedJobIndex.value == 'all') {
        expandedJobIndex.value = null
        return
    }
    expandedJobIndex.value = 'all'
}



// add/edit
const showJobForm = ref(false);
const jobFormData = reactive({
    jobFormType: "",
    jobId: null,

    jobTitle: null,
    jobDescription: null,
    jobDuties: null,
    jobQualifications: null,
    jobExperience: null,
    jobDepartmentId: null,
    jobLocationId: null,
    jobManagerId: null,
    jobMinSalary: null,
    jobMaxSalary: null,
});
const handleAddJobButtonClick = () => {
    clearJobForm();
    jobFormData.jobFormType = "add";
    showJobForm.value = !showJobForm.value;
};
const handleEditJobButtonClick = (job) => {
    populateJobFormData(job)
    showJobForm.value = !showJobForm.value;
};
const populateJobFormData = (job) => {
    jobFormData.jobFormType = "edit";
    jobFormData.jobId = job.ID;

    jobFormData.jobTitle = job.title;
    jobFormData.jobDescription = job.description;
    jobFormData.jobDuties = job.duties;
    jobFormData.jobQualifications = job.qualifications;
    jobFormData.jobExperience = job.experience;
    jobFormData.jobDepartmentId = job.departmentId;
    jobFormData.jobLocationId = job.locationId;
    jobFormData.jobManagerId = job.managerId;
    jobFormData.jobMinSalary = job.minSalary;
    jobFormData.jobMaxSalary = job.maxSalary;
};
const clearJobForm = () => {
    jobFormData.jobFormType = "";
    jobFormData.jobId = null;

    jobFormData.jobTitle = null;
    jobFormData.jobDescription = null;
    jobFormData.jobDuties = null;
    jobFormData.jobQualifications = null;
    jobFormData.jobExperience = null;
    jobFormData.jobDepartmentId = null;
    jobFormData.jobLocationId = null;
    jobFormData.jobManagerId = null;
    jobFormData.jobMinSalary = null;
    jobFormData.jobMaxSalary = null;
};
const handleFormSubmit = () => {
    if (jobFormData.jobFormType == "add") {
        companyStore.createJob({ companyId: activeCompanyId, jobFormData: jobFormData });
    } else if (jobFormData.jobFormType == "edit") {
        companyStore.updateJob({ companyId: activeCompanyId, jobId: jobFormData.jobId, jobFormData: jobFormData });
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
        await companyStore.deleteJob({ companyId: activeCompanyId, jobId: jobId });

        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
// watchers
watch(() => companyStore.getCompanyJobs, (newJobList) => {
    if (!newJobList || newJobList.length < 1) {
        showJobForm.value = true;
        jobFormData.jobFormType = "add"
    }
}, { immediate: true });
</script>