import { error } from '@sveltejs/kit';

export async function load({ fetch }) {
    try {
        const res = await fetch('http://localhost:8080/api/users');

        if (!res.ok) {
            throw error(res.status, `Failed to fetch users: ${res.statusText}`);
        }

        const users = await res.json();
        return { users };
    } catch (err) {
        // fallback: return empty array for UI resilience, but surface in logs
        console.error('users load error', err);
        return { users: [] };
    }
}