const API_BASE = 'http://localhost:8081';

// --- CHAT ENDPOINTS ---

// Получить чат по subject_id
export async function getChatBySubject(subjectId, token) {
  const res = await fetch(`${API_BASE}/chat/subject/${subjectId}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to fetch chat by subject');
  return res.json();
}

// Добавить чат по subject_id
export async function addChat(subjectId, token) {
  const res = await fetch(`${API_BASE}/chat/subject/${subjectId}`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to add chat');
  return res.json();
}

// открыть чат по chat_id
export async function openChatByID(chatId, token, limit) {
  const params = new URLSearchParams();
  if (limit) params.append('limit', limit);

  const res = await fetch(`${API_BASE}/chat/${chatId}?${params.toString()}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to fetch chat by ID');
  return res.json();
}

// Получить все чаты
export async function getChats(token) {
  const res = await fetch(`${API_BASE}/chats`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Failed to fetch chats');
  return res.json();
}

// --- MESSAGES ENDPOINTS ---

export async function getMessages(
  token,
  { chat_id, before, after, limit } = {}
) {
  const params = new URLSearchParams();

  if (chat_id) params.append('chat_id', chat_id);
  if (before !== undefined) params.append('before', before);
  if (after !== undefined) params.append('after', after);
  if (limit) params.append('limit', limit);

  const url = `${API_BASE}/messages?${params.toString()}`;
  console.log('[getMessages] Fetching:', url);

  const res = await fetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!res.ok) {
    const errorText = await res.text();
    console.error('[getMessages] API Error:', res.status, errorText);
    throw new Error(`Failed to fetch messages: ${res.status} ${errorText}`);
  }

  // 🔥 защита от пустого ответа
  const text = await res.text();
  if (!text) return [];

  return JSON.parse(text);
}


// Добавить сообщение
export async function addMessage(chatId, content, token) {
  const res = await fetch(`${API_BASE}/message`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ chat_id: chatId, content }),
  });
  if (!res.ok) throw new Error('Failed to add message');
  return res.json();
}

// Обновить сообщение
export async function updateMessage(messageId, content, version, token) {
  const res = await fetch(`${API_BASE}/message`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ message_id: messageId, content, version }),
  });
  if (!res.ok) throw new Error('Failed to update message');
  return res.json();
}

// --- LAST READ ENDPOINT ---

// Обновить последний прочитанный message_id в чате
export async function updateLastRead(chatId, messageId, token) {
  const res = await fetch(`${API_BASE}/lastread`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ chat_id: chatId, message_id: messageId }),
  });
  if (!res.ok) throw new Error('Failed to update last read');
  return res.json();
}
