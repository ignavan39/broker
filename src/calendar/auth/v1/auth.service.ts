import { repository } from "../../common/repository/user.repository";

export class AuthService {
    public signIn = async (login: string,password: string) => {
        const user = await repository.User.save({
            email: login,
            password,
        })
    }
}