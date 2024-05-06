<template>
    <Card class="">
        <form @submit.prevent="submitForm" class="flex flex-col gap-4 p-4">
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
                <input name="password" type="password" v-model="form.password" placeholder="********" minlength="6"
                    required
                    class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                    :disabled="form.submitting">
            </div>

            <!-- confirm password -->
            <div v-auto-animate>
                <label for="confirmPassword">Confirm Password</label>
                <div v-if="form.confirmPasswordErr" class="text-xs text-destructive"> {{ form.confirmPasswordErr }}
                </div>
                <input name="confirmPassword" type="password" v-model="form.confirmPassword" placeholder="********"
                    minlength="6" required
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
</template>

<script setup>
import { reactive } from 'vue';

import { Icon } from '@iconify/vue';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

import { useUserStore } from '@/stores/user';

const userStore = useUserStore();
const form = reactive({
    email: '',
    emailErr: '',
    password: '',
    confirmPassword: '',
    passwordErr: '',
    confirmPasswordErr: '',
    submitting: false,
});

const submitForm = () => {
    form.submitting = true;
    validateEmail(form.email).catch((error) => {
        console.error(error);
        form.emailErr = error.message;
        return;
    });
    validatePassword(form.password, form.confirmPassword).catch((error) => {
        console.error(error);
        form.passwordErr = error.message;
        form.confirmPasswordErr = error.message;
        return;
    });

    console.log(form);

    try {
        userStore.signup(form);
    } catch (error) {
        console.error(error);
    } finally {
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

const validatePassword = (password, confirmPassword, minLength = 6) => {
    let error = '';

    switch (true) {
        case !password && !confirmPassword:
            error = 'Password required';
            break;
        case !password || !confirmPassword:
            error = 'Both password fields must be filled';
            break;
        case password.length < minLength || confirmPassword.length < minLength:
            error = 'Password must be at least ' + minLength + ' characters';
            break;
        case password !== confirmPassword:
            error = 'Passwords do not match';
            break;
        default:
            return true;
    }

    if (error) {
        throw new Error(error);
    }
};

onBeforeMount(() => {
    try {
        userStore.getCsrfToken();
    } catch (error) {
        console.error(error);
    }
});
</script>