<template>
    <div class="border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] bg-secondary-bg">
        <div v-auto-animate class="px-2 gap-2 flex justify-between items-center p-1 md:p-2">
            <button @click="navigateTo('/company')" class="hover:text-primary">
                <client-only>
                    <div v-if="companyName" class="flex gap-2 font-bold">
                        <AppLogo :src="companyLogoUrl || '/defaultLogo.png'" shape="square" class="w-6 h-6" />
                        <div>{{ companyName }}</div>
                    </div>
                    <div v-else> No Company Data </div>
                </client-only>
            </button>
            <DropdownMenu>
                <DropdownMenuTrigger as-child>
                    <Button variant="outline" class="relative h-8 w-8 rounded-full">
                        <Avatar class="h-8 w-8">
                            <AvatarImage src="userStore.getCurrentUserAvatarUrl" alt="avatar" />
                            <AvatarFallback>SC</AvatarFallback>
                        </Avatar>
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent class="w-56" align="end">
                    <DropdownMenuLabel class="font-normal flex">
                        <div class="flex flex-col space-y-1">
                            <p class="text-sm font-medium leading-none">
                                John Doe
                            </p>
                            <p class="text-xs leading-none text-muted-foreground">
                                john@doe.com
                            </p>
                        </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuGroup>
                        <DropdownMenuItem @click="navigateTo('/')">
                            Home
                            <DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
                        </DropdownMenuItem>
                        <DropdownMenuItem @click="navigateTo('/company')">
                            Company
                            <DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
                        </DropdownMenuItem>
                        <DropdownMenuItem @click="navigateTo('/user/list')">
                            Users
                            <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
                        </DropdownMenuItem>
                    </DropdownMenuGroup>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem @click="userStore.signout">
                        Sign out
                        <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    </div>
</template>

<script setup>
const companyName = persistedState.sessionStorage.getItem('companyName');
const companyLogoUrl = persistedState.sessionStorage.getItem('companyLogoUrl');

const userStore = useUserStore();
</script>