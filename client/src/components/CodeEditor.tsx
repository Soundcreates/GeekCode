import React, { useEffect, useMemo, useRef, useState, useCallback } from "react";
import Editor from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react";
import LanguageSelector from "./LanguageSelector";
import Output from "./Output";
import ClientCounter from "./ClientCounter";
import { useCode, type FileType } from "../context/globalCode";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router";
interface ClientInfo {
  user: string;
  userId: string;
  joinedAt: string;
  isOnline: boolean;
}

interface WebSocketMessage {
  action: string;
  room?: string;
  user?: string;
  userId?: string;
  change?: any;
  clients?: ClientInfo[];
  clientCount?: number;
  timestamp?: string;
}

type CodeEditorProps = {
  roomId: string | undefined;
}

const CodeEditor: React.FC<CodeEditorProps> = ({ roomId: roomName }: CodeEditorProps) => {
  const navigate = useNavigate();
  const { file, setFile, files, setFiles } = useCode();
  const { user } = useAuth();
  const editorRef = useRef<any | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // WebSocket state
  const [clients, setClients] = useState<ClientInfo[]>([]);
  const [clientCount, setClientCount] = useState(0);
  const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'disconnected' | 'error'>('disconnected');
  const [isReconnecting, setIsReconnecting] = useState(false);

  const currentUser = user;

  // Single WebSocket connection function
  const connectWebSocket = useCallback(() => {
    if (!roomName || !currentUser) {
      console.log('Missing roomName or currentUser');
      return;
    }

    console.log('Attempting to connect to WebSocket for room:', roomName);
    setConnectionStatus('connecting');
    setIsReconnecting(false);

    const ws = new WebSocket(`ws://localhost:8080/api/ws/${roomName}`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log(' WebSocket connected to room:', roomName);
      setConnectionStatus('connected');

      // Join the room
      const joinMessage: WebSocketMessage = {
        action: "join",
        room: roomName,
        user: currentUser.username,
        userId: currentUser.id?.toString() || currentUser.username
      };

      console.log('Sending join message:', joinMessage);
      ws.send(JSON.stringify(joinMessage));

      // Request current room info
      setTimeout(() => {
        const roomInfoMessage: WebSocketMessage = {
          action: "get_room_info",
          room: roomName
        };
        console.log('Requesting room info');
        ws.send(JSON.stringify(roomInfoMessage));
      }, 100);
    };

    ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        console.log(' WebSocket message received:', message);

        switch (message.action) {
          case 'room_info':
            if (message.clients) {
              console.log('Room info received:', message.clients);
              setClients(message.clients);
              setClientCount(message.clientCount || message.clients.length);
            }
            break;

          case 'room_update':
            if (message.clients) {
              console.log('Room update received:', message.clients);
              setClients(message.clients);
              setClientCount(message.clientCount || message.clients.length);
            }
            break;

          case 'system':
            console.log(' System message:', message.change);
            break;

          case 'edit':
            console.log(' Edit message received:', message.change);
            // Handle code edits from other users
            if (message.change && message.user !== currentUser.username) {
              // Update the editor with changes from other users
              // You can implement this based on your needs
            }
            break;

          default:
            console.log(' Unknown message action:', message.action);
        }
      } catch (error) {
        console.error(' Error parsing WebSocket message:', error);
      }
    };

    ws.onclose = (event) => {
      console.log('🔌 WebSocket connection closed:', event.code, event.reason);
      setConnectionStatus('disconnected');
      setClients([]);
      setClientCount(0);

      // Auto-reconnect if not intentionally closed
      if (event.code !== 1000 && roomName) {
        setIsReconnecting(true);
        reconnectTimeoutRef.current = setTimeout(() => {
          console.log('🔄 Attempting to reconnect...');
          connectWebSocket();
        }, 3000);
      }
    };

    ws.onerror = (error) => {
      console.error(' WebSocket error:', error);
      setConnectionStatus('error');
    };

    // Cleanup function
    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        const leaveMessage: WebSocketMessage = {
          action: "leave",
          room: roomName
        };
        ws.send(JSON.stringify(leaveMessage));
      }
      ws.close(1000);
    };
  }, [roomName, currentUser]);

  // WebSocket connection effect
  useEffect(() => {
    if (!roomName || !currentUser) return;

    const cleanup = connectWebSocket();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (cleanup) {
        cleanup();
      }
    };
  }, [connectWebSocket]);

  // Send WebSocket message
  const sendWebSocketMessage = useCallback((message: WebSocketMessage) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log(' Sending WebSocket message:', message);
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn(' WebSocket is not connected');
    }
  }, []);

  const handleEditorDidMount: OnMount = (editor) => {
    editorRef.current = editor;
    editor.focus();
    console.log(' Editor mounted');
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

    // Send code change notification
    sendWebSocketMessage({
      action: 'edit',
      room: roomName,
      user: currentUser?.username,
      userId: currentUser?.id?.toString() || currentUser?.username,
      change: {
        code: newValue,
        language: file.language,
        timestamp: new Date().toISOString()
      }
    });

    console.log(" Code changed by:", currentUser?.username);
  };

  const handleLanguageSelect = (language: string) => {
    if (!file) return;
    setFiles((prev: FileType[]) =>
      prev.map((f: FileType) => (f === file ? { ...f, language } : f))
    );
    setFile((prev) => (prev ? { ...prev, language } : prev));
  };

  // Move the conditional rendering AFTER all hooks
  if (!file) {
    return (
      <div className="h-screen w-full flex items-center justify-center">
        <div className="text-white">Loading...</div>
      </div>
    );
  }

  const handleLeave = () => {
    sendWebSocketMessage({
      action: "leave",
    });
    navigate("/dashboard")
  };
  return (
    <div className="h-screen w-full flex flex-col p-4 gap-4">
      {/* Top Row - Editor and Client Info */}
      <div className="flex gap-4 h-3/5">
        {/* Left Side - Editor */}
        <div className="w-3/4 flex flex-col">
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
            className="rounded-md mt-1 flex-1"
            height="100%"
            width="100%"
            value={file.value}
            language={file.language}
            theme="vs-dark"
            onMount={handleEditorDidMount}
            onChange={handleCodeChange}
            options={{
              tabSize: 4,
              automaticLayout: true,
              minimap: { enabled: false },
              fontSize: 14,
              wordWrap: 'on',
            }}
          />
        </div>

        {/* Right Side - Client Info */}
        <div className="w-1/4 flex flex-col">
          {/* Client Counter */}
          <ClientCounter
            clientCount={clientCount}
            clients={clients}
            connectionStatus={isReconnecting ? 'connecting' : connectionStatus}
          />

          {/* Room Info */}
          <div className="bg-gray-800 p-3 rounded-lg border border-gray-700 mb-4">
            <h4 className="text-white text-sm font-semibold mb-2">Room Info</h4>
            <div className="text-xs text-gray-400 space-y-1">
              <div>Room ID: {roomName}</div>
              <div>Your Username: {currentUser?.username}</div>
              <div>Connection: <span className={`${connectionStatus === 'connected' ? 'text-green-400' :
                connectionStatus === 'connecting' ? 'text-yellow-400' : 'text-red-400'
                }`}>{connectionStatus}</span></div>
            </div>
          </div>

          {/* WebSocket Test Panel */}
          <div className="bg-gray-800 p-3 rounded-lg border border-gray-700">
            <h4 className="text-white text-sm font-semibold mb-2">WebSocket Test</h4>
            <button
              onClick={() => sendWebSocketMessage({ action: 'get_room_info', room: roomName })}
              className="bg-blue-600 text-white px-2 py-1 text-xs rounded w-full"
            >
              Refresh Room Info
            </button>
            <button onClick={handleLeave} className="bg-red-500 text-white px-2 py-1 text-s rounded-md mt-2 cursor-pointer w-full hover:bg-red-700 transition-all duration-300">
              Leave
            </button>
          </div>
        </div>
      </div>

      {/* Bottom Row - Output Terminal */}
      <div className="h-2/5">
        <Output value={file.value} language={file.language} />
      </div>
    </div>
  );
};

export default CodeEditor;