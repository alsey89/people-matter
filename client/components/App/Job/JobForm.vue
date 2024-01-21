<template>
    <div class="w-full h-full flex flex-col gap-4">
        <NBCard class="flex flex-col p-4">
            <div class="text-lg font-bold"> Create New Job </div>
            <form @submit.prevent="handleSubmit" class="flex flex-col gap-4">
                <!-- Input for Job Title -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="title"> Job Title </label>
                    <input v-model="formData.title" type="text" name="title" id="title"
                        class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500" />
                </div>
                <!-- Dropdown Selector for Department -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="department"> Department </label>
                    <AppDropdown :options="companyStore.getCompanyDepartments" @selected="handleDepartmentSelection" />
                </div>
                <!-- Dropdown Selector for Location -->
                <div v-auto-animate class="flex flex-col gap-2">
                    <label for="location"> Location </label>
                    <AppDropdown :options="companyStore.getCompanyLocations" @selected="handleLocationSelection" />
                </div>
                <!-- Submit Button -->
                <div class="w-full mt-4">
                    <NBButtonSquare type="submit" size="sm" textSize="md"
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

const { formData } = defineProps({
    formData: Object
});

const handleTitleSelection = (newSelection) => {
    formData.title = newSelection;
};
const handleDepartmentSelection = (newSelection) => {
    formData.department = newSelection;
};
const handleLocationSelection = (newSelection) => {
    formData.location = newSelection;
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
  