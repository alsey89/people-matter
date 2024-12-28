<template>
    <div class="w-full h-full flex flex-col gap-2">
        <div class="h-8 flex justify-between items-center text-lg font-semibold px-4 py-2">
            <div> Departments </div>
            <button @click="toggleAddDepartmentCard"
                class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                <Icon v-if="!showAddDepartmentCard" icon="mdi:plus" class="text-xl" />
                <Icon v-else icon="mdi:close" class="text-xl" />
            </button>
        </div>
        <hr class="border-t mx-2">
        <div v-auto-animate class="flex flex-col gap-2 px-4">
            <Card v-if="showAddDepartmentCard" ref="addDepartmentCard" class="w-full p-4">
                <form @submit.prevent="onSubmitAddForm" class="flex flex-col gap-2">
                    <div class="flex flex-col">
                        <label for="addDepartmentName"> Name </label>
                        <input type="text" id="addDepartmentName" v-model.trim="addDepartmentForm.name"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addDepartmentDescription"> Description </label>
                        <textarea id="addDepartmentDescription" v-model.trim="addDepartmentForm.description"
                            class="bg-background border px-2"></textarea>
                    </div>
                    <Button type="submit" variant="border" class="mt-2"> Submit </Button>
                </form>
            </Card>
        </div>
        <div v-for="department in companyStore.departments" :key="department.id">
            <div class="flex justify-between items-center px-4">
                <div class="h-8 flex items-center gap-2 text-lg font-semibold py-2">
                    {{ department.name }}
                </div>
                <button @click="toggleEditDepartmentCard(department)"
                    class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                    <Icon v-if="!(showEditDepartmentCard && editDepartmentForm.id === department.id)" icon="mdi:pencil"
                        class="text-xl" />
                    <Icon v-else icon="mdi:close" class="text-xl" />
                </button>
            </div>
            <div v-auto-animate class="flex flex-col gap-2 px-4">
                <Card v-if="showEditDepartmentCard && editDepartmentForm.id === department.id" ref="editDepartmentCard"
                    class="w-full p-4">
                    <form @submit.prevent="onSubmitEditForm" class="flex flex-col gap-2">
                        <div class="flex flex-col">
                            <label for="editDepartmentName"> Name </label>
                            <input type="text" id="editDepartmentName" v-model.trim="editDepartmentForm.name"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editDepartmentDescription"> Description </label>
                            <textarea id="editDepartmentDescription" v-model.trim="editDepartmentForm.description"
                                class="bg-background border px-2"></textarea>
                        </div>
                        <Button @click="onDeleteDepartment(department)" type="button" variant="border-destructive"
                            class="mt-2"> Delete
                        </Button>
                        <Button type="submit" variant="border" class="mt-2"> Submit </Button>
                    </form>
                </Card>
                <Card class="w-full flex flex-col gap-4 p-4">
                    <div>
                        <p class="font-bold">Description</p>
                        <p>{{ department.description }}</p>
                    </div>
                </Card>
            </div>
        </div>
    </div>
    <!-- confirmation modal -->
    <div v-auto-animate v-if="showConfirmationDialog" @click="showConfirmationDialog = false"
        class="fixed ml-0 md:ml-[250px] inset-0 bg-black bg-opacity-50 flex justify-center items-center">
        <Card class="w-[90%] md:w-[25%] p-4">
            <p>Are you sure you want to delete this department?</p>
            <div class="flex justify-end gap-2 mt-4">
                <button @click.stop="showConfirmationDialog = false"
                    class="border rounded-md px-2 py-1 hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                    Cancel
                </button>
                <button @click.stop="onConfirmDeleteDepartment"
                    class="border rounded-md px-2 py-1 hover:cursor-pointer hover:bg-destructive hover:text-destructive-foreground">
                    Delete
                </button>
            </div>
        </Card>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useCompanyStore } from '@/stores/Company';
import { Icon } from '@iconify/vue';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

const companyStore = useCompanyStore();

const showEditDepartmentCard = ref(false);
const showAddDepartmentCard = ref(false);
const showConfirmationDialog = ref(false);

const editDepartmentForm = reactive({
    id: '',
    name: '',
    description: ''
});

const addDepartmentForm = reactive({
    name: '',
    description: ''
});

const toggleEditDepartmentCard = (department) => {
    if (showEditDepartmentCard.value && editDepartmentForm.id === department.id) {
        showEditDepartmentCard.value = false;
        cleanEditForm();
        return;
    }
    populateEditForm(department);
    showEditDepartmentCard.value = true;
};

const toggleAddDepartmentCard = () => {
    showAddDepartmentCard.value = !showAddDepartmentCard.value;
    cleanAddForm();
};

const cleanEditForm = () => {
    editDepartmentForm.id = '';
    editDepartmentForm.name = '';
    editDepartmentForm.description = '';
};

const cleanAddForm = () => {
    addDepartmentForm.name = '';
    addDepartmentForm.description = '';
};

const populateEditForm = (department) => {
    editDepartmentForm.id = department.id;
    editDepartmentForm.name = department.name;
    editDepartmentForm.description = department.description;
};

const onSubmitEditForm = async () => {
    await companyStore.updateDepartment(editDepartmentForm);
    showEditDepartmentCard.value = false;
};

const onSubmitAddForm = async () => {
    await companyStore.addDepartment(addDepartmentForm);
    showAddDepartmentCard.value = false;
};

const onDeleteDepartment = async () => {
    showConfirmationDialog.value = true;
};

const onConfirmDeleteDepartment = async () => {
    await companyStore.deleteDepartment(editDepartmentForm.id);
    showEditDepartmentCard.value = false;
    showConfirmationDialog.value = false;
};

onMounted(() => {
    companyStore.getCompany();
});
</script>
