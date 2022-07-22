import axios from "axios";
import { Host } from "../config";
import { User } from "../types/User";

export type AuthorizationResponseDto = User;

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

export type SendVerifyCodePayload = {
  email: string;
}

export type AuthorizationService = {
  signIn: (
    payload: SignInPayload
  ) => Promise<AuthorizationResponseDto>;
  signUp: (
    payload: SignUpPayload
  ) => Promise<AuthorizationResponseDto>;
  sendVerifyCode: (
    payload: SendVerifyCodePayload
  ) => void
};

export const authorizationService: AuthorizationService = {
  signIn: async (
    payload: SignInPayload
  ): Promise<AuthorizationResponseDto> => {
    const url = Host + "/auth/signIn";
    const user = await axios.post<AuthorizationResponseDto>(url, payload);
    axios.defaults.headers.common["Authorization"] =
      "Bearer " + user.data.auth.access.token;
    return user.data;
  },

  signUp: async (
    payload: SignUpPayload
  ): Promise<AuthorizationResponseDto> => {
    const url = Host + "/auth/signUp";
    const user = await axios.post<AuthorizationResponseDto>(url, payload);
    axios.defaults.headers.common["Authorization"] =
      "Bearer " + user.data.auth.access.token;
    return user.data;
  },

  sendVerifyCode: async (
    payload: SendVerifyCodePayload
  ) => {
    const url = Host + "/auth/sendVerifyCode";
    await axios.post(url, payload);
  }
};
