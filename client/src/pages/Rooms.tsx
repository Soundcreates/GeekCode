import React, { useState, useEffect, useMemo } from 'react'
import * as Y from "yjs";
import { WebsocketProvider } from "y-websocket";
import { MonacoBinding } from "y-monaco";

import { Editor } from '@monaco-editor/react';
import { useParams } from 'react-router';

function Rooms() {

  const { roomId: roomName } = useParams().

  const ydoc = useMemo(() => new Y.Doc(), []);
  const [editor, setEditor] = useState<any | null>(null);

  useEffect(() => {
    const provider = new WebsocketProvider()
  })
  return (
    <div>Rooms</div>
  )
}

export default Rooms