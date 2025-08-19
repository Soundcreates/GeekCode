import React from "react";

import { BrowserRouter, Routes, Route } from "react-router";
import CodeLayout from "./pages/CodeLayout.tsx";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<CodeLayout />} />
        
      </Routes>
    </BrowserRouter>
  );
}

export default App;
