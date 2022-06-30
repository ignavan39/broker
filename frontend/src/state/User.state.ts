import axios from "axios";
import { atom, selector } from "recoil";
import { User } from "../types/User";

const getDefaultUser = (): User => {
  const cache = localStorage.getItem("user");
  if (!cache) {
    return {
      auth : {
        access : {
          token : "",
          expireAt : 0,
        },
        refresh : {
          token : "",
          expireAt : 0,
        }
      },
      user : {
        firstName: "",
        lastName: "",
        email: "",
        password: "",
        nickname: "",
        avatarUrl: "",
      }
    }
  } else {
    const user = JSON.parse(cache) as User;
    axios.defaults.headers.common["Authorization"] =
      "Bearer " + user.auth.access.token;
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

    return user?.auth.access.token?.length;
  },
});
