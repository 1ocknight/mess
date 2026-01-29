import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useKeycloak } from '@react-keycloak/web';
import ProfilePage from './pages/CreateProfilePage';
import HomePage from './pages/HomePage';  // ← импортируем HomePage
import { getProfile } from './api/profile';
import ChatPage from './pages/ChatPage';

function App() {
  const { keycloak, initialized } = useKeycloak();

  if (!initialized) return <div>Loading...</div>;

  return (
    <Router>
      <Routes>
        <Route path="/profile" element={<ProfilePage />} />
        <Route
          path="/"
          element={
            <RequireProfile>
              <HomePage />  {/* ← отображаем HomePage на главной */}
            </RequireProfile>
          }
        />
        <Route
          path="/chat/:id"
          element={
            <RequireProfile>
              <ChatPage />
            </RequireProfile>
          }
        />
      </Routes>
    </Router>
  );
}

// Компонент для проверки, есть ли профиль
function RequireProfile({ children }) {
  const { keycloak } = useKeycloak();
  const [profileExists, setProfileExists] = useState(null);

  useEffect(() => {
    if (!keycloak.authenticated) return;

    getProfile(keycloak.token)
      .then(() => setProfileExists(true))
      .catch(() => setProfileExists(false));
  }, [keycloak]);

  if (profileExists === null) return <div>Checking profile...</div>;
  if (profileExists === false) return <Navigate to="/profile" replace />;

  return children;
}

export default App;
