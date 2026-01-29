import React, { useEffect, useState } from 'react';
import { useKeycloak } from '@react-keycloak/web';
import { useParams, useNavigate } from 'react-router-dom';
import { WSProvider } from '../context/WebSocketContext';
import ChatWindow from '../components/ChatWindow';
import { getProfile } from '../api/profile';

export default function ChatPage() {
  const { keycloak } = useKeycloak();
  const { id: chatId } = useParams(); // берём ID чата из URL
  const navigate = useNavigate();
  const [profile, setProfile] = useState(null);

  // Загружаем профиль пользователя
  useEffect(() => {
    if (!keycloak?.authenticated) return;
    getProfile(keycloak.token).then(setProfile).catch(() => console.log('No profile'));
  }, [keycloak]);

  if (!profile) return <div>Loading profile...</div>;

  return (
    <div style={{ 
        display: 'flex', 
        height: '100vh', 
        width: '100vw', 
        gap: 20, 
        justifyContent: 'center', 
        alignItems: 'center', 
        background: 'linear-gradient(135deg, rgba(126, 87, 194, 0.7), rgba(156, 39, 176, 0.7))'
        }}>
    <div
        style={{
            display: 'flex',
            flexDirection: 'column',
            width: '500px',
            height: '750px', // высота карточки
            padding: 20,
            borderRadius: 20,
            backgroundColor: 'rgba(255, 255, 255, 0.05)',
            boxShadow: '0 8px 30px rgba(0,0,0,0.3)',
            color: '#fff',
        }}
        >
        {/* Кнопка назад */}
        <div style={{ marginBottom: 12 }}>
            <button
            onClick={() => navigate('/')}
            style={{
                padding: '6px 12px',
                borderRadius: 6,
                border: 'none',
                backgroundColor: '#6a0dad',
                color: '#fff',
                cursor: 'pointer',
                fontWeight: 'bold',
            }}
            >
            ← Назад к чатам
            </button>
        </div>

        {/* ChatWindow / ChatsList */}
        <div
            style={{
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            minHeight: 0, // важно для корректного scroll внутри flex
            overflowY: 'auto', // включаем скролл
            }}
        >
            <WSProvider token={keycloak.token}>
            <ChatWindow chatId={chatId} token={keycloak.token} myProfile={profile} />
            {/* или <ChatsList token={keycloak.token} /> */}
            </WSProvider>
        </div>
        </div>
        </div>
  );
}
