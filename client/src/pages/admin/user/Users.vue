<template>
    <div class="w-full h-full flex flex-col gap-2 p-4">
        <div class="h-8 flex gap-2 text-lg font-semibold px-4 py-2">
            <div>Company Users</div>
        </div>
        <hr class="border-t mx-2">
        <Table>
            <TableCaption>A list of all company users.</TableCaption>
            <TableHeader>
                <TableRow>
                    <TableHead>Avatar</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Email</TableHead>
                    <TableHead>Role(s)</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody v-for="user in userStore.users" :key="user.id">
                <TableRow @click="onSelectUser(user)">
                    <TableCell>
                        <Avatar>
                            <AvatarImage :src="user.avatar" />
                            <AvatarFallback>{{ user.firstName }}</AvatarFallback>
                        </Avatar>
                    </TableCell>
                    <TableCell>{{ formatName(user) }}</TableCell>
                    <TableCell>{{ user.email }}</TableCell>
                    <TableCell class="flex flex-col justify-between">
                        <div v-for="userRole in user.userRoles">
                            {{ userRole.role.name }} | {{ userRole.location.name }}
                        </div>
                    </TableCell>
                </TableRow>
            </TableBody>
        </Table>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useUserStore } from '@/stores/User';
import { Icon } from '@iconify/vue';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'

const userStore = useUserStore();

const onSelectUser = (user) => {
    userStore.setSelectedUser(user);
};

const formatName = (user) => {
    return `${user.firstName} ${user.middleName} ${user.lastName}`;
};

onMounted(() => {
    userStore.getUsers();
});
</script>
