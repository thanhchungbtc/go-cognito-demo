<template>
  <div class="container fluid mt-5 col-md-6 offset-md-3">
    <h1 class="pb-4 mb-4 border-bottom">Profile</h1>
    <div class="row">
      <div class="col-12 d-flex justify-content-between">
        <div v-if="user">
          Hello,
          <span class="font-weight-bold">{{ user.user_name }}</span>
        </div>
        <a href="#" @click="logout">Logout</a>
      </div>
    </div>
  </div>
</template>
<script>
// interceptor to automate refresh token
const createInterceptor = () => {
  const interceptor = axios.interceptors.response.use(
    res => res,
    err => {
      // not an unauthorized request, continue with error
      if (err.response.status !== 401) {
        return Promise.reject(err);
      }

      // get refresh token
      const data = window.localStorage.getItem("auth");
      const obj = JSON.parse(data);
      const { refresh_token } = obj;
      if (!data) {
        return Promise.reject(err);
      }
      // eject interceptor to prevent loop in case refresh token itself returns 401
      axios.interceptors.response.eject(interceptor);
      return axios
        .post("/api/login", {
          refresh_token
        })
        .then(({ data }) => {
          // update the authenticated token
          window.localStorage.setItem("auth", JSON.stringify(data));
        })
        .catch(err => {
          // remove the invalid authorization info
          window.localStorage.removeItem("auth");
          return Promise.reject(err);
        })
        .finally(createInterceptor);
    }
  );
};

createInterceptor();

module.exports = {
  data() {
    return {
      auth: null,
      user: null
    };
  },
  async created() {
    const data = window.localStorage.getItem("auth");
    if (!data) {
      this.$router.push("/login");
      return;
    }
    const obj = JSON.parse(data);
    this.auth = obj;
    try {
      const { data } = await axios.get("/api/me", {
        headers: { Authorization: `Bearer ${obj.id_token}` }
      });
      this.user = data;
    } catch (e) {}
  },
  methods: {
    async logout() {
      try {
        const { data } = await axios.get("/api/logout", {
          headers: { Authorization: `Bearer ${this.auth.access_token}` }
        });
        window.localStorage.removeItem("auth");
        this.$router.push("/login");
      } catch (e) {
        console.error(e);
      }
    }
  }
};
</script>