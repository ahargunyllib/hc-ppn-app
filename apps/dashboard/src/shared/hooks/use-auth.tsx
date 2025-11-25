/**
 * Authentication Context and Hook
 * Provides global authentication state and methods
 */

import type { ReactNode } from "react";
import { createContext, useContext, useEffect, useState } from "react";
import type {
  AuthState,
  LoginCredentials,
  StoredSession,
  User,
} from "../types/auth";

type AuthContextType = AuthState & {
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthProviderProps = {
  children: ReactNode;
};

/**
 * AuthProvider component that wraps the application
 * Manages authentication state and provides auth methods
 */
export default function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  // Restore session on mount
  useEffect(() => {
    const getCurrentUser = (): User | null => {
      const session = localStorage.getItem("authSession");
      if (!session) {
        return null;
      }

      const storedSession = JSON.parse(session) as StoredSession;
      if (Date.now() >= storedSession.expiresAt) {
        localStorage.removeItem("authSession");
        return null;
      }

      return storedSession.user;
    };

    const restoreSession = (): void => {
      const currentUser = getCurrentUser();
      setUser(currentUser);
      setIsLoading(false);
    };

    restoreSession();
  }, []);

  const login = (credentials: LoginCredentials): Promise<void> => {
    if (
      credentials.email !== "admin@gmail.com" ||
      credentials.password !== "password"
    ) {
      return Promise.reject(new Error("Invalid email or password"));
    }

    const loggedInUser: User = {
      id: "1",
      email: credentials.email,
      name: "Admin User",
    };

    const expiresAt =
      Date.now() +
      (credentials.rememberMe ? 7 * 24 * 60 * 60 * 1000 : 60 * 60 * 1000);

    const session: StoredSession = {
      user: loggedInUser,
      token: "dummy-token",
      expiresAt,
    };

    localStorage.setItem("authSession", JSON.stringify(session));
    setUser(loggedInUser);
    return Promise.resolve();
  };

  /**
   * Logout current user
   */
  const logout = (): void => {
    localStorage.removeItem("authSession");
    setUser(null);
  };

  const value: AuthContextType = {
    user,
    isAuthenticated: user !== null,
    isLoading,
    login,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * Hook to access authentication context
 * Must be used within AuthProvider
 */
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
}
