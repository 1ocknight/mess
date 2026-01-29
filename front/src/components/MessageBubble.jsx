import React from 'react';

export default function MessageBubble({ message, isMine, isReadByOther }) {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: isMine ? 'flex-end' : 'flex-start', // выравниваем блок влево/вправо
        width: '100%',
      }}
    >
      <div
        style={{
          backgroundColor: isMine ? '#7e57c2' : '#e0e0e0',
          color: isMine ? '#fff' : '#000',
          padding: '8px 12px',
          borderRadius: 12,
          maxWidth: '70%',
          wordBreak: 'break-word',
        }}
      >
        <div style={{ marginBottom: 4 }}>{message.content}</div>
        <div
          style={{
            fontSize: 10,
            opacity: 0.7,
            textAlign: 'right',
          }}
        >
          {new Date(message.created_at).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
          })}
          {isMine && (isReadByOther ? ' ✔✔' : ' ✔')}
        </div>
      </div>
    </div>
  );
}
