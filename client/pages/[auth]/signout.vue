<template>
  <div>
    FOR USE WHEN THERE'S NO SIGNOUT BUTTON
  </div>
</template>

<script setup>
const userState = useUserStore();
await userState.signout();
</script>
