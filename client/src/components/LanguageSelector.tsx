import React, { useState } from "react";
import { LANGUAGE_OPTIONS } from "../utils/constants";
import { useCode } from "../context/globalCode";

const languages = Object.keys(LANGUAGE_OPTIONS);

const DISPLAY_NAMES: Record<string, string> = {
  javascript: "JavaScript",
  python: "Python",
  java: "Java",
  cpp: "C++",
};

type Props = {
  language: string;
  onSelect: (lang: string) => void;
  setLanguage: React.Dispatch<React.SetStateAction<string>>;
};

const FILE_MAP: Record<string, string> = {
  ".js": "javascript",
  ".cpp": "cpp",
  ".py": "python",
  ".java": "java",
};

function LanguageSelector({ language, onSelect, setLanguage }: Props) {
  const [fileName, setFileName] = useState<string>("");

  const codeCtx = useCode();
  if (!codeCtx) return null;
  const { AddFile } = codeCtx;

  const templates: Record<string, string> = {
    javascript: "// Write your JavaScript code here...",
    python: "# Write your Python code here...",
    java: "// Write your Java code here...",
    cpp: "// Write your C++ code here...",
    plaintext: "",
  };

  const handleAddFile = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const name = fileName.trim();
    if (!name) return;

    const ext = Object.keys(FILE_MAP).find((ext) => name.endsWith(ext));
    if (!ext) {
      alert("Unsupported File extension! Use .js, .py, .java or .cpp");
      return;
    }

    const detectedLanguage = FILE_MAP[ext];
    const content = templates[detectedLanguage] ?? templates.plaintext;

    AddFile({
      name,
      value: content,
      language: detectedLanguage,
    });

    setFileName("");
  };

  return (
    <div className="bg-slate-900 text-white mt-5 p-3 rounded">
      {/* <div className="mb-3">
        <label className="block text-sm text-gray-300 mb-1">Language</label>
        <select
          value={language}
          onChange={(e) => onSelect(e.target.value)}
          className="bg-gray-700 text-white p-2 rounded-md w-full"
        >
          {languages.map((lang) => (
            <option key={lang} value={lang}>
              {DISPLAY_NAMES[lang] ?? lang}
            </option>
          ))}
        </select>
      </div> */}

      <form onSubmit={handleAddFile} className="flex gap-2">
        <input
          type="text"
          value={fileName}
          onChange={(e) => setFileName(e.target.value)}
          placeholder="filename.ext (e.g. index.js)"
          className="flex-1 p-2 rounded bg-gray-800 text-white"
        />
        <button
          type="submit"
          className="px-3 py-2 bg-blue-600 text-white rounded"
        >
          Add
        </button>
      </form>
    </div>
  );
}

export default LanguageSelector;
