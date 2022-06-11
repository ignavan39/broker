import axios from "axios";
import { Host } from "../config";
import { User } from "../types/User";

export type SignUpDto = {
  user: {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    password: string;
  };
  auth: {
    accessToken: string;
    refreshToken: string;
  }
};

export const signUp = async (payload: Omit<User, "auth">) => {
  const user = await axios.post<SignUpDto>(Host + "/auth/signUp", {
    ...payload,
  });
  console.log(user);
  axios.defaults.headers.common["Authorization"] = "Bearer " + user.data.auth.accessToken;
  return user.data;
};