<template>
    <div class="w-full h-full flex flex-col gap-4">
        <Card class="flex flex-col p-4">
            <form @submit.prevent="handleSubmit" class="flex flex-col gap-4">
                <!-- Input for Title -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="title"> Title* </label>
                    <input v-model.trim="jobFormData.jobTitle" type="text" name="title" id="title"
                        placeholder="Title of the Position"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary"
                        required />
                </div>
                <div class="w-full flex flex-wrap md:flex-nowrap gap-2">
                    <!-- Dropdown Selector for Department -->
                    <div v-auto-animate class="w-full md:w-1/2 flex flex-col gap-2">
                        <label for="department"> Department* </label>
                        <AppDropdown v-model="jobFormData.jobDepartmentId" :options="companyStore.getCompanyDepartments"
                            :required="true" @selected="handleDepartmentSelection" />
                    </div>
                    <!-- Dropdown Selector for Location -->
                    <div v-auto-animate class="w-full md:w-1/2 flex flex-col gap-2">
                        <label for="location"> Location* </label>
                        <AppDropdown v-model="jobFormData.jobLocationId" :options="companyStore.getCompanyLocations"
                            :required="true" @selected="handleLocationSelection" />
                    </div>
                </div>
                <!-- Input for Manager -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="manager"> Manager </label>
                    <AppDropdown v-model="jobFormData.jobManagerId" :options="companyStore.getCompanyJobs"
                        :required="false" @selected="handleManagerSelection" />
                </div>
                <!-- Input for Salary Range -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="salaryRange"> Salary Range [optional] </label>
                    <div class="flex items-center gap-2">
                        <input v-model.number="jobFormData.jobMinSalary" type="number" name="salaryRange"
                            id="salaryRange" placeholder="Min"
                            class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                        -
                        <input v-model.number="jobFormData.jobMaxSalary" type="number" name="salaryRange"
                            id="salaryRange" placeholder="Max"
                            class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                    </div>
                </div>
                <!-- Input for Description -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="description"> Description [optional] </label>
                    <textarea v-model.trim="jobFormData.jobDescription" type="text" name="description" id="description"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Input for Duties -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="duties"> Duties [optional] </label>
                    <textarea v-model.trim="jobFormData.jobDuties" type="text" name="duties" id="duties"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Input for Qualifications -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="qualifications"> Qualifications [optional] </label>
                    <textarea v-model.trim="jobFormData.jobQualifications" type="text" name="qualifications"
                        id="qualifications"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Input for Experience -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="experience"> Experience [optional] </label>
                    <textarea v-model.trim="jobFormData.jobExperience" type="text" name="experience" id="experience"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Submit Button -->
                <div class="w-full mt-4">
                    <Button type="submit" size="md" textSize="md"
                        class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">
                        Submit
                    </Button>
                </div>
            </form>
        </Card>
    </div>
</template>


<script setup>
const companyStore = useCompanyStore()

const { jobFormData } = defineProps({
    jobFormData: Object
});

const handleDepartmentSelection = (newSelection) => {
    jobFormData.departmentId = newSelection;
};
const handleLocationSelection = (newSelection) => {
    jobFormData.locationId = newSelection;
};
const handleManagerSelection = (newSelection) => {
    jobFormData.managerId = newSelection;
};

const emit = defineEmits(['submit']);
const handleSubmit = () => {
    emit('submit');
};
</script>