import { error } from '@sveltejs/kit';

export async function load({ fetch, params, url }) {
    const { c2 } = params;
    const limit = url.searchParams.get('limit') || '50';
    const offset = url.searchParams.get('offset') || '0';
    
    try {
        const res = await fetch(`http://localhost:8080/api/events?c2=${c2}&limit=${limit}&offset=${offset}`);

        if (!res.ok) {
            throw error(res.status, `Failed to fetch events for C2 ${c2}: ${res.statusText}`);
        }

        const events = await res.json();
        return { events, c2, limit: parseInt(limit), offset: parseInt(offset) };
    } catch (err) {
        console.error('c2 events load error', err);
        return { events: [], c2, limit: parseInt(limit), offset: parseInt(offset) };
    }
}
