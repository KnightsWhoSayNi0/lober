<script lang="ts">
    const { data } = $props<{ data: { teams: Array<{ name: string; color: string; lead: string }> }}>();
</script>

<h1>Teams</h1>

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
                    <td>{team.name}</td>
                    <td style="background-color: #{team.color};"></td>
                    <td><a href="/users/{team.lead}">{team.lead}</a></td>
                    <td>
                        <button>Edit</button>
                        <button>Show Logs</button>
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <p>No teams available yet.</p>
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