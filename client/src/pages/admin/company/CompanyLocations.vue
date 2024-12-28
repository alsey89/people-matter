<template>
    <div class="w-full h-full flex flex-col gap-2">
        <div class="h-8 flex justify-between items-center text-lg font-semibold px-4 py-2">
            <div> Locations </div>
            <button @click="toggleAddLocationCard"
                class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                <Icon v-if="!showAddLocationCard" icon="mdi:plus" class="text-xl" />
                <Icon v-else icon="mdi:close" class="text-xl" />
            </button>
        </div>
        <hr class="border-t mx-2">
        <div v-auto-animate class="flex flex-col gap-2 px-4">
            <Card v-if="showAddLocationCard" ref="addLocationCard" class="w-full p-4">
                <form @submit.prevent="onSubmitAddForm" class="flex flex-col gap-2">
                    <div class="flex flex-col">
                        <label for="addLocationName"> Name </label>
                        <input type="text" id="addLocationName" v-model.trim="addLocationForm.name"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex items-center gap-2">
                        <input type="checkbox" id="addIsHeadquarters" v-model="addLocationForm.isHeadOffice" />
                        <label for="addIsHeadquarters"> This location is the headquarters </label>
                    </div>
                    <div class="flex flex-col">
                        <label for="addLocationPhone"> Phone </label>
                        <input type="tel" id="addLocationPhone" v-model.trim="addLocationForm.phone"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addLocationAddress"> Address </label>
                        <input type="text" id="addLocationAddress" v-model.trim="addLocationForm.address"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addCompanyCity"> City </label>
                        <input type="text" id="addCompanyCity" v-model.trim="addLocationForm.city"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addCompanyState"> State </label>
                        <input type="text" id="addCompanyState" v-model.trim="addLocationForm.state"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addLocationCountry"> Country </label>
                        <input type="text" id="addLocationCountry" v-model.trim="addLocationForm.country"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="addLocationPostalCode"> Postal Code </label>
                        <input type="text" id="addLocationPostalCode" v-model.trim="addLocationForm.postalCode"
                            class="bg-background border px-2" />
                    </div>
                    <button type="submit"
                        class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground mt-2">
                        Submit
                    </button>
                </form>
            </Card>
        </div>
        <div v-for="location in companyStore.locations" :key="location.id">
            <div class="flex justify-between items-center px-4">
                <div class="h-8 flex items-center gap-2 text-lg font-semibold py-2">
                    {{ location.name }}
                    <Icon v-if="location.isHeadOffice" icon="mdi:home" class="text-primary text-xl border rounded-md" />
                </div>
                <button @click="toggleEditLocationCard(location)"
                    class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                    <Icon v-if="!(showEditLocationCard && editLocationForm.id === location.id)" icon="mdi:pencil"
                        class="text-xl" />
                    <Icon v-else icon="mdi:close" class="text-xl" />
                </button>
            </div>
            <div v-auto-animate class="flex flex-col gap-2 px-4">
                <Card v-if="showEditLocationCard && editLocationForm.id === location.id" ref="editCompanyCard"
                    class="w-full p-4">
                    <form @submit.prevent="onSubmitEditForm" class="flex flex-col gap-2">
                        <div class="flex flex-col">
                            <label for="editLocationName"> Name </label>
                            <input type="text" id="editLocationName" v-model.trim="editLocationForm.name"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex items-center gap-2">
                            <input type="checkbox" id="editIsHeadquarters" v-model="editLocationForm.isHeadOffice" />
                            <label for="editIsHeadquarters"> This location is the headquarters </label>
                        </div>
                        <div class="flex flex-col">
                            <label for="editLocationPhone"> Phone </label>
                            <input type="tel" id="editLocationPhone" v-model.trim="editLocationForm.phone"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editLocationAddress"> Address </label>
                            <input type="text" id="editLocationAddress" v-model.trim="editLocationForm.address"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editCompanyCity"> City </label>
                            <input type="text" id="editCompanyCity" v-model.trim="editLocationForm.city"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editCompanyState"> State </label>
                            <input type="text" id="editCompanyState" v-model.trim="editLocationForm.state"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editLocationCountry"> Country </label>
                            <input type="text" id="editLocationCountry" v-model.trim="editLocationForm.country"
                                class="bg-background border px-2" />
                        </div>
                        <div class="flex flex-col">
                            <label for="editLocationPostalCode"> Postal Code </label>
                            <input type="text" id="editLocationPostalCode" v-model.trim="editLocationForm.postalCode"
                                class="bg-background border px-2" />
                        </div>
                        <Button @click="onDeleteLocation(location)" type="button" variant="border-destructive"
                            class="mt-2"> Delete
                        </Button>
                        <Button type="submit" variant="border" class="mt-2"> Submit </Button>
                    </form>
                </Card>
                <Card class="w-full flex flex-col gap-4 p-4">
                    <div>
                        <p class="font-bold">Phone</p>
                        <p>{{ location.phone }}</p>
                    </div>
                    <div>
                        <p class="font-bold">Address</p>
                        <p>{{ formattedAddress(location) }}</p>
                    </div>
                </Card>
            </div>
        </div>
    </div>
    <!-- confirmation modal -->
    <div v-auto-animate v-if="showConfirmationDialog" @click="showConfirmationDialog = false"
        class="fixed ml-0 md:ml-[250px] inset-0 bg-black bg-opacity-50 flex justify-center items-center">
        <Card class="w-[90%] md:w-[25%] p-4">
            <p>Are you sure you want to delete this location?</p>
            <div class="flex justify-end gap-2 mt-4">
                <button @click.stop="showConfirmationDialog = false"
                    class="border rounded-md px-2 py-1 hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                    Cancel
                </button>
                <button @click.stop="onConfirmDeleteLocation"
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

const showEditLocationCard = ref(false);
const showAddLocationCard = ref(false);
const showConfirmationDialog = ref(false);

const editLocationForm = reactive({
    id: '',
    name: '',
    isHeadOffice: false,
    phone: '',
    address: '',
    city: '',
    state: '',
    country: '',
    postalCode: ''
});

const addLocationForm = reactive({
    name: '',
    isHeadOffice: false,
    phone: '',
    address: '',
    city: '',
    state: '',
    country: '',
    postalCode: ''
});

const toggleEditLocationCard = (location) => {
    if (showEditLocationCard.value && editLocationForm.id === location.id) {
        showEditLocationCard.value = false;
        cleanEditForm();
        return;
    }
    populateEditForm(location);
    showEditLocationCard.value = true;
};

const toggleAddLocationCard = () => {
    showAddLocationCard.value = !showAddLocationCard.value;
    cleanAddForm();
};

const cleanEditForm = () => {
    editLocationForm.id = '';
    editLocationForm.name = '';
    editLocationForm.isHeadOffice = false;
    editLocationForm.phone = '';
    editLocationForm.address = '';
    editLocationForm.city = '';
    editLocationForm.state = '';
    editLocationForm.country = '';
    editLocationForm.postalCode = '';
};

const cleanAddForm = () => {
    addLocationForm.name = '';
    addLocationForm.isHeadOffice = false;
    addLocationForm.phone = '';
    addLocationForm.address = '';
    addLocationForm.city = '';
    addLocationForm.state = '';
    addLocationForm.country = '';
    addLocationForm.postalCode = '';
};

const populateEditForm = (location) => {
    editLocationForm.id = location.id;
    editLocationForm.name = location.name;
    editLocationForm.isHeadOffice = location.isHeadOffice;
    editLocationForm.phone = location.phone;
    editLocationForm.address = location.address;
    editLocationForm.city = location.city;
    editLocationForm.state = location.state;
    editLocationForm.country = location.country;
    editLocationForm.postalCode = location.postalCode;
};

const onSubmitEditForm = async () => {
    await companyStore.updateLocation(editLocationForm);
    showEditLocationCard.value = false;
};

const onSubmitAddForm = async () => {
    await companyStore.addLocation(addLocationForm);
    showAddLocationCard.value = false;
};

const onDeleteLocation = async () => {
    showConfirmationDialog.value = true;
};

const onConfirmDeleteLocation = async () => {
    await companyStore.deleteLocation(editLocationForm.id);
    showEditLocationCard.value = false;
    showConfirmationDialog.value = false;
};

const formattedAddress = (location) => {
    return `${location.address}, ${location.city}, ${location.state}, ${location.country}, ${location.postalCode}`;
};

onMounted(() => {
    companyStore.getCompany();
});
</script>
