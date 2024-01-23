<template>
    <NBCard>
        <form @submit.prevent="handleSubmit" class="flex flex-col gap-4">
            <div class="w-full flex gap-2">
                <!-- Name Field -->
                <div class="w-1/2 flex flex-col gap-2">
                    <label for="name">Name*</label>
                    <input v-model="formData.locationName" type="text" id="name" class="focus:outline-primary p-1"
                        required />
                </div>

                <!-- Phone Field -->
                <div class="w-1/2 flex flex-col gap-2">
                    <label for="phone">Phone</label>
                    <input v-model="formData.locationPhone" type="text" id="phone" class="focus:outline-primary p-1" />
                </div>
            </div>

            <!-- Head Office Checkbox -->
            <div class="flex items-center gap-2">
                <input v-model="formData.locationIsHeadOffice" type="checkbox" id="isHeadOffice"
                    class="focus:outline-primary p-1" />
                <label for="headOffice"> this location is the head office </label>
            </div>

            <!-- Address Field -->
            <div class="flex flex-col gap-2">
                <label for="address">Address*</label>
                <input v-model="formData.locationAddress" type="text" id="address" class="focus:outline-primary p-1"
                    required />
            </div>

            <div class="w-full flex gap-2">
                <!-- City Field -->
                <div class="w-1/2 flex flex-col gap-2">
                    <label for="city">City*</label>
                    <input v-model="formData.locationCity" type="text" id="city" class="focus:outline-primary p-1"
                        required />
                </div>

                <!-- State Field -->
                <div class="w-1/2 flex flex-col gap-2">
                    <label for="state">State*</label>
                    <input v-model="formData.locationState" type="text" id="state" class="focus:outline-primary p-1" />
                </div>
            </div>

            <div class="w-full flex gap-2">
                <!-- Country Field -->
                <div class="w-1/2 flex flex-col gap-2">
                    <label for="country">Country*</label>
                    <input v-model="formData.locationCountry" type="text" id="country" class="focus:outline-primary p-1"
                        required />
                </div>

                <!-- Postal Code Field -->
                <div class="w-1/2 flex flex-col gap-2">
                    <label for="postalCode">Postal Code</label>
                    <input v-model="formData.locationPostalCode" type="text" id="postalCode"
                        class="focus:outline-primary p-1" />
                </div>
            </div>
            <!-- Submit Button -->
            <div class="w-full">
                <NBButtonSquare type="submit" size="xs" textSize="md"
                    class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">
                    Submit
                </NBButtonSquare>
            </div>
        </form>
    </NBCard>
</template>

<script setup>
const { formData } = defineProps({
    formData: Object
});

const emit = defineEmits(['submit']);

const handleSubmit = () => {
    emit('submit', formData);
};
</script>
