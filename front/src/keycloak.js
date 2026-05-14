import Keycloak from 'keycloak-js';

const keycloak = new Keycloak({
  url: '/auth',
  realm: 'main',
  clientId: 'front',
});

export default keycloak;
