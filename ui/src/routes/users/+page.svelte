<script lang="ts">
    import { invalidateAll } from '$app/navigation';
    import { getContext } from 'svelte';

    const { data } = $props<{ data: { users: Array<{ username: string; team: string }>, teams: Array<{ name: string }> }}>();
    const getMasterToken = getContext<() => string>('masterToken');

    let showModal = $state(false);
    let formData = $state({ username: '', password: '', team: '' });

    async function handleSubmit(e: Event) {
        e.preventDefault();

        const response = await fetch('/api/users', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getMasterToken()}`
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            showModal = false;
            formData = { username: '', password: '', team: '' };
            await invalidateAll();
        }
    }

    async function handleDelete(name: string) {
        if (!confirm(`Are you sure you want to remove user "${name}"?`)) return;
        
        const response = await fetch(`/api/users/${name}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${getMasterToken()}`
            }
        });

        if (response.ok) {
            await invalidateAll();
        }
    }
</script>

<title>Lober Users</title>

<div class="header-container">
    <h1>Users</h1>
    <button onclick={() => showModal = true} >New User</button>
</div>

{#if showModal}
    <div class="modal-backdrop">
        <div class="modal">
            <h2>New User</h2>
            <form onsubmit={handleSubmit}>
                <label>
                    Username:
                    <input type="text" bind:value={formData.username} required />
                </label>
                <label>
                    Password:
                    <input type="password" bind:value={formData.password} required />
                </label>
                <label>
                    Team:
                    <select bind:value={formData.team} required>
                        <option value="">Select Team</option>
                        {#each data.teams as team}
                            <option value={team.name}>{team.name}</option>
                        {/each}
                    </select>
                </label>
                <div class="actions">
                    <button type="button" onclick={() => showModal = false}>Cancel</button>
                    <button type="submit">Create</button>
                </div>
            </form>
        </div>
    </div>
{/if}

<div class="table-container">
    {#if data.users.length > 0}
        <table>
            <thead>
                <tr>
                    <th>Username</th>
                    <th>Team</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each data.users as user}
                    <tr>
                        <td><a href="/users/{user.username}">{user.username}</a></td>
                        <td><a href="/teams/{user.team}">{user.team}</a></td>
                        <td>
                            {#if user.username !== 'admin'}
                                <button onclick={() => handleDelete(user.username)} class="danger"><i class="fa-solid fa-trash"></i> Remove</button>
                            {/if}
                            <button>Edit</button>
                            <button>Update Password</button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:else}
        <p class="empty-msg">No users available yet.</p>
    {/if}
</div>

<style>
    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }

    .modal {
        background: white;
        padding: 20px;
        border-radius: 8px;
        min-width: 300px;
        color: black;
    }

    :global(html.dark) .modal {
        background: #2d2d2d;
        color: white;
    }

    form {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    label {
        display: flex;
        flex-direction: column;
    }

    .actions {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 10px;
    }

    input {
        padding: 4px;
    }

    :global(html.dark) input {
        background: #1a1a1a;
        color: white;
        border: 1px solid #444;
    }

    .empty-msg {
        text-align: center;
        padding: 4rem;
        color: #888;
    }
</style>
