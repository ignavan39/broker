import axios from "axios";
import { Host } from "../config";
import { User } from "../types/User";

export type SignDto = User;

export const sign = async (
  payload:
    | {
        operation: "sign_up";
        email: string;
        password: string;
        nickname: string;
        lastName: string;
        firstName: string;
        code: number;
      }
    | {
        operation: "sign_in";
        email: string;
        password: string;
      }
) => {
  const url =
    payload.operation === "sign_in"
      ? Host + "/auth/signIn"
      : Host + "/auth/signUp";
  const user = await axios.post<SignDto>(url, {
    ...payload,
  });
  axios.defaults.headers.common["Authorization"] =
    "Bearer " + user.data.auth.access.token;
  return user.data;
};
