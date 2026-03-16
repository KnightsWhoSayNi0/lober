<script lang="ts">
    const { data } = $props<{ data: { users: Array<{ username: string; email: string; team: string }> }}>();
</script>

<h1>Users</h1>

{#if data.users.length > 0}
    <table>
        <thead>
            <tr>
                <th>username</th>
                <th>email</th>
                <th>team</th>
            </tr>
        </thead>
        <tbody>
            {#each data.users as user}
                <tr>
                    <td><a href="/users/{user.username}">{user.username}</a></td>
                    <td>{user.email}</td>
                    <td>{user.team}</td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <p>No users available yet.</p>
{/if}

<style>
    table {
        width: 100%;
        border-collapse: collapse;
    }

    th, td {
        border: 1px solid #ddd;
        padding: 8px;
    }

    th {
        background-color: #f2f2f2;
        text-align: left;
    }
</style>