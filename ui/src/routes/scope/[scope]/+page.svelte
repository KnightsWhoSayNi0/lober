<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/state';
    import EventTable from '$lib/components/EventTable.svelte';
    import Pagination from '$lib/components/Pagination.svelte';

    const { data } = $props<{ data: { events: Array<any>; scope: string; limit: number; offset: number } }>();

    let colorize = $state(false);

    function handlePageChange(newOffset: number, newLimit: number) {
        const url = new URL(page.url);
        url.searchParams.set('offset', newOffset.toString());
        url.searchParams.set('limit', newLimit.toString());
        goto(url.toString());
    }
</script>

<h1>Events for Scope: {data.scope}</h1>

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

<EventTable events={data.events} {colorize} showScope={false} />

<style>
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
