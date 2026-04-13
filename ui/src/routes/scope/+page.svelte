<script lang="ts">
    import { invalidateAll } from '$app/navigation';
    import { getContext } from 'svelte';

    const { data } = $props<{ data: { scopes: Array<{ name: string; ipaddr: string; team: number }> }}>();
    const getMasterToken = getContext<() => string>('masterToken');

    let showModal = $state(false);
    let formData = $state({ name: '', ipaddr: '', team: 1 });

    const groupedScopes = $derived(
        data.scopes.reduce((acc, scope) => {
            if (!acc[scope.name]) {
                acc[scope.name] = [];
            }
            acc[scope.name].push(scope);
            return acc;
        }, {} as Record<string, typeof data.scopes>)
    );

    async function handleSubmit(e: Event) {
        e.preventDefault();

        const response = await fetch('/api/scope', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getMasterToken()}`
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            showModal = false;
            formData = { name: '', ipaddr: '', team: 1 };
            await invalidateAll();
        }
    }

    async function handleDelete(name: string) {
        if (!confirm(`Are you sure you want to remove ALL scopes named "${name}"?`)) return;
        
        const response = await fetch(`/api/scope/${name}`, {
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

<title>Lober Scope</title>

<div class="header-container">
    <h1>Scope</h1>
    <button onclick={() => showModal = true}>New Scope</button>
</div>

{#if showModal}
    <div class="modal-backdrop">
        <div class="modal">
            <h2>New Scope</h2>
            <form onsubmit={handleSubmit}>
                <label>
                    Name:
                    <input type="text" bind:value={formData.name} placeholder="e.g. Windows 1" required />
                </label>
                <label>
                    IP Address:
                    <input type="text" bind:value={formData.ipaddr} placeholder="e.g. 10.1.1.70" required />
                </label>
                <label>
                    Team:
                    <input type="number" bind:value={formData.team} min="1" required />
                </label>
                <div class="actions">
                    <button type="button" onclick={() => showModal = false}>Cancel</button>
                    <button type="submit">Create</button>
                </div>
            </form>
        </div>
    </div>
{/if}

<div class="table-scroll-container">
    {#if Object.keys(groupedScopes).length > 0}
        <table>
            <thead>
                <tr>
                    <th>Box Name</th>
                    <th>Hosts (IPs)</th>
                    <th>Teams</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each Object.entries(groupedScopes) as [name, hosts]}
                    <tr>
                        <td class="box-name"><strong>{name}</strong></td>
                        <td>
                            <div class="ip-list">
                                {#each hosts as host}
                                    <code>{host.ipaddr}</code>
                                {/each}
                            </div>
                        </td>
                        <td>
                            <div class="team-list">
                                {#each [...new Set(hosts.map(h => h.team))] as team}
                                    <span class="team-badge">Team {team}</span>
                                {/each}
                            </div>
                        </td>
                        <td class="actions-cell">
                            <a href="/scope/{name}" class="btn-link">View Events</a>
                            <button onclick={() => handleDelete(name)} class="danger sm"><i class="fa-solid fa-trash"></i> Delete Box</button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:else}
        <p class="empty-msg">No scopes available yet.</p>
    {/if}
</div>

<style>
    .table-scroll-container {
        max-height: calc(100vh - 250px);
        overflow-y: auto;
        background: #fff;
        border: 1px solid #ddd;
        border-radius: 8px;
    }

    :global(html.dark) .table-scroll-container {
        background: #252525;
        border-color: #333;
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    th {
        position: sticky;
        top: 0;
        background: #f8f9fa;
        z-index: 10;
        text-align: left;
        padding: 1rem;
        font-size: 0.85rem;
        text-transform: uppercase;
        color: #666;
        border-bottom: 2px solid #eee;
    }

    :global(html.dark) th {
        background: #1a1a1a;
        color: #aaa;
        border-bottom-color: #333;
    }

    td {
        padding: 1rem;
        border-bottom: 1px solid #eee;
        vertical-align: middle;
    }

    :global(html.dark) td {
        border-bottom-color: #333;
    }

    .box-name {
        font-size: 1.1rem;
    }

    .ip-list {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
    }

    .team-list {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
    }

    .team-badge {
        background: #e9ecef;
        padding: 0.2rem 0.5rem;
        border-radius: 4px;
        font-size: 0.8rem;
        font-weight: 600;
    }

    :global(html.dark) .team-badge {
        background: #333;
        color: #ddd;
    }

    .btn-link {
        display: inline-flex;
        align-items: center;
        padding: 0.4rem 0.75rem;
        background: #eee;
        color: #333;
        text-decoration: none;
        border-radius: 4px;
        font-size: 0.85rem;
        font-weight: 600;
    }

    :global(html.dark) .btn-link {
        background: #333;
        color: #ccc;
    }

    .actions-cell {
        display: flex;
        gap: 0.75rem;
        align-items: center;
    }

    code {
        font-family: "IBM Plex Mono", monospace;
        background: #f4f4f4;
        padding: 0.2rem 0.4rem;
        border-radius: 3px;
        font-size: 0.9rem;
    }

    button.danger.sm {
        padding: 0.4rem 0.75rem;
        font-size: 0.85rem;
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
    }

    input {
        padding: 0.6rem;
        border: 1px solid #ddd;
        border-radius: 4px;
    }

    :global(html.dark) input {
        background: #1a1a1a;
        color: white;
        border-color: #444;
    }

    .actions {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 0.5rem;
    }

    .empty-msg {
        text-align: center;
        padding: 4rem;
        color: #888;
    }
</style>
