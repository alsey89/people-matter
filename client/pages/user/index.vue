<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full" />
        <!-- !Users Header -->
        <div class="min-h-max flex justify-between border-b-2 border-black py-2 pr-4">
            <h1 class="text-lg font-bold"> All Users </h1>
            <!-- !add user button -->
            <NBButtonSquare @click="handleAddUserButtonClick" size="xs">
                <Icon v-if="showUserForm" name="material-symbols:close" class="text-primary-dark h-6 w-6" />
                <Icon v-else name="material-symbols:add" class="h-6 w-6" />
            </NBButtonSquare>
        </div>
        <div v-auto-animate class="w-full flex flex-col gap-4">
            <!-- !New Company Form -->
            <AppUsersForm v-if="showUserForm" :userFormData="userFormData" @submit="handleUserFormSubmit" />
            <!-- !Company List -->
            <div v-if="userStore.getUserList && userStore.getUserList.length > 0" v-for="user in userStore.getUserList"
                :key="user.ID">
                <NBCard class="w-full flex justify-between">
                    <div class="">
                        {{ user }}
                    </div>
                    <button @click="handleEditUserButtonClick(user)"> EDIT </button>
                </NBCard>
            </div>
            <div v-else>
                <NBCard>
                    <div v-if="userStore.isLoading">
                        Loading...
                    </div>
                    <div v-else class="m-auto">
                        No Data
                    </div>
                </NBCard>
            </div>
        </div>
    </div>
</template>

<script setup>
//meta
definePageMeta({
    title: "Users",
    description: "Users",
    layout: "users",
});
//stores
const userStore = useUserStore()
const messageStore = useMessageStore()

// add/edit user
const showUserForm = ref(false)
const userFormData = reactive({
    userFormType: null,
    userId: null,
    email: '',
    firstName: '',
    middleName: '',
    lastName: '',
    isAdmin: false,
})

const handleAddUserButtonClick = () => {
    clearUserForm();
    userFormData.userFormType = "add"
    showUserForm.value = !showUserForm.value
};
const handleEditUserButtonClick = (user) => {
    populateUserForm(user)
    userFormData.userFormType = "edit"
    showUserForm.value = true
};
const clearUserForm = () => {
    userFormData.userFormType = null
    userFormData.userId = null
    userFormData.email = ''
    userFormData.firstName = ''
    userFormData.middleName = ''
    userFormData.lastName = ''
    userFormData.isEmployee = false
    userFormData.isAdmin = false
};
const populateUserForm = (user) => {
    userFormData.userFormType = "edit"
    userFormData.userId = user.ID
    userFormData.email = user.email
    userFormData.firstName = user.firstName
    userFormData.middleName = user.middleName
    userFormData.lastName = user.lastName
    userFormData.isEmployee = user.isEmployee
    userFormData.isAdmin = user.isAdmin
};
const handleUserFormSubmit = async () => {
    showUserForm.value = false;
    if (userFormData.userFormType === "edit") {
        await userStore.updateUser({ userFormData: userFormData });
    } else if (userFormData.userFormType === "add") {
        await userStore.createUser({ companyFormData: companyFormData });
    } else {
        console.error("No user form type")
        messageStore.setError("Error submitting user form")
        return
    }
};
// delete
const showConfirmationModal = ref(false)
const confirmationModalMessage = ref("")
const handleModalConfirmEvent = ref(null)

const handleModalCancelEvent = () => {
    showConfirmationModal.value = false;
    handleModalConfirmEvent.value = null; //! clear the stored function
};
const handleDeleteUserButtonClick = (user) => {
    const userId = user.ID
    if (!userId) {
        console.error("No user ID")
        messageStore.setError("Error deleting user")
        return
    }
    confirmationModalMessage.value = `Are you sure you want to delete ${user.firstName} ${user.lastName}? All data associated with this user will be deleted. This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showCompanyForm.value = false;
        await userStore.deleteUser({ userId: userId });
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
};
</script>

