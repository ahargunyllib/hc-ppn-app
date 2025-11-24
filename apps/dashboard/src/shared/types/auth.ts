/**
 * Authentication type definitions
 */

export type User = {
  id: string;
  email: string;
  name: string;
};

export type LoginCredentials = {
  email: string;
  password: string;
  rememberMe?: boolean;
};

export type AuthState = {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
};

export type StoredSession = {
  user: User;
  token: string;
  expiresAt: number;
};
