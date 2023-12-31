<template>
  <div class="flex items-center justify-center h-screen">
    <div class="w-80 flex flex-col justfiy-center items-center">
      <NBCardVertical class=" bg-secondary-bg">
        <NBCardHeader>
          <NBText>
            Sign In To Verve HR
          </NBText>
        </NBCardHeader>
        <!-- Form -->
        <form class="w-full flex flex-col justify-center items-center gap-4" @submit.prevent="submitForm">
          <NBFormField>
            <template v-slot:icon>
              <Icon name="material-symbols:attach-email-outline" />
            </template>
            <input type="email" v-model="email" class="h-8 w-full focus:outline-none px-1" maxlength="256" name="name"
              placeholder="Email Address" required />
          </NBFormField>
          <NBFormField>
            <template v-slot:icon>
              <Icon name="material-symbols:lock-outline" />
            </template>
            <input type="password" v-model="password" class="h-10 w-full focus:outline-none px-1" maxlength="256"
              name="password" placeholder="Enter Password" required />
          </NBFormField>
          <!-- Button -->
          <NBButtonSquare class="w-full bg-sky-700 hover:bg-sky-900 text-white">
            Sign In
          </NBButtonSquare>
        </form>
        <NBCardFooter>
          Already have an account?
          <NuxtLink to="/signup">
            <button class="text-sky-700 hover:text-sky-900 hover:underline">Sign Up</button>
          </NuxtLink>
        </NBCardFooter>
      </NBCardVertical>
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
// const remember = ref(false);

const userStore = useUserStore();

const submitForm = async () => {
  await userStore.signin({ email: email.value, password: password.value });
};
</script>