<template>
  <div class="login-container">
    <h2>Login</h2>
    <form @submit.prevent="loginUser">
      <div class="mb-3">
        <label for="username" class="form-label d-block">Username</label>
        <input 
          type="text" 
          id="username" 
          class="form-control-md" 
          v-model="username" 
          placeholder="Enter your username" 
          required
        />
      </div>
      <button type="submit" class="btn btn-primary">Login</button>
      <p v-if="msg" :class="{ 'text-success': success, 'text-danger': !success }">{{ msg }}</p>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      msg: "",
      securityKey: null,
      userId: null,
      success: false,
    };
  },
  created() {
    // Clear session data when accessing the login page
    localStorage.removeItem("securityKey");
    localStorage.removeItem("userId");
    console.log("User has been logged out.");
  },
  methods: {
    async loginUser() {
      this.success = false;
      this.msg = "";
      try {
        // Login to get the security key and user ID
        let response = await this.$axios.post("/session", { name: this.username });

        // Store security key and user ID
        this.securityKey = response.data.apiKey;
        this.userId = response.data.userId;

        // Store the security key and user ID locally
        localStorage.setItem("securityKey", this.securityKey);
        localStorage.setItem("userId", this.userId);
        console.log("User data:", this.securityKey, this.userId);

        // Fetch username using the user ID and security key
        let userResponse = await this.$axios.get(`/users/${this.userId}/username`, {
          headers: { Authorization: `Bearer ${this.securityKey}` }
        });

        this.username = userResponse.data.username;
        this.msg = `Logged in successfully. Welcome back ${this.username}`;
        this.success = true;
      } catch (e) {
        this.msg = "Login failed: " + (e.response?.data?.error || e.message);
      }
    }
  }
};
</script>

<style>
</style>