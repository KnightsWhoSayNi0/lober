<script lang="ts">
    const { data } = $props();
</script>

<h1>Users / {data.user.username}</h1>

<ul>
    <li>Email: <a href="mailto:{data.user.email}">{data.user.email}</a></li>
    <li>Team: {data.user.team}</li>
</ul>