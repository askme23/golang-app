<template>
  <div class="d-flex justify-content-center">
    <b-form id="login-form" @submit="onSubmit" v-if="show">
      <b-form-group
        id="input-group-1"
        label="Email address:"
        label-for="input-1"
        description="We'll never share your email with anyone else."
      >
        <b-form-input
          id="input-1"
          v-model="form.email"
          type="email"
          required
          placeholder="Enter email"
        ></b-form-input>
      </b-form-group>

      <b-form-group id="input-group-2" label="Password:" label-for="input-2">
        <b-input
          v-model="form.password"
          type="password"
          id="text-password"
          placeholder="Enter password"
        ></b-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Log In</b-button>

      <br><b-link :to="'sign-in'">Sign In</b-link>
    </b-form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Auth',
  data() {
    return {
      form: {
        email: '',
        password: '',
      },
      show: true,
    };
  },
  methods: {
    onSubmit(evt) {
      evt.preventDefault();

      const path = 'http://localhost:8081/login';
      const self = this;
      axios({
        url: path,
        method: 'POST',
        withCredentials: true,
        data: {
          email: this.form.email,
          password: this.form.password,
        },
      })
        .then((response) => {
          console.log(response);
          self.$router.push({ name: 'Home'})
        })
        .catch((error) => {
          console.log(error);
        });
    },
    // onReset(evt) {
    //   evt.preventDefault();
    //   this.form.email = '';
    //   this.form.password = '';
    //   this.show = false;
    //   this.$nextTick(() => {
    //     this.show = true;
    //   });
    // },
  },
};
</script>

<style scoped>
  #login-form {
    width: 400px;
    background-color: #fff;
    padding: 25px;
    box-shadow: 0 5px 30px 0px rgba(0, 0, 0, 0.1);
    border-radius: 10px;
  }
</style>
