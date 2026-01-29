import React, { useEffect, useState } from 'react';
import { getProfileById } from '../api/profile';

export default function ChatListItem({ chat, token, onClick }) {
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    const subjectId = chat.second_subject_id;
    if (!subjectId) return;

    getProfileById(token, subjectId)
      .then(setProfile)
      .catch(console.error);
  }, [chat.second_subject_id, token]);

  const unreadCount = chat.unread_count || 0;

  const lastMessage = chat.last_message;

  // --- логика галочек ---
  const isMyLastMessage =
    lastMessage && lastMessage.sender_id != chat.second_subject_id;

  const showSingleCheck =
    isMyLastMessage && !chat.is_last_message_read;

  const showDoubleCheck =
    isMyLastMessage && chat.is_last_message_read;

  return (
    <div
      onClick={() => onClick(chat)}
      style={{
        display: 'flex',
        alignItems: 'center',
        cursor: 'pointer',
        padding: 10,
        borderRadius: 12,
        marginBottom: 8,
        backgroundColor: unreadCount > 0 ? '#6a0dad' : '#7e57c2',
        color: '#fff',
      }}
    >
      <img
        src={profile?.avatar_url}
        alt="Avatar"
        style={{
          width: 40,
          height: 40,
          borderRadius: '50%',
          marginRight: 10,
          border: unreadCount > 0 ? '2px solid #fff' : 'none',
        }}
      />

      <div style={{ flex: 1 }}>
        <div style={{ fontWeight: 'bold' }}>
          {profile?.alias || 'Loading...'}
        </div>

        <div
          style={{
            fontSize: 12,
            opacity: 0.8,
            display: 'flex',
            alignItems: 'center',
            gap: 6,
          }}
        >
          <span>
            {lastMessage
              ? lastMessage.content.length > 40
                ? lastMessage.content.slice(0, 40) + ' .....'
                : lastMessage.content
              : 'No messages'}
          </span>

          {/* ✔ / ✔✔ */}
          {showSingleCheck && <span>✔</span>}
          {showDoubleCheck && <span>✔✔</span>}
        </div>
      </div>

      {/* непрочитанные — оставляем как есть */}
      {unreadCount > 0 && (
        <div
          style={{
            backgroundColor: '#ff5252',
            borderRadius: '50%',
            width: 20,
            height: 20,
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            fontSize: 12,
            fontWeight: 'bold',
          }}
        >
          {unreadCount}
        </div>
      )}
    </div>
  );
}
