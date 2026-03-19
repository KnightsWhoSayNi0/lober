import { error } from '@sveltejs/kit'

export async function load({ fetch, params }) {
    try {
        const res = await fetch(`http://localhost:8080/api/users/${params.user}`);

        if (!res.ok) {
            throw error(res.status, `Failed to fetch users: ${res.statusText}`);
        }

        const user = await res.json();
        return { user: user };
    } catch (err) {
        // fallback: return empty array for UI resilience, but surface in logs
        console.error('users load error', err);
        return { user: {} };
    }
}