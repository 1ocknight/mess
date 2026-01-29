import React, { useState, useEffect, useRef } from 'react';
import { useKeycloak } from '@react-keycloak/web';
import { getProfile } from '../api/profile';
import ProfileCard from '../components/ProfileCard';
import ProfileModal from '../components/ProfileModal';
import ProfileSearchModal from '../components/ProfileSearchModal';
import ChatsList from '../components/ChatsList'; // <-- импортируем компонент чатов
import { WSProvider } from '../context/WebSocketContext';

export default function HomePage() {
  const { keycloak, initialized } = useKeycloak();
  const [profile, setProfile] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [showSearch, setShowSearch] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    if (!initialized || !keycloak.authenticated) return;
    getProfile(keycloak.token).then(setProfile).catch(() => console.log('No profile found'));
  }, [initialized, keycloak]);

  if (!initialized || !profile) return <div>Loading profile...</div>;

  const handleLogout = () => keycloak.logout();
  const handleAccountPage = () => {
    if (!keycloak || !keycloak.authServerUrl || !keycloak.realm) return;
    const baseUrl = keycloak.authServerUrl.replace(/\/$/, '');
    window.open(`${baseUrl}/realms/${keycloak.realm}/account`, '_blank');
  };

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
      <div style={{
        width: '500px',
        height: '750px',
        display: 'flex',
        flexDirection: 'column',
        padding: 20,
        backgroundColor: 'rgba(255, 255, 255, 0.05)', // слегка прозрачная карточка
        boxShadow: '0 8px 30px rgba(0,0,0,0.3)',
        color: '#fff',
        borderRadius: '12px', // верхний левый и нижний левый углы без скругления
        boxSizing: 'border-box',
      }}>
        {/* Меню + Поиск сверху */}
        <div style={{ marginBottom: 20 }}>
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <button
              onClick={() => setMenuOpen(!menuOpen)}
              style={{
                padding: '8px',
                width: 40,
                height: 40,
                borderRadius: 8,
                border: 'none',
                backgroundColor: '#6a0dad',
                color: '#fff',
                cursor: 'pointer',
                fontWeight: 'bold',
                fontSize: 18,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                flexShrink: 0,
              }}
            >
              ☰
            </button>
            <button
              onClick={() => setShowSearch(true)}
              style={{
                padding: '8px 12px',
                borderRadius: 8,
                border: 'none',
                backgroundColor: '#4a148c',
                color: '#fff',
                cursor: 'pointer',
                fontWeight: 'bold',
                flex: 1,
              }}
            >
              🔎 Поиск профилей
            </button>
          </div>
          {menuOpen && (
            <div style={{ marginTop: 5, backgroundColor: '#7b1fa2', borderRadius: 8, overflow: 'hidden' }}>
              <button
                onClick={handleAccountPage}
                style={{
                  display: 'block',
                  width: '100%',
                  padding: '10px 20px',
                  border: 'none',
                  backgroundColor: 'transparent',
                  color: '#fff',
                  textAlign: 'left',
                  cursor: 'pointer',
                }}
              >
                Настройки аккаунта
              </button>
              <button
                onClick={handleLogout}
                style={{
                  display: 'block',
                  width: '100%',
                  padding: '10px 20px',
                  border: 'none',
                  backgroundColor: 'transparent',
                  color: '#fff',
                  textAlign: 'left',
                  cursor: 'pointer',
                }}
              >
                Выйти
              </button>
            </div>
          )}
        </div>

        {showSearch && (
          <ProfileSearchModal
            token={keycloak.token}
            onClose={() => setShowSearch(false)}
          />
        )}

        {/* --- Список чатов --- */}
        <WSProvider token={keycloak.token}>
          <ChatsList token={keycloak.token} />
        </WSProvider>

        {/* Профиль снизу по центру */}
        <div style={{ display: 'flex', justifyContent: 'center', marginTop: 'auto' }}>
          <ProfileCard profile={profile} onClick={() => setShowModal(true)} />
          {showModal && (
            <ProfileModal
              profile={profile}
              token={keycloak.token}
              onClose={() => setShowModal(false)}
              onUpdate={setProfile}
              keycloak={keycloak}
            />
          )}
        </div>
      </div>
    </div>
  );
}
