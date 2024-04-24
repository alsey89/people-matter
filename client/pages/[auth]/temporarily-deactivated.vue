<template>
  <div class="w-full h-screen flex items-center justify-center">
    <div v-auto-animate class="w-full flex flex-col justify-center items-center">
      <Card v-if="showForm" class="bg-gray-100">
        <div class="flex flex-col gap-4 p-4">
          <CardHeader>
            <UIText>
              Sign Up for Verve HR
            </UIText>
          </CardHeader>
          <!-- ----------------------------------------- -->
          <form class="w-full flex flex-col justify-center items-center gap-4" @submit.prevent="submitForm">
            <UIFormField>
              <template v-slot:icon>
                <Icon name="material-symbols:attach-email-outline" />
              </template>
              <input type="email" v-model="email" class="h-8 w-full focus:outline-primary px-1" maxlength="256"
                placeholder="Email Address" required />
            </UIFormField>

            <UIFormField>

              <template v-slot:icon>
                <Icon name="material-symbols:lock-outline" />
              </template>
              <input type="password" v-model="password" class="h-10 w-full focus:outline-primary px-1" maxlength="256"
                name="password" placeholder="Password (min 6 characters)" required />
            </UIFormField>
            <UIFormField>

              <template v-slot:icon>
                <Icon name="material-symbols:lock-outline" />
              </template>
              <input type="password" v-model="confirmPassword" class="h-10 w-full focus:outline-primary px-1"
                maxlength="256" name="confirmPassword" placeholder="Repeat Password" required />
            </UIFormField>
            <Button size="lg" class="w-full bg-primary hover:bg-primary-dark text-white">
              Sign Up
            </Button>
          </form>
          <!-- ----------------------------------------- -->
          <CardFooter>
            <div class="flex gap-2">
              Already have an account?
              <button @click="navigateTo('/auth/signin')"
                class="text-primary hover:text-primary-dark hover:underline">Sign
                In</button>
            </div>
          </CardFooter>
        </div>
      </Card>
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

// render form on mounted, to trigger animation
const showForm = ref(false);
onMounted(() => {
  showForm.value = true;
  userStore.getCsrfToken();
});
</script>