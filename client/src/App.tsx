import React from "react";

import { BrowserRouter, Routes, Route } from "react-router";
import CodeLayout from "./pages/CodeLayout";
import { CodeProvider } from "./context/globalCode";
import Landing from "./pages/Landing.js";
import Login from "./pages/auth/Login";
import Register from "./pages/auth/Register";
import Dashboard from "./pages/Dashboard";
import Rooms from "./pages/Rooms";
import { RoomProvider } from "./context/RoomContext";
import { AuthProvider } from "./context/AuthContext";
function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <CodeProvider>
          <RoomProvider>
            <Routes>
              <Route path="/" element={<Landing />} />
              <Route path="/login" element={<Login />} />
              <Route path="/get-started" element={<Register />} />
              <Route path="/code/:roomId" element={<CodeLayout />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/rooms/:roomId" element={<Rooms />} />
            </Routes>
          </RoomProvider>

        </CodeProvider>
      </AuthProvider>

    </BrowserRouter>
  );
}

export default App;
