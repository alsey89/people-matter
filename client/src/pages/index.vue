<template>
    <div v-auto-animate>
        <!-- Company Created -->
        <Card class="w-[450px] flex flex-col gap-2 p-4">
            <p class="text-xl font-bold"> This is the dashboard! </p>
            <p> This is a placeholder. The actual dashboard is coming soon. </p>
            <Button @click="this.$router.push('/auth/signin')" class="w-full"> Take me back to the login screen!
            </Button>
        </Card>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

const router = useRouter();
</script>