import React, { useEffect, useRef, useState, useCallback } from "react";
import Editor from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react";
import * as monaco from "monaco-editor";
import LanguageSelector from "./LanguageSelector";
import Output from "./Output";
import ClientCounter from "./ClientCounter";
import { useCode, type FileType } from "../context/globalCode";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router";
import { Play, Terminal, X } from "lucide-react";
import GeminiAssistant from "./GeminiAssistant";
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
  change?: unknown;
  clients?: ClientInfo[];
  clientCount?: number;
  timestamp?: string;
  code?: string;
  language?: string;
  fileName?: string;
}

type CodeEditorProps = {
  roomId: string | undefined;
}

const CodeEditor: React.FC<CodeEditorProps> = ({ roomId: roomName }: CodeEditorProps) => {
  const navigate = useNavigate();
  const { file, setFile, files, setFiles } = useCode();
  const { user } = useAuth();
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // WebSocket state
  const [clients, setClients] = useState<ClientInfo[]>([]);
  const [clientCount, setClientCount] = useState(0);
  const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'disconnected' | 'error'>('disconnected');
  const [isReconnecting, setIsReconnecting] = useState(false);
  const [isTerminalOpen, setIsTerminalOpen] = useState(false);
  const [terminalHeight, setTerminalHeight] = useState(0);
  const [isGeminiMinimized, setIsGeminiMinimized] = useState(false);

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
      console.log('✅ WebSocket connected to room:', roomName);
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
        console.log('📨 WebSocket message received:', message);

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
            console.log('📢 System message:', message.change);
            break;

          case 'edit':
            console.log('✏️ Edit message received:', message.change);
            // Handle code edits from other users
            if (message.change && message.user !== currentUser.username) {
              // Update the editor with changes from other users
              // You can implement this based on your needs
            }
            break;

          default:
            console.log('❓ Unknown message action:', message.action);
        }
      } catch (error) {
        console.error('❌ Error parsing WebSocket message:', error);
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
      console.error('❌ WebSocket error:', error);
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
    if (!roomName || !currentUser) {
      console.log('Waiting for roomName or currentUser...', { roomName, currentUser });
      return;
    }

    console.log('Both roomName and currentUser are available, connecting...');
    const cleanup = connectWebSocket();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (cleanup) {
        cleanup();
      }
    };
  }, [roomName, currentUser, connectWebSocket]);

  // Send WebSocket message
  const sendWebSocketMessage = useCallback((message: WebSocketMessage) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log('📤 Sending WebSocket message:', message);
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn('⚠️ WebSocket is not connected');
    }

  }, []);

  const handleEditorDidMount: OnMount = (editor) => {
    editorRef.current = editor;
    editor.focus();
    console.log('📝 Editor mounted');
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

    console.log("✏️ Code changed by:", currentUser?.username);
  };

  const handleLanguageSelect = (language: string) => {
    if (!file) return;
    setFiles((prev: FileType[]) =>
      prev.map((f: FileType) => (f === file ? { ...f, language } : f))
    );
    setFile((prev) => (prev ? { ...prev, language } : prev));
  };

  // Move the conditional rendering AFTER all hooks
  if (!file || !currentUser) {
    return (
      <div className="h-screen w-full flex items-center justify-center bg-gray-900">
        <div className="text-center">
          <div className="text-white text-lg mb-2">Loading CodeSpace...</div>
          <div className="text-gray-400 text-sm">
            {!currentUser && "Authenticating user..."}
            {!file && currentUser && "Initializing editor..."}
          </div>
          <div className="mt-4">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          </div>
        </div>
      </div>
    );
  }

  const handleLeave = () => {
    sendWebSocketMessage({
      action: "leave",
    });
    navigate("/dashboard")
  };

  const toggleTerminal = () => {
    setIsTerminalOpen(!isTerminalOpen);
    setTerminalHeight(isTerminalOpen ? 0 : 300);
  };

  const runCode = () => {
    // Send run code message via WebSocket
    sendWebSocketMessage({
      action: 'run_code',
      room: roomName,
      user: currentUser?.username,
      userId: currentUser?.id?.toString() || currentUser?.username,
      code: file.value,
      language: file.language,
      fileName: file.name || 'main'
    });

    // Open terminal to show output
    if (!isTerminalOpen) {
      toggleTerminal();
    }
  };

  return (
    <div className="h-screen w-full flex flex-col bg-gray-900">
      {/* Top Bar */}
      <div className="flex items-center justify-between p-2 bg-gray-800 border-b border-gray-700">
        <div className="flex items-center gap-4">
          <div className="flex gap-2">
            {files.map((f: FileType, idx: number) => (
              <button
                key={idx}
                className={`px-3 py-1 rounded-t text-sm ${f === file
                  ? "bg-gray-900 text-white border-b-2 border-blue-500"
                  : "bg-gray-700 text-gray-300 hover:bg-gray-600"
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
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={runCode}
            className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-md transition-colors"
          >
            <Play className="w-4 h-4" />
            Run
          </button>

          <button
            onClick={toggleTerminal}
            className={`flex items-center gap-2 px-4 py-2 rounded-md transition-colors ${isTerminalOpen
              ? 'bg-blue-600 hover:bg-blue-700 text-white'
              : 'bg-gray-700 hover:bg-gray-600 text-gray-300'
              }`}
          >
            <Terminal className="w-4 h-4" />
            Terminal
          </button>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex flex-1 overflow-hidden">
        {/* Editor Area */}
        <div className="flex-1 flex flex-col">
          <Editor
            className="flex-1"
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
              scrollBeyondLastLine: false,
              renderLineHighlight: 'line',
              cursorBlinking: 'blink',
              cursorStyle: 'line',
            }}
          />
        </div>

        {/* Right Sidebar */}
        <div className="w-80 bg-gray-800 border-l border-gray-700 flex flex-col">
          {/* Client Counter */}
          <div className="p-4 border-b border-gray-700">
            <ClientCounter
              clientCount={clientCount}
              clients={clients}
              connectionStatus={isReconnecting ? 'connecting' : connectionStatus}
            />
          </div>

          {/* Room Info */}
          <div className="p-4 border-b border-gray-700">
            <h4 className="text-white text-sm font-semibold mb-3">Room Info</h4>
            <div className="text-xs text-gray-400 space-y-2">
              <div className="flex justify-between">
                <span>Room ID:</span>
                <span className="text-gray-300">{roomName}</span>
              </div>
              <div className="flex justify-between">
                <span>Username:</span>
                <span className="text-gray-300">{currentUser?.username}</span>
              </div>
              <div className="flex justify-between items-center">
                <span>Status:</span>
                <span className={`px-2 py-1 rounded text-xs ${connectionStatus === 'connected' ? 'bg-green-600 text-white' :
                  connectionStatus === 'connecting' ? 'bg-yellow-600 text-white' :
                    'bg-red-600 text-white'
                  }`}>
                  {connectionStatus}
                </span>
              </div>
            </div>
          </div>

          {/* AI Assistant */}
          <div className="p-4 border-b border-gray-700">
            <GeminiAssistant
              isMinimized={isGeminiMinimized}
              onToggleMinimize={() => setIsGeminiMinimized(!isGeminiMinimized)}
            />
          </div>

          {/* Actions */}
          <div className="p-4 mt-auto">
            <button
              onClick={() => sendWebSocketMessage({ action: 'get_room_info', room: roomName })}
              className="w-full bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 text-sm rounded mb-2 transition-colors"
            >
              Refresh Room Info
            </button>
            <button
              onClick={handleLeave}
              className="w-full bg-red-600 hover:bg-red-700 text-white px-3 py-2 text-sm rounded transition-colors"
            >
              Leave Room
            </button>
          </div>
        </div>
      </div>

      {/* Slide-up Terminal */}
      <div
        className="bg-gray-900 border-t border-gray-700 transition-all duration-300 ease-in-out overflow-hidden"
        style={{ height: `${terminalHeight}px` }}
      >
        <div className="h-full flex flex-col">
          <div className="flex items-center justify-between p-2 bg-gray-800 border-b border-gray-700">
            <div className="flex items-center gap-2">
              <Terminal className="w-4 h-4 text-gray-400" />
              <span className="text-sm text-gray-300">Terminal</span>
            </div>
            <button
              onClick={toggleTerminal}
              className="text-gray-400 hover:text-white transition-colors"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
          <div className="flex-1 p-4">
            <Output value={file.value} language={file.language} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default CodeEditor;