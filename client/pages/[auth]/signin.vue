<template>
  <div class="w-full h-screen flex items-center justify-center">
    <div v-auto-animate class="w-full flex flex-col justfiy-center items-center">
      <Card v-if="showForm" class="min-w-64 bg-gray-100">
        <CardHeader>
          <div class="text-center font-bold">
            Sign In To Verve HR
          </div>
        </CardHeader>
        <!-- ----------------------------------------- -->
        <CardContent>
          <div class="flex flex-col gap-4">
            <div class="border border-primary rounded-sm p-4 text-sm text-gray-500">
              <span class="font-bold">
                Email:
              </span>
              <span class="text-primary">
                test@michaelchen.me
              </span>
              <br />
              <span class="font-bold">
                Password:
              </span>
              <span class="text-primary">
                testtesttest
              </span>
            </div>
            <div>
              <form class="flex flex-col justify-start gap-4" @submit.prevent="submitForm">
                <div class="flex items-center gap-2">
                  <Icon name="ic:twotone-email" class="h-6 w-6" />
                  <input type="email" v-model="email" class="p-2 h-10 w-full focus:outline-secondary" maxlength="256"
                    name="name" placeholder="Email Address" required />
                </div>
                <div class="flex items-center gap-2">
                  <Icon name="ic:twotone-password" class="h-6 w-6" />
                  <input type="password" v-model="password" class="p-2 h-10 w-full focus:outline-secondary"
                    maxlength="256" name="password" placeholder="Enter Password" required />
                </div>
                <Button class="w-full bg-primary hover:bg-primary-dark text-white">
                  Sign In
                </Button>
              </form>
            </div>
          </div>
        </CardContent>
        <!-- <CardFooter>
            <div class="flex gap-2">
              Already have an account?
              <button @click="navigateTo('/auth/signup')"
                class="text-primary hover:text-primary-dark hover:underline">Sign
                Up</button>
            </div>
          </CardFooter> -->
      </Card>
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
  userStore.getCsrfToken();
});
</script>