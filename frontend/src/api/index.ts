import axios from "axios";
import { User } from "../types/User";

const Host = process.env.HOST || "http://127.0.0.1:8080/api/v1";

export type SignUpDto = {
  user:{
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    password: string;
  };
  token: string;
}

export namespace Api {
  export const signUp = async (payload: Omit<User, "auth">) => {
    const user = await axios.post<SignUpDto>(Host + "/users/signUp", {
      ...payload,
    });
    console.log(user);
    axios.defaults.headers.common["Authorization"] = 'Bearer ' + user.data.token;
    return user.data;
  };
}
