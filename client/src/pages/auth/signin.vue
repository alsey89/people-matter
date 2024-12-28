<template>
    <div v-auto-animate class="w-full">
        <Card v-if="showCard" class="w-[400px] md:w-[480px]">
            <form @submit.prevent="submitForm" v-auto-animate class="flex flex-col gap-4 p-4">
                <!-- Error -->
                <section v-if="authStore.getError">
                    <p class="text-base font-bold text-destructive">Error: {{ authStore.getError }}</p>
                    <hr v-if="authStore.getError" class="border-t border-gray-300 my-4">
                </section>
                <!-- email -->
                <div v-auto-animate>
                    <label for="email">Email</label>
                    <div v-if="form.emailErr" class="text-xs text-destructive"> {{ form.emailErr }} </div>
                    <input name="email" type="email" v-model.trim="form.email" placeholder="email@domain.com" required
                        class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                        :disabled="form.submitting">
                </div>
                <!-- password -->
                <div v-auto-animate>
                    <label for="password">Password</label>
                    <div v-if="form.passwordErr" class="text-xs text-destructive"> {{ form.passwordErr }} </div>
                    <input name="password" type="password" v-model="form.password" placeholder="********" required
                        class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                        :disabled="form.submitting">
                </div>
                <!-- submit -->
                <Button v-auto-animate>
                    <Icon v-if="form.submitting" icon="svg-spinners:6-dots-rotate" class="h-6 w-6" />
                    <div v-else> Sign In </div>
                </Button>
            </form>
        </Card>
    </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import { onBeforeMount } from 'vue';
import { useRouter } from 'vue-router';

import { Icon } from '@iconify/vue';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

import { useAuthStore } from '@/stores/Auth';

const router = useRouter();
const authStore = useAuthStore();
const form = reactive({
    email: '',
    emailErr: '',
    password: '',
    passwordErr: '',
    submitting: false,
});
const showCard = ref(false);

const submitForm = async () => {
    form.submitting = true;

    try {
        validateEmail(form.email);
    } catch (error) {
        form.emailErr = error.message;
        form.submitting = false;
        return;
    }

    try {
        const success = await authStore.signin(form, router);
    } catch (error) {
        console.error(error);
        form.submitting = false;
    }
};

const validateEmail = (email) => {
    if (!email) {
        throw new Error('Email required');
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        throw new Error('Invalid email');
    }

    return true;
};

onMounted(() => {
    authStore.getCsrfToken();
    showCard.value = true;
});
</script>