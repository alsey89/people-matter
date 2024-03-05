<template>
    <Card class="w-full">
        <div class="flex flex-col gap-4 p-2">
            <div v-if="userFormData.userFormType == 'add'">
                <h1 class="font-bold underline"> Create New User </h1>
                <p>
                    This will create a new user account. An email will be sent to the designated email, asking the user
                    to
                    set
                    their password and complete their profile.
                </p>
            </div>
            <form @submit.prevent="handleSubmit" class="w-full flex flex-col gap-4">
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
                    </div>
                </div>
                <div class="w-full flex justify-end">
                    <Button type="submit" size="xs" textSize="md"
                        class="min-w-full items-center text-lg font-bold bg-primary hover:bg-primary-dark">
                        Submit
                    </Button>
                </div>
            </form>
        </div>
    </Card>
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