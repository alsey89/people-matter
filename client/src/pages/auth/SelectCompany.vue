<template>
    <div v-auto-animate class="w-full">
        <Card v-if="showCard" class="w-[400px] md:w-[480px] flex flex-col gap-2 p-4">
            <div>
                <h1 class="text-2xl font-bold">Select Company</h1>
                <p class="text-sm text-gray-600">Select a company to sign in</p>
            </div>
            <div v-for="userRole in authStore.userRoles" :key="userRole.company.id">
                <Button @click="onSelect(userRole.company.id)"
                    class="w-full flex justify-between items-center px-4 py-2">
                    <div>{{ userRole.company.name }}</div>
                    <div class="border-r border-red-500 h-full mx-2"></div>
                    <div>{{ userRole.role.name }}
                        <span v-if="userRole.role.name == 'root'"> user</span>
                    </div>
                </Button>
            </div>
        </Card>
    </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';

import { Icon } from '@iconify/vue';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

import { useAuthStore } from '@/stores/Auth';

const router = useRouter();
const authStore = useAuthStore();

const showCard = ref(false);

//get email from query params
const email = router.currentRoute.value.query.email;

const onSelect = async (companyId) => {
    if (!email) {
        console.error('Email not found');
        return;
    }
    const success = await authStore.getJwt(companyId, email, router);
};

onMounted(() => {
    authStore.getCsrfToken();
    showCard.value = true;
});
</script>