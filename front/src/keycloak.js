import Keycloak from 'keycloak-js';

const keycloak = new Keycloak({
  url: 'http://localhost:7070', 
  realm: 'main',
  clientId: 'front',
});

export default keycloak;
