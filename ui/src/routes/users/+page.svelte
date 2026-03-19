<script lang="ts">
    const { data } = $props<{ data: { users: Array<{ username: string; email: string; team: string }> }}>();

    let showModal = $state(false);
    let formData = $state({ username: '', email: '', team: '' });

    async function handleSubmit(e: Event) {
        e.preventDefault();
        
        const response = await fetch('/api/users', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            showModal = false;
            formData = { username: '', email: '', team: '' };
        }
    }
</script>

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
                    Email:
                    <input type="email" bind:value={formData.email} required />
                </label>
                <label>
                    Team:
                    <input type="text" bind:value={formData.team} required />
                </label>
                <div class="actions">
                    <button type="button" onclick={() => showModal = false}>Cancel</button>
                    <button type="submit">Create</button>
                </div>
            </form>
        </div>
    </div>
{/if}

{#if data.users.length > 0}
    <table>
        <thead>
            <tr>
                <th>Username</th>
                <th>Email</th>
                <th>Team</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            {#each data.users as user}
                <tr>
                    <td><a href="/users/{user.username}">{user.username}</a></td>
                    <td>{user.email}</td>
                    <td>{user.team}</td>
                    <td>
                        <button>Edit</button>
                        <button>Show Logs</button>
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <p>No users available yet.</p>
{/if}

<style>
    .header-container {
        display: flex;
        align-items: center;
        gap: 20px;
        margin-bottom: 20px;
    }

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
    }

    .modal {
        background: white;
        padding: 20px;
        border-radius: 8px;
        min-width: 300px;
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