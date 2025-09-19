import React from "react";
import CodeEditor from "../components/CodeEditor";
import { useParams } from "react-router";

function CodeLayout() {
  const { roomId } = useParams<{ roomId: string }>();
  return (
    <div className="w-full h-full flex justify-between gap-2 p-5 bg-slate-950">
      {/* Code editor box */}
      <div className="w-full  rounded-md">
        <CodeEditor roomId={roomId} />
      </div>
    </div>
  );
}

export default CodeLayout;
