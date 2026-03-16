<script lang="ts">
    const { data } = $props<{ data: { scopes: Array<{ name: string }> }}>();
</script>

<h1>Scope</h1>

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
                    <td>{scope.name}</td>
                    <td>
                        <button>Edit</button>
                        <button>Show Logs</button>
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <p>No scopes available yet.</p>
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