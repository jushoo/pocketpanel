<script lang="ts">
  let count = $state(0)

  function increment() {
    count += 1
  }
</script>

<button onclick={increment}>
  Count is {count}
</button>

<style>
  button {
    padding: 0.5rem 1rem;
    font-size: 1rem;
    cursor: pointer;
  }
</style>
