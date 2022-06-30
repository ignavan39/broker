export type User = {
  auth: { 
    access: {
      token: string;
      expireAt: number;
    };
    refresh: {
      token: string;
      expireAt: number;
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
