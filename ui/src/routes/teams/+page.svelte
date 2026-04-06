<script lang="ts">
    import { invalidateAll } from '$app/navigation';
    import { getContext } from 'svelte';

    const { data } = $props<{ data: { teams: Array<{ name: string; color: string; lead: string }> }}>();
    const getMasterToken = getContext<() => string>('masterToken');

    let showModal = $state(false);
    let formData = $state({ name: '', color: '000000' });

    async function handleSubmit(e: Event) {
        e.preventDefault();

        const response = await fetch('/api/teams', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getMasterToken()}`
            },
            body: JSON.stringify({
                ...formData,
                color: formData.color.replace('#', '')
            })
        });

        if (response.ok) {
            showModal = false;
            formData = { name: '', color: '000000' };
            await invalidateAll();
        }
    }

    async function handleDelete(name: string) {
        if (!confirm(`Are you sure you want to remove team "${name}"?`)) return;
        
        const response = await fetch(`/api/teams/${name}`, {
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

<title>Lober Teams</title>

<div class="header-container">
    <h1>Teams</h1>
    <button onclick={() => showModal = true}>New Team</button>
</div>

{#if showModal}
    <div class="modal-backdrop">
        <div class="modal">
            <h2>New Team</h2>
            <form onsubmit={handleSubmit}>
                <label>
                    Name:
                    <input type="text" bind:value={formData.name} required />
                </label>
                <label>
                    Color (Hex, e.g. FF0000):
                    <input type="text" bind:value={formData.color} required pattern="#?[0-9A-Fa-f]{6}" title="6-character hex color code (with optional #)" />
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
    {#if data.teams.length > 0}
        <table>
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Color</th>
                    <th>Lead</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each data.teams as team}
                    <tr>
                        <td><a href="/teams/{team.name}">{team.name}</a></td>
                        <td style="background-color: #{team.color};"></td>
                        <td><a href="/users/{team.lead}">{team.lead}</a></td>
                        <td>
                            <button onclick={() => handleDelete(team.name)} class="danger"><i class="fa-solid fa-trash"></i> Remove</button>
                            <button>Edit</button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:else}
        <p class="empty-msg">No teams available yet.</p>
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
