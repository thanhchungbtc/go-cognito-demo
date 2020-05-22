<template>
  <div class="container fluid mt-5 col-md-6 offset-md-3">
    <h1 class="pb-4 mb-4 border-bottom">Register</h1>

    <div v-text="error" v-if="error" class="alert alert-danger"></div>

    <div class="row">
      <div class="col-md-3">
        <label>Email</label>
      </div>

      <div class="form-group col col-md-9">
        <input v-model="email" type="text" class="form-control" />
      </div>
    </div>

    <div class="row">
      <div class="col-md-3">
        <label>Username</label>
      </div>

      <div class="form-group col col-md-9">
        <input v-model="username" type="text" class="form-control" />
      </div>
    </div>

    <div class="row">
      <div class="col-md-3">
        <label>Password</label>
      </div>
      <div class="form-group col col-md-9">
        <input v-model="password" type="password" class="form-control" />
      </div>
    </div>

    <div class="row">
      <div class="col-md-3">
        <label>Confirm</label>
      </div>
      <div class="form-group col col-md-9">
        <input v-model="confirmPassword" type="password" class="form-control" />
      </div>
    </div>

    <div class="row my-3">
      <div class="col-md-9 offset-md-3">
        <span>Already have an account?</span>
        <a href="#/login">Login</a>
      </div>
    </div>

    <div class="row">
      <div class="col-md-9 offset-md-3">
        <button :disabled="loading" @click="submit" class="btn btn-block btn-primary">Register</button>
      </div>
    </div>
  </div>
</template>
<script>
module.exports = {
  data() {
    return {
      email: "",
      username: "",
      password: "",
      confirmPassword: "",
      error: "",
      loading: false
    };
  },
  methods: {
    async submit() {
      this.loading = true;
      this.error = "";
      try {
        const { email, username, password, confirmPassword } = this;
        if (password !== confirmPassword) {
          this.error = "password did not match";
          return;
        }
        const { data } = await axios.post("/api/register", {
          email,
          username,
          password
        });
        console.log(data);
        this.$router.push({
          name: "confirm_auth",
          params: { username, password }
        });
      } catch (e) {
        this.error = e.response.data.error;
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>