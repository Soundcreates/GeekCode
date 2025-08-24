import React, { useRef } from "react";
import Editor from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react";
import LanguageSelector from "./LanguageSelector";
import Output from "./Output";
import { useCode, type FileType } from "../context/globalCode";

const CodeEditor: React.FC = () => {
  const { file, setFile, files, setFiles } = useCode();
  const editorRef = useRef<any | null>(null);

  if (!file) return <div>Loading...</div>;

  // Save editor instance on mount
  const handleEditorDidMount: OnMount = (editor) => {
    editorRef.current = editor;
    editor.focus();
  };

  // Switch active file (do not call editor.setModel here)
  const handleTabClick = (f: FileType) => {
    setFile(f);
  };

  // Update code for the active file (controlled editor)
  const handleCodeChange = (value: string | undefined) => {
    if (!file) return;
    const newValue = value ?? "";

    setFiles((prev: FileType[]) =>
      prev.map((f: FileType) => (f === file ? { ...f, value: newValue } : f))
    );
    setFile((prev) => (prev ? { ...prev, value: newValue } : prev));
  };

  // Update language for the active file (update context only)
  const handleLanguageSelect = (language: string) => {
    if (!file) return;
    setFiles((prev: FileType[]) =>
      prev.map((f: FileType) => (f === file ? { ...f, language } : f))
    );
    setFile((prev) => (prev ? { ...prev, language } : prev));
  };

  return (
    <div className="h-screen w-full rounded-md flex justify-between gap-2">
      <div className="w-1/2">
        {/* Tabs for each file */}
        <div className="flex gap-2 mb-2">
          {files.map((f: FileType, idx: number) => (
            <button
              key={idx}
              className={`px-3 py-1 rounded-t ${
                f === file
                  ? "bg-blue-600 text-white"
                  : "bg-gray-700 text-gray-300"
              }`}
              onClick={() => handleTabClick(f)}
            >
              {f.name || `File ${idx + 1}`}
            </button>
          ))}
        </div>

        {/* Language Selector */}
        <LanguageSelector
          language={file.language}
          onSelect={handleLanguageSelect}
          setLanguage={() => {}}
        />

        {/* Monaco Editor (controlled) */}
        <Editor
          className="rounded-md mt-1"
          height="75vh"
          width="100%"
          value={file.value}
          language={file.language}
          theme="vs-dark"
          onMount={handleEditorDidMount}
          onChange={handleCodeChange}
          options={{
            tabSize: 4,
            automaticLayout: true,
          }}
        />
      </div>

      {/* Output */}
      <div className="w-1/2">
        <Output value={file.value} language={file.language} />
      </div>
    </div>
  );
};

export default CodeEditor;
