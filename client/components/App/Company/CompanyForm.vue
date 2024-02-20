<template>
    <UICard class="w-full">
        <form @submit.prevent="handleSubmit" class="w-full flex flex-col gap-4 p-2">
            <div class="flex flex-wrap md:flex-nowrap justify-between gap-4">
                <div class="w-full md:w-3/4 flex-flex-col gap-4">
                    <div class="flex flex-col">
                        <label for="name">Name*</label>
                        <input v-model="formData.companyName" type="text" name="name" id="name"
                            class="focus:outline-primary p-1" required />
                    </div>
                    <div class="flex flex-col">
                        <label for="phone">Phone</label>
                        <input v-model="formData.companyPhone" type="text" name="phone" id="phone"
                            class="focus:outline-primary p-1" />
                    </div>
                    <div class="flex flex-col">
                        <label for="email">Email</label>
                        <input v-model="formData.companyEmail" type="email" name="email" id="email"
                            class="focus:outline-primary p-1" />
                    </div>
                    <div class="flex flex-col">
                        <label for="website"> Website </label>
                        <input v-model="formData.companyWebsite" type="text" name="name" id="website"
                            class="focus:outline-primary p-1" />
                    </div>
                </div>
                <div class="w-full md:w-1/4 flex justify-center">
                    <!-- todo: logo upload -->
                    <UIButtonSquare type="button" size="xs" textsize="xs" class="w-full block md:hidden">
                        <div class="flex justify-center items-center gap-2">
                            <Icon name="material-symbols:upload" class="h-6 w-6" />
                            Upload Photo
                        </div>
                    </UIButtonSquare>
                    <AppLogo :src="`${formData.companyLogoUrl}`" shape="square"
                        class="hidden md:block w-full border-2 border-black hover:cursor-pointer" />
                </div>
            </div>
            <div class="flex flex-col">
                <label for="address">Address</label>
                <input v-model="formData.companyAddress" type="text" name="address" id="address"
                    class="focus:outline-primary p-1" />
            </div>
            <div class="flex flex-wrap md:flex-nowrap gap-2 md:gap-4">
                <div class="w-full md:w-1/2 flex flex-col gap-4">
                    <div class="flex flex-col">
                        <label for="city">City</label>
                        <input v-model="formData.companyCity" type="text" name="city" id="city"
                            class="focus:outline-primary p-1" />
                    </div>
                    <div class="flex flex-col">
                        <label for="state">State</label>
                        <input v-model="formData.companyState" type="text" name="state" id="state"
                            class="focus:outline-primary p-1" />
                    </div>
                </div>
                <div class="w-full md:w-1/2 flex flex-col gap-4">
                    <div class="flex flex-col">
                        <label for="country">Country</label>
                        <input v-model="formData.companyCountry" type="text" name="country" id="country"
                            class="focus:outline-primary p-1" />
                    </div>
                    <div class="flex flex-col">
                        <label for="postal_code">Postal Code</label>
                        <input v-model="formData.companyPostalCode" type="text" name="postal_code" id="postal_code"
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
