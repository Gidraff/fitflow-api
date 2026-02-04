# --- Stage 1: Build Stage ---
FROM node:24-alpine AS build

WORKDIR /app

# Install dependencies first (better caching)
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Pass build-time arguments (Vite needs these to bake them in)
ARG VITE_API_URL
ARG VITE_AUTH_URL
ARG VITE_REALM
ARG VITE_CLIENT_ID

# Set them as ENV so the 'npm run build' command sees them
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_AUTH_URL=$VITE_AUTH_URL
ENV VITE_REALM=$VITE_REALM
ENV VITE_CLIENT_ID=$VITE_CLIENT_ID

RUN npm run build

# --- Stage 2: Production Stage ---
FROM nginx:stable-alpine

# Copy the build output to Nginx's serving directory
COPY --from=build /app/dist /usr/share/nginx/html

# Copy a custom Nginx config to handle SPA routing (crucial for React/Vite)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]