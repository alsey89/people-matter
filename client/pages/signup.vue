<template>
  <div class="flex items-center justify-center h-screen">
    <div class="w-80 flex flex-col justify-center items-center">
      <NBCard class="bg-gray-100">
        <NBCardHeader>
          <NBText>
            Sign Up for Verve HR
          </NBText>
        </NBCardHeader>
        <!-- Form -->
        <form class="w-full flex flex-col justify-center items-center gap-4" @submit.prevent="submitForm">
          <!-- Email Field -->
          <NBFormField>
            <template v-slot:icon>
              <Icon name="material-symbols:attach-email-outline" />
            </template>
            <input type="email" v-model="email" class="h-8 w-full focus:outline-none px-1" maxlength="256"
              placeholder="Email Address" required />
          </NBFormField>
          <!-- Password Field -->
          <NBFormField>
            <template v-slot:icon>
              <Icon name="material-symbols:lock-outline" />
            </template>
            <input type="password" v-model="password" class="h-10 w-full focus:outline-none px-1" maxlength="256"
              name="password" placeholder="Password (min 6 characters)" required />
          </NBFormField>
          <NBFormField>
            <template v-slot:icon>
              <Icon name="material-symbols:lock-outline" />
            </template>
            <input type="password" v-model="confirmPassword" class="h-10 w-full focus:outline-none px-1" maxlength="256"
              name="confirmPassword" placeholder="Repeat Password" required />
          </NBFormField>
          <!-- Sign Up Button -->
          <NBButtonSquare class="w-full bg-secondary hover:bg-secondary-dark text-white">
            Sign Up
          </NBButtonSquare>
        </form>
        <NBCardFooter>
          Already have an account?
          <NuxtLink to="/signin">
            <button class="text-secondary hover:text-secondary-dark hover:underline">Sign In</button>
          </NuxtLink>
        </NBCardFooter>
      </NBCard>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  title: 'Sign Up',
  layout: 'fullscreen',
});

const name = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');

const userStore = useUserStore();
const messageStore = useMessageStore();

const submitForm = async () => {
  const isValidEmail = (email) => {
    const re = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
    return re.test(email);
  };
  if (!isValidEmail(email.value)) {
    messageStore.setError("Invalid email address.");
    return;
  }

  if (password.value !== confirmPassword.value) {
    messageStore.setError("Passwords do not match.");
    return;
  }
  if (password.value.length < 6) {
    messageStore.setError("Password must be at least 6 characters.");
    return;
  }

  try {
    const res = await userStore.signup({
      username: name.value,
      email: email.value,
      password: password.value,
      confirmPassword: confirmPassword.value
    });
  } catch (error) {
    console.error(error);
    messageStore.setError("Something went wrong.");
  }
};
</script>