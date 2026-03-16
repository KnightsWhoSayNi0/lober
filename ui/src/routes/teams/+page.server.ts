import { error } from '@sveltejs/kit';

export async function load({ fetch }) {
    try {
        const res = await fetch('http://localhost:8080/api/teams');

        if (!res.ok) {
            throw error(res.status, `Failed to fetch teams: ${res.statusText}`);
        }

        const teams = await res.json();
        return { teams };
    } catch (err) {
        // fallback: return empty array for UI resilience, but surface in logs
        console.error('teams load error', err);
        return { teams: [] };
    }
}