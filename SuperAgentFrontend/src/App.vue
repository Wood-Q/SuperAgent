<template>
  <div id="app">
    <ChatRoom />
  </div>
</template>

<script>
import ChatRoom from './components/ChatRoom.vue'

export default {
  name: 'App',
  components: {
    ChatRoom
  }
}
</script>

<style scoped>
#app {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
</style> 