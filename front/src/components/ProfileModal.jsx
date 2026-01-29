import React, { useState } from 'react';
import { updateProfile, getUploadAvatarUrl } from '../api/profile';

export default function ProfileModal({ profile, token, onClose, onUpdate, keycloak }) {
  const [alias, setAlias] = useState(profile.alias);
  const [avatarFile, setAvatarFile] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(profile.avatar_url);
  const [loading, setLoading] = useState(false);

  const handleAvatarChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setAvatarFile(file);
      setAvatarPreview(URL.createObjectURL(file));
    }
  };

  const handleSave = async () => {
    setLoading(true);
    try {
      const updated = await updateProfile(token, alias, profile.version);

      if (avatarFile) {
        const { upload_url } = await getUploadAvatarUrl(token);
        await fetch(upload_url, {
          method: 'PUT',
          body: avatarFile,
        });
      }

      onUpdate(updated);
      onClose();
    } catch (err) {
      console.error(err);
      alert('Failed to update profile or upload avatar');
    } finally {
      setLoading(false);
    }
  };

  return (
      <div
    style={{
      position: 'fixed',
      top: 0, left: 0, right: 0, bottom: 0,
      backgroundColor: 'rgba(75, 0, 130, 0.6)', // мягкий фиолетовый фон
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      zIndex: 1000,
    }}
  >
    <div
      style={{
        backgroundColor: '#6a0dad', // насыщенный фиолетовый
        padding: '30px 40px',
        borderRadius: 15,
        minWidth: 350,
        boxShadow: '0 10px 30px rgba(0,0,0,0.3)',
        color: '#fff',
        fontFamily: 'Arial, sans-serif',
      }}
    >
      <h2 style={{ textAlign: 'center', marginBottom: 20 }}>Профиль</h2>

      {/* Аватар */}
      <div style={{ marginBottom: 20, textAlign: 'center' }}>
        <img
          src={avatarPreview}
          alt="Avatar"
          style={{
            width: 100,
            height: 100,
            borderRadius: '50%',
            marginBottom: 10,
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

      {/* ID */}
      <div style={{ marginBottom: 15, textAlign: 'center' }}>
        <label>ID: <strong>{profile.subject_id}</strong></label>
      </div>

      {/* Alias */}
      <div style={{ marginBottom: 20, textAlign: 'center' }}>
        <label>
          Alias:
          <input
            value={alias}
            onChange={(e) => setAlias(e.target.value)}
            style={{
              marginLeft: 10,
              padding: '5px 10px',
              borderRadius: 8,
              border: 'none',
              outline: 'none',
              fontSize: 14,
            }}
          />
        </label>
      </div>

      {/* Кнопки */}
      <div style={{ display: 'flex', justifyContent: 'center', flexWrap: 'wrap', gap: 10 }}>
        <button
          onClick={handleSave}
          disabled={loading}
          style={{
            padding: '8px 16px',
            borderRadius: 8,
            border: 'none',
            backgroundColor: '#9b59b6',
            color: '#fff',
            cursor: 'pointer',
            fontWeight: 'bold',
          }}
        >
          {loading ? 'Сохранение...' : 'Сохранить'}
        </button>

        <button
          onClick={onClose}
          style={{
            padding: '8px 16px',
            borderRadius: 8,
            border: 'none',
            backgroundColor: '#8e44ad',
            color: '#fff',
            cursor: 'pointer',
          }}
        >
          Назад
        </button>
      </div>
    </div>
  </div>

  );
}
