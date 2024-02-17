<template>
    <NBCard v-if="userStore.getSingleUserData">
        <div class="flex flex-col gap-2 px-2">
            <h1 class="text-lg font-bold">
                Contact Information
            </h1>
            <div>
                Email: {{ userStore.getSingleUserEmail || "No Data" }}
            </div>
            <div>
                Mobile: {{ userStore.getSingleUserMobile || "No Data" }}
            </div>
            <div>
                Address: {{ userStore.getSingleUserFullAddress || "No Data" }}
            </div>
        </div>
    </NBCard>
    <NBCard v-else>
        <div v-if="userStore.isLoading">
            Loading...
        </div>
        <div v-else class="m-auto">
            No Data
        </div>
    </NBCard>
    <NBCard v-if="userStore.getSingleUserEmergencyContactFullName">
        <div class="flex flex-col gap-2 px-2">
            <h1 class="text-lg font-bold">
                Emergency Contact Information
            </h1>
            <div>
                Name: {{ userStore.getSingleUserEmergencyContactFullName }}
            </div>
            <div>
                Email: {{ userStore.getSingleUserEmergencyContactEmail || "No Data" }}
            </div>
            <div>
                Mobile: {{ userStore.getSingleUserEmergencyContactMobile || "No Data" }}
            </div>
            <div>
                Relation: {{ userStore.getSingleUserEmergencyContactRelation || "No Data" }}
            </div>
        </div>
    </NBCard>
    <NBCard v-else>
        <div class="m-auto">
            No Data
        </div>
    </NBCard>
</template>

<script setup>
definePageMeta({
    title: "User Details",
    description: "User details",
    layout: "user-details",
});

const userStore = useUserStore();

const route = useRouter().currentRoute.value
</script>

type ContactInfo struct {
	gorm.Model
	UserID     uint   `json:"userId"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
}

type EmergencyContact struct {
	gorm.Model
	UserID     uint    `json:"userId"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Relation   string  `json:"relation"`
	Mobile     string  `json:"mobile"`
	Email      string  `json:"email"`
}