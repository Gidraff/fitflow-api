import Keycloak from "keycloak-js";


const keycloak = new Keycloak({
    url: import.meta.env.VITE_AUTH_URL,
    realm: import.meta.env.VITE_REALM,
    clientId: import.meta.env.VITE_CLIENT_ID,
});

export const initKeycloak = (onAuthenticated: () => void) => {
    keycloak.init({ 
        onLoad: 'login-required', 
        checkLoginIframe: false,
        // redirectUri: 'https://localhost/' 
    }).then((authenticated) => {
        if (authenticated) {
            onAuthenticated();
        }
    }).catch(console.error);
};

export default keycloak;