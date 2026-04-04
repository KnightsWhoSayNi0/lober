<script lang="ts">
    import { onMount } from 'svelte';
    import EventTable from '$lib/components/EventTable.svelte';
    import Pagination from '$lib/components/Pagination.svelte';

    const { data } = $props<{ data: { events: Array<any> } }>();
    
    let events = $state(data.events);
    let filter = $state("");
    let socket: WebSocket | null = $state(null);
    let colorize = $state(false);
    
    let pageSize = $state(50);
    let offset = $state(0);

    onMount(() => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.host;
        const s = new WebSocket(`${protocol}//${host}/api/ws`);

        s.onopen = () => {
            s.send(filter);
            socket = s;
        };

        s.onmessage = (event) => {
            const receivedData = JSON.parse(event.data);
            if (Array.isArray(receivedData)) {
                events = receivedData;
            } else {
                if (offset === 0) {
                    events = [receivedData, ...events];
                    if (events.length > pageSize) events = events.slice(0, pageSize);
                }
            }
        };

        return () => s.close();
    });

    $effect(() => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(filter);
        }
    });

    async function fetchPage(newOffset: number, newLimit: number) {
        const res = await fetch(`/api/events?filter=${filter}&limit=${newLimit}&offset=${newOffset}`);
        if (res.ok) {
            events = await res.json();
        }
    }

    function handleFilter(e: Event) {
        filter = (e.target as HTMLInputElement).value;
        offset = 0;
        fetchPage(0, pageSize);
    }
</script>

<title>Lober</title>

<h1>Events</h1>

<div class="header-actions">
    <div class="search-group">
        <input 
            type="text" 
            placeholder="Search events..." 
            value={filter}
            oninput={handleFilter}
            class="search-input"
        />
        <label class="checkbox-label">
            <input type="checkbox" bind:checked={colorize} />
            Colorize
        </label>
    </div>

    <Pagination 
        bind:limit={pageSize} 
        bind:offset={offset} 
        totalCount={events.length} 
        onchange={fetchPage} 
    />
</div>

<EventTable {events} {colorize} />

<style>
    .search-input {
        padding: 0.6rem 1rem;
        border: 1px solid #ddd;
        border-radius: 6px;
        width: 300px;
        background: #fff;
    }

    :global(html.dark) .search-input {
        background: #2d2d2d;
        border-color: #444;
        color: #fff;
    }

    .search-group {
        display: flex;
        align-items: center;
        gap: 1.5rem;
    }

    .checkbox-label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9rem;
        cursor: pointer;
        user-select: none;
    }
</style>
