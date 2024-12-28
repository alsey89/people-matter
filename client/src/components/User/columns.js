import { h } from "vue";

export const columns = [
  {
    accessorKey: "avatar",
    header: () => h("div", { class: "text-center" }, "Avatar"),
    cell: ({ row }) => {
      const avatarUrl = row.getValue("avatar");
      return h("div", { class: "text-center" }, [
        h("img", {
          src: avatarUrl,
          alt: "Avatar",
          class: "w-10 h-10 rounded-full mx-auto",
        }),
      ]);
    },
  },
  {
    accessorKey: "name",
    header: () => h("div", { class: "text-left" }, "Name"),
    cell: ({ row }) =>
      h("div", { class: "text-left font-medium" }, row.getValue("name")),
  },
  {
    accessorKey: "email",
    header: () => h("div", { class: "text-left" }, "Email"),
    cell: ({ row }) => h("div", { class: "text-left" }, row.getValue("email")),
  },
  {
    accessorKey: "roles",
    header: () => h("div", { class: "text-left" }, "Role(s)"),
    cell: ({ row }) => {
      const roles = row.getValue("roles");
      return h("div", { class: "text-left" }, roles.join(", "));
    },
  },
  {
    accessorKey: "location",
    header: () => h("div", { class: "text-left" }, "Location"),
    cell: ({ row }) =>
      h("div", { class: "text-left" }, row.getValue("location")),
  },
];
