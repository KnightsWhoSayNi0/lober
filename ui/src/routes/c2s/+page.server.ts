import { error } from '@sveltejs/kit';

export async function load({ fetch }) {
    try {
        const res = await fetch('http://localhost:8080/api/c2s');

        if (!res.ok) {
            throw error(res.status, `Failed to fetch c2s: ${res.statusText}`);
        }

        const c2s = await res.json();
        return { c2s };
    } catch (err) {
        // fallback: return empty array for UI resilience, but surface in logs
        console.error('c2s load error', err);
        return { c2s: [] };
    }
}