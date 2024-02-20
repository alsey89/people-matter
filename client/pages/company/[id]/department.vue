<template>
    <div v-auto-animate class="w-full flex flex-col gap-4">
        <!-- !Confirmation Modal -->
        <AppConfirmationModal v-if="showConfirmationModal" :confirmationModalMessage="confirmationModalMessage"
            @confirm="handleModalConfirmEvent" @cancel="handleModalCancelEvent" class="w-full" />
        <!-- !Department Header -->
        <div class="min-h-max flex justify-between border-b-2 border-black py-2 pr-4">
            <div class=" items-center  text-lg font-bold"> Departments </div>
            <div v-auto-animate class="flex gap-4">
                <!-- show/hide details Button -->
                <UIButtonSquare size="xs" @click="handleToggleAllDepartmentsButtonClick">
                    <Icon v-if="expandedDepartmentIndex == 'all'" name="solar:list-arrow-up-bold"
                        class="text-primary-dark h-6 w-6" />
                    <Icon v-else name="solar:list-arrow-down-bold" class="h-6 w-6" />
                </UIButtonSquare>
                <!-- !add department button -->
                <UIButtonSquare @click="handleAddDepartmentButtonClick" size="xs">
                    <Icon v-if="showDepartmentForm" name="material-symbols:close" class="text-primary-dark h-6 w-6" />
                    <Icon v-else name="material-symbols:add" class="h-6 w-6" />
                </UIButtonSquare>
            </div>
        </div>
        <!-- !department form -->
        <div v-auto-animate>
            <AppCompanyDepartmentForm v-if="showDepartmentForm" @submit="handleDepartmentFormSubmit"
                :formData="departmentFormData" />
        </div>
        <!-- !department list -->
        <div v-if="companyStore.getCompanyDepartments?.length > 0"
            v-for="(department, index) in companyStore.companyDepartments" :key="department.ID">
            <UICard v-auto-animate>
                <UICardHeader>
                    <div v-auto-animate @click="handleExpandDepartmentButtonClick(index)"
                        class="flex justify-between items-center px-2 hover:cursor-pointer">
                        <p :class="{ 'text-primary text-xl': expandedDepartmentIndex === index }">
                            {{ department.name }}
                        </p>
                        <div class="flex gap-4">
                            <UIButtonSquare size="xs" @click.stop="handleEditDepartmentButtonClick(department)">
                                <Icon name="material-symbols:edit" class="h-6 w-6 hover:text-primary" />
                            </UIButtonSquare>
                            <UIButtonSquare size="xs" @click.stop="handleDeleteDepartmentButtonClick(department)">
                                <Icon name="material-symbols:delete" class="h-6 w-6 hover:text-primary" />
                            </UIButtonSquare>
                        </div>
                    </div>
                </UICardHeader>
                <div v-if="department.description && expandedDepartmentIndex == index || department.description && expandedDepartmentIndex == 'all'"
                    class="px-2">
                    <MDRender :content="department.description" />
                </div>
            </UICard>
        </div>
        <div v-else>
            <UICard>
                <div class="m-auto">
                    No Data
                </div>
            </UICard>
        </div>
    </div>
</template>

<script setup>
definePageMeta({
    title: "Departments",
    description: "Manage departments",
    layout: "company",
})

const companyStore = useCompanyStore()
const messageStore = useMessageStore()

// show/hide details
const expandedDepartmentIndex = ref("all")
const handleExpandDepartmentButtonClick = (index) => {
    if (expandedDepartmentIndex.value == index) {
        expandedDepartmentIndex.value = null
        return
    }
    expandedDepartmentIndex.value = index
}
const handleToggleAllDepartmentsButtonClick = () => {
    if (expandedDepartmentIndex.value == 'all') {
        expandedDepartmentIndex.value = null
        return
    }
    expandedDepartmentIndex.value = 'all'
}

// add/edit department form
const showDepartmentForm = ref(false)
const departmentFormData = reactive({
    departmentFormType: null,
    departmentId: null,
    departmentName: null,
    departmentDescription: null,
})
const handleDepartmentFormSubmit = async () => {
    showDepartmentForm.value = false;
    if (departmentFormData.departmentFormType == "edit") {
        await companyStore.updateDepartment({ companyId: companyStore.getCompanyId, departmentFormData: departmentFormData });
    } else if (departmentFormData.departmentFormType == "add") {
        await companyStore.createDepartment({ companyId: companyStore.getCompanyId, departmentFormData: departmentFormData });
    } else {
        console.error("No department form type")
        messageStore.setError("Error submitting department form")
        return
    }
};
const handleAddDepartmentButtonClick = () => {
    clearDepartmentForm()
    departmentFormData.departmentFormType = "add"
    showDepartmentForm.value = !showDepartmentForm.value
};
const handleEditDepartmentButtonClick = (department) => {
    departmentFormData.departmentFormType = "edit"
    populateDepartmentForm(department)
    showDepartmentForm.value = !showDepartmentForm.value
};
const clearDepartmentForm = () => {
    departmentFormData.departmentId = null
    departmentFormData.departmentName = ''
    departmentFormData.departmentDescription = ''
};
const populateDepartmentForm = (department) => {
    departmentFormData.departmentId = department.ID
    departmentFormData.departmentName = department.name
    departmentFormData.departmentDescription = department.description
};

// delete department
const showConfirmationModal = ref(false)
const confirmationModalMessage = ref("")
const handleModalConfirmEvent = ref(null) //! stored function to be called when confirmation modal is confirmed
const handleModalCancelEvent = () => {
    showConfirmationModal.value = false
};
const handleDeleteDepartmentButtonClick = (department) => {
    const departmentId = department.ID
    if (!departmentId) {
        console.error("No department ID")
        messageStore.setError("Error deleting department")
        return
    }

    confirmationModalMessage.value = `Are you sure you want to delete ${department.name}? This action cannot be undone.`
    showConfirmationModal.value = true

    //* store the function to be called when confirmation modal is confirmed, along with its arguments
    handleModalConfirmEvent.value = async () => {
        showConfirmationModal.value = false;
        showDepartmentForm.value = false;
        await companyStore.deleteDepartment({ companyId: companyStore.getCompanyId, departmentId: departmentId });
        handleModalConfirmEvent.value = null; //! clear the stored function
    };
}
</script>
