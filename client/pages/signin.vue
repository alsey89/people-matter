<template>
  <div class="flex justify-center items-center min-h-screen">
    <div class="w-full max-w-md px-6 py-10 space-y-8 bg-white rounded-lg shadow-xl dark:bg-gray-800">
      <!-- Form Starts Here -->
      <form class="space-y-6" @submit.prevent="submitForm">
        <!-- Email Input -->
        <div>
          <label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your
            email</label>
          <input type="email" v-model="email" id="email"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white"
            required>
        </div>

        <!-- Password Input -->
        <div>
          <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your
            password</label>
          <input type="password" v-model="password" id="password"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white"
            required>
        </div>

        <!-- Remember Me and Lost Password -->
        <div class="flex items-start justify-between">
          <div class="flex items-center">
            <input type="checkbox" v-model="remember" id="remember"
              class="w-4 h-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:bg-gray-700 dark:border-gray-600">
            <label for="remember" class="ml-2 text-sm font-medium text-gray-500 dark:text-gray-400">Remember
              this device</label>
          </div>
          <a href="#" class="text-sm font-medium text-blue-600 hover:underline dark:text-blue-500">Lost
            Password?</a>
        </div>

        <!-- Submit Button -->
        <button type="submit"
          class="w-full px-5 py-3 text-base font-medium text-center text-white bg-blue-700 rounded-lg hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Login
          to your account</button>

        <!-- Registration Link -->
        <p class="text-sm font-medium text-gray-900 dark:text-white">
          Not registered yet? <NuxtLink to="/signup" class="text-blue-600 hover:underline dark:text-blue-500">
            Sign Up </NuxtLink>
        </p>
      </form>
    </div>
  </div>
</template>


<script setup>
definePageMeta({
  title: 'Sign In',
  layout: 'fullscreen',
});

const email = ref('');
const password = ref('');
const remember = ref(false);

const userStore = useUserStore();

const submitForm = async () => {
  await userStore.signin({ store: userStore, email: email.value, password: password.value });
};

</script>
