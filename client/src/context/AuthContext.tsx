import React, { useEffect, useState, useContext, createContext } from "react";
import { fetchData } from "../services/backendApi";

type User = {
  id: string;
  firstname: string;
  lastname: string;
  username: string;
  email: string;
  createdAt: string;
  updatedAt: string;
};

type AuthContextType = {
  user: User | null;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
  isLoading: boolean;
  setIsLoading: React.Dispatch<React.SetStateAction<boolean>>;
};

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  useEffect(() => {
    const fetchUser = async () => {
      setIsLoading(true);
      try {
        const response = await fetchData.get("/api/me", {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        });

        if (response.status === 200) {
          setUser(response.data.currentUser);
          console.log("Current user: ", response.data.currentUser);
        }
      } catch (err: any) {
        console.error(err.message);
      } finally {
        setIsLoading(false);
      }
    };
    fetchUser();
  }, []);

  return (
    <AuthContext.Provider value={{ user, setUser, isLoading, setIsLoading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  return useContext(AuthContext);
};
