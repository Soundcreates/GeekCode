import React from "react";
import { useAuth } from "../context/AuthContext";

function Dashboard() {
  const { user } = useAuth();

  return <div>{user ? `Welcome, ${user.firstname}` : "Please log in"}</div>;
}

export default Dashboard;
