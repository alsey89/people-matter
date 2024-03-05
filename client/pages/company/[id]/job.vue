<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full max-h-32" />
        <!-- !Job Header -->
        <div class="min-h-max flex justify-between border-b-2 border-black py-2 pr-6">
            <div class="text-lg font-bold"> Jobs </div>
            <div v-auto-animate class="flex gap-4">
                <!-- !toggle job action buttons -->
                <Button @click="handleToggleAllJobsButtonClick" size="xs">
                    <Icon v-if="expandedJobIndex == 'all'" name="solar:list-arrow-up-bold"
                        class="text-primary-dark h-6 w-6" />
                    <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                </Button>
                <!-- !add job button -->
                <Button @click="handleAddJobButtonClick" size="xs">
                    <Icon v-if="showJobForm" name="material-symbols:close" class="text-primary-dark h-6 w-6" />
                    <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                </Button>
            </div>
        </div>
        <!-- !Job Form -->
        <div v-auto-animate>
            <AppCompanyJobForm v-if="showJobForm" :jobFormData="jobFormData" @submit="handleFormSubmit" />
        </div>
        <!-- !Job List -->
        <div v-if="companyStore.getCompanyJobs?.length > 0" v-for="(job, index) in companyStore.getCompanyJobs"
            :key="job.ID">
            <Card v-auto-animate class="p-4">
                <CardHeader>
                    <div v-auto-animate @click="handleExpandJobButtonClick(index)"
                        class="flex justify-between items-center px-2 hover:cursor-pointer">
                        <p :class="{ 'text-primary text-xl': expandedJobIndex === index }">
                            {{ job.title }}
                        </p>
                        <div class="flex gap-4">
                            <Button size="xs" @click.stop="handleEditJobButtonClick(job)">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </Button>
                            <Button size="xs" @click.stop="handleDeleteJobButtonClick(job)">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </Button>
                        </div>
                    </div>
                </CardHeader>
                <CardContent v-auto-animate>
                    <div v-if="expandedJobIndex == index || expandedJobIndex == 'all'" class="px-2">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div class="flex flex-col flex-wrap gap-2">
                                <p class="font-semibold border-b py-1">Description:</p>
                                <p v-if="job.description">
                                    <MDRender :content="job.description" />
                                </p>
                                <p v-else>
                                    No Data
                                </p>
                            </div>
                            <div class="flex flex-col gap-2">
                                <p class="font-semibold border-b py-1">Duties:</p>
                                <p v-if="job.duties">
                                    <MDRender :content="job.duties" />
                                </p>
                                <p v-else>
                                    No Data
                                </p>
                            </div>
                            <div class="flex flex-col gap-2">
                                <p class="font-semibold border-b py-1">Qualifications:</p>
                                <p v-if="job.qualifications">
                                    <MDRender :content="job.qualifications" />
                                </p>
                                <p v-else>
                                    No Data
                                </p>
                            </div>
                            <div class="flex flex-col gap-2">
                                <p class="font-semibold border-b py-1">Experience:</p>
                                <p v-if="job.experience">
                                    <MDRender :content="job.experience" />
                                </p>
                                <p v-else>
                                    No Data
                                </p>
                            </div>
                            <div>
                                <p class="font-semibold">Manager:</p>
                                <p v-if="companyStore.getManagerJobById(job.managerId)">{{
            companyStore.getManagerJobById(job.managerId) }}</p>
                                <p v-else>No Data</p>
                            </div>
                            <div>
                                <p class="font-semibold">Salary Range:</p>
                                <p>{{ job.minSalary }} - {{ job.maxSalary }}</p>
                            </div>
                            <div>
                                <p class="font-semibold">Location:</p>
                                <p v-if="companyStore.getLocationNameById(job.locationId)">{{
            companyStore.getLocationNameById(job.locationId) }}</p>
                                <p v-else>No Data</p>
                            </div>
                            <div>
                                <p class="font-semibold">Department:</p>
                                <p v-if="companyStore.getDepartmentNameById(job.departmentId)">{{
            companyStore.getDepartmentNameById(job.departmentId) }}</p>
                                <p v-else>No Data</p>
                            </div>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
        <div v-else>
            <Card>
                <div class="flex justify-center items-center">
                    No Data
                </div>
            </Card>
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
        companyStore.createJob({ companyId: companyStore.companyId, jobFormData: jobFormData });
    } else if (jobFormData.jobFormType == "edit") {
        companyStore.updateJob({ companyId: companyStore.companyId, jobId: jobFormData.jobId, jobFormData: jobFormData });
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
        await companyStore.deleteJob({ companyId: companyStore.companyId, jobId: jobId });

        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
</script>