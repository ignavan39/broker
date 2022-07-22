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
  profile : {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
    nickname: string;
    avatarUrl: string;
  };
};
