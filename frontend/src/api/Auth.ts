import axios from "axios";
import { Host } from "../config";
import { User } from "../types/User";

export type AuthorizationResponseDto = User;
export enum AuthorizationActionType {
  signIn = "signIn",
  signUp = "signUp",
}

export type SignUpPayload = {
  email: string;
  password: string;
  nickname: string;
  lastName: string;
  firstName: string;
  code: number;
};

export type SignInPayload = {
  email: string;
  password: string;
};

export type AuthorizationService = {
  [AuthorizationActionType.signIn]: (
    payload: SignInPayload
  ) => Promise<AuthorizationResponseDto>;
  [AuthorizationActionType.signUp]: (
    payload: SignUpPayload
  ) => Promise<AuthorizationResponseDto>;
};

export const authorizationService: AuthorizationService = {
  [AuthorizationActionType.signIn]: async (
    payload: SignInPayload
  ): Promise<AuthorizationResponseDto> => {
    const url = Host + "/auth/signIn";
    const user = await axios.post<AuthorizationResponseDto>(url, payload);
    axios.defaults.headers.common["Authorization"] =
      "Bearer " + user.data.auth.access.token;
    return user.data;
  },

  [AuthorizationActionType.signUp]: async (
    payload: SignUpPayload
  ): Promise<AuthorizationResponseDto> => {
    const url = Host + "/auth/signUp";
    const user = await axios.post<AuthorizationResponseDto>(url, payload);
    axios.defaults.headers.common["Authorization"] =
      "Bearer " + user.data.auth.access.token;
    return user.data;
  },
};
