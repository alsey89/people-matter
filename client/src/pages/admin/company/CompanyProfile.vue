<template>
    <div class="w-full h-full flex flex-col gap-2">
        <div class="h-8 flex gap-2 text-lg font-semibold px-4 py-2">
            <div> Company Profile </div>
        </div>
        <hr class="border-t mx-2">
        <div class="flex justify-between items-center px-4">
            <div class="h-8 flex items-center gap-2 text-lg font-semibold py-2">
                <Avatar class="w-7 h-7">
                    <AvatarImage :src="companyStore.logoUrl" alt="logo" />
                    <AvatarFallback> Logo </AvatarFallback>
                </Avatar>
                <span>{{ companyStore.company.name }}</span>
            </div>
            <button @click="onClickEditCompanyButton"
                class="border rounded-md hover:cursor-pointer hover:bg-primary hover:text-primary-foreground">
                <Icon v-if="!showEditCompanyCard" icon="mdi:pencil" class="text-xl" />
                <Icon v-else icon="mdi:close" class="text-xl" />
            </button>
        </div>
        <div v-auto-animate class="flex flex-col gap-2 px-4">
            <Card v-if="showEditCompanyCard" ref="editCompanyCard" class="w-full p-4">
                <form @submit.prevent="onSubmitForm" class="flex flex-col gap-2">
                    <div class="flex flex-col">
                        <label for="companyName"> Name </label>
                        <input type="text" id="companyName" v-model.trim="editCompanyForm.name"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="companyEmail"> Email </label>
                        <input type="email" id="companyEmail" v-model.trim="editCompanyForm.email"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="companyPhone"> Phone </label>
                        <input type="tel" id="companyPhone" v-model.trim="editCompanyForm.phone"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="companyAddress"> Address </label>
                        <input type="text" id="companyAddress" v-model.trim="editCompanyForm.address"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="companyCity"> City </label>
                        <input type="text" id="companyCity" v-model.trim="editCompanyForm.city"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="companyState"> State </label>
                        <input type="text" id="companyState" v-model.trim="editCompanyForm.state"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="companyCountry"> Country </label>
                        <input type="text" id="companyCountry" v-model.trim="editCompanyForm.country"
                            class="bg-background border px-2" />
                    </div>
                    <div class="flex flex-col">
                        <label for="companyPostalCode"> Postal Code </label>
                        <input type="text" id="companyPostalCode" v-model.trim="editCompanyForm.postalCode"
                            class="bg-background border px-2" />
                    </div>
                    <Button type="submit" variant="border" class="mt-2"> Submit </Button>
                </form>
            </Card>
            <Card class="w-full flex flex-col gap-4 p-4">
                <div class="grid grid-cols-1 md:grid-cols-12 gap-4">
                    <div class="col-span-6">
                        <p class="font-bold">Email</p>
                        <p>{{ companyStore.company.email }}</p>
                    </div>
                    <div class="col-span-6">
                        <p class="font-bold">Phone</p>
                        <p>{{ companyStore.company.phone }}</p>
                    </div>
                </div>
                <div>
                    <p class="font-bold">Address</p>
                    <p>{{ companyStore.getAddress }}</p>
                </div>
            </Card>
        </div>
    </div>
</template>

<script setup>
import { onBeforeMount, onMounted, reactive, ref } from 'vue';
import { useCompanyStore } from '@/stores/Company';
import { Icon } from '@iconify/vue';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'

const companyStore = useCompanyStore();

const showEditCompanyCard = ref(false);
const editCompanyForm = reactive({
    name: '',
    email: '',
    phone: '',
    address: '',
    city: '',
    state: '',
    country: '',
    postalCode: ''
});
const onClickEditCompanyButton = () => {
    if (showEditCompanyCard.value) {
        showEditCompanyCard.value = false;
        resetForm();
        return;
    }
    populateForm();
    showEditCompanyCard.value = true;
};
const resetForm = () => {
    editCompanyForm.name = '';
    editCompanyForm.email = '';
    editCompanyForm.phone = '';
    editCompanyForm.address = '';
    editCompanyForm.city = '';
    editCompanyForm.state = '';
    editCompanyForm.country = '';
    editCompanyForm.postalCode = '';
};
const populateForm = () => {
    editCompanyForm.name = companyStore.company.name;
    editCompanyForm.email = companyStore.company.email;
    editCompanyForm.phone = companyStore.company.phone;
    editCompanyForm.address = companyStore.company.address;
    editCompanyForm.city = companyStore.company.city;
    editCompanyForm.state = companyStore.company.state;
    editCompanyForm.country = companyStore.company.country;
    editCompanyForm.postalCode = companyStore.company.postalCode;
};
const onSubmitForm = async () => {
    await companyStore.updateCompany(editCompanyForm);
    showEditCompanyCard.value = false;
}

onMounted(() => {
    companyStore.getCompany();
})
</script>