<script lang="ts">
    const { data } = $props<{ data: { c2s: Array<{ name: string }> }}>();
</script>

<h1>C2s</h1>

{#if data.c2s.length > 0}
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            {#each data.c2s as c2}
                <tr>
                    <td>{c2.name}</td>
                    <td>
                        <button>Edit</button>
                        <button>Show Logs</button>
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <p>No c2s available yet.</p>
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