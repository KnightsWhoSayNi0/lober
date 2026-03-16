import { error } from '@sveltejs/kit';

export async function load({ fetch }) {
    try {
        const res = await fetch('http://localhost:8080/events');

        if (!res.ok) {
            throw error(res.status, `Failed to fetch events: ${res.statusText}`);
        }

        const events = await res.json();
        return { events };
    } catch (err) {
        // fallback: return empty array for UI resilience, but surface in logs
        console.error('events load error', err);
        return { events: [] };
    }
}