<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/state';
    import EventTable from '$lib/components/EventTable.svelte';
    import Pagination from '$lib/components/Pagination.svelte';

    const { data } = $props<{ data: { user: { username: string; team: string; email?: string }; events: Array<any>; limit: number; offset: number } }>();

    let colorize = $state(false);

    function handlePageChange(newOffset: number, newLimit: number) {
        const url = new URL(page.url);
        url.searchParams.set('offset', newOffset.toString());
        url.searchParams.set('limit', newLimit.toString());
        goto(url.toString());
    }
</script>

<title>Lober User {data.user.username}</title>

<h1>User: {data.user.username}</h1>

<div class="header-actions">
    <div class="search-group">
        <label class="checkbox-label">
            <input type="checkbox" bind:checked={colorize} />
            Colorize
        </label>
    </div>

    <Pagination 
        limit={data.limit} 
        offset={data.offset} 
        totalCount={data.events.length} 
        onchange={handlePageChange} 
    />
</div>

<EventTable events={data.events} {colorize} showUser={false} />

<style>
    .header-actions {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
        gap: 1rem;
    }

    .search-group {
        display: flex;
        align-items: center;
        gap: 1.5rem;
        flex-grow: 1;
    }

    h2 {
        margin: 0;
        font-size: 1.25rem;
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
