import axios from "axios";
import { Host } from "../config";
import { User } from "../types/User";

export type SignDto = User;

export const sign = async (
  payload: {
    email: string;
    password: string;
    nickname?: string;
    lastName?: string;
    firstName?: string;
  },
  operation: "signIn" | "signUp"
) => {
  const url =
    operation === "signIn" ? Host + "/auth/signIn" : Host + "/auth/signUp";
  const user = await axios.post<SignDto>(url, {
    ...payload,
  });
  axios.defaults.headers.common["Authorization"] =
    "Bearer " + user.data.auth.accessToken;
  return user.data;
};
