import React, { useEffect, useRef, useState } from "react";
import { useParams } from "react-router";
import { Editor } from "@monaco-editor/react";

const WS_BASE = "ws://localhost:8080/api/ws/";

const Rooms: React.FC = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const ws = useRef<WebSocket | null>(null);
  const [code, setCode] = useState("");
  const [connected, setConnected] = useState(false);
  const [language, setLanguage] = useState("javascript");
  const editorRef = useRef<any>(null);

  useEffect(() => {
    if (!roomId) return;

    // Connect to WebSocket
    ws.current = new WebSocket(`${WS_BASE}${roomId}`);

    ws.current.onopen = () => {
      setConnected(true);
      console.log("Connected to room:", roomId);
    };

    ws.current.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (message.type === "code_change") {
          setCode(message.payload.code);
          // Update editor without triggering onChange
          if (editorRef.current) {
            editorRef.current.setValue(message.payload.code);
          }
        }
      } catch (error) {
        // Handle plain text messages
        setCode(event.data);
      }
    };

    ws.current.onclose = () => {
      setConnected(false);
      console.log("Disconnected from room");
    };

    return () => {
      ws.current?.close();
    };
  }, [roomId]);

  const handleEditorChange = (value: string | undefined) => {
    if (value !== undefined) {
      setCode(value);
      // Send code changes to other users
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        const message = {
          type: "code_change",
          payload: { code: value, language },
        };
        ws.current.send(JSON.stringify(message));
      }
    }
  };

  const handleEditorDidMount = (editor: any) => {
    editorRef.current = editor;
  };

  return (
    <div className="flex flex-col h-screen bg-gray-900 text-white">
      {/* Header */}
      <div className="flex items-center justify-between p-4 bg-gray-800 border-b border-gray-700">
        <div className="flex items-center gap-4">
          <h1 className="text-xl font-bold">Room: {roomId}</h1>
          <div className="flex items-center gap-2">
            <div
              className={`w-3 h-3 rounded-full ${
                connected ? "bg-green-500" : "bg-red-500"
              }`}
            />
            <span className="text-sm">
              {connected ? "Connected" : "Disconnected"}
            </span>
          </div>
        </div>

        <select
          value={language}
          onChange={(e) => setLanguage(e.target.value)}
          className="bg-gray-700 text-white px-3 py-1 rounded"
        >
          <option value="javascript">JavaScript</option>
          <option value="python">Python</option>
          <option value="java">Java</option>
          <option value="cpp">C++</option>
          <option value="go">Go</option>
        </select>
      </div>

      {/* Code Editor */}
      <div className="flex-1">
        <Editor
          height="100%"
          language={language}
          value={code}
          onChange={handleEditorChange}
          onMount={handleEditorDidMount}
          theme="vs-dark"
          options={{
            fontSize: 14,
            wordWrap: "on",
            minimap: { enabled: false },
            scrollBeyondLastLine: false,
          }}
        />
      </div>
    </div>
  );
};

export default Rooms;
