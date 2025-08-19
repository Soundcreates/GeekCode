import React, { useRef, useState } from "react";
import Editor from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react";
import LanguageSelector from "./LanguageSelector";
import Output from "./Output";

const CodeEditor: React.FC = () => {
  const [value, setValue] = useState<string>("");
  const [language, setLanguage] = useState<string>("javascript");

  const editorRef = useRef<any>(null);

  const handleEditorDidMount: OnMount = (editor) => {
    editorRef.current = editor;
    editor.focus();

    console.log("Monaco Instance: ", editor);
  };

  const onSelect = (language) => {
    setLanguage(language);
  };

  return (
    <div className="h-screen w-full rounded-md flex justify-between gap-2  ">
      <div className="w-1/2">
        <LanguageSelector language={language} onSelect={onSelect} />
        <Editor
          className="rounded-md mt-1 "
          height="75vh"
          width="100%"
          language={language}
          value={`${
            language === "python" ? "#" : "//"
          } Write your ${language} code here...`}
          theme="vs-dark"
          onMount={handleEditorDidMount}
          options={{
            tabSize: 4,
          }}
          onChange={() => setValue(editorRef.current?.getValue())}
        />
      </div>

      {/* Output */}
      <div className="w-1/2">
        <Output editorRef={editorRef} language={language} />
      </div>
    </div>
  );
};
export default CodeEditor;
