import { error } from '@sveltejs/kit';

export async function load({ fetch, params, url }) {
    const { team } = params;
    const limit = url.searchParams.get('limit') || '50';
    const offset = url.searchParams.get('offset') || '0';
    
    try {
        const res = await fetch(`http://localhost:8080/api/events?team=${team}&limit=${limit}&offset=${offset}`);

        if (!res.ok) {
            throw error(res.status, `Failed to fetch events for team ${team}: ${res.statusText}`);
        }

        const events = await res.json();
        return { events, team, limit: parseInt(limit), offset: parseInt(offset) };
    } catch (err) {
        console.error('team events load error', err);
        return { events: [], team, limit: parseInt(limit), offset: parseInt(offset) };
    }
}
