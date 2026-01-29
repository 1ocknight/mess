import React, { useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { getChats } from '../api/chat';
import { useWS } from '../context/WebSocketContext';
import ChatListItem from './ChatListItem';

export default function ChatsList({ token, pageSize = 20 }) {
  const [chats, setChats] = useState([]);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  const afterRef = useRef(null);
  const containerRef = useRef(null);

  const { messages } = useWS();
  const navigate = useNavigate(); // <-- для перехода

  // ---------- Загрузка чатов ----------
  const fetchChats = async (after = null) => {
    if (!token || loading || !hasMore) return;
    setLoading(true);

    try {
      const params = { limit: pageSize };
      if (after) params.after = after;

      const newChats = await getChats(token, params);
      if (!Array.isArray(newChats)) return;

      setChats(prev => {
        const ids = new Set(prev.map(c => c.chat_id));
        return [...prev, ...newChats.filter(c => !ids.has(c.chat_id))];
      });

      if (newChats.length > 0) {
        afterRef.current = newChats[newChats.length - 1].chat_id;
      }

      if (newChats.length < pageSize) setHasMore(false);

    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchChats();
  }, [token]);

  // ---------- Скролл ----------
  const handleScroll = () => {
    if (!containerRef.current || loading || !hasMore) return;
    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    if (scrollTop + clientHeight >= scrollHeight - 50) {
      fetchChats(afterRef.current);
    }
  };

  // ---------- WS подписка ----------
  useEffect(() => {
    messages.forEach(msg => {
      const chatMsg = msg.data;

      setChats(prev => {
        const map = new Map(prev.map(c => [c.chat_id, c]));
        const existing = map.get(chatMsg.chat_id);

        if (existing) {
          const currentLast = existing.last_message;

          if (msg.type == 'update_last_read' && chatMsg.message_id === currentLast.id ){
            map.set(chatMsg.chat_id, {
                ...existing,
                is_last_message_read: true,
              });
          }

          if (msg.type === 'send_message') {
            if (!currentLast || chatMsg.id > currentLast.id) {
              map.set(chatMsg.chat_id, {
                ...existing,
                last_message: chatMsg,
                unread_count: (existing.unread_count || 0) + 1,
              });
            } else if (chatMsg.id === currentLast.id) {
              map.set(chatMsg.chat_id, {
                ...existing,
                last_message: chatMsg,
              });
            }
          }

          if (msg.type === 'update_message' && currentLast?.id === chatMsg.id) {
            map.set(chatMsg.chat_id, {
              ...existing,
              last_message: chatMsg,
            });
          }
        } else {
          map.set(chatMsg.chat_id, {
            chat_id: chatMsg.chat_id,
            second_subject_id: chatMsg.sender_id,
            last_message: chatMsg,
            unread_count: msg.type === 'send_message' ? 1 : 0,
          });
        }

        return Array.from(map.values()).sort(
          (a, b) =>
            new Date(b.last_message?.createdAt) -
            new Date(a.last_message?.createdAt)
        );
      });
    });
  }, [messages]);

  return (
    <div
      ref={containerRef}
      onScroll={handleScroll}
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: 8,
        overflowY: 'auto',
        height: '100%',
      }}
    >
      {chats.map(chat => (
        <ChatListItem
          key={chat.chat_id}
          chat={chat}
          token={token}
          onClick={() => navigate(`/chat/${chat.chat_id}`)} // <-- переход на страницу чата
        />
      ))}
    </div>
  );
}
