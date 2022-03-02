<template>
  <div id="app">
    <IndexPage />
  </div>
</template>

<script>
import IndexPage from "./components/IndexPage.vue";

export default {
  name: "App",
  components: {
    IndexPage,
  },
};
</script>

<style>
#app {
  font-family: Ubuntu, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  background-color: #fafafa;
}
</style>
