<script lang="ts">
    let { limit = $bindable(50), offset = $bindable(0), totalCount = 0, onchange } = $props();

    function changePage(newOffset: number) {
        offset = newOffset;
        if (onchange) onchange(offset, limit);
    }

    function handleLimitChange(e: Event) {
        const val = parseInt((e.target as HTMLInputElement).value);
        if (!isNaN(val) && val > 0) {
            limit = val;
            offset = 0;
            if (onchange) onchange(offset, limit);
        }
    }

    let currentPage = $derived(Math.floor(offset / limit) + 1);
</script>

<div class="pagination-group">
    <div class="page-size">
        Rows:
        <input type="number" value={limit} onchange={handleLimitChange} min="1" step="1" class="rows-input">
    </div>

    <div class="nav-buttons">
        <button onclick={() => changePage(Math.max(0, offset - limit))} disabled={offset === 0}>Prev</button>
        <span class="page-num">Pg {currentPage}</span>
        <button onclick={() => changePage(offset + limit)} disabled={totalCount > 0 && totalCount < limit}>Next</button>
    </div>
</div>

<style>
    .pagination-group {
        display: flex;
        align-items: center;
        gap: 1.5rem;
        margin-left: auto;
    }

    .page-size {
        font-size: 0.9rem;
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .rows-input {
        width: 80px;
        padding: 0.4rem 0.6rem;
        border: 1px solid #ddd;
        border-radius: 6px;
        background: #fff;
    }

    :global(html.dark) .rows-input {
        background: #2d2d2d;
        border-color: #444;
        color: #fff;
    }

    .nav-buttons {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .page-num {
        font-size: 0.9rem;
        font-weight: 600;
        min-width: 40px;
        text-align: center;
    }

    button {
        padding: 0.4rem 0.8rem;
    }
</style>
