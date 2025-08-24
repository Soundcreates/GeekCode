import { useContext, createContext, ReactNode, useEffect } from "react";
import { useState } from "react";
import * as monaco from "monaco-editor";

export type FileType = {
  name: string;
  value: string;
  language: string;
  model: monaco.editor.ITextModel | null;
};

type CodeContextType = {
  code: string;
  setCode: React.Dispatch<React.SetStateAction<string>>;
  files: FileType[];
  setFiles: React.Dispatch<React.SetStateAction<FileType[]>>;
  file: FileType | null;
  setFile: React.Dispatch<React.SetStateAction<FileType | null>>;
  FileTemplates: Omit<FileType, "model">[];
  AddFile: (newFile: Omit<FileType, "model">) => void;
};

export const CodeContext = createContext<CodeContextType | undefined>(
  undefined
);

export const CodeProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [code, setCode] = useState("");
  const [files, setFiles] = useState<FileType[]>([]);
  const [file, setFile] = useState<FileType | null>(null);

  const FileTemplates: Omit<FileType, "model">[] = [
    {
      name: "main.js",
      value: "// Write your JavaScript code here...",
      language: "javascript",
    },
    // { value: "# Write your Python code here...", language: "python" },
    // { value: "// Write your Java code here...", language: "java" },
    // { value: "// Write your C++ code here...", language: "cpp" },
  ];

  const AddFile = (newFile: Omit<FileType, "model">) => {
    //this ensures that every file has its own model
    const model = monaco.editor.createModel(newFile.value, newFile.language); //createmodel only needs value and language (Sourcecode, and language)
    const fileWithModel: FileType = { ...newFile, model }; //basically adding the model to the file
    setFiles((prevFiles) => [...prevFiles, fileWithModel]); //so basically what it does it it doesnt write over the existing files, it just appends
    setFile(fileWithModel); //to switch to the new file and make it active, we use 'filewithmodel' because it has all the fields in one object so we can just destructure it when actually using it in the editor component
  };

  useEffect(() => {
    const models = FileTemplates.map((tpl) => ({
      ...tpl,
      model: monaco.editor.createModel(tpl.value, tpl.language), //this means for every obj in the filetemplate obj u will create a new model wrt to the language and the value(sourcecode)
    }));

    setFiles(models);
    setFile(models[0]);

    //cleanup func
    return () => {
      models.forEach((f) => f.model.dispose());
    };
  }, []);

  return (
    <CodeContext.Provider
      value={{
        code,
        setCode,
        files,
        setFiles,
        file,
        setFile,
        FileTemplates,
        AddFile,
      }}
    >
      {children}
    </CodeContext.Provider>
  );
};

export const useCode = () => {
  const content = useContext(CodeContext);
  if (!content) return null;
  return content;
};
