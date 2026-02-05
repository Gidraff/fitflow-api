import axios from 'axios';
import keycloak from './auth';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
});

api.interceptors.request.use(async (config) => {
    try {
        await keycloak.updateToken(30);
    } catch (error) {
        console.error("Token refresh failed", error);
        keycloak.login();
    }
    if (keycloak.token) { 
        config.headers.Authorization = `Bearer ${keycloak.token}`;
    } 
    return config;
});

export const getClients = async () => {
    const response = await api.get('/clients');
    return response.data;
};

export const getClientById = async (id: string) => {
    const response = await api.get(`/clients/${id}`);
    return response.data;
};

export const createClient = async (data: { name: string; goal: string; email: string; profile: string }) => {
    const response = await api.post('/clients', data);
    return response.data;
};

export const updateClient = async (id: string, data: { name: string; goal: string; email: string; profile: string }) => {
    const response = await api.put(`/clients/${id}`, data);
    return response.data;
};

export const deleteClient = async (id: string) => {
    await api.delete(`/clients/${id}`);
};

export const getClientSessions = async (clientId: string) => {
    const response = await api.get(`/clients/${clientId}/sessions`);
    return response.data;
};

export const addSessionFeedback = async (sessionId: string, feedback: { feedback: string; performance_rating: number }) => {
    const response = await api.patch(`/sessions/${sessionId}/feedback`, feedback);
    return response.data;
};

export { api };