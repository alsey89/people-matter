<template>
    <div v-auto-animate class="w-full">
        <!-- Company Created -->
        <Card v-if="!form.submitting && companyCreated" class="w-[450px] flex flex-col gap-2 p-4">
            <p class="text-xl font-bold"> Company Created! </p>
            <p> A confirmation email has been sent to your admin email. Please click on the confirmation link to begin
                setting up the HRMS. </p>
            <Button @click="this.$router.push('/auth/signin')" class="w-full"> Take me back to the login screen!
            </Button>
        </Card>

        <!-- Company Creation -->
        <Card v-else class="w-[400px] md:w-[480px]">

            <form v-auto-animate @submit.prevent="submitForm" class="flex flex-col gap-4 p-4">
                <!-- Error -->
                <section v-if="authStore.getError">
                    <p class="text-base font-bold text-destructive">Error: {{ authStore.getError }}</p>
                    <hr v-if="authStore.getError" class="border-t border-gray-300 my-4">
                </section>
                <!-- Company Information -->
                <section class="flex flex-col gap-2">
                    <p>
                        First, tell us a bit about your company.
                    </p>
                    <!-- company name -->
                    <div v-auto-animate>
                        <label for="companyName">Company </label>
                        <div v-if="form.companyNameErr" class="text-xs text-destructive"> {{ form.companyNameErr }}
                        </div>
                        <input name="companyName" type="text" v-model.trim="form.companyName" placeholder="Company Name"
                            required
                            class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                            :disabled="form.submitting">
                    </div>

                    <!-- company size -->
                    <div v-auto-animate>
                        <label for="companySize">Company Size</label>
                        <select name="companySize" v-model="form.companySize" required
                            class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                            :disabled="form.submitting">
                            <option value="" disabled selected>Select Company Size</option>
                            <option value="1-10">1-10</option>
                            <option value="11-25">11-25</option>
                            <option value="26-50">26-50</option>
                            <option value="51-200">51-200</option>
                            <option value=">200"> >200</option>
                        </select>
                    </div>
                </section>

                <!-- divider -->
                <hr v-if="form.companyName && form.companySize && !form.companyNameErr"
                    class="border-t border-gray-300 my-4">

                <!-- Admin Information -->
                <section v-if="form.companyName && form.companySize && !form.companyNameErr"
                    class="flex flex-col gap-2">
                    <p>
                        Great, now let's set up the admin account.
                    </p>
                    <!-- email -->
                    <div v-auto-animate>
                        <label for="adminEmail">Email</label>
                        <div v-if="form.adminEmailErr" class="text-xs text-destructive"> {{ form.adminEmailErr }} </div>
                        <input name="adminEmail" type="email" v-model.trim="form.adminEmail"
                            placeholder="email@domain.com" required
                            class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                            :disabled="form.submitting">
                    </div>

                    <!-- password -->
                    <div v-auto-animate>
                        <label for="password">Password</label>
                        <div v-if="form.passwordErr" class="text-xs text-destructive"> {{ form.passwordErr }} </div>
                        <input name="password" type="password" v-model="form.password" placeholder="********"
                            minlength="6" required
                            class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                            :disabled="form.submitting">
                    </div>

                    <!-- confirm password -->
                    <div v-auto-animate>
                        <label for="confirmPassword">Confirm Password</label>
                        <div v-if="form.confirmPasswordErr" class="text-xs text-destructive"> {{ form.confirmPasswordErr
                            }}
                        </div>
                        <input name="confirmPassword" type="password" v-model="form.confirmPassword"
                            placeholder="********" minlength="6" required
                            class="appearance-none border-2 border-secondary w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:border-ring"
                            :disabled="form.submitting">
                    </div>
                </section>

                <!-- divider -->
                <hr v-if="formIsValid" class="border-t border-gray-300 my-4">

                <!-- submit -->
                <section v-auto-animate v-if="formIsValid" class="flex flex-col gap-2">
                    <p>
                        All set, let's do this!
                    </p>
                    <Button v-auto-animate>
                        <Icon v-if="form.submitting" icon="svg-spinners:6-dots-rotate" class="h-6 w-6" />
                        <div v-else> Sign In </div>
                    </Button>
                </section>
            </form>
        </Card>
    </div>
</template>

<script setup>
import { ref, reactive, onBeforeMount } from 'vue';
import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { Icon } from '@iconify/vue';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

import { useAuthStore } from '@/stores/Auth';
const router = useRouter();
const authStore = useAuthStore();
const companyCreated = ref(false);
const form = reactive({
    companyName: '',
    companyNameErr: '',
    companySize: '',

    adminEmail: '',
    adminEmailErr: '',
    password: '',
    confirmPassword: '',
    passwordErr: '',
    confirmPasswordErr: '',

    submitting: false,
});

const formIsValid = computed(() => {
    return form.companyName && !form.companyNameErr &&
        form.companySize &&
        form.adminEmail && !form.adminEmailErr &&
        form.password && !form.passwordErr &&
        form.confirmPassword && !form.confirmPasswordErr &&
        form.password === form.confirmPassword;
});

const submitForm = async () => {
    form.submitting = true;

    try {
        validateEmail(form.adminEmail);
    } catch (error) {
        console.error(error);
        form.emailErr = error.message;
        form.submitting = false;
        return;
    }

    try {
        validatePassword(form.password, form.confirmPassword);
    } catch (error) {
        console.error(error);
        form.passwordErr = error.message;
        form.confirmPasswordErr = error.message;
        form.submitting = false;
        return;
    }

    try {
        const success = await authStore.createCompany(form, router);
        if (success) {
            companyCreated.value = true;
            form.submitting = false;
        }
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
        authStore.getCsrfToken();
    } catch (error) {
        console.error(error);
    }
});
</script>