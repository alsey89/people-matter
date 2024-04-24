<script setup lang="ts">
import {
    Avatar,
    AvatarFallback,
    AvatarImage,
} from '@/components/ui/avatar'

// Mock data for demonstration purposes
const people = [
    { name: 'Olivia Martin', email: 'olivia.martin@email.com', avatar: '/avatars/01.png', infoLabel: 'claims:', infoValue: '+$1,999.00' },
    { name: 'Jackson Lee', email: 'jackson.lee@email.com', avatar: '/avatars/02.png', infoLabel: 'claims:', infoValue: '+$39.00' },
    { name: 'Isabella Nguyen', email: 'isabella.nguyen@email.com', avatar: '/avatars/03.png', infoLabel: 'Leave', infoValue: '2024/12/31' },
    { name: 'William Kim', email: 'will@email.com', avatar: '/avatars/04.png', infoLabel: 'claims:', infoValue: '+$39.00' },
    { name: 'Sofia Davis', email: 'sofia.davis@email.com', avatar: '/avatars/05.png', infoLabel: 'Leave', infoValue: '2024/06/06' },
]
</script>

<template>
    <div class="space-y-4">
        <div v-for="person in people" :key="person.email" class="flex items-center">
            <Avatar class="h-9 w-9">
                <AvatarImage :src="person.avatar" :alt="person.name" />
                <AvatarFallback>{{ person.name[0] }}</AvatarFallback>
            </Avatar>
            <div class="ml-4 flex-1 min-w-0">
                <p class="text-sm font-medium truncate">{{ person.name }}</p>
                <p class="text-sm text-muted-foreground truncate">{{ person.email }}</p>
            </div>
            <div class="ml-4 flex-none text-right">
                <div class="text-sm font-semibold">{{ person.infoLabel }}</div>
                <div class="text-sm">{{ person.infoValue }}</div>
            </div>
        </div>
    </div>
</template>
