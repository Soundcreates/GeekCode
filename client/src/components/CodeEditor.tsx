import React, { useEffect, useMemo, useRef } from "react";
import Editor from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react";
import LanguageSelector from "./LanguageSelector";
import Output from "./Output";
import { useCode, type FileType } from "../context/globalCode";
import * as Y from "yjs";
import { WebsocketProvider } from "y-websocket";
import { MonacoBinding } from "y-monaco";

type CodeEditorProps = {
  roomId: string | undefined;
}

const CodeEditor: React.FC<CodeEditorProps> = ({ roomId: roomName }: CodeEditorProps) => {
  const { file, setFile, files, setFiles } = useCode();
  const editorRef = useRef<any | null>(null);
  const bindingRef = useRef<MonacoBinding | null>(null);

  if (!file) return <div>Loading...</div>;

  const ydoc = useMemo(() => new Y.Doc(), []);

  useEffect(() => {
    if (!editorRef.current || !roomName) return;

    const provider = new WebsocketProvider(
      `ws://localhost:8080/api/ws/${roomName}`,
      `room-${roomName}`,
      ydoc
    );

    const ytext = ydoc.getText("monaco");

    const binding = new MonacoBinding(
      ytext,
      editorRef.current.getModel(),
      new Set([editorRef.current]),
      provider?.awareness
    );

    bindingRef.current = binding;

    return () => {
      binding.destroy();
      provider.destroy();
    };
  }, [roomName, ydoc]);

  const handleEditorDidMount: OnMount = (editor) => {
    editorRef.current = editor;
    editor.focus();
  };

  const handleTabClick = (f: FileType) => {
    setFile(f);
  };

  const handleCodeChange = (value: string | undefined) => {
    if (!file) return;
    const newValue = value ?? "";

    setFiles((prev: FileType[]) =>
      prev.map((f: FileType) => (f === file ? { ...f, value: newValue } : f))
    );
    setFile((prev) => (prev ? { ...prev, value: newValue } : prev));
  };

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
        <div className="flex gap-2 mb-2">
          {files.map((f: FileType, idx: number) => (
            <button
              key={idx}
              className={`px-3 py-1 rounded-t ${f === file
                ? "bg-blue-600 text-white"
                : "bg-gray-700 text-gray-300"
                }`}
              onClick={() => handleTabClick(f)}
            >
              {f.name || `File ${idx + 1}`}
            </button>
          ))}
        </div>

        <LanguageSelector
          language={file.language}
          onSelect={handleLanguageSelect}
          setLanguage={() => { }}
        />

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

      <div className="w-1/2">
        <Output value={file.value} language={file.language} />
      </div>
    </div>
  );
};

export default CodeEditor;