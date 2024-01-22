<template>
    <div class="w-full h-full flex flex-col gap-4">
        <NBCard class="flex flex-col p-4">
            <form @submit.prevent="handleSubmit" class="flex flex-col gap-4">
                <!-- Input for Title -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="title"> Title* </label>
                    <input v-model="jobFormData.title" type="text" name="title" id="title" placeholder="Title of the Job"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary"
                        required />
                </div>
                <div class="w-full flex gap-2">
                    <!-- Dropdown Selector for Department -->
                    <div v-auto-animate class="w-1/2 flex flex-col gap-2">
                        <label for="department"> Department* </label>
                        <AppDropdown :options="companyStore.getCompanyDepartments" :required="true"
                            @selected="handleDepartmentSelection" />
                    </div>
                    <!-- Dropdown Selector for Location -->
                    <div v-auto-animate class="w-1/2 flex flex-col gap-2">
                        <label for="location"> Location* </label>
                        <AppDropdown :options="companyStore.getCompanyLocations" :required="true"
                            @selected="handleLocationSelection" />
                    </div>
                </div>
                <!-- Input for Manager -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="manager"> Manager </label>
                    <AppDropdown :options="jobStore.getAllJobs" :required="false" @selected="handleManagerSelection" />
                </div>
                <!-- Input for Salary Range -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="salaryRange"> Salary Range [optional] </label>
                    <div class="flex items-center gap-2">
                        <input v-model="jobFormData.minSalary" type="number" name="salaryRange" id="salaryRange"
                            placeholder="Min"
                            class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                        -
                        <input v-model="jobFormData.maxSalary" type="number" name="salaryRange" id="salaryRange"
                            placeholder="Max"
                            class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                    </div>
                </div>
                <!-- Input for Description -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="description"> Description [optional] </label>
                    <textarea v-model="jobFormData.description" type="text" name="description" id="description"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Input for Duties -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="duties"> Duties [optional] </label>
                    <textarea v-model="jobFormData.duties" type="text" name="duties" id="duties"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Input for Qualifications -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="qualifications"> Qualifications [optional] </label>
                    <textarea v-model="jobFormData.qualifications" type="text" name="qualifications" id="qualifications"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Input for Experience -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="experience"> Experience [optional] </label>
                    <textarea v-model="jobFormData.experience" type="text" name="experience" id="experience"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-primary" />
                </div>
                <!-- Submit Button -->
                <div class="w-full mt-4">
                    <NBButtonSquare type="submit" size="md" textSize="md"
                        class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">
                        Submit
                    </NBButtonSquare>
                </div>
            </form>
        </NBCard>
    </div>
</template>

  
<script setup>
const companyStore = useCompanyStore()
const jobStore = useJobStore()

const { jobFormData } = defineProps({
    jobFormData: Object
});

const handleDepartmentSelection = (newSelection) => {
    jobFormData.departmentId = newSelection.ID;
};
const handleLocationSelection = (newSelection) => {
    jobFormData.locationId = newSelection.ID;
};
const handleManagerSelection = (newSelection) => {
    jobFormData.managerId = newSelection.ID;
};

const emit = defineEmits(['submit']);
const handleSubmit = () => {
    emit('submit');
};

const activeCompanyId = persistedState.sessionStorage.getItem('activeCompanyId')
onBeforeMount(() => {
    companyStore.fetchCompanyListAndExpandById(activeCompanyId)
});
</script>
  