export type User = {
  auth: {
    access: {
      token: string;
      expireAt: number | null;
    };
    refresh: {
      token: string;
      expireAt: number | null;
    };
  };
  user : {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
    nickname: string;
    avatarUrl: string;
  };
};
