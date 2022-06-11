export type User = {
  auth: { accessToken: string; refreshToken: string };
  firstName: string;
  lastName: string;
  email: string;
  password: string;
};
