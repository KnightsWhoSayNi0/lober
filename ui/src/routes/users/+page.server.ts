import { error } from '@sveltejs/kit';

export async function load({ fetch }) {
    try {
        const [usersRes, teamsRes] = await Promise.all([
            fetch('http://localhost:8080/api/users'),
            fetch('http://localhost:8080/api/teams')
        ]);

        if (!usersRes.ok) throw error(usersRes.status, `Failed to fetch users: ${usersRes.statusText}`);
        if (!teamsRes.ok) throw error(teamsRes.status, `Failed to fetch teams: ${teamsRes.statusText}`);

        const users = await usersRes.json();
        const teams = await teamsRes.json();
        
        return { users, teams };
    } catch (err) {
        console.error('users load error', err);
        return { users: [], teams: [] };
    }
}
