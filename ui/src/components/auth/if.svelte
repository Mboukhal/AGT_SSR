<script>
    let { condition, children } = $props();
</script>

{`{{if ${condition}}}`}
{@render children()}

{"{{end}}"}
