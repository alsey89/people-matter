<template>
  <div class="flex items-center justify-center h-screen">
    <div v-auto-animate class="w-80 flex flex-col justfiy-center items-center">
      <NBCard v-if="showForm" class="bg-gray-100">
        <div class="flex flex-col gap-4 p-4">
          <NBCardHeader>
            <NBText>
              Sign In To Verve HR
            </NBText>
          </NBCardHeader>
          <!-- ----------------------------------------- -->
          <form class="w-full flex flex-col justify-center items-center gap-4" @submit.prevent="submitForm">
            <NBFormField>
              <template v-slot:icon>
                <Icon name="material-symbols:attach-email-outline" />
              </template>
              <input type="email" v-model="email" class="h-8 w-full focus:outline-primary px-1" maxlength="256"
                name="name" placeholder="Email Address" required />
            </NBFormField>
            <NBFormField>
              <template v-slot:icon>
                <Icon name="material-symbols:lock-outline" />
              </template>
              <input type="password" v-model="password" class="h-10 w-full focus:outline-primary px-1" maxlength="256"
                name="password" placeholder="Enter Password" required />
            </NBFormField>
            <NBButtonSquare size="lg" class="w-full bg-primary hover:bg-primary-dark text-white">
              Sign In
            </NBButtonSquare>
          </form>
          <!-- ----------------------------------------- -->
          <NBCardFooter>
            Already have an account?
            <button @click="navigateTo('/signup')" class="text-primary hover:text-primary-dark hover:underline">Sign
              Up</button>
          </NBCardFooter>
        </div>
      </NBCard>
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


const userStore = useUserStore();

const submitForm = async () => {
  await userStore.signin({ email: email.value, password: password.value });
};

// render form on mounted, to trigger animation
const showForm = ref(false);
onMounted(() => {
  showForm.value = true;
});
</script>