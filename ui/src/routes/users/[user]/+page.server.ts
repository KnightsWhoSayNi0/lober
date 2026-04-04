import { error } from '@sveltejs/kit';

export async function load({ fetch, params, url }) {
    const { user: username } = params;
    const limit = url.searchParams.get('limit') || '50';
    const offset = url.searchParams.get('offset') || '0';

    try {
        const [userRes, eventsRes] = await Promise.all([
            fetch(`http://localhost:8080/api/users/${username}`),
            fetch(`http://localhost:8080/api/events?user=${username}&limit=${limit}&offset=${offset}`)
        ]);

        if (!userRes.ok) {
            throw error(userRes.status, `Failed to fetch user ${username}: ${userRes.statusText}`);
        }

        const user = await userRes.json();
        const events = eventsRes.ok ? await eventsRes.json() : [];

        return { user, events, limit: parseInt(limit), offset: parseInt(offset) };
    } catch (err) {
        console.error('user load error', err);
        throw error(500, 'Internal Server Error');
    }
}
