<template>
    <Card v-auto-animate>
        <div v-if="confirmed" class="flex text-base font-bold p-4">
            Your account has been confirmed. You can now sign in.
        </div>
        <div v-else class="flex gap-2 text-base font-bold p-4">
            <Icon icon="svg-spinners:blocks-scale" class="h-6 w-6 text-primary" />
            Please hold on while we confirm your account...
        </div>
    </Card>
</template>

<script setup>

import { useRouter } from 'vue-router';
import { computed, onMounted, ref } from 'vue';

import { Card } from '@/components/ui/card';
import { Icon } from '@iconify/vue';

import { redirectAfterDelay } from "@/util.js"

import axios from 'axios';

//get query parameters from router
const router = useRouter();
const query = computed(() => router.currentRoute.value.query);

const api_url = import.meta.env.VITE_API_URL;

const confirmed = ref(false);

onMounted(async () => {
    if (!query.value.token) {
        console.error('No token provided');
        router.push({ name: 'onboarding' });
        return;
    }
    try {
        axios.defaults.headers.post["Content-Type"] = "application/json";
        axios.defaults.headers.post["Accept"] = "application/json";
        axios.defaults.withCredentials = true;
        const tokenParam = encodeURIComponent(query.value.token);
        const response = await axios.get(`${api_url}/auth/confirmation?token=${tokenParam}`);
        if (response.status === 200) {
            setTimeout(() => {
                confirmed.value = true;
            }, 3000);
            redirectAfterDelay(router, "/auth/signin", 6000);
        } else {
            console.error('Failed to confirm account');
            setTimeout(() => {
                router.push("/onboarding");
            }, 3000);
        }
    } catch (error) {
        console.error(error);
        setTimeout(() => {
            router.push("/onboarding");
        }, 3000);
    }
});

</script>