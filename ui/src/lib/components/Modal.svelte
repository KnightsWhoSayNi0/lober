<script lang="ts">
    let { show = $bindable(false), title, children } = $props();
</script>

{#if show}
    <div class="modal-backdrop" onclick={() => show = false}>
        <div class="modal" onclick={(e) => e.stopPropagation()}>
            <h2>{title}</h2>
            {@render children()}
        </div>
    </div>
{/if}

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
        padding: 2rem;
        border-radius: 8px;
        min-width: 350px;
        color: black;
        box-shadow: 0 4px 20px rgba(0,0,0,0.2);
    }

    :global(html.dark) .modal {
        background: #2d2d2d;
        color: white;
    }

    :global(form) {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
    }

    :global(label) {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        font-weight: 600;
        font-size: 0.9rem;
    }

    :global(.modal-actions) {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 0.5rem;
    }
</style>
