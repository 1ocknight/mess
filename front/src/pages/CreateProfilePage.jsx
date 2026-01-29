import React, { useState, useEffect } from 'react';
import { useKeycloak } from '@react-keycloak/web';
import { addProfile, getUploadAvatarUrl, getProfile } from '../api/profile';
import { useNavigate } from 'react-router-dom';
import defaultAvatar from '../../public/vite.svg'; // дефолтный аватар

export default function CreateProfilePage() {
  const { keycloak, initialized } = useKeycloak();
  const navigate = useNavigate();

  const [alias, setAlias] = useState('');
  const [avatarFile, setAvatarFile] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(defaultAvatar);
  const [loading, setLoading] = useState(false);
  const [checkingProfile, setCheckingProfile] = useState(true);

  // Проверяем, есть ли уже профиль
  useEffect(() => {
    if (!initialized || !keycloak.authenticated) return;

    getProfile(keycloak.token)
      .then(() => navigate('/', { replace: true }))
      .catch(() => setCheckingProfile(false));
  }, [initialized, keycloak, navigate]);

  if (!initialized || !keycloak.authenticated || checkingProfile) {
    return <div>Loading...</div>;
  }

  const handleAvatarChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (file.size > 5 * 1024 * 1024) {
        alert('Файл слишком большой. Максимум 5MB.');
        return;
      }
      if (!file.type.startsWith('image/')) {
        alert('Можно загружать только изображения.');
        return;
      }

      setAvatarFile(file);

      const reader = new FileReader();
      reader.onload = () => setAvatarPreview(reader.result);
      reader.readAsDataURL(file);
    }
  };

  const handleCreate = async () => {
    if (!alias) {
      alert('Alias is required');
      return;
    }

    try {
      setLoading(true);

      await addProfile(keycloak.token, alias);

      const { upload_url } = await getUploadAvatarUrl(keycloak.token);

      if (avatarFile) {
        await fetch(upload_url, { method: 'PUT', body: avatarFile });
      } else {
        const response = await fetch(defaultAvatar);
        const blob = await response.blob();
        await fetch(upload_url, { method: 'PUT', body: blob });
      }

      alert('Profile created!');
      navigate('/', { replace: true });
    } catch (err) {
      console.error(err);
      alert('Failed to create profile or upload avatar');
    } finally {
      setLoading(false);
    }
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
      <div
        style={{
          width: '500px',
          padding: 40,
          borderRadius: 20,
          backgroundColor: 'rgba(255, 255, 255, 0.05)', // слегка прозрачная карточка
          boxShadow: '0 8px 30px rgba(0,0,0,0.3)',
          color: '#fff',
          textAlign: 'center',
        }}
      >
        <h1 style={{ marginBottom: 30 }}>Create Profile</h1>

        {/* Аватар */}
        <div style={{ marginBottom: 20 }}>
          <img
            src={avatarPreview}
            alt="Avatar"
            style={{
              width: 120,
              height: 120,
              borderRadius: '50%',
              cursor: 'pointer',
              border: '3px solid #d1b3ff',
              transition: 'transform 0.2s',
            }}
            onClick={() => document.getElementById('avatarInput').click()}
            onMouseOver={(e) => (e.currentTarget.style.transform = 'scale(1.05)')}
            onMouseOut={(e) => (e.currentTarget.style.transform = 'scale(1)')}
          />
          <input
            id="avatarInput"
            type="file"
            accept="image/*"
            style={{ display: 'none' }}
            onChange={handleAvatarChange}
          />
        </div>

        {/* Alias */}
        <div style={{ marginBottom: 30 }}>
          <input
            type="text"
            placeholder="Enter alias"
            value={alias}
            onChange={(e) => setAlias(e.target.value)}
            style={{
              width: '80%',
              padding: '10px 15px',
              borderRadius: 12,
              border: 'none',
              outline: 'none',
              fontSize: 16,
            }}
          />
        </div>

        {/* Кнопка создания */}
        <button
          onClick={handleCreate}
          disabled={loading}
          style={{
            padding: '12px 25px',
            borderRadius: 12,
            border: 'none',
            backgroundColor: '#9b59b6',
            color: '#fff',
            fontSize: 16,
            fontWeight: 'bold',
            cursor: 'pointer',
            transition: 'background-color 0.2s',
          }}
          onMouseOver={(e) => (e.currentTarget.style.backgroundColor = '#8e44ad')}
          onMouseOut={(e) => (e.currentTarget.style.backgroundColor = '#9b59b6')}
        >
          {loading ? 'Creating...' : 'Create Profile'}
        </button>
      </div>
    </div>
  );
}
