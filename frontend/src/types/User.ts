export type User = {
  auth: {
    accessToken: string;
    refreshToken: string;
  };
  user: {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
    nickname: string;
    avatarUrl: string;
  };
};
