import { error } from '@sveltejs/kit';

export async function load({ fetch }) {
    try {
        const res = await fetch('http://localhost:8080/api/scope');

        if (!res.ok) {
            throw error(res.status, `Failed to fetch scopes: ${res.statusText}`);
        }

        const scopes = await res.json();
        return { scopes };
    } catch (err) {
        // fallback: return empty array for UI resilience, but surface in logs
        console.error('scopes load error', err);
        return { scopes: [] };
    }
}