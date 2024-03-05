<template>
    <Card>
        <form @submit.prevent="handleSubmit" class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="name">Name*</label>
                <input v-model="formData.departmentName" type="text" id="name" class="focus:outline-primary p-1"
                    required />
            </div>
            <div class="flex flex-col">
                <label for="description">Description</label>
                <textarea v-model="formData.departmentDescription" id="description"
                    class="focus:outline-primary p-1"></textarea>
            </div>
            <!-- todo: add more fields if needed -->
            <div class="w-full">
                <Button type="submit" size="xs" textSize="md"
                    class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">Submit
                </Button>
            </div>
        </form>
    </Card>
</template>

<script setup>
const { formData } = defineProps({
    formData: Object
});

const emit = defineEmits(['submit']);

const handleSubmit = () => {
    emit('submit');
};
</script>