const API_BASE = 'http://localhost:8080';

export async function getProfile(token) {
  const res = await fetch(`${API_BASE}/profile`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to fetch profile');
  return res.json();
}

export async function getProfileById(token, id) {
  const res = await fetch(`${API_BASE}/profile/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to fetch profile by ID');
  return res.json();
}

export async function getProfiles(token, params = {}) {
  const query = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/profiles?${query}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to fetch profiles');
  return res.json();
}

export async function addProfile(token, alias) {
  const res = await fetch(`${API_BASE}/profile`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ alias }),
  });
  if (!res.ok) throw new Error('Failed to add profile');
  return res.json();
}

export async function updateProfile(token, alias, version) {
  const res = await fetch(`${API_BASE}/profile`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ alias, version }),
  });
  if (!res.ok) throw new Error('Failed to update profile');
  return res.json();
}

export async function getUploadAvatarUrl(token) {
  const res = await fetch(`${API_BASE}/avatar`, {
    method: 'PUT',
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to get upload URL');
  return res.json();
}

export async function deleteAvatar(token) {
  const res = await fetch(`${API_BASE}/avatar`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to delete avatar');
  return res.json();
}
