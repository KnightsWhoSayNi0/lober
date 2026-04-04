<script lang="ts">
    import { invalidateAll } from '$app/navigation';

    const { data } = $props<{ data: { scopes: Array<{ name: string }> }}>();

    let showModal = $state(false);
    let formData = $state({ name: '' });

    async function handleSubmit(e: Event) {
        e.preventDefault();

        const response = await fetch('/api/scope', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            showModal = false;
            formData = { name: '' };
            await invalidateAll();
        }
    }

    async function handleDelete(name: string) {
        if (!confirm(`Are you sure you want to remove scope "${name}"?`)) return;
        
        const response = await fetch(`/api/scope/${name}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            await invalidateAll();
        }
    }
</script>

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
                    <input type="text" bind:value={formData.name} required />
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
    {#if data.scopes.length > 0}
        <table>
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each data.scopes as scope}
                    <tr>
                        <td><a href="/scope/{scope.name}">{scope.name}</a></td>
                        <td>
                            <button onclick={() => handleDelete(scope.name)} class="danger"><i class="fa-solid fa-trash"></i> Remove</button>
                            <button>Edit</button>
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
