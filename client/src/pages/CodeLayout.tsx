import React from "react";
import CodeEditor from "../components/CodeEditor";

function CodeLayout() {
  return (
    <div className="w-full h-full flex justify-between gap-2 p-5 bg-slate-950">
      {/* Code editor box */}
      <div className="w-full  rounded-md">
        <CodeEditor />
      </div>
    </div>
  );
}

export default CodeLayout;
