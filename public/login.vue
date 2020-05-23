<template>
  <div class="container fluid mt-5 col-md-6 offset-md-3">
    <h1 class="pb-4 mb-4 border-bottom">Login</h1>

    <div v-text="error" v-if="error" class="alert alert-danger"></div>

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

    <div class="row my-3">
      <div class="col-md-9 offset-md-3">
        <span>Not have an account yet?</span>
        <a href="#/register">Register</a>
      </div>
    </div>

    <div class="row">
      <div class="col-md-9 offset-md-3">
        <button @click="submit" :disabled="loading" class="btn btn-block btn-primary">Login</button>
      </div>
    </div>
  </div>
</template>
<script>
module.exports = {
  data() {
    return {
      username: "",
      password: "",
      error: "",
      loading: false
    };
  },
  methods: {
    async submit() {
      const { username, password } = this;
      this.loading = true;
      try {
        this.error = "";
        const { data } = await axios.post("/api/login", {
          username,
          password
        });
        window.localStorage.setItem("auth", JSON.stringify(data));
        this.$router.push("/profile");
      } catch (e) {
        this.error = e.response.data.error;
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>