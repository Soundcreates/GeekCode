import { useContext, createContext, ReactNode, useState, useMemo } from "react";
import * as monaco from "monaco-editor";

export type FileType = {
  name: string;
  value: string;
  language: string;
  model: monaco.editor.ITextModel;
};

type CodeContextType = {
  files: FileType[];
  setFiles: React.Dispatch<React.SetStateAction<FileType[]>>;
  file: FileType | null;
  setFile: React.Dispatch<React.SetStateAction<FileType | null>>;
  FileTemplates: Omit<FileType, "model">[];
  AddFile: (newFile: Omit<FileType, "model">) => void;
  RemoveFile: (name: string) => void;
};

export const CodeContext = createContext<CodeContextType | undefined>(undefined);

// helper to generate Monaco model for a file
const createFileWithModel = (file: Omit<FileType, "model">): FileType => ({
  ...file,
  model: monaco.editor.createModel(file.value, file.language),
});

const FileTemplates: Omit<FileType, "model">[] = [
  {
    name: "main.js",
    value: "// Write your JavaScript code here...",
    language: "javascript",
  },
  // { name: "main.py", value: "# Write your Python code here...", language: "python" },
  // { name: "Main.java", value: "// Write your Java code here...", language: "java" },
  // { name: "main.cpp", value: "// Write your C++ code here...", language: "cpp" },
];

export const CodeProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  // Initialize files immediately with templates
  const initialFiles = useMemo(() => {
    return FileTemplates.map(createFileWithModel);
  }, []);

  const [files, setFiles] = useState<FileType[]>(initialFiles);
  const [file, setFile] = useState<FileType | null>(initialFiles[0] || null);

  const AddFile = (newFile: Omit<FileType, "model">) => {
    // prevent duplicates
    if (files.some((f) => f.name === newFile.name)) return;

    const fileWithModel = createFileWithModel(newFile);
    setFiles((prev) => [...prev, fileWithModel]);
    setFile(fileWithModel);
  };

  const RemoveFile = (name: string) => {
    setFiles((prev) => {
      const target = prev.find((f) => f.name === name);
      target?.model.dispose(); // cleanup Monaco model
      return prev.filter((f) => f.name !== name);
    });

    if (file?.name === name) {
      setFile(null);
    }
  };


  return (
    <CodeContext.Provider
      value={{ files, setFiles, file, setFile, FileTemplates, AddFile, RemoveFile }}
    >
      {children}
    </CodeContext.Provider>
  );
};

export const useCode = () => {
  const context = useContext(CodeContext);
  if (!context) {
    throw new Error("useCode must be used within a CodeProvider");
  }
  return context;
};
