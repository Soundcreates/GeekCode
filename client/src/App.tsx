import React from "react";

import { BrowserRouter, Routes, Route } from "react-router";
import CodeLayout from "./pages/CodeLayout.tsx";
import { CodeProvider } from "./context/globalCode.tsx";
import Landing from "./pages/Landing.js";
import Login from "./pages/auth/Login.tsx";
import Register from "./pages/auth/Register.tsx";

function App() {
  return (
    <BrowserRouter>
      <CodeProvider>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/login" element={<Login />} />
          <Route path="/get-started" element={<Register />} />
          <Route path="/code/:room" element={<CodeLayout />} />
        </Routes>
      </CodeProvider>
    </BrowserRouter>
  );
}

export default App;
