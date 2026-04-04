<script lang="ts">
    import { invalidateAll } from '$app/navigation';

    const { data } = $props<{ data: { tokens: Array<any>, users: Array<any>, c2s: Array<any> }}>();

    let showModal = $state(false);
    let formData = $state({ username: '', c2: '', expires_days: 7 });
    let lastGeneratedToken = $state('');

    async function handleSubmit(e: Event) {
        e.preventDefault();

        const response = await fetch('/api/tokens', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            const result = await response.json();
            lastGeneratedToken = result.token; // Show the raw token once
            showModal = false;
            formData = { username: '', c2: '', expires_days: 7 };
            await invalidateAll();
        }
    }

    async function handleDelete(id: number) {
        if (!confirm('Remove this token?')) return;
        const res = await fetch(`/api/tokens/${id}`, { method: 'DELETE' });
        if (res.ok) await invalidateAll();
    }
</script>

<title>Lober Tokens</title>

<div class="header-container">
    <h1>Authentication Tokens</h1>
    <button onclick={() => { showModal = true; lastGeneratedToken = ''; }}>New Token</button>
</div>

{#if lastGeneratedToken}
    <div class="token-reveal">
        <p><strong>Token Generated!</strong> Copy this now, it won't be shown again:</p>
        <code>{lastGeneratedToken}</code>
        <button onclick={() => lastGeneratedToken = ''}>Dismiss</button>
    </div>
{/if}

<div class="table-container">
    {#if data.tokens.length > 0}
        <table>
            <thead>
                <tr>
                    <th>User</th>
                    <th>C2</th>
                    <th>Created</th>
                    <th>Expires</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each data.tokens as token}
                    <tr>
                        <td>{token.username}</td>
                        <td>{token.c2}</td>
                        <td>{new Date(token.created).toLocaleString()}</td>
                        <td>{new Date(token.expires).toLocaleString()}</td>
                        <td>
                            <button class="danger" onclick={() => handleDelete(token.id)}>Remove</button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:else}
        <p class="empty-msg">No active tokens found.</p>
    {/if}
</div>

{#if showModal}
    <div class="modal-backdrop">
        <div class="modal">
            <h2>Generate New Token</h2>
            <form onsubmit={handleSubmit}>
                <label>
                    User:
                    <select bind:value={formData.username} required>
                        <option value="">Select User</option>
                        {#each data.users as user}
                            <option value={user.username}>{user.username}</option>
                        {/each}
                    </select>
                </label>
                <label>
                    C2:
                    <select bind:value={formData.c2} required>
                        <option value="">Select C2</option>
                        {#each data.c2s as c2}
                            <option value={c2.name}>{c2.name}</option>
                        {/each}
                    </select>
                </label>
                <label>
                    Expires in (days):
                    <input type="number" bind:value={formData.expires_days} min="1" max="365" />
                </label>
                <div class="actions">
                    <button type="button" onclick={() => showModal = false}>Cancel</button>
                    <button type="submit">Generate</button>
                </div>
            </form>
        </div>
    </div>
{/if}

<style>
    .header-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
    }

    .token-reveal {
        background: #fff3cd;
        border: 1px solid #ffeeba;
        padding: 1.5rem;
        border-radius: 8px;
        margin-bottom: 1.5rem;
        color: #856404;
    }

    :global(html.dark) .token-reveal {
        background: #332701;
        border-color: #4d3d02;
        color: #ffda6a;
    }

    code {
        display: block;
        background: #fff;
        padding: 1rem;
        border-radius: 4px;
        margin: 1rem 0;
        font-family: "IBM Plex Mono", monospace;
        font-size: 1.1rem;
        word-break: break-all;
    }

    :global(html.dark) code {
        background: #000;
    }

    .empty-msg {
        text-align: center;
        padding: 4rem;
        color: #888;
    }

    /* Modal Styles (Consistent with Teams/Users) */
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
        padding: 2rem;
        border-radius: 8px;
        min-width: 350px;
        color: black;
    }

    :global(html.dark) .modal {
        background: #2d2d2d;
        color: white;
    }

    form {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
    }

    label {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        font-weight: 600;
        font-size: 0.9rem;
    }

    select, input {
        padding: 0.6rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        background: #fff;
        color: inherit;
    }

    :global(html.dark) select, :global(html.dark) input {
        background: #1a1a1a;
        border-color: #444;
    }

    .actions {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 0.5rem;
    }
</style>
