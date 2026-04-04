import { error } from '@sveltejs/kit';

export async function load({ fetch, params, url }) {
    const { scope } = params;
    const limit = url.searchParams.get('limit') || '50';
    const offset = url.searchParams.get('offset') || '0';
    
    try {
        const res = await fetch(`http://localhost:8080/api/events?scope=${scope}&limit=${limit}&offset=${offset}`);

        if (!res.ok) {
            throw error(res.status, `Failed to fetch events for Scope ${scope}: ${res.statusText}`);
        }

        const events = await res.json();
        return { events, scope, limit: parseInt(limit), offset: parseInt(offset) };
    } catch (err) {
        console.error('scope events load error', err);
        return { events: [], scope, limit: parseInt(limit), offset: parseInt(offset) };
    }
}
