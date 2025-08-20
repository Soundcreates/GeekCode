import React, { useState, useEffect } from "react";
import { executeCode, getLanguages } from "../services/api";

interface OutputProps {
  value: string;
  language: string;
}
function Output({ value, language }: OutputProps) {
  const [output, setOutput] = useState<string | null>(null);

  useEffect(() => {
    getLanguages();
  }, []);

  const runCode = async () => {
    setOutput(null);
    const sourceCode = value;
    if (!sourceCode) return;
    try {
      const { run: result } = await executeCode(language, sourceCode);
      setOutput(result.output);
    } catch (err) {
      console.error("Error at output.tsx in runcode function: ", err.message);
    }
  };
  return (
    <div className="w-full h-[90vh] bg-gray-800 text-white p-4 overflow-auto">
      {/* Output content will go here */}
      <button
        onClick={runCode}
        className="cursor-pointer p-2 bg-blue-500 text-white rounded mt-2"
      >
        Run code
      </button>
      <div className="h-[79vh] bg-slate-950 rounded-md w-full mt-2 text-white py-2 shadow-md">
        <p className="p-5">
          {output ? output : "Run code to see the output here"}
        </p>
      </div>
    </div>
  );
}

export default Output;
