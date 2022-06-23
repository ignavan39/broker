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
    return JSON.parse(cache) as User;
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
