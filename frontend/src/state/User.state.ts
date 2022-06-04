import { atom, selector } from "recoil";
import { User } from "../types/User";

const getDefaultUser = () => {
  const cache = localStorage.getItem("user");
  if (!cache) {
    return {
      email: "",
      name: "",
      auth: {
        expiresAt: null,
        accessToken: "",
      },
    };
  } else {
    return JSON.parse(cache) as User;
  }
};

export const userState = atom({
  key: "userState",
  default: getDefaultUser(),
});

export const userIsLoggined = selector({
  key: "isLoggined",
  get: ({ get }) => {
    const user = get(userState);

    return user.auth.accessToken.length;
  },
});
