<template>
    <Card>
        <CardHeader>
            <h1> Welcome! </h1>
            <p> This is a starter template example. </p>
        </CardHeader>
        <CardContent>
            <Button @click="isConfirmed = true"> 
                <p v-if="isConfirmed"> Good luck! </p>
                <p v-else> Got it! </p>    
            </Button>
        </CardContent>
    </Card>
</template>

<script setup>
const isConfirmed = ref(false)
</script>