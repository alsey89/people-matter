<template>
  <div class="flex justify-center items-center min-h-screen">
    <div class="w-full max-w-md px-6 py-12 space-y-8 bg-white rounded-lg shadow-xl dark:bg-gray-800">
      <form @submit.prevent="submitForm" class="space-y-6">

        <!-- Email Input -->
        <div>
          <label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your
            email</label>
          <input type="email" v-model="email" id="email"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white"
            required>
          <span v-if="emailError" class="text-red-500 text-sm mt-1">{{ emailError }}</span>
        </div>

        <!-- Password Input -->
        <div>
          <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your
            password</label>
          <input type="password" v-model="password" id="password"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white"
            required>
        </div>

        <!-- Confirm Password Input -->
        <div>
          <label for="confirmPassword" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Confirm your
            password</label>
          <input type="password" v-model="confirmPassword" id="confirmPassword"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white"
            required>
          <span v-if="passwordError" class="text-red-500 text-sm mt-1">{{ passwordError }}</span>
        </div>

        <!-- Submit Button -->
        <button type="submit"
          class="w-full px-5 py-3 text-base font-medium text-center text-white bg-blue-700 rounded-lg hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Create
          your account</button>

        <!-- Login Link -->
        <p class="text-sm font-medium text-gray-900 dark:text-white">
          Already have an account? <NuxtLink to="/signin" class="text-blue-600 hover:underline dark:text-blue-500"> Sign
            In </NuxtLink>
        </p>
      </form>
    </div>
  </div>
</template>


<script setup>
definePageMeta({
  title: 'Sign Up',
  layout: 'fullscreen',
});

const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const emailError = ref('');
const passwordError = ref('');

const userStore = useUserStore();

const submitForm = async () => {
  emailError.value = validateEmail(email.value);
  passwordError.value = validatePassword(password.value, confirmPassword.value);

  if (emailError.value || passwordError.value) {
    return;
  }

  const res = await userStore.signup({ store: userStore, email: email.value, password: password.value });
  if (res.status == 409) {
    emailError.value = res.data.message;
  }
};

const validateEmail = (email) => {
  const re = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
  if (re.test(email)) {
    return null
  } else {
    return "Invalid email address.";
  }
};

const validatePassword = (password1, password2) => {
  if (password1 !== password2) {
    return "Passwords do not match.";
  }
  if (password1.length < 6) {
    return "Password must be at least 6 characters long.";
  }
  return null;
};
</script>

