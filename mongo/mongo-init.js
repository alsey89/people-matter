db = db.getSiblingDB("verve");
db.createCollection("users");

//todo: add sample data to populate at startup

db.users.insertMany([
  {
    _id: {
      $oid: "65707fa26d8bba7914f919ad",
    },
    userId: {
      $binary: {
        base64: "fuFncJRAEe6Z5gJCrBIAAw==",
        subType: "00",
      },
    },
    username: "Michael",
    email: "test@mailnesia.com",
    password: "$2a$10$36VCf6OAqJXRy4cVvl7CCOUlhVDIY9Y3Mt6nDvhdFa8dm6uh3U3cC",
    isAdmin: true,
    avatarUrl:
      "https://static.vecteezy.com/system/resources/previews/000/439/863/original/vector-users-icon.jpg",
    createdAt: {
      $date: "2023-12-06T14:05:22.258Z",
    },
    updatedAt: {
      $date: "2023-12-06T14:05:22.258Z",
    },
  },
  {
    _id: {
      $oid: "657087f17b80f68cb0efacbf",
    },
    userId: {
      $binary: {
        base64: "cqPG4pRFEe6LdAJCrBIAAw==",
        subType: "00",
      },
    },
    username: "Michael the 2nd",
    email: "test2@mailnesia.com",
    password: "$2a$10$fjMQ911jIYEj77fREldrEuI4QBtGU6z.8ke4b36RCymo64BDXvlWG",
    isAdmin: false,
    avatarUrl:
      "https://cdn3.vectorstock.com/i/1000x1000/30/97/flat-business-man-user-profile-avatar-icon-vector-4333097.jpg",
    createdAt: {
      $date: "2023-12-06T14:40:49.206Z",
    },
    updatedAt: {
      $date: "2023-12-06T14:40:49.206Z",
    },
  },
]);
