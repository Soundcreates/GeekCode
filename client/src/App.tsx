import React from "react";

import { BrowserRouter, Routes, Route } from "react-router";
import CodeLayout from "./pages/CodeLayout.tsx";
import { CodeProvider } from "./context/globalCode.tsx";
import Landing from "./pages/Landing.js";
import Login from "./pages/auth/Login.tsx";
import Register from "./pages/auth/Register.tsx";
import Dashboard from "./pages/Dashboard.tsx";
import Rooms from "./pages/Rooms.tsx";
import { RoomProvider } from "./context/RoomContext.tsx";
function App() {
  return (
    <BrowserRouter>
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
    </BrowserRouter>
  );
}

export default App;
