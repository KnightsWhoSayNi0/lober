export async function load({ fetch }) {
    try {
        const [tokensRes, usersRes, c2sRes] = await Promise.all([
            fetch('http://localhost:8080/api/tokens'),
            fetch('http://localhost:8080/api/users'),
            fetch('http://localhost:8080/api/c2s')
        ]);

        const tokens = tokensRes.ok ? await tokensRes.json() : [];
        const users = usersRes.ok ? await usersRes.json() : [];
        const c2s = c2sRes.ok ? await c2sRes.json() : [];

        return { tokens, users, c2s };
    } catch (err) {
        console.error('tokens load error', err);
        return { tokens: [], users: [], c2s: [] };
    }
}
