<template>
  <div class="hello">
    <h1 v-if="user">Welcome {{ user }} to Your App.</h1>
    <h1 v-else>You need to authenticate, please click to <a href="/login">Login</a></h1>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Home',
  data() {
    return {
      user: null,
    };
  },
  methods: {
    authenticate() {
      const path = 'http://localhost:8081';
      const self = this;
      axios({
        url: path,
        method: 'GET',
        // Инача куки не отправляются на другой домен и не выставляется на клиенте
        withCredentials: true,
      })
        .then((response) => {
          console.log(response);
          self.user = response.data;
        })
        .catch((error) => {
          self.user = null;
          console.log(error);
        });
    },
  },
  created() {
    this.authenticate();
  },
};
</script>

<style scoped>
  h1, h2 {
    font-weight: normal;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    display: inline-block;
    margin: 0 10px;
  }
  a {
    color: #42b983;
  }
</style>
