<template>
    <UICard class="w-full">
        <form @submit.prevent="handleSubmit" class="w-full flex flex-col gap-4 p-2">
            <div class="flex flex-wrap md:flex-nowrap justify-between gap-4">
                <div class="w-full md:w-1/2 flex flex-col">
                    <label for="name">Name*</label>
                    <input v-model="formData.locationName" type="text" name="name" id="name"
                        class="focus:outline-primary p-1" required />
                </div>
                <div class="w-full md:w-1/2 flex flex-col">
                    <label for="phone">Phone</label>
                    <input v-model="formData.locationPhone" type="text" name="phone" id="phone"
                        class="focus:outline-primary p-1" />
                </div>
            </div>
            <!-- Head Office Checkbox -->
            <div class="flex items-center gap-2">
                <input v-model="formData.locationIsHeadOffice" type="checkbox" id="isHeadOffice"
                    class="focus:outline-primary p-1" />
                <label for="headOffice"> this location is the head office </label>
            </div>

            <div class="flex flex-col">
                <label for="address">Address</label>
                <input v-model="formData.locationAddress" type="text" name="address" id="address"
                    class="focus:outline-primary p-1" />
            </div>
            <div class="flex flex-wrap md:flex-nowrap gap-2 md:gap-4">
                <div class="w-full md:w-1/2 flex flex-col gap-4">
                    <div class="flex flex-col">
                        <label for="city">City</label>
                        <input v-model="formData.locationCity" type="text" name="city" id="city"
                            class="focus:outline-primary p-1" />
                    </div>
                    <div class="flex flex-col">
                        <label for="state">State</label>
                        <input v-model="formData.locationState" type="text" name="state" id="state"
                            class="focus:outline-primary p-1" />
                    </div>
                </div>
                <div class="w-full md:w-1/2 flex flex-col gap-4">
                    <div class="flex flex-col">
                        <label for="country">Country</label>
                        <input v-model="formData.locationCountry" type="text" name="country" id="country"
                            class="focus:outline-primary p-1" />
                    </div>
                    <div class="flex flex-col">
                        <label for="postal_code">Postal Code</label>
                        <input v-model="formData.locationPostalCode" type="text" name="postal_code" id="postal_code"
                            class="focus:outline-primary p-1" />
                    </div>
                </div>
            </div>
            <div class="w-full flex justify-end">
                <UIButtonSquare type="submit" size="xs" textSize="md"
                    class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">
                    Submit
                </UIButtonSquare>
            </div>
        </form>
    </UICard>
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

