export type User = {
  auth: { accessToken: string; expiresAt: Date | null };
  name: string;
  email: string;
};

