<template>
    <NBCard>
        <form @submit.prevent="handleSubmit" class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="name">Name*</label>
                <input v-model="formData.departmentName" type="text" id="name" class="focus:outline-primary p-1" required />
            </div>
            <!-- Searchable Dropdown Selector for Title -->
            <!-- Searchable Dropdown Selector for Department -->
            <!-- Searchable Dropdown Selector for Location -->
            <!-- Searchable Dropdown Selector for Manager -->
            <!-- todo: add more fields if needed -->
            <div class="w-full">
                <NBButtonSquare type="submit" size="sm" textSize="md"
                    class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">Submit
                </NBButtonSquare>
            </div>
        </form>
    </NBCard>
</template>
  
<script setup>
const companyStore = useCompanyStore()

const { formData } = defineProps({
    formData: Object
});

const emit = defineEmits(['submit']);

const handleSubmit = () => {
    emit('submit');
};

onBeforeMount(() => {

});
</script>
  