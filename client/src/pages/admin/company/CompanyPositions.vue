<template>
    <div class="w-full h-full flex flex-col gap-2">
        <div class="h-8 flex justify-between items-center text-lg font-semibold px-4 py-2">
            <div>Positions</div>
            <button @click="toggleAddPositionCard"
                class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                <Icon v-if="!showAddPositionCard" icon="mdi:plus" class="text-xl" />
                <Icon v-else icon="mdi:close" class="text-xl" />
            </button>
        </div>
        <hr class="border-t mx-2">
        <div v-auto-animate class="flex flex-col gap-2 px-4">
            <Card v-if="showAddPositionCard" ref="addPositionCard" class="w-full p-4">
                <form @submit.prevent="onSubmitAddForm" class="flex flex-col gap-2">
                    <div class="flex flex-col">
                        <label for="addPositionName">Name</label>
                        <input type="text" id="addPositionName" v-model.trim="addPositionForm.name"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addPositionDepartment">Department</label>
                        <select id="addPositionDepartment" v-model.number="addPositionForm.departmentId"
                            class="bg-background border px-2">
                            <option v-for="department in companyStore.company.departments" :key="department.id"
                                :value="department.id">{{ department.name }}</option>
                        </select>
                    </div>
                    <div class="flex flex-col">
                        <label for="addPositionDescription">Description</label>
                        <textarea id="addPositionDescription" v-model.trim="addPositionForm.description"
                            class="bg-background border px-2"></textarea>
                    </div>
                    <div class="flex flex-col">
                        <label for="addPositionDuties">Duties</label>
                        <textarea id="addPositionDuties" v-model.trim="addPositionForm.duties"
                            class="bg-background border px-2"></textarea>
                    </div>
                    <div class="flex flex-col">
                        <label for="addPositionQualifications">Qualifications</label>
                        <textarea id="addPositionQualifications" v-model.trim="addPositionForm.qualifications"
                            class="bg-background border px-2"></textarea>
                    </div>
                    <div class="flex flex-col">
                        <label for="addPositionExperience">Experience</label>
                        <textarea id="addPositionExperience" v-model.trim="addPositionForm.experience"
                            class="bg-background border px-2"></textarea>
                    </div>
                    <div class="flex flex-col">
                        <label for="addPositionMinSalary">Min Salary</label>
                        <input type="number" id="addPositionMinSalary" v-model.number="addPositionForm.minSalary"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addPositionMaxSalary">Max Salary</label>
                        <input type="number" id="addPositionMaxSalary" v-model.number="addPositionForm.maxSalary"
                            class="bg-background border px-2" />
                    </div>
                    <Button type="submit" variant="border" class="mt-2"> Submit </Button>
                </form>
            </Card>
        </div>
        <div v-for="position in companyStore.positions" :key="position.id">
            <div class="flex justify-between items-center px-4">
                <div class="h-8 flex items-center gap-2 text-lg font-semibold py-2">
                    {{ position.name }}
                </div>
                <button @click="toggleEditPositionCard(position)"
                    class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                    <Icon v-if="!(showEditPositionCard && editPositionForm.id === position.id)" icon="mdi:pencil"
                        class="text-xl" />
                    <Icon v-else icon="mdi:close" class="text-xl" />
                </button>
            </div>
            <div v-auto-animate class="flex flex-col gap-2 px-4">
                <Card v-if="showEditPositionCard && editPositionForm.id === position.id" ref="editPositionCard"
                    class="w-full p-4">
                    <form @submit.prevent="onSubmitEditForm" class="flex flex-col gap-2">
                        <div class="flex flex-col">
                            <label for="editPositionName">Name</label>
                            <input type="text" id="editPositionName" v-model.trim="editPositionForm.name"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editPositionDepartment">Department</label>
                            <select id="editPositionDepartment" v-model.number="editPositionForm.departmentId"
                                class="bg-background border px-2">
                                <option v-for="department in companyStore.company.departments" :key="department.id"
                                    :value="department.id">{{ department.name }}</option>
                                <option :value=null class="text-destructive font-bold"> Unassign Department </option>
                            </select>
                        </div>
                        <div class="flex flex-col">
                            <label for="editPositionDescription">Description</label>
                            <textarea id="editPositionDescription" v-model.trim="editPositionForm.description"
                                class="bg-background border px-2"></textarea>
                        </div>
                        <div class="flex flex-col">
                            <label for="editPositionDuties">Duties</label>
                            <textarea id="editPositionDuties" v-model.trim="editPositionForm.duties"
                                class="bg-background border px-2"></textarea>
                        </div>
                        <div class="flex flex-col">
                            <label for="editPositionQualifications">Qualifications</label>
                            <textarea id="editPositionQualifications" v-model.trim="editPositionForm.qualifications"
                                class="bg-background border px-2"></textarea>
                        </div>
                        <div class="flex flex-col">
                            <label for="editPositionExperience">Experience</label>
                            <textarea id="editPositionExperience" v-model.trim="editPositionForm.experience"
                                class="bg-background border px-2"></textarea>
                        </div>
                        <div class="flex flex-col">
                            <label for="editPositionMinSalary">Min Salary</label>
                            <input type="number" id="editPositionMinSalary" v-model.number="editPositionForm.minSalary"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editPositionMaxSalary">Max Salary</label>
                            <input type="number" id="editPositionMaxSalary" v-model.number="editPositionForm.maxSalary"
                                class="bg-background border px-2" />
                        </div>
                        <Button @click="onDeletePosition(position)" type="button" variant="border-destructive"
                            class="mt-2"> Delete
                        </Button>
                        <Button type="submit" variant="border" class="mt-2"> Submit </Button>
                    </form>
                </Card>
                <Card class="w-full flex flex-col gap-4 p-4">
                    <div>
                        <p class="font-bold">Department</p>
                        <p>{{ getDepartmentName(position.departmentId) }}</p>
                    </div>
                    <div>
                        <p class="font-bold">Description</p>
                        <p>{{ position.description }}</p>
                    </div>
                    <div>
                        <p class="font-bold">Duties</p>
                        <p>{{ position.duties }}</p>
                    </div>
                    <div>
                        <p class="font-bold">Qualifications</p>
                        <p>{{ position.qualifications }}</p>
                    </div>
                    <div>
                        <p class="font-bold">Experience</p>
                        <p>{{ position.experience }}</p>
                    </div>
                    <div>
                        <p class="font-bold">Salary Range</p>
                        <p>{{ position.minSalary }} - {{ position.maxSalary }}</p>
                    </div>
                </Card>
            </div>
        </div>
    </div>
    <!-- confirmation modal -->
    <div v-if="showConfirmationDialog" @click="showConfirmationDialog = false"
        class="fixed ml-0 md:ml-[250px] inset-0 bg-black bg-opacity-50 flex justify-center items-center">
        <Card class="w-[90%] md:w-[25%] p-4">
            <p>Are you sure you want to delete this position?</p>
            <div class="flex justify-end gap-2 mt-4">
                <button @click.stop="showConfirmationDialog = false"
                    class="border rounded-md px-2 py-1 hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                    Cancel
                </button>
                <button @click.stop="onConfirmDeletePosition"
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

const showEditPositionCard = ref(false);
const showAddPositionCard = ref(false);
const showConfirmationDialog = ref(false);

const editPositionForm = reactive({
    id: '',
    departmentId: '',
    name: '',
    description: '',
    duties: '',
    qualifications: '',
    experience: '',
    minSalary: '',
    maxSalary: ''
});

const addPositionForm = reactive({
    departmentId: '',
    name: '',
    description: '',
    duties: '',
    qualifications: '',
    experience: '',
    minSalary: '',
    maxSalary: ''
});

const toggleEditPositionCard = (position) => {
    if (showEditPositionCard.value && editPositionForm.id === position.id) {
        showEditPositionCard.value = false;
        cleanEditForm();
        return;
    }
    populateEditForm(position);
    showEditPositionCard.value = true;
};

const toggleAddPositionCard = () => {
    showAddPositionCard.value = !showAddPositionCard.value;
    cleanAddForm();
};

const cleanEditForm = () => {
    editPositionForm.id = '';
    editPositionForm.departmentId = '';
    editPositionForm.name = '';
    editPositionForm.description = '';
    editPositionForm.duties = '';
    editPositionForm.qualifications = '';
    editPositionForm.experience = '';
    editPositionForm.minSalary = '';
    editPositionForm.maxSalary = '';
};

const cleanAddForm = () => {
    addPositionForm.departmentId = '';
    addPositionForm.name = '';
    addPositionForm.description = '';
    addPositionForm.duties = '';
    addPositionForm.qualifications = '';
    addPositionForm.experience = '';
    addPositionForm.minSalary = '';
    addPositionForm.maxSalary = '';
};

const populateEditForm = (position) => {
    editPositionForm.id = position.id;
    editPositionForm.departmentId = position.departmentId;
    editPositionForm.name = position.name;
    editPositionForm.description = position.description;
    editPositionForm.duties = position.duties;
    editPositionForm.qualifications = position.qualifications;
    editPositionForm.experience = position.experience;
    editPositionForm.minSalary = position.minSalary;
    editPositionForm.maxSalary = position.maxSalary;
};

const onSubmitEditForm = async () => {
    await companyStore.updatePosition(editPositionForm);
    showEditPositionCard.value = false;
};

const onSubmitAddForm = async () => {
    await companyStore.addPosition(addPositionForm);
    showAddPositionCard.value = false;
};

const onDeletePosition = async () => {
    showConfirmationDialog.value = true;
};

const onConfirmDeletePosition = async () => {
    await companyStore.deletePosition(editPositionForm.id);
    showEditPositionCard.value = false;
    showConfirmationDialog.value = false;
};

const getDepartmentName = (departmentId) => {
    const department = companyStore.company.departments.find(dept => dept.id === departmentId);
    return department ? department.name : 'Unassigned';
};

onMounted(() => {
    companyStore.getCompany();
});
</script>
