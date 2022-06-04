export type User = {
  auth: { accessToken: string; expiresAt: Date | null };
  firstName: string,
  lastName: string,
  email: string;
  password: string;
};

