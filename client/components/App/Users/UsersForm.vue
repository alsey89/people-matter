<template>
    <NBCard class="w-full">
        <form @submit.prevent="handleSubmit" class="w-full flex flex-col gap-4 p-2">
            <div class="flex flex-wrap justify-between gap-4">
                <div class="w-full flex flex-col gap-4">
                    <div class="flex flex-col">
                        <label for="email">Email*</label>
                        <input v-model="userFormData.email" type="email" name="email" id="email"
                            class="focus:outline-primary p-1" required />
                    </div>
                    <div class="flex flex-col">
                        <label for="firstName">First Name*</label>
                        <input v-model="userFormData.firstName" type="text" name="firstName" id="firstName"
                            class="focus:outline-primary p-1" required />
                    </div>
                    <div class="flex flex-col">
                        <label for="middleName">Middle Name</label>
                        <input v-model="userFormData.middleName" type="text" name="middleName" id="middleName"
                            class="focus:outline-primary p-1" />
                    </div>
                    <div class="flex flex-col">
                        <label for="lastName">Last Name*</label>
                        <input v-model="userFormData.lastName" type="text" name="lastName" id="lastName"
                            class="focus:outline-primary p-1" required />
                    </div>
                    <div class="flex gap-2 items-center">
                        <NBButtonSquare class="max-w-max" size="xs" id="isAdmin" name="isAdmin" type="button"
                            @click="toggleIsAdmin">
                            <Icon name="material-symbols:check" class="h-6 w-6"
                                :class="userFormData.isAdmin === true ? 'bg-primary text-white' : 'text-black hover:text-primary'" />
                        </NBButtonSquare>
                        <label for="isAdmin"> Administrator </label>
                    </div>
                </div>
            </div>
            <div class="w-full flex justify-end">
                <NBButtonSquare type="submit" size="xs" textSize="md"
                    class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">
                    Submit
                </NBButtonSquare>
            </div>
        </form>
    </NBCard>
</template>

<script setup>
const { userFormData } = defineProps({
    userFormData: Object
});

const toggleIsAdmin = () => {
    userFormData.isAdmin = !userFormData.isAdmin;
};

const emit = defineEmits(['submit']);

const handleSubmit = () => {
    emit('submit');
};
</script>