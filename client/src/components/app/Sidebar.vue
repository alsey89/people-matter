<template>
    <aside ref="sidebar"
        :class="['fixed top-0 left-0 h-full flex flex-col px-2 gap-2 border-r-2 border-accent bg-background select-none', isMobile ? 'w-3/4' : 'w-[250px]', 'z-50']">
        <h2 class="h-8 text-lg font-semibold px-4 py-2"> People Matter </h2>
        <div v-for="section in menuSections" :key="section.label">
            <h2 class="text-xs text-gray-500 p-2 border-t">{{ section.label }}</h2>
            <ul class="py-2">
                <li v-for="item in section.items" :key="item.name" class="group">
                    <div v-if="item.disabledUntilSelected && !selectedUser"
                        class="flex items-center p-2 border-2 border-background rounded-md opacity-50 cursor-not-allowed">
                        <Icon :icon="item.disabledIcon" class="text-xl" />
                        <span class="ml-3">{{ item.disabledText }}</span>
                    </div>
                    <div v-else-if="item.children" @click="toggle(item, $event)"
                        class="flex items-center p-2 border-2 border-background rounded-md cursor-pointer hover:border-accent">
                        <Icon :icon="item.icon" class="text-xl" />
                        <span class="ml-3">{{ item.name }}</span>
                        <span class="ml-auto transition-transform duration-200"
                            :class="item.isOpen ? 'rotate-180' : ''">
                            <Icon icon="material-symbols:expand-more" />
                        </span>
                    </div>
                    <div v-else @click="navigateTo(item.path)" :class="['flex items-center p-2 border-2 border-background rounded-md cursor-pointer hover:border-accent',
                        item.path === $route.path ? 'border-primary/20 bg-gray-100' : '']">
                        <Icon :icon="item.icon" class="text-xl" />
                        <span class="ml-3">{{ item.name }}</span>
                    </div>
                    <ul ref="submenu" class="ml-8 overflow-hidden">
                        <li @click="navigateTo(child.path)" v-for="child in item.children" :key="child.name" :class="['p-2 border-2 border-background rounded-md cursor-pointer hover:border-accent',
                            child.path === $route.path ? 'border-primary/20 bg-gray-100' : '']">
                            <div class="flex items-center">
                                <Icon :icon="child.icon" class="text-xl" />
                                <span class="ml-2">{{ child.name }}</span>
                            </div>
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
        <router-link to="/auth/signout"
            class="flex items-center p-2 border-2 border-background rounded-md cursor-pointer hover:border-accent">
            <Icon icon="mdi:logout" class="text-xl" />
            <span class="ml-3"> signout </span>
        </router-link>
    </aside>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, nextTick } from 'vue';
import { gsap } from 'gsap';
import { Icon } from '@iconify/vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/User';

const userStore = useUserStore();

const props = defineProps({
    showSidebar: Boolean,
    isMobile: Boolean,
    selectedUser: Object
});
const emit = defineEmits(['closeSidebar']);

const router = useRouter();
const navigateTo = (path) => {
    if (router.currentRoute.value.path === path) {
        return;
    }
    router.push(path);
    //emit close sidebar event
    if (props.isMobile) {
        emit('closeSidebar');
    }
};
const selectedUser = computed(() => userStore.selectedUser);
// const selectedUser = {
//     id: 1,
//     firstName: 'John',
//     lastName: 'Doe'
// };

const menuSections = ref([
    {
        label: 'Administrator',
        items: [
            { name: 'Dashboard', icon: 'material-symbols:dashboard', path: '/admin' },
            {
                name: 'Company',
                icon: 'mdi:office-building',
                isOpen: false,
                children: [
                    { name: 'Profile', icon: 'ooui:view-details-ltr', path: '/admin/company' },
                    { name: 'Locations', icon: 'material-symbols:map', path: '/admin/company/locations' },
                    { name: 'Departments', icon: 'system-uicons:hierarchy', path: '/admin/company/departments' },
                    { name: 'Positions', icon: 'hugeicons:job-link', path: '/admin/company/positions' },
                ]
            },
            {
                name: 'Policies',
                icon: 'material-symbols:policy',
                isOpen: false,
                children: [
                    { name: 'Attendance', icon: 'material-symbols:punch-clock', path: '/admin/attendance' },
                    { name: 'Leave', icon: 'material-symbols:sick', path: '/admin/leave' },
                    { name: 'Salary', icon: 'material-symbols:payments', path: '/admin/salary' },
                    { name: 'Claims', icon: 'material-symbols:money', path: '/admin/claims' },
                    { name: 'Compliance', icon: 'octicon:law', path: '/admin/compliance' },
                ]
            },
            { name: 'Users', icon: 'ph:user-list', path: '/admin/user' },
            {
                name: selectedUser.value ? `${selectedUser.value.firstName} ${selectedUser.value.lastName}` : 'No User Selected',
                icon: 'material-symbols:account-circle',
                path: `/admin/user/${selectedUser.id}`,
                disabledUntilSelected: true,
                disabledIcon: 'material-symbols:account-circle-off',
                disabledText: 'No user selected',
                children: [
                    { name: 'Profile', icon: 'fluent:slide-text-person-16-filled', path: `/admin/user/${selectedUser.id}/profile` },
                    { name: 'Leave', icon: 'material-symbols:sick-outline', path: `/admin/user/${selectedUser.id}/leave` },
                    { name: 'Salary', icon: 'material-symbols:payments', path: `/admin/user/${selectedUser.id}/salary` },
                    { name: 'Claims', icon: 'material-symbols:money-outline', path: `/admin/user/${selectedUser.id}/claims` },
                ]
            },
        ]
    },
    {
        label: 'Location: {location}',
        items: [
            { name: 'Dashboard', icon: 'material-symbols:dashboard', path: '/manager' },
            {
                name: 'Location',
                icon: 'material-symbols:map',
                isOpen: false,
                children: [
                    { name: 'Users', icon: 'ph:user-list', path: '/manager/users' },
                    { name: 'Leave', icon: 'material-symbols:sick-outline', path: '/manager/leave' },
                    { name: 'Salary', icon: 'material-symbols:payments', path: '/manager/salary' },
                    { name: 'Claims', icon: 'material-symbols:money-outline', path: '/manager/claims' },
                ]
            }
        ]
    },
    {
        label: 'User: {userId}',
        items: [
            { name: 'Dashboard', icon: 'material-symbols:dashboard', path: '/user' },
            {
                name: 'User',
                icon: 'material-symbols:frame-person-outline',
                isOpen: false,
                children: [
                    { name: 'Profile', icon: 'fluent:slide-text-person-16-filled', path: '/user/profile' },
                    { name: 'Leave', icon: 'material-symbols:sick-outline', path: '/user/leave' },
                    { name: 'Salary', icon: 'material-symbols:payments', path: '/user/salary' },
                    { name: 'Claims', icon: 'material-symbols:money-outline', path: '/user/claims' },
                ]
            }
        ]
    }
]);

const submenu = ref([]);

const toggle = (item, event) => {
    if (item.disabled) return;
    if (!item.children) return;
    item.isOpen = !item.isOpen;
    const submenu = event.currentTarget.nextElementSibling;
    if (item.isOpen) {
        gsap.fromTo(submenu, { height: 0 }, { height: 'auto', duration: 0.5 });
    } else {
        gsap.to(submenu, { height: 0, duration: 0.5 });
    }
};

const setInitialState = () => {
    submenu.value.forEach((item, index) => {
        if (!item.children) return;
        if (item.isOpen) {
            gsap.set(submenu.value[index], { height: 'auto' });
        } else {
            gsap.set(submenu.value[index], { height: 0 });
        }
    });
};

const sidebar = ref(null);

const animateSidebar = async () => {
    await nextTick();  // Ensure DOM updates are completed
    const sidebarWidth = sidebar.value.offsetWidth;
    gsap.to(sidebar.value, { x: props.showSidebar ? 0 : -sidebarWidth, duration: 0.5 });
};

watch(() => props.showSidebar, () => {
    animateSidebar();
});

onMounted(async () => {
    animateSidebar();
    await nextTick();
    setInitialState();
});
</script>