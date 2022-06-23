import axios from "axios";
import { atom, selector } from "recoil";
import { User } from "../types/User";

const getDefaultUser = (): User => {
  const cache = localStorage.getItem("user");
  if (!cache) {
    return {
      user: {
        email: "",
        firstName: "",
        lastName: "",
        avatarUrl: "",
        password: "",
        nickname: "",
      },
      auth: {
        accessToken: "",
        refreshToken: "",
      },
    };
  } else {
    const user = JSON.parse(cache) as User;
    axios.defaults.headers.common["Authorization"] =
      "Bearer " + user.auth.accessToken;
    return user;
  }
};

export const userState = atom<User>({
  key: "userState",
  default: getDefaultUser(),
});

export const userIsLoggined = selector({
  key: "isLoggined",
  get: ({ get }) => {
    const user = get(userState);

    return user?.auth.accessToken?.length;
  },
});
