<template>
  <div class="container fluid mt-5 col-md-6 offset-md-3">
    <h1 class="pb-4 mb-4 border-bottom">Verification</h1>
    <p>Check your email for verification code</p>

    <div v-text="error" v-if="error" class="alert alert-danger"></div>

    <div class="row">
      <div class="col-md-3">
        <label>Code</label>
      </div>

      <div class="form-group col col-md-9">
        <input v-model="code" type="text" class="form-control" />
      </div>
    </div>
    <div class="row">
      <div class="col-md-9 offset-md-3">
        <button :disabled="loading" @click="submit" class="btn btn-block btn-primary">Send</button>
      </div>
    </div>
  </div>
</template>
<script>
module.exports = {
  props: ["username", "password"],
  data() {
    return {
      loading: false,
      error: "",
      code: ""
    };
  },
  created() {
    if (!this.username) {
      this.$router.push("/login");
    }
  },
  methods: {
    async submit() {
      this.loading = true;
      this.error = "";
      try {
        const { code, username, password } = this;
        await axios.post("/api/confirm", {
          code,
          username
        });
        // ok, try to login
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