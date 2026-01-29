import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { ReactKeycloakProvider } from '@react-keycloak/web'
import './index.css'
import App from './App.jsx'
import keycloak from './keycloak'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <ReactKeycloakProvider 
    authClient={keycloak}
    initOptions={{
      onLoad: 'login-required', 
      checkLoginIframe: false,
      pkceMethod: 'S256', 
      enableLogging: true,
    }}
    >
      <App />
    </ReactKeycloakProvider>
  </StrictMode>,
)
