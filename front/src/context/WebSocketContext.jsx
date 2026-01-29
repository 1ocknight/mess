import React, { createContext, useContext, useState, useCallback } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

const WSContext = createContext();

export const WSProvider = ({ token, children }) => {
  const [messages, setMessages] = useState([]);

  const wsUrl = token
    ? `ws://localhost:8082/ws?token=${encodeURIComponent(token)}`
    : null;

  const {
    sendMessage,
    lastMessage,
    readyState,
  } = useWebSocket(wsUrl, {
    onOpen: () => console.log('WS connected'),
    onMessage: (event) => {
      const raw = event.data;
      const msgs = raw
        .split('\n')
        .map(str => str.trim())
        .filter(Boolean);

      msgs.forEach(msgStr => {
        try {
          const msg = JSON.parse(msgStr);
          setMessages(prev => [...prev, msg]);
        } catch (e) {
          console.warn('WS ignored non-JSON message:', msgStr);
        }
      });
    },
    onError: (err) => {
      console.error('WS error:', err);
    },
    shouldReconnect: (closeEvent) => true, // всегда переподключаемся
    reconnectAttempts: 10, // макс попыток
    reconnectInterval: 3000, // пауза между попытками (мс)
  });

  const addListener = useCallback((cb) => {
    if (!lastMessage) return () => {};
    cb(JSON.parse(lastMessage.data));
    return () => {};
  }, [lastMessage]);

  const connected = readyState === ReadyState.OPEN;

  return (
    <WSContext.Provider value={{ connected, messages, addListener, sendMessage }}>
      {children}
    </WSContext.Provider>
  );
};

export const useWS = () => useContext(WSContext);
