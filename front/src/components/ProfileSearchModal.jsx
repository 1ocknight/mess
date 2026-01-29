import React, { useState } from 'react';
import { getProfiles, getProfileById } from '../api/profile';
import { getChatBySubject, addChat } from '../api/chat';
import ProfileCard from './ProfileCard';
import { useNavigate } from 'react-router-dom';

export default function ProfileSearchModal({ token, onClose }) {
  const [query, setQuery] = useState('');
  const [profiles, setProfiles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [limit] = useState(8);
  const [after, setAfter] = useState(null);
  const [before, setBefore] = useState(null);
  const [hasMore, setHasMore] = useState(false);

  const navigate = useNavigate();

  const search = async (dir = null) => {
    if (!token) return;
    setLoading(true);
    try {
      // если начинается с @ — ищем по alias (поддерживаем пагинацию)
      if (query.startsWith('@')) {
        const alias = query.slice(1);
        const params = { limit };
        if (dir === 'after' && after) params.after = after;
        if (dir === 'before' && before) params.before = before;

        try {
          const res = await getProfiles(token, { alias, ...params });
          const arr = (res && res.profiles) || [];
          setProfiles(arr);
          setHasMore(arr.length === limit);
          if (arr.length > 0) {
            setAfter(arr[arr.length - 1].subject_id);
            setBefore(arr[0].subject_id);
          }
        } catch (err) {
          // 204 No Content или другие ошибки — считаем пустым
          setProfiles([]);
          setHasMore(false);
        }
      } else {
        // обычный ввод — трактуем как ID
        try {
          const p = await getProfileById(token, query);
          setProfiles([p]);
          setHasMore(false);
        } catch (err) {
          setProfiles([]);
        }
      }
    } finally {
      setLoading(false);
    }
  };

  const openChatFor = async (profile) => {
    try {
      const ch = await getChatBySubject(profile.subject_id, token);
      if (ch && ch.chat_id) {
        onClose();
        navigate(`/chat/${ch.chat_id}`);
        return;
      }
    } catch (err) {
      // continue to create
    }

    try {
      const created = await addChat(profile.subject_id, token);
      if (created && created.chat_id) {
        onClose();
        navigate(`/chat/${created.chat_id}`);
      }
    } catch (err) {
      console.error('Failed to open or create chat', err);
      alert('Не удалось открыть или создать чат');
    }
  };

  return (
    <div
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(75, 0, 130, 0.6)',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        zIndex: 1000,
      }}
    >
      <div
        style={{
          backgroundColor: '#6a0dad',
          padding: 24,
          borderRadius: 12,
          minWidth: 420,
          maxWidth: '90%',
          color: '#fff',
        }}
      >
        <h3 style={{ textAlign: 'center' }}>Поиск профилей</h3>
        <div style={{ display: 'flex', gap: 8, marginBottom: 12 }}>
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Введите ID или @alias"
            style={{ flex: 1, padding: '8px 10px', borderRadius: 8, border: 'none' }}
            onKeyDown={(e) => { if (e.key === 'Enter') search(); }}
          />
          <button
            onClick={() => search()}
            disabled={loading}
            style={{ padding: '8px 12px', borderRadius: 8, background: '#4a148c', color: '#fff', border: 'none' }}
          >
            {loading ? 'Поиск...' : 'Найти'}
          </button>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: 8, maxHeight: 320, overflowY: 'auto' }}>
          {profiles.map((p) => (
            <div key={p.subject_id} onClick={() => openChatFor(p)}>
              <ProfileCard profile={{ alias: p.alias, avatar_url: p.avatar_url }} />
            </div>
          ))}
          {profiles.length === 0 && <div style={{ textAlign: 'center', opacity: 0.9 }}>Ничего не найдено</div>}
        </div>

        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 12 }}>
          <div>
            <button
              onClick={() => search('before')}
              disabled={!before}
              style={{ padding: '6px 10px', borderRadius: 8, border: 'none', background: '#7b1fa2', color: '#fff' }}
            >
              Назад
            </button>
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
            <button
              onClick={onClose}
              style={{ padding: '6px 10px', borderRadius: 8, border: 'none', background: '#8e44ad', color: '#fff' }}
            >
              Закрыть
            </button>
            <button
              onClick={() => search('after')}
              disabled={!hasMore}
              style={{ padding: '6px 10px', borderRadius: 8, border: 'none', background: '#9c27b0', color: '#fff' }}
            >
              Далее
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
