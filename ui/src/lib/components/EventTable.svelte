<script lang="ts">
    import { getContrastColor } from '$lib/utils';

    let { 
        events = [], 
        colorize = false,
        showUser = true,
        showC2 = true,
        showScope = true,
        emptyMessage = "No events found."
    } = $props();
</script>

<div class="table-container">
    {#if events.length > 0}
        <table>
            <thead>
                <tr>
                    <th style="width: auto">Command</th>
                    {#if showUser}<th style="width: 120px">User</th>{/if}
                    {#if showC2}<th style="width: 120px">C2</th>{/if}
                    {#if showScope}<th style="width: 120px">Scope</th>{/if}
                    <th style="width: 180px">Time</th>
                </tr>
            </thead>
            <tbody>
                {#each events as event}
                    <tr style={colorize ? `background-color: #${event.team_color}; color: ${getContrastColor(event.team_color)}` : ''}>
                        <td><div class="cell-scroll">{event.command}</div></td>
                        {#if showUser}
                            <td><a href="/users/{event.user}">{event.user}</a></td>
                        {/if}
                        {#if showC2}
                            <td><a href="/c2s/{event.c2}">{event.c2}</a></td>
                        {/if}
                        {#if showScope}
                            <td><a href="/scope/{event.scope}">{event.scope}</a></td>
                        {/if}
                        <td>{new Date(event.time).toLocaleString()}</td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:else}
        <p class="empty-msg">{emptyMessage}</p>
    {/if}
</div>

<style>
    .empty-msg {
        text-align: center;
        padding: 4rem;
        color: #888;
    }

    tr[style*="background-color"] a {
        color: inherit;
    }
</style>
