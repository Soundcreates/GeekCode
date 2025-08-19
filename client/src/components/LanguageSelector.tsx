import React from "react";
import { LANGUAGE_OPTIONS } from "../utils/constants";

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
};

function LanguageSelector({ language, onSelect }: Props) {
  return (
    <div className="p-4 bg-gray-800 text-white rounded-md">
      <h2 className="text-lg font-semibold mb-2">{language}</h2>
      <select
        onChange={(e) => onSelect(e.target.value)}
        className="bg-gray-700 text-white p-2 rounded-md"
      >
        {languages.map((lang) => (
          <option key={lang} value={lang}>
            {DISPLAY_NAMES[lang] || lang}
          </option>
        ))}
      </select>
    </div>
  );
}

export default LanguageSelector;
