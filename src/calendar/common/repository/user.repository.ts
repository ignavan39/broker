import { getRepository } from "typeorm";
import { User as UserEntity } from "../entities/user";

export namespace repository {
  export const User = getRepository(UserEntity);
}
