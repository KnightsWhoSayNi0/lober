import type { HandleFetch } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const handleFetch: HandleFetch = async ({ request, fetch }) => {
    if (request.url.startsWith('http://localhost:8080/api')) {
        request = new Request(
            request.url.replace('http://localhost:8080/api', 'http://server:8080/api'),
            request
        );
    }

    if (request.url.startsWith('http://server:8080/api')) {
        const token = env.MASTER_TOKEN || 'dev-master-token';
        request.headers.set('Authorization', `Bearer ${token}`);
    }

    return fetch(request);
};
