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
import api from "@/plugins/axios";  // Import your custom axios instance

const router = useRouter();
const query = computed(() => router.currentRoute.value.query);
const confirmed = ref(false);

onMounted(async () => {
    if (!query.value.token) {
        console.error('No token provided');
        router.push({ name: 'onboarding' });
        return;
    }
    try {
        const tokenParam = encodeURIComponent(query.value.token);
        const response = await api.get(`/auth/confirmation?token=${tokenParam}`);  // Use the custom axios instance
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