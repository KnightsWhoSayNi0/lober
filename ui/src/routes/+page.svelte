<script lang="ts">
    const { data } = $props<{ data: { events: Array<{ ID: number; command: string; user: string; c2: string; scope: string; time: string }> }}>();
</script>

<h1>events</h1>

{#if data.events.length > 0}
    <table>
        <thead>
            <tr>
                <th>command</th>
                <th>user</th>
                <th>c2</th>
                <th>scope</th>
                <th>time</th>
            </tr>
        </thead>
        <tbody>
            {#each data.events as event}
                <tr>
                    <td>{event.command}</td>
                    <td>{event.user}</td>
                    <td>{event.c2}</td>
                    <td>{event.scope}</td>
                    <td>{new Date(event.time).toTimeString().substring(0, 8)}</td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <p>No events available yet.</p>
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