import React, { useEffect, useRef, useState, useCallback } from 'react';
import { openChatByID, getMessages, addMessage } from '../api/chat';
import { getProfileById } from '../api/profile';
import { useWS } from '../context/WebSocketContext';
import MessageBubble from './MessageBubble';

const PAGE_SIZE = 20;

export default function ChatWindow({ chatId, token, myProfile }) {
  const [chat, setChat] = useState(null);
  const [otherUser, setOtherUser] = useState(null);
  const [otherUserLastRead, setOtherUserLastRead] = useState(0);
  const [messages, setMessages] = useState([]);
  const [loadingUp, setLoadingUp] = useState(false);
  const [hasMoreUp, setHasMoreUp] = useState(true);
  const [hasMoreDown, setHasMoreDown] = useState(false);
  const [text, setText] = useState('');
  const [isDown, setIsDown] = useState(true);

  const containerRef = useRef(null);
  const bottomRef = useRef(null);
  const isUserAtBottomRef = useRef(true);
  const lastProcessedWsRef = useRef(0);

  const { messages: wsMessages } = useWS();

  // ---------- УТИЛИТА ОБНОВЛЕНИЯ (ЕДИНСТВЕННЫЙ ИСТОЧНИК ПРАВДЫ) ----------
  const upsertMessages = useCallback((newMsgs) => {
    setMessages(prev => {
      const map = new Map(prev.map(m => [m.id, m]));
      newMsgs.forEach(msg => {
        map.set(msg.id, msg);
      });
      return Array.from(map.values()).sort(
        (a, b) => new Date(a.created_at) - new Date(b.created_at)
      );
    });
  }, []);

  // ---------- INITIAL LOAD ----------
  useEffect(() => {
    let isMounted = true;
    const fetchChat = async () => {
      try {
        const fullChat = await openChatByID(chatId, token, PAGE_SIZE);
        if (!isMounted) return;

        setChat(fullChat);

        // Загружаем профиль другого юзера
        if (fullChat.second_subject_id) {
          try {
            const otherUserProfile = await getProfileById(token, fullChat.second_subject_id);
            if (isMounted) {
              setOtherUser(otherUserProfile);
            }
          } catch (err) {
            console.error('Failed to load other user profile', err);
          }
        }

        // Инициализируем last_read другого юзера
        const otherSubjectId = fullChat.second_subject_id;
        const otherLastReadId = fullChat.last_reads?.[otherSubjectId] || 0;
        if (isMounted) {
          setOtherUserLastRead(otherLastReadId);
        }

        // Вычисляем last_read для текущего пользователя
        const lastReadId = fullChat.last_reads?.[myProfile.subject_id] || 0;

        let finalMessages = [];
        let hasMoreUpFlag = false;
        let hasMoreDownFlag = false;
        let isDownFlag = false;

        try {
          if (lastReadId) {
            // Есть last_read: сначала загружаем сообщения ПОСЛЕ него (after)
            const newerMessages = await getMessages(token, {
              chat_id: fullChat.chat_id,
              after: lastReadId,
              limit: PAGE_SIZE,
            });

            if (newerMessages && newerMessages.length > 0) {
              // Есть новые сообщения после last_read
              finalMessages = newerMessages.sort(
                (a, b) => new Date(a.created_at) - new Date(b.created_at)
              );
              hasMoreDownFlag = newerMessages.length >= PAGE_SIZE;
              isDownFlag = false; // Пока не дойдем до конца

              // Загружаем и предыдущие сообщения (until last_read)
              if (finalMessages.length > 0) {
                const oldestId = finalMessages[0].id;
                try {
                  const olderMessages = await getMessages(token, {
                    chat_id: fullChat.chat_id,
                    before: oldestId,
                    limit: PAGE_SIZE,
                  });
                  if (olderMessages && olderMessages.length > 0) {
                    finalMessages = [...olderMessages, ...finalMessages].sort(
                      (a, b) => new Date(a.created_at) - new Date(b.created_at)
                    );
                    hasMoreUpFlag = olderMessages.length >= PAGE_SIZE;
                  }
                } catch (err) {
                  console.error('Failed to load messages before last_read', err);
                }
              }
            } else {
              // Нет сообщений после last_read: грузим before последние сообщения
              const latestMessages = await getMessages(token, {
                chat_id: fullChat.chat_id,
                before: 999999999,
                limit: PAGE_SIZE,
              });
              if (latestMessages && latestMessages.length > 0) {
                finalMessages = latestMessages.sort(
                  (a, b) => new Date(a.created_at) - new Date(b.created_at)
                );
                hasMoreUpFlag = latestMessages.length >= PAGE_SIZE;
              }
              isDownFlag = true; // Нет новых — ставим isDown = true
              hasMoreDownFlag = false;
            }
          } else {
            // Нет last_read: грузим последние сообщения и ставим isDown = true
            const latestMessages = await getMessages(token, {
              chat_id: fullChat.chat_id,
              before: 999999999,
              limit: PAGE_SIZE,
            });
            if (latestMessages && latestMessages.length > 0) {
              finalMessages = latestMessages.sort(
                (a, b) => new Date(a.created_at) - new Date(b.created_at)
              );
              hasMoreUpFlag = latestMessages.length >= PAGE_SIZE;
            }
            isDownFlag = true;
            hasMoreDownFlag = false;
          }
        } catch (err) {
          console.error('Failed to load messages', err);
        }

        if (isMounted) {
          setMessages(finalMessages);
          setHasMoreUp(hasMoreUpFlag);
          setHasMoreDown(hasMoreDownFlag);
          setIsDown(isDownFlag);
          isUserAtBottomRef.current = isDownFlag;

          // Если есть last_read и он в списке, скроллим к нему
          if (lastReadId && finalMessages.some(m => m.id === lastReadId)) {
            setTimeout(() => {
              const el = document.getElementById(`msg-${lastReadId}`);
              if (el) {
                el.scrollIntoView({ block: 'start' });
              }
            }, 100);
          }
        }
      } catch (err) {
        console.error('Failed to open chat', err);
      }
    };

    fetchChat();
    return () => { isMounted = false; };
  }, [chatId, token, myProfile.subject_id]);

  // ---------- ОБРАБОТКА WEBSOCKET (ОТРИСОВКА ТУТ) ----------
  useEffect(() => {
    if (!chat || !wsMessages.length) return;

    // Берем только необработанные сообщения
    const newWsEvents = wsMessages.slice(lastProcessedWsRef.current);
    lastProcessedWsRef.current = wsMessages.length;

    const toUpsert = [];

    newWsEvents.forEach(msg => {
      console.log(msg)
      const { type, data } = msg;
      if (data.chat_id !== chat.chat_id) return;

      if (type === 'send_message') {
        toUpsert.push(data);

        // Скроллим вниз только если сообщение отправлено текущим пользователем
        if (data.sender_id === myProfile.subject_id) {
          isUserAtBottomRef.current = true;
          setIsDown(true);
        }
      } else if (type === 'update_last_read') {
        // Обновляем last_read другого юзера
        if (data.subject_id !== myProfile.subject_id) {
          setOtherUserLastRead(data.message_id);
        }
      }
    });

    if (toUpsert.length > 0) {
      // Если isDown (пользователь у низа) — добавляем сообщения в конец
      if (isDown) {
        upsertMessages(toUpsert);
      } else {
        // Если NOT isDown — запрашиваем сообщения до первого текущего
        // Это имитирует пагинацию вверх автоматически
        if (messages.length > 0) {
          const firstId = messages[0].id;
          (async () => {
            try {
              const older = await getMessages(token, {
                chat_id: chat.chat_id,
                before: firstId,
                limit: PAGE_SIZE,
              });
              if (older && older.length > 0) {
                setMessages(prev => {
                  const combined = [...older, ...prev].sort(
                    (a, b) => new Date(a.created_at) - new Date(b.created_at)
                  );
                  return combined;
                });
              }
            } catch (err) {
              console.error('Failed to load messages before in WS handler', err);
            }
          })();
        }
      }
    }
  }, [wsMessages, chat, upsertMessages, myProfile.subject_id, isDown, messages, token, otherUserLastRead]);

  // ---------- ПАГИНАЦИЯ ВВЕРХ ----------
  const loadHistoryUp = useCallback(async () => {
    if (loadingUp || !hasMoreUp || !messages.length || !chat) return;
    setLoadingUp(true);

    const oldestId = messages[0].id;

    try {
      const older = await getMessages(token, {
        chat_id: chat.chat_id,
        before: oldestId,
        limit: PAGE_SIZE,
      });

      if (older.length < PAGE_SIZE) setHasMoreUp(false);

      if (older.length > 0) {
        const prevHeight = containerRef.current.scrollHeight;
        
        // Обновляем messages напрямую
        setMessages(prev => {
          const combined = [...older, ...prev].sort(
            (a, b) => new Date(a.created_at) - new Date(b.created_at)
          );
          return combined;
        });

        requestAnimationFrame(() => {
          if (containerRef.current) {
            containerRef.current.scrollTop = containerRef.current.scrollHeight - prevHeight;
          }
        });
      }
    } finally {
      setLoadingUp(false);
    }
  }, [loadingUp, hasMoreUp, messages, chat, token]);

  const loadMoreDown = useCallback(async () => {
    if (hasMoreDown === false) return;
    if (!messages.length || !chat) return;
    const lastId = messages[messages.length - 1].id;
    try {
      const newer = await getMessages(token, {
        chat_id: chat.chat_id,
        after: lastId,
        limit: PAGE_SIZE,
      });
      // Если получили пустой массив — значит достигли конца, выставляем isDown = true
      if (newer.length === 0) {
        setHasMoreDown(false);
        setIsDown(true);
      } else {
        if (newer.length < PAGE_SIZE) setHasMoreDown(false);
        upsertMessages(newer);
      }
    } catch (err) {
      console.error('Failed to load more down', err);
    }
  }, [hasMoreDown, messages, chat, token, upsertMessages]);

  const handleScroll = useCallback(() => {
    const el = containerRef.current;
    if (!el) return;

    const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight;
    isUserAtBottomRef.current = distanceFromBottom < 100;
    setIsDown(distanceFromBottom < 100);

    if (el.scrollTop < 50) {
      loadHistoryUp();
    }

    // Подгрузка вниз, если скроллим к низу
    if (distanceFromBottom < 100) {
      loadMoreDown();
    }
  }, [loadHistoryUp, loadMoreDown]);

  // ---------- ОТПРАВКА СООБЩЕНИЯ (БЕЗ ОБНОВЛЕНИЯ СТЕЙТА) ----------
  const handleSend = async () => {
    const messageText = text.trim();
    if (!messageText || !chat) return;

    setText('');

    try {
      await addMessage(chat.chat_id, messageText, token);
    } catch (err) {
      console.error('API Error:', err);
    }
  };

  // Авто-скролл
  useEffect(() => {
    if (isUserAtBottomRef.current) {
      bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  if (!chat) return <div style={{ padding: 20 }}>Загрузка...</div>;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div style={{ padding: 12, borderBottom: '1px solid #ddd', display: 'flex', alignItems: 'center', gap: 12 }}>
        {otherUser?.avatar_url && <img src={otherUser.avatar_url} alt={otherUser.alias} style={{ width: 40, height: 40, borderRadius: '50%', objectFit: 'cover' }} />}
        <span style={{ fontSize: 16, fontWeight: 500 }}>{otherUser?.alias || 'Загрузка...'}</span>
      </div>

      <div
        ref={containerRef}
        onScroll={handleScroll}
        style={{ flex: 1, overflowY: 'auto', padding: 12, display: 'flex', flexDirection: 'column' }}
      >
        {loadingUp && <div style={{ textAlign: 'center', color: '#999' }}>Загрузка истории...</div>}
        
        {messages.map(m => (
          <div id={`msg-${m.id}`} key={m.id}>
            <MessageBubble
              message={m}
              isMine={m.sender_id === myProfile.subject_id}
              isReadByOther={m.id <= otherUserLastRead}
            />
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      <div style={{ padding: 12, display: 'flex', gap: 8 }}>
        <textarea
          value={text}
          onChange={e => setText(e.target.value)}
          style={{ flex: 1, padding: 8, borderRadius: 8, border: '1px solid #ccc' }}
          onKeyDown={e => {
            if (e.key === 'Enter' && !e.shiftKey) {
              e.preventDefault();
              handleSend();
            }
          }}
        />
        <button onClick={handleSend} style={{ background: '#6a0dad', color: '#fff', border: 'none', borderRadius: 8, padding: '0 15px' }}>
          ➤
        </button>
      </div>
    </div>
  );
}